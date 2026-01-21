package services

import (
	"fmt"

	"main.go/database"
)

type IUrlResponseProcessor interface {
	ProcessAndStoreUrlResponse(url string, content string) error
}

type UrlResponseProcessor struct {
	repo database.IRepository
}

func NewUrlResponseProcessor(repo database.IRepository) IUrlResponseProcessor {
	return &UrlResponseProcessor{repo: repo}
}

func (p *UrlResponseProcessor) ProcessAndStoreUrlResponse(url string, content string) error {
	record, err := p.repo.CreateUrlRecord(url, content)

	if err != nil {
		fmt.Println("Error processing and storing URL response:", err)
		return err
	}
	
	fmt.Println("Successfully stored URL record with ID:", record.ID)
	
	return nil
}
