package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/view_db.go <db_path>")
		fmt.Println("Example: go run cmd/view_db.go url_fetcher.db")
		os.Exit(1)
	}

	dbPath := os.Args[1]
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// List all tables
	fmt.Println("=== Tables ===")
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		fmt.Printf("Error querying tables: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Printf("Error scanning: %v\n", err)
			continue
		}
		fmt.Printf("  - %s\n", tableName)
	}

	// Show records from url_records table
	fmt.Println("\n=== URL Records ===")
	recordRows, err := db.Query(`
		SELECT id, url, LENGTH(content) as content_length, created_at 
		FROM url_records 
		ORDER BY created_at DESC
	`)
	if err != nil {
		fmt.Printf("Error querying records: %v\n", err)
		return
	}
	defer recordRows.Close()

	fmt.Printf("%-5s %-50s %-15s %-25s\n", "ID", "URL", "Content Length", "Created At")
	fmt.Println(strings.Repeat("-", 100))

	for recordRows.Next() {
		var id int64
		var url string
		var contentLength int
		var createdAt string

		if err := recordRows.Scan(&id, &url, &contentLength, &createdAt); err != nil {
			fmt.Printf("Error scanning record: %v\n", err)
			continue
		}

		// Truncate URL if too long
		if len(url) > 47 {
			url = url[:44] + "..."
		}

		fmt.Printf("%-5d %-50s %-15d %-25s\n", id, url, contentLength, createdAt)
	}

	// Count records
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM url_records").Scan(&count)
	if err == nil {
		fmt.Printf("\nTotal records: %d\n", count)
	}
}

