package enricher

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
	"golang.org/x/sync/errgroup"
)

// TODO do not process the four sources in separated goroutines, use a loop instead

const (
	listOfCastlesInEngland        = "https://medievalbritain.com/medieval-castles-of-england"
	listOfCastlesInScotland       = "https://medievalbritain.com/medieval-castles-of-scotland"
	listOfCastlesInWales          = "https://medievalbritain.com/medieval-castles-of-wales"
	listOfCastlesInNorthenIreland = "https://medievalbritain.com/medieval-castles-of-northern-ireland"

	workersToExtractCastlesFromHTML = 3
)

type medievalbritainEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewMedievalBritainEnricher(httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &medievalbritainEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (be *medievalbritainEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errorsChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errorsChan)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("UK got done")
				return
			default:
				castles, err := be.collect(ctx)
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

func (be *medievalbritainEnricher) collect(ctx context.Context) ([]castle.Model, error) {
	sources := []string{
		listOfCastlesInEngland,
		listOfCastlesInScotland,
		listOfCastlesInWales,
		listOfCastlesInNorthenIreland,
	}
	collectedHTMLs, err := be.collectHTMLPagesToExtractCastlesInfo(ctx, sources)
	if err != nil {
		return nil, err
	}
	return be.extractTheListOfCastlesToEnrich(ctx, collectedHTMLs, workersToExtractCastlesFromHTML)
}

func (be *medievalbritainEnricher) collectHTMLPagesToExtractCastlesInfo(ctx context.Context, sources []string) ([][]byte, error) {
	var rawHTMLs [][]byte
	var mutex sync.Mutex
	errs, errCtx := errgroup.WithContext(ctx)
	for _, source := range sources {
		s := source
		errs.Go(func() error {
			rawHTML, err := be.fetchHTML(errCtx, s, be.httpClient)
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

func (be *medievalbritainEnricher) extractTheListOfCastlesToEnrich(ctx context.Context, rawHTMLs [][]byte, workers int) ([]castle.Model, error) {
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
						collectedCastlesToEnrich, err := be.extractTheListOfCastlesFromPage(html)
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

func (be *medievalbritainEnricher) extractTheListOfCastlesFromPage(rawHTML []byte) ([]castle.Model, error) {
	var castles []castle.Model
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %v", err)
	}
	doc.Find(".elementor-post .elementor-post__title a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, _ := s.Attr("href")
		castles = append(castles, castle.Model{
			Name:                    strings.ReplaceAll(strings.ReplaceAll(title, "\t", ""), "\n", ""),
			CurrentEnrichmentLink:   link,
			Country:                 castle.UK,
			CurrentEnrichmentSource: MedievalBritain.String(),
		})
	})
	return castles, nil
}

func (be *medievalbritainEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := be.fetchHTML(ctx, c.CurrentEnrichmentLink, be.httpClient)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled, err := be.extractDataOfUKCastle(castlePage, c)
	if err != nil {
		return castle.Model{}, err
	}
	enrichedCastled.CleanFields()
	return enrichedCastled, nil
}

/*
About the UK state

From inspection we could see that the query '.elementor-widget-container div.elementor-text-editor.elementor-clearfix p a”
provides the state as its first item.

# An example of chunk in which we parse is bellow one

<div class="elementor-widget-container">

		<div
			class="elementor-text-editor elementor-clearfix">
			<p>Berkshire, <a
					href="https://medievalbritain.com/category/locations/england/greater-london/">Greater
					London</a><br />(<a
					class="external text"
					href="https://geohack.toolforge.org/geohack.php?pagename=Windsor_Castle&amp;params=51_29_0_N_00_36_15_W_region:GB_type:landmark"
					target="_blank"
					rel="nofollow noopener"><span
						class="geo-default"><span
							class="geo-dms"
							title="Maps, aerial photos, and other data for this location"><span
								class="latitude">51°29′0″N</span> <span
								class="longitude">00°36′15″W</span></span></span></a>)
			</p>
		</div>
	</div>
*/
func (be *medievalbritainEnricher) extractUkState(doc *goquery.Document, c castle.Model) (string, error) {
	state := doc.Find(".elementor-widget-container div.elementor-text-editor.elementor-clearfix p a").First().Parent().Text()
	state = strings.ReplaceAll(strings.ReplaceAll(state, "\n", ""), "\t", "")
	before, _, found := strings.Cut(state, "(")
	if found {
		state = before
	}
	return state, nil
}

/*
About the UK city

From inspection we could see that the query '.elementor-widget-container div.elementor-text-editor.elementor-clearfix p a”
provides the state as its first item.

# An example of chunk in which we parse is bellow one
<div class="elementor-widget-container">

		<div
			class="elementor-text-editor elementor-clearfix">
			<p><strong>Address</strong></p>
			<p>Windsor SL4 1NJ</p>
		</div>
	</div>
*/
func (be *medievalbritainEnricher) extractUkCity(doc *goquery.Document, c castle.Model) (string, error) {
	var city string
	doc.Find(".elementor-text-editor.elementor-clearfix p").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Address" {
			city = s.Next().Text()
			before, _, found := strings.Cut(city, ",")
			if found {
				city = before
			}
		}
	})
	return city, nil
}

func (be *medievalbritainEnricher) extractDataOfUKCastle(rawHTML []byte, c castle.Model) (castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return castle.Model{}, fmt.Errorf("failed to create reader for castle [%s], got %v", c.Name, err)
	}
	state, err := be.extractUkState(doc, c)
	if err != nil {
		return castle.Model{}, err
	}
	city, err := be.extractUkCity(doc, c)
	if err != nil {
		return castle.Model{}, err
	}
	return castle.Model{
		Name:                  c.Name,
		Country:               c.Country,
		CurrentEnrichmentLink: c.CurrentEnrichmentLink,
		State:                 state,
		City:                  city,
		PictureURL:            be.collectImage(doc),
		Coordinates:           be.collectCoordinates(doc),
		Contact:               be.collectContactInfo(doc),
		Sources:               []string{c.CurrentEnrichmentLink},
		VisitingInfo:          be.collectVisitingInfo(doc),
		PropertyCondition:     castle.Unknown,
	}, nil
}

func (be *medievalbritainEnricher) collectImage(doc *goquery.Document) string {
	var imageSrc string
	metaTag := doc.Find("meta[property='og:image']")
	imageSrc, _ = metaTag.Attr("content")
	return imageSrc
}

func (be *medievalbritainEnricher) collectCoordinates(doc *goquery.Document) string {
	replacer := strings.NewReplacer(
		`′`, `'`,
		`″`, `"`,
		`\n`, "",
		` `, ",",
	)

	latitude := doc.Find(".geo-default .latitude").First().Text()
	longitude := doc.Find(".geo-default .longitude").First().Text()

	if latitude != "" && longitude != "" {
		return fmt.Sprintf("%s,%s", replacer.Replace(latitude), replacer.Replace(longitude))
	}

	geoDec := doc.Find(".geo-default .geo-dec").First().Text()
	if geoDec != "" {
		return replacer.Replace(geoDec)
	}

	// fallback for when .geo-default is not present and it uses tools.wmflabs.org
	targetReferenceHost := "tools.wmflabs.org"

	latitudeLongitude := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, targetReferenceHost) {
			latLng := strings.Split(s.Text(), " ")
			latitudeLongitude = fmt.Sprintf("%s,%s", latLng[0], latLng[1])
			return
		}
	})

	// fallback for when .geo-default is not present and it uses https://goo.gl/maps/
	targetReferenceHost = "https://goo.gl/maps/"

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, targetReferenceHost) {
			latitudeLongitude = strings.ReplaceAll(s.Text(), " ", "")
			return
		}
	})

	return latitudeLongitude
}

func (be *medievalbritainEnricher) collectContactInfo(doc *goquery.Document) *castle.Contact {
	phoneElement := doc.Find(".elementor-text-editor .w8qArf").First()
	collectedPhone := phoneElement.Parent().Next().Text()
	if collectedPhone != "" {
		return &castle.Contact{
			Phone: collectedPhone,
		}
	}
	return nil
}

func (be *medievalbritainEnricher) collectVisitingInfo(doc *goquery.Document) *castle.VisitingInfo {
	return &castle.VisitingInfo{
		WorkingHours: be.collectHorkingHours(doc),
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

// the searching key for this one is <strong>Hours</strong>
func (be *medievalbritainEnricher) collectHorkingHours(doc *goquery.Document) string {
	var hours []string

	doc.Find(".elementor-text-editor p").Each(func(i int, s *goquery.Selection) {
		if s.Find("strong").Text() == "Hours" {
			s.NextAll().Each(func(j int, p *goquery.Selection) {
				text := p.Text()
				if strings.Contains(text, ":") {
					hours = append(hours, text)
				}
			})
			return
		}
	})

	replacer := strings.NewReplacer(
		`–`, "-",
	)

	workingHours := replacer.Replace(strings.Join(hours, ","))
	if workingHours != "" {
		return workingHours
	}

	// given hours as raw text
	doc.Find(".elementor-text-editor").Each(func(i int, s *goquery.Selection) {
		if s.Find("strong").Text() == "Hours" {
			wh := s.Find("p").Text()
			if wh != "" {
				workingHours = strings.ReplaceAll(wh, "Hours", "")
				return
			}
		}
	})

	if workingHours != "" {
		return workingHours
	}

	// given hours as table
	doc.Find(".elementor-text-editor").Each(func(i int, s *goquery.Selection) {
		if s.Find("p:contains('Hours')").Length() > 0 {
			// get the table after the <strong>Hours</strong> element
			s.Find("table.WgFkxc tr").Each(func(i int, tr *goquery.Selection) {
				season := strings.TrimSpace(tr.Find("td.SKNSIb").Text())
				time := strings.TrimSpace(tr.Find("td").Last().Text())
				if season != "" && time != "" {
					hours = append(hours, fmt.Sprintf("%s: %s", season, time))
				}
			})
		}
	})

	return replacer.Replace(strings.Join(hours, ", "))

}
