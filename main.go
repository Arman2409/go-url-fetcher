package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"main.go/constants"
	"main.go/database"
	"main.go/services"
	appTypes "main.go/types"
)

func main() {
	fmt.Println("Starting the URL fetcher with SQLite database")

	// Initialize SQLite database
	dbPath := "url_fetcher.db"
	repo, err := database.NewSQLiteRepository(dbPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer repo.Close()

	fmt.Println("Database initialized successfully!")

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	urlFetcher := services.NewUrlFetcher(httpClient)
	urlProcessor := services.NewUrlResponseProcessor(repo)

	var wg sync.WaitGroup
	jobs := make(chan string)
	results := make(chan appTypes.FetchResult)
	workerCount := 3

	var collectorWG sync.WaitGroup
	statsRunner := services.NewStatsRunner()

	collectorWG.Add(1)
	go func() {
		defer collectorWG.Done()
		statsRunner.Consume(results)
	}()

	for i := range workerCount {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for url := range jobs {
				fmt.Printf("Worker %d fetching URL: %s\n", workerID, url)
				start := time.Now()

				// Allow time for multiple fetch attempts (retries + http.Client timeout per attempt).
				ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
				content, statusCode, err := urlFetcher.FetchUrlWithContext(ctx, url)
				if err != nil {
					fmt.Printf("Error fetching the URL: %v\n", err)
					cancel()
					results <- appTypes.FetchResult{
						URL:        url,
						StatusCode: statusCode,
						Duration:   time.Since(start),
						FetchErr:   err,
					}
					continue
				}

				err = urlProcessor.ProcessAndStoreUrlResponseWithContext(ctx, url, statusCode, content)
				cancel()
				if err != nil {
					fmt.Printf("Error processing and storing URL response: %v\n", err)
					results <- appTypes.FetchResult{
						URL:        url,
						StatusCode: statusCode,
						Duration:   time.Since(start),
						StoreErr:   err,
					}
					continue
				}

				fmt.Printf("Successfully processed and stored URL: %s\n", url)
				results <- appTypes.FetchResult{
					URL:        url,
					StatusCode: statusCode,
					Duration:   time.Since(start),
					Stored:     true,
				}
			}
		}(i + 1)
	}

	uniqueUrls := make(map[string]bool)

	for _, url := range constants.URLS_TO_FETCH {
		if _, exists := uniqueUrls[url]; exists {
			fmt.Printf("Skipping duplicate URL: %s\n", url)
			results <- appTypes.FetchResult{
				URL:          url,
				WasDuplicate: true,
			}
			continue
		}
		uniqueUrls[url] = true
		jobs <- url
	}
	close(jobs)

	wg.Wait()
	close(results)
	collectorWG.Wait()
	stats := statsRunner.Stats()

	fmt.Printf("Run summary: queued=%d stored=%d duplicates=%d fetch_failed=%d store_failed=%d non_2xx=%d avg_duration=%s\n",
		stats.Queued,
		stats.StoredSuccess,
		stats.Duplicates,
		stats.FetchFailed,
		stats.StoreFailed,
		stats.Non2xxFailures,
		stats.AvgDuration(),
	)

	// Get all records
	fmt.Println("All records in database:")
	allRecords, err := repo.GetAllUrlRecords()
	if err != nil {
		fmt.Printf("Error retrieving all records: %v\n", err)
		return
	}

	for _, rec := range allRecords {
		fmt.Printf("  ID: %d, URL: %s, Status: %d, Content Length: %d bytes, Created: %s\n",
			rec.ID, rec.URL, rec.StatusCode, len(rec.Content), rec.CreatedAt)
	}
}
