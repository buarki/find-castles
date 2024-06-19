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
	castelosdeportugalHost          = "https://www.castelosdeportugal.pt"
	castelosdeportugalCastlesSource = castelosdeportugalHost + "/castelos/SiteMap.html"
)

type castelosDePortugalEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewCastelosDePortugalEnricher(
	httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &castelosDePortugalEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (p *castelosDePortugalEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("castelosdeportugal received done!")
				return
			default:
				htmlWithCastlesToCollect, err := p.fetchHTML(ctx, castelosdeportugalCastlesSource, p.httpClient)
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

func (p *castelosDePortugalEnricher) collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML of portugal: %v", err)
	}
	var castles []castle.Model
	doc.Find("#indice div a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		name := s.Text()
		castles = append(castles, castle.Model{
			Name:    name,
			Country: castle.Portugal,
			Link:    fmt.Sprintf("%s/castelos/%s", castelosdeportugalHost, link),
		})
	})
	return castles, nil
}

func (p *castelosDePortugalEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := p.fetchHTML(ctx, c.Link, p.httpClient)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled, err := p.extractCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled.CleanFields()
	return enrichedCastled, nil
}

func (p *castelosDePortugalEnricher) extractCastleInfo(c castle.Model, rawHTMLPage []byte) (castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTMLPage))
	if err != nil {
		return castle.Model{}, fmt.Errorf("failed to load page, got %v", err)
	}

	var tableData = make(map[string]string)

	rowsToExtract := []string{
		"Distrito",
		"Concelho",
		"Freguesia",
		"Construção",
		"Conservação",
	}

	doc.Find("#info-table tbody tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		if p.contains(rowsToExtract, key) {
			value := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
			tableData[key] = value
		}
	})

	// we execute it again because some castles have the city name as district
	// ex: https://www.castelosdeportugal.pt/castelos/Castelos(pre)SECXII/velhoAlcoutim.html
	rowsToExtract = append(rowsToExtract, tableData["Concelho"])

	doc.Find("#info-table tbody tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		if p.contains(rowsToExtract, key) {
			value := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
			tableData[key] = value
		}
	})

	district := tableData["Freguesia"]
	if district == "" {
		district = tableData[tableData["Concelho"]]
	}

	return castle.Model{
		Name:              c.Name,
		Country:           c.Country,
		Link:              c.Link,
		City:              tableData["Concelho"],
		State:             tableData["Distrito"],
		District:          district,
		FoundationPeriod:  tableData["Construção"],
		PropertyCondition: p.parseCondition(tableData["Conservação"]),
	}, nil
}

func (p castelosDePortugalEnricher) parseCondition(rawCondition string) castle.PropertyCondition {
	filtered := strings.ToLower(rawCondition)
	switch filtered {
	case "boa":
		return castle.Intact
	case "razoável":
		return castle.Damaged
	case "mau":
		return castle.Ruins
	case "submerso":
		return castle.Ruins
	default:
		return castle.Unknown
	}
}

func (p *castelosDePortugalEnricher) contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
