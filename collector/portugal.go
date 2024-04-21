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
	castlesSource = "https://www.castelosdeportugal.pt"
)

func collectHomePageHTML(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
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

func collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML of portugal: %v", err)
	}
	var castles []castle.Model
	doc.Find("#div-list-alfa-cast p a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		name := s.Text()
		castles = append(castles, castle.Model{
			Name:     name,
			Country:  castle.Portugal,
			Link:     fmt.Sprintf("%s/castelos/%s", castlesSource, strings.ReplaceAll(link, "../", "")),
			FlagLink: "/pt-flag.webp",
		})
	})
	return castles, nil
}

func collectForPotugal(ctx context.Context, httpClient *http.Client) (chan castle.Model, chan error) {
	castlesToEnrich := make(chan castle.Model)
	errChan := make(chan error)
	go func() {
		defer close(castlesToEnrich)
		defer close(errChan)
		htmlWithCastlesToCollect, err := collectHomePageHTML(ctx, castlesSource, httpClient)
		if err != nil {
			errChan <- err
		}
		castles, err := collectCastleNameAndLinks(htmlWithCastlesToCollect)
		if err != nil {
			errChan <- err
		}
		for _, c := range castles {
			castlesToEnrich <- c
		}
	}()
	return castlesToEnrich, errChan
}
