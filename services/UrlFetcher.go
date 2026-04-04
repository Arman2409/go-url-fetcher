package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	maxFetchAttempts = 4
	retryBaseDelay   = 200 * time.Millisecond
	maxRetryBackoff  = 2 * time.Second
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
	var lastErr error
	var lastCode int
	backoff := retryBaseDelay

	for attempt := range maxFetchAttempts {
		if attempt > 0 {
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return "", lastCode, ctx.Err()
			}
			backoff *= 2
			if backoff > maxRetryBackoff {
				backoff = maxRetryBackoff
			}
		}

		body, code, err := f.fetchOnce(ctx, url)
		if err == nil {
			return body, code, nil
		}
		lastErr, lastCode = err, code
		if !isRetryable(err, code) {
			return "", code, err
		}
	}

	return "", lastCode, lastErr
}

func (f *UrlFetcher) fetchOnce(ctx context.Context, url string) (string, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", 0, err
	}
	resp, err := f.client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		_, _ = io.Copy(io.Discard, resp.Body)
		return "", resp.StatusCode, fmt.Errorf("unexpected status code %d for url %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}

	return string(body), resp.StatusCode, nil
}

func isRetryable(err error, statusCode int) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}
	if statusCode == 0 {
		return true
	}
	if statusCode == http.StatusTooManyRequests {
		return true
	}
	return statusCode >= http.StatusInternalServerError
}
