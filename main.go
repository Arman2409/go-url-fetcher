package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"main.go/constants"
	"main.go/database"
	"main.go/services"
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

	// 	fmt.Println("Database initialized successfully")

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	urlFetcher := services.NewUrlFetcher(httpClient)

	var wg sync.WaitGroup

	for _, url := range constants.URLS_TO_FETCH {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Printf("Fetching URL: %s\n", url)

			content, err := urlFetcher.FetchUrl(url)
			if err != nil {
				fmt.Printf("Error fetching the URL: %v\n", err)
				return
			}

			urlProcessor := services.NewUrlResponseProcessor(repo)
			err = urlProcessor.ProcessAndStoreUrlResponse(url, content)

			if err != nil {
				fmt.Printf("Error processing and storing URL response: %v\n", err)
				return
			}

			fmt.Printf("Successfully processed and stored URL: %s\n", url)
		}(url)
	}

	wg.Wait()

	// Get all records
	fmt.Println("All records in database:")
	allRecords, err := repo.GetAllUrlRecords()
	if err != nil {
		fmt.Printf("Error retrieving all records: %v\n", err)
		return
	}

	for _, rec := range allRecords {
		fmt.Printf("  ID: %d, URL: %s, Content Length: %d bytes, Created: %s\n",
			rec.ID, rec.URL, len(rec.Content), rec.CreatedAt)
	}
}
