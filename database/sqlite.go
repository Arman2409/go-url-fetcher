package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"main.go/models"
)

// SQLiteRepository implements IRepository using SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository instance
func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	// Open database connection (lazy - doesn't connect until first query)
	// _busy_timeout=5000 sets a 5 second timeout if database is locked
	db, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Ping forces an immediate connection test to catch errors early
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &SQLiteRepository{db: db}

	// Initialize the database schema
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates the necessary tables if they don't exist
func (r *SQLiteRepository) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS url_records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// CreateUrlRecord inserts a new URL record into the database
func (r *SQLiteRepository) CreateUrlRecordWithContext(ctx context.Context, url, content string) (*models.UrlRecord, error) {
	createdAt := time.Now().Format(time.RFC3339)

	query := `
	INSERT INTO url_records (url, content, created_at)
	VALUES (?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, url, content, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.UrlRecord{
		ID:        id,
		URL:       url,
		Content:   content,
		CreatedAt: createdAt,
	}, nil
}

// GetUrlRecordByID retrieves a URL record by its ID
func (r *SQLiteRepository) GetUrlRecordByID(id int64) (*models.UrlRecord, error) {
	query := `
	SELECT id, url, content, created_at
	FROM url_records
	WHERE id = ?
	`

	var record models.UrlRecord
	err := r.db.QueryRow(query, id).Scan(
		&record.ID,
		&record.URL,
		&record.Content,
		&record.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	return &record, nil
}

// GetUrlRecordByURL retrieves a URL record by its URL
func (r *SQLiteRepository) GetUrlRecordByURL(url string) (*models.UrlRecord, error) {
	query := `
	SELECT id, url, content, created_at
	FROM url_records
	WHERE url = ?
	`

	var record models.UrlRecord
	err := r.db.QueryRow(query, url).Scan(
		&record.ID,
		&record.URL,
		&record.Content,
		&record.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("record not found")
		}
		return nil, fmt.Errorf("failed to get record: %w", err)
	}

	return &record, nil
}

// GetAllUrlRecords retrieves all URL records
func (r *SQLiteRepository) GetAllUrlRecords() ([]*models.UrlRecord, error) {
	query := `
	SELECT id, url, content, created_at
	FROM url_records
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query records: %w", err)
	}
	defer rows.Close()

	var records []*models.UrlRecord
	for rows.Next() {
		var record models.UrlRecord
		if err := rows.Scan(
			&record.ID,
			&record.URL,
			&record.Content,
			&record.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan record: %w", err)
		}
		records = append(records, &record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating records: %w", err)
	}

	return records, nil
}

// DeleteUrlRecord deletes a URL record by ID
func (r *SQLiteRepository) DeleteUrlRecord(id int64) error {
	query := `DELETE FROM url_records WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}

// Close closes the database connection
func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
