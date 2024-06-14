package enricher

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
)

const (
	castlesSource = "https://www.castelosdeportugal.pt"
)

type portugueseEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewPortugueseEnricher(
	httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &portugueseEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (p *portugueseEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Portugal received done!")
				return
			default:
				htmlWithCastlesToCollect, err := p.fetchHTML(ctx, castlesSource, p.httpClient)
				if err != nil {
					errChan <- err
					return
				}
				castles, err := p.collectCastleNameAndLinks(htmlWithCastlesToCollect)
				if err != nil {
					errChan <- err
					return
				}
				for _, c := range castles {
					castlesToEnrichChan <- c
				}
				return
			}
		}

	}()

	return castlesToEnrichChan, errChan
}

func (p *portugueseEnricher) collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
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

func (p *portugueseEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := p.fetchHTML(ctx, c.Link, p.httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := p.extractCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	return enrichedCastled, nil
}

func (p *portugueseEnricher) extractCastleInfo(c castle.Model, rawHTMLPage []byte) (castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTMLPage))
	if err != nil {
		return castle.Model{}, fmt.Errorf("failed to load page, got %v", err)
	}

	var tableData = make(map[string]string)

	rowsToExtract := []string{"Distrito", "Concelho", "Freguesia", "Construção"}

	doc.Find("#info-table tbody tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		if p.contains(rowsToExtract, key) {
			value := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
			tableData[key] = value
		}
	})

	return castle.Model{
		Name:             c.Name,
		Country:          c.Country,
		Link:             c.Link,
		City:             tableData["Concelho"],
		State:            tableData["Distrito"],
		District:         tableData["Freguesia"],
		FoundationPeriod: tableData["Construção"],
		FlagLink:         "/pt-flag.webp",
	}, nil
}

func (p *portugueseEnricher) contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
