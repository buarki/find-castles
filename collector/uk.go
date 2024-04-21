package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
	"golang.org/x/sync/errgroup"
)

const (
	listOfCastlesInEngland        = "https://medievalbritain.com/medieval-castles-of-england"
	listOfCastlesInScotland       = "https://medievalbritain.com/medieval-castles-of-scotland"
	listOfCastlesInWales          = "https://medievalbritain.com/medieval-castles-of-wales"
	listOfCastlesInNorthenIreland = "https://medievalbritain.com/medieval-castles-of-northern-ireland"

	workersToExtractCastlesFromHTML = 3
)

func getHTMLHavingTheListOfInARegionOfUk(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GET request for [%s] to get the list of castles, got %v", link, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get the HTML with the list of castles by doing GET at [%s], got %v", link, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body of GET at [%s], got %v", link, err)
	}
	return rawBody, nil
}

func extractTheListOfCastlesFromPage(rawHTML []byte) ([]castle.Model, error) {
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

func collectHTMLPagesToExtractCastlesInfo(ctx context.Context, sources []string, httpClient *http.Client) ([][]byte, error) {
	var rawHTMLs [][]byte
	var mutex sync.Mutex
	errs, errCtx := errgroup.WithContext(ctx)
	for _, source := range sources {
		s := source
		errs.Go(func() error {
			rawHTML, err := getHTMLHavingTheListOfInARegionOfUk(errCtx, s, httpClient)
			if err != nil {
				return err
			}
			mutex.Lock()
			rawHTMLs = append(rawHTMLs, rawHTML)
			mutex.Unlock()
			return nil
		})
	}
	return rawHTMLs, errs.Wait()
}

func extractTheListOfCastlesToEnrich(ctx context.Context, rawHTMLs [][]byte, workers int) ([]castle.Model, error) {
	errs, errCtx := errgroup.WithContext(ctx)
	htmlsChan := make(chan []byte)
	var collectedCastles []castle.Model
	var mutex sync.Mutex
	for i := 0; i < workers; i++ {
		errs.Go(func() error {
			for {
				select {
				case <-errCtx.Done():
					return nil
				case html, ok := <-htmlsChan:
					if ok {
						collectedCastlesToEnrich, err := extractTheListOfCastlesFromPage(html)
						if err != nil {
							return err
						}
						mutex.Lock()
						collectedCastles = append(collectedCastles, collectedCastlesToEnrich...)
						mutex.Unlock()
					} else {
						return nil
					}
				}
			}
		})
	}
	errs.Go(func() error {
		defer close(htmlsChan)
		for _, html := range rawHTMLs {
			htmlsChan <- html
		}
		return nil
	})
	return collectedCastles, errs.Wait()
}

func collectCastlesFromUK(ctx context.Context, httpClient *http.Client) (chan castle.Model, chan error) {
	castlesToEnrich := make(chan castle.Model)
	errChan := make(chan error)
	sources := []string{
		listOfCastlesInEngland,
		listOfCastlesInScotland,
		listOfCastlesInWales,
		listOfCastlesInNorthenIreland,
	}

	go func() {
		defer close(castlesToEnrich)
		defer close(errChan)

		collectedHTMLs, err := collectHTMLPagesToExtractCastlesInfo(ctx, sources, httpClient)
		if err != nil {
			errChan <- err
		}
		englandCastles, err := extractTheListOfCastlesToEnrich(ctx, collectedHTMLs, workersToExtractCastlesFromHTML)
		if err != nil {
			errChan <- err
		}
		for _, c := range englandCastles {
			castlesToEnrich <- c
		}
	}()

	return castlesToEnrich, errChan
}
