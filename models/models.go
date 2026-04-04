package models

// UrlRecord represents a URL record in the database
type UrlRecord struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
