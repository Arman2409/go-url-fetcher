package services

import (
	"io"
	"net/http"
)

type IUrlFetcher interface {
   FetchUrl(url string) (string, error)
}

type UrlFetcher struct {
	client *http.Client
}

func NewUrlFetcher(client *http.Client) IUrlFetcher {
	return &UrlFetcher{client: client}
}

func (f *UrlFetcher) FetchUrl(url string) (string, error) {
	resp, err := f.client.Get(url);
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