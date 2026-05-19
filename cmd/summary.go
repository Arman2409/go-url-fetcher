package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type statusStat struct {
	code  int
	count int
}

func main() {
	dbPath := flag.String("db", "url_fetcher.db", "path to the SQLite database file")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var total int
	var avgContentSize float64
	var totalContentSize int64

	err = db.QueryRow("SELECT COUNT(*), COALESCE(SUM(LENGTH(content)), 0) FROM url_records").
		Scan(&total, &totalContentSize)
	if err != nil {
		fmt.Printf("Error querying summary: %v\n", err)
		os.Exit(1)
	}

	if total > 0 {
		avgContentSize = float64(totalContentSize) / float64(total)
	}

	rows, err := db.Query(`
		SELECT status_code, COUNT(*) as count
		FROM url_records
		GROUP BY status_code
		ORDER BY count DESC
	`)
	if err != nil {
		fmt.Printf("Error querying status codes: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	var statusStats []statusStat
	for rows.Next() {
		var s statusStat
		if err := rows.Scan(&s.code, &s.count); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}
		statusStats = append(statusStats, s)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterating rows: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Database Summary ===")
	fmt.Printf("Total URLs stored : %d\n", total)
	fmt.Printf("Avg content size  : %.0f bytes\n", avgContentSize)
	fmt.Printf("Total content size: %d bytes\n", totalContentSize)

	fmt.Println("\n=== Status Code Breakdown ===")
	for _, s := range statusStats {
		fmt.Printf("  %d : %d URL(s)\n", s.code, s.count)
	}
}
