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

const (
	listOfCastlesInEngland        = "https://medievalbritain.com/medieval-castles-of-england"
	listOfCastlesInScotland       = "https://medievalbritain.com/medieval-castles-of-scotland"
	listOfCastlesInWales          = "https://medievalbritain.com/medieval-castles-of-wales"
	listOfCastlesInNorthenIreland = "https://medievalbritain.com/medieval-castles-of-northern-ireland"

	workersToExtractCastlesFromHTML = 3
)

type britishEnricher struct {
	httpClient *http.Client
	fetchHTML  func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)
}

func NewBritishEnricher(httpClient *http.Client,
	fetchHTML func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error)) Enricher {
	return &britishEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (be *britishEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
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

func (be *britishEnricher) collect(ctx context.Context) ([]castle.Model, error) {
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

func (be *britishEnricher) collectHTMLPagesToExtractCastlesInfo(ctx context.Context, sources []string) ([][]byte, error) {
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

func (be *britishEnricher) extractTheListOfCastlesToEnrich(ctx context.Context, rawHTMLs [][]byte, workers int) ([]castle.Model, error) {
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

func (be *britishEnricher) extractTheListOfCastlesFromPage(rawHTML []byte) ([]castle.Model, error) {
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

func (be *britishEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	castlePage, err := be.fetchHTML(ctx, c.Link, be.httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := be.extractDataOfUKCastle(castlePage, c)
	if err != nil {
		return castle.Model{}, err
	}
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
func (be *britishEnricher) extractUkState(rawHTML []byte, c castle.Model) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return "", fmt.Errorf("failed to create reader for castle [%s], got %v", c.Name, err)
	}
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
func (be *britishEnricher) extractUkCity(rawHTML []byte, c castle.Model) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return "", fmt.Errorf("failed to create reader for castle [%s], got %v", c.Name, err)
	}
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

func (be *britishEnricher) extractDataOfUKCastle(rawHTML []byte, c castle.Model) (castle.Model, error) {
	state, err := be.extractUkState(rawHTML, c)
	if err != nil {
		return castle.Model{}, err
	}
	city, err := be.extractUkCity(rawHTML, c)
	if err != nil {
		return castle.Model{}, err
	}
	return castle.Model{
		Name:     c.Name,
		Country:  c.Country,
		FlagLink: c.FlagLink,
		Link:     c.Link,
		State:    state,
		City:     city,
	}, nil
}
