package models

// UrlRecord represents a URL record in the database
type UrlRecord struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}
