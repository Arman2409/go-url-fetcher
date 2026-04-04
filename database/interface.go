package database

import (
	"context"

	"main.go/models"
)

// IRepository defines the interface for database operations
type IRepository interface {
	// CreateUrlRecord inserts a new URL record into the database
	CreateUrlRecordWithContext(ctx context.Context, url string, statusCode int, content string) (*models.UrlRecord, error)

	// GetUrlRecordByID retrieves a URL record by its ID
	GetUrlRecordByID(id int64) (*models.UrlRecord, error)

	// GetUrlRecordByURL retrieves a URL record by its URL
	GetUrlRecordByURL(url string) (*models.UrlRecord, error)

	// GetAllUrlRecords retrieves all URL records
	GetAllUrlRecords() ([]*models.UrlRecord, error)

	// DeleteUrlRecord deletes a URL record by ID
	DeleteUrlRecord(id int64) error

	// Close closes the database connection
	Close() error
}
