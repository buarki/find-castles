package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
)

const (
	EnglandCastlesListPage = "https://medievalbritain.com/medieval-castles-of-england"
)

func getHTMLHavingTheListOfCastlesInEngland(ctx context.Context, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", EnglandCastlesListPage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GET request for [%s] to get the list of castles in england, got %v", EnglandCastlesListPage, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get the HTML with the list of castles in england by doing GET at [%s], got %v", EnglandCastlesListPage, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body of GET at [%s], got %v", EnglandCastlesListPage, err)
	}
	return rawBody, nil
}

func extractTheListOfCastlesInEngland(rawHTML []byte) ([]castle.Model, error) {
	var castles []castle.Model
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %v", err)
	}
	doc.Find(".elementor-post .elementor-post__title a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, _ := s.Attr("href")
		castles = append(castles, castle.Model{
			Name:     strings.ReplaceAll(strings.ReplaceAll(title, "\t", ""), "\n", ""),
			Link:     link,
			Country:  castle.UK,
			FlagLink: "/uk-flag.webp",
		})
	})
	return castles, nil
}

func collectForUK(ctx context.Context, httpClient *http.Client) (chan castle.Model, chan error) {
	castlesToEnrich := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(castlesToEnrich)
		defer close(errChan)
		htmlOfCastlesInEngland, err := getHTMLHavingTheListOfCastlesInEngland(ctx, httpClient)
		if err != nil {
			errChan <- err
		}
		englandCastles, err := extractTheListOfCastlesInEngland(htmlOfCastlesInEngland)
		if err != nil {
			errChan <- err
		}
		for _, c := range englandCastles {
			castlesToEnrich <- c
		}
	}()

	return castlesToEnrich, errChan
}
