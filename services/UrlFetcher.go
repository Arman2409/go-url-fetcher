package services

import (
	"context"
	"io"
	"net/http"
)

type IUrlFetcher interface {
   FetchUrlWithContext(ctx context.Context, url string) (string, error)
}

type UrlFetcher struct {
	client *http.Client
}

func NewUrlFetcher(client *http.Client) IUrlFetcher {
	return &UrlFetcher{client: client}
}

func (f *UrlFetcher) FetchUrlWithContext(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}