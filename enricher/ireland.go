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
	irelandCastlesURL = "https://heritageireland.ie/visit/castles/"
)

type irishEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewIrishEnricher(httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &irishEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (ie *irishEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errorsChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errorsChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Ireland got done")
				return
			default:
				htmlWithCastlesToCollect, err := ie.fetchHTML(ctx, irelandCastlesURL, ie.httpClient)
				if err != nil {
					errorsChan <- err
				}
				castles, err := ie.collectCastleNameAndLinks(htmlWithCastlesToCollect)
				if err != nil {
					errorsChan <- err
				}
				for _, c := range castles {
					castlesToEnrichChan <- c
				}
				return
			}
		}

	}()

	return castlesToEnrichChan, errorsChan
}

func (ie *irishEnricher) collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
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

func (ie *irishEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := ie.fetchHTML(ctx, c.Link, ie.httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := ie.extractCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	return enrichedCastled, nil
}

func (ie *irishEnricher) extractCastleInfo(c castle.Model, castlePage []byte) (castle.Model, error) {
	rawAddress, err := ie.extractContact(castlePage)
	if err != nil {
		return castle.Model{}, err
	}

	district, city, state := ie.get(rawAddress)

	return castle.Model{
		Name:     c.Name,
		Country:  castle.Ireland,
		Link:     c.Link,
		City:     city,
		State:    state,
		District: district,
		FlagLink: c.FlagLink,
	}, nil
}

func (ie *irishEnricher) extractContact(rawHTML []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return "", fmt.Errorf("error loading HTML: %v", err)
	}

	address, err := doc.Find("#place--contact div p.address").First().Html()
	if err != nil {
		return "", err
	}

	return address, nil
}

func (ie *irishEnricher) get(raw string) (string, string, string) {
	parts := strings.Split(raw, "<br/>")
	for i := range parts {
		parts[i] = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(parts[i], ",", ""), "\n", ""))
	}
	// ex: [Ross Castle, Ross Road, Killarney, Co. Kerry, V93 V304]
	placeInsideDistrictInformed := len(parts) == 5
	if placeInsideDistrictInformed {
		return parts[1], parts[2], parts[len(parts)-2]
	}
	// ex: [Trim, Co Meath, C15 HN90]
	districtAndCityAreEqual := len(parts) < 4
	if districtAndCityAreEqual {
		return parts[0], parts[0], parts[len(parts)-2]
	}
	// ex: [Adare Heritage Centre, Adare, Co. Limerick, V94 DWV7]
	return parts[0], parts[1], parts[len(parts)-2]
}
