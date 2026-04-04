package services

import (
	"context"
	"fmt"

	"main.go/database"
)

type IUrlResponseProcessor interface {
	ProcessAndStoreUrlResponseWithContext(ctx context.Context, url string, statusCode int, content string) error
}

type UrlResponseProcessor struct {
	repo database.IRepository
}

func NewUrlResponseProcessor(repo database.IRepository) IUrlResponseProcessor {
	return &UrlResponseProcessor{repo: repo}
}

func (p *UrlResponseProcessor) ProcessAndStoreUrlResponseWithContext(ctx context.Context, url string, statusCode int, content string) error {
	record, err := p.repo.CreateUrlRecordWithContext(ctx, url, statusCode, content)

	if err != nil {
		fmt.Println("Error processing and storing URL response:", err)
		return err
	}

	fmt.Println("Successfully stored URL record with ID:", record.ID)

	return nil
}
