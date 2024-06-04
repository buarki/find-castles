package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
)

const (
	irelandCastles = "https://heritageireland.ie/visit/castles/"
)

func collectIrelanHomePageHTML(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request [%s], got %v", link, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET at [%s], got %v", link, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content of GET [%s], got %v", link, err)
	}
	return rawBody, nil
}

func extract(rawHTML []byte) ([]castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %v", err)
	}

	var castles []castle.Model
	doc.Find("#placesgrid ul li a").Each(func(i int, s *goquery.Selection) {
		name := s.Find("header div h3").Text()
		link, exists := s.Attr("href")
		if exists {
			castles = append(castles, castle.Model{
				Name:     name,
				Link:     link,
				Country:  castle.Ireland,
				FlagLink: "/ir-flag.jpeg",
			})
		}
	})

	return castles, nil
}

func collectCastlesFromIreland(ctx context.Context, httpClient *http.Client) (chan castle.Model, chan error) {
	castlesToEnrich := make(chan castle.Model)
	errChan := make(chan error)
	go func() {
		defer close(castlesToEnrich)
		defer close(errChan)
		htmlWithCastlesToCollect, err := collectIrelanHomePageHTML(ctx, irelandCastles, httpClient)
		if err != nil {
			errChan <- err
		}
		castles, err := extract(htmlWithCastlesToCollect)
		if err != nil {
			errChan <- err
		}
		for _, c := range castles {
			castlesToEnrich <- c
		}
	}()
	return castlesToEnrich, errChan
}
