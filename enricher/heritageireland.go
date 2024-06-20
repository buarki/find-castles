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
	heritageIrelandHost = "https://heritageireland.ie"
	herirageIrelandURL  = heritageIrelandHost + "/visit/castles/"
)

type heritageirelandEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewHeritageIreland(httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &heritageirelandEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (ie *heritageirelandEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
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
				htmlWithCastlesToCollect, err := ie.fetchHTML(ctx, herirageIrelandURL, ie.httpClient)
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

func (ie *heritageirelandEnricher) collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
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
				Name:    name,
				Link:    link,
				Country: castle.Ireland,
			})
		}
	})

	return castles, nil
}

func (ie *heritageirelandEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := ie.fetchHTML(ctx, c.Link, ie.httpClient)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled, err := ie.extractCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled.CleanFields()
	return enrichedCastled, nil
}

func (ie *heritageirelandEnricher) extractCastleInfo(c castle.Model, castlePage []byte) (castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(castlePage))
	if err != nil {
		return castle.Model{}, fmt.Errorf("error loading HTML: %v", err)
	}
	rawAddress, err := ie.extractContact(doc)
	if err != nil {
		return castle.Model{}, err
	}

	district, city, state := ie.get(rawAddress)

	return castle.Model{
		Name:        c.Name,
		Country:     castle.Ireland,
		Link:        c.Link,
		City:        city,
		State:       state,
		District:    district,
		PictureLink: ie.collectImage(doc),
	}, nil
}

func (ie *heritageirelandEnricher) extractContact(doc *goquery.Document) (string, error) {
	address, err := doc.Find("#place--contact div p.address").First().Html()
	if err != nil {
		return "", err
	}

	return address, nil
}

func (ie *heritageirelandEnricher) get(raw string) (string, string, string) {
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

func (ie *heritageirelandEnricher) collectImage(doc *goquery.Document) string {
	var srcset string
	picture := doc.Find("section.gallery ul.hi_gallery li a figure picture").First()
	if picture.Find("source").Length() > 0 {
		srcset, _ = picture.Find("source").First().Attr("srcset")
	}
	parts := strings.Split(srcset, " ")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
