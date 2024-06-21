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
				Name:                  name,
				CurrentEnrichmentLink: link,
				Country:               castle.Ireland,
			})
		}
	})

	return castles, nil
}

func (ie *heritageirelandEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := ie.fetchHTML(ctx, c.CurrentEnrichmentLink, ie.httpClient)
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
		Name:                  c.Name,
		Country:               castle.Ireland,
		CurrentEnrichmentLink: c.CurrentEnrichmentLink,
		City:                  city,
		State:                 state,
		District:              district,
		PictureURL:            ie.collectImage(doc),
		Contact:               ie.collectContactInfo(doc),
		Sources:               []string{c.CurrentEnrichmentLink},
		VisitingInfo:          ie.collectVisitingInfo(doc),
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

func (ie *heritageirelandEnricher) collectContactInfo(doc *goquery.Document) *castle.Contact {
	collectedPhone := doc.Find("#place--contact div .phone").First().Text()
	collectedEmail := doc.Find("#place--contact div .email").First().Text()
	if collectedEmail != "" || collectedPhone != "" {
		return &castle.Contact{
			Email: collectedEmail,
			Phone: collectedPhone,
		}
	}
	return nil
}

func (ie *heritageirelandEnricher) collectVisitingInfo(doc *goquery.Document) *castle.VisitingInfo {
	return &castle.VisitingInfo{
		WorkingHours: ie.collectHorkingHours(doc),
		Facilities: &castle.Facilities{
			AssistanceDogsAllowed: doc.Find(".fa-dog").Length() > 0,
			Giftshops:             doc.Find(".fa-shopping-bag").Length() > 0,
			WheelchairSupport:     doc.Find(".fa-wheelchair").Length() > 0,
			Restrooms:             doc.Find(".fa-toilet").Length() > 0,
			PinicArea:             doc.Find(".fa-tree").Length() > 0,
			Exhibitions:           doc.Find(".fa-vector-square").Length() > 0,
			Cafe:                  doc.Find(".fa-coffee").Length() > 0,
			Parking:               doc.Find(".fa-car-alt").Length() > 0,
		},
	}
}

// we can use id place--opening
func (ie heritageirelandEnricher) collectHorkingHours(doc *goquery.Document) string {
	replacer := strings.NewReplacer(
		`–`, "-",
	)
	var dateRange, timeRange string
	doc.Find("section#place--opening").Each(func(i int, s *goquery.Selection) {
		dateRange = strings.TrimSpace(s.Find("p strong").Text())
		timeRange = strings.TrimSpace(s.Find("p").Next().Text())
	})
	if dateRange != "" && timeRange != "" {
		return replacer.Replace(fmt.Sprintf("%s - %s", dateRange, timeRange))
	}

	//in case of a accordion used
	var openingDates string

	doc.Find("section#place--opening").Each(func(i int, s *goquery.Selection) {
		openingDates = strings.TrimSpace(s.Find("div p").First().Text())
	})

	return openingDates
}
