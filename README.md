# go-url-fetcher

A concurrent URL fetcher written in Go as a learning project. Fetches a list of URLs in parallel, stores the responses in a local SQLite database, and prints a run summary at the end.

## Stack

- **Language:** Go
- **Database:** SQLite (via `github.com/mattn/go-sqlite3`)
- **Concurrency:** goroutines, channels, `sync.WaitGroup`

## How it works

1. A configurable worker pool (3 workers by default) picks URLs from a jobs channel.
2. Each worker fetches a URL using an HTTP client with a 10s timeout.
3. Failed requests are retried up to 4 times with exponential backoff (capped at 2s).
4. Successful responses are stored in SQLite (`url_records` table).
5. A separate collector goroutine consumes results and tracks run statistics.
6. Duplicate URLs in the input list are detected and skipped before hitting the network.

## Features

- Worker pool with configurable concurrency
- Retry logic with exponential backoff
- Non-2xx responses treated as errors (with retries for 429 and 5xx)
- Duplicate URL detection
- Per-run stats: queued, stored, duplicates, fetch failures, store failures, avg duration
- SQLite persistence with context-aware queries

## Run

```bash
go run main.go
```

Results are saved to `url_fetcher.db` in the project root.
