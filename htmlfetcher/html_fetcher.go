package htmlfetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type HTMLFetcher func(ctx context.Context, url string, httpClient *http.Client) ([]byte, error)

func Fetch(ctx context.Context, url string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request [%s], got %v", url, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET at [%s], got %v", url, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content of GET [%s], got %v", url, err)
	}
	return rawBody, nil
}
