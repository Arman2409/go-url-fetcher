package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type IUrlFetcher interface {
	FetchUrlWithContext(ctx context.Context, url string) (string, int, error)
}

type UrlFetcher struct {
	client *http.Client
}

func NewUrlFetcher(client *http.Client) IUrlFetcher {
	return &UrlFetcher{client: client}
}

func (f *UrlFetcher) FetchUrlWithContext(ctx context.Context, url string) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", resp.StatusCode, fmt.Errorf("unexpected status code %d for url %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	return string(body), resp.StatusCode, nil
}
