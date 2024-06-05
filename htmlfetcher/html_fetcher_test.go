package htmlfetcher_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/buarki/find-castles/htmlfetcher"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) *http.Response
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req), nil
}

func newMockClient(roundTripFunc func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &mockTransport{roundTripFunc: roundTripFunc},
	}
}

type errReader struct{}

func (e *errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func TestFetch(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		statusCode int
		body       string
		expectErr  bool
	}{
		{
			name:       "Successful fetch",
			url:        "http://example.com",
			statusCode: http.StatusOK,
			body:       "<html>Example</html>",
			expectErr:  false,
		},
		{
			name:       "Error creating request",
			url:        "http://[invalid-url",
			statusCode: http.StatusOK,
			body:       "<html>Example</html>",
			expectErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockClient(func(req *http.Request) *http.Response {
				if tt.expectErr && tt.body == "" {
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(&errReader{}),
						Header:     make(http.Header),
					}
				}
				return &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
					Header:     make(http.Header),
				}
			})

			ctx := context.Background()
			_, err := htmlfetcher.Fetch(ctx, tt.url, mockClient)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
