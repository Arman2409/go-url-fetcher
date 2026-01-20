package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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

	// Fetch URL
	url := "https://www.google.com"
	urlFetcher := services.NewUrlFetcher(httpClient)

	fmt.Printf("Fetching URL: %s\n", url)
	content, err := urlFetcher.FetchUrl(url)
	if err != nil {
		fmt.Printf("Error fetching the URL: %v\n", err)
		return
	}

	fmt.Printf("Fetched content length: %d bytes\n", len(content))

	// Save to database
	fmt.Println("Saving to database...")
	record, err := repo.CreateUrlRecord(url, content)
	if err != nil {
		fmt.Printf("Error saving to database: %v\n", err)
		fmt.Println("Continuing with the error...")
	}

	if record != nil {
		fmt.Printf("Saved record with ID: %d\n", record.ID)

		// Retrieve from database
		fmt.Println("Retrieving from database...")
		retrieved, err := repo.GetUrlRecordByID(record.ID)
		if err != nil {
			fmt.Printf("Error retrieving record: %v\n", err)
			return
		}

		fmt.Printf("Retrieved record - ID: %d, URL: %s, Created: %s\n",
			retrieved.ID, retrieved.URL, retrieved.CreatedAt)
	}

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
