package httpclient

import (
	"net/http"
	"time"
)

const (
	timeout = 10 * time.Second
)

func New() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConns = 100

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}
