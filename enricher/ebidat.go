package enricher

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/htmlfetcher"
	"golang.org/x/net/html"
)

// TODO support the other countries available other than Slovakia

type ebidatCountrySource struct {
	country   castle.Country
	sourceURL string
}

const (
	ebidatHost   = "www.ebidat.de"
	slovakSource = "https://" + ebidatHost + "/cgi-bin/ebidat.pl?a=a&te53=6"
	danishSource = "https://" + ebidatHost + "/cgi-bin/ebidat.pl?a=a&te53=2"
)

var (
	countriesSources = []ebidatCountrySource{
		{
			country:   castle.Denmark,
			sourceURL: danishSource,
		},
		{
			country:   castle.Slovakia,
			sourceURL: slovakSource,
		},
	}

	stateCollector = map[castle.Country]func(doc *goquery.Document) string{
		castle.Slovakia: func(doc *goquery.Document) string {
			var state string
			doc.Find("li.daten").Each(func(i int, s *goquery.Selection) {
				if s.Find(".gruppe").Text() == "Bundesland:" {
					state = s.Find(".gruppenergebnis").Text()
				}
			})
			return state
		},
		castle.Denmark: func(doc *goquery.Document) string {
			var state string
			referencePoints := []string{
				"Region:",
				"Kreis:",
				"Stadt / Gemeinde:", // fallback for when neither region or kreis is provided
			}
			doc.Find("li.daten").Each(func(i int, s *goquery.Selection) {
				foundText := s.Find(".gruppe").Text()
				if slices.Contains(referencePoints, foundText) {
					state = s.Find(".gruppenergebnis").Text()
				}
			})
			return state
		},
	}
)

type ebidatEnricher struct {
	httpClient *http.Client
	fetchHTML  htmlfetcher.HTMLFetcher
}

func NewEbidatEnricher(
	httpClient *http.Client,
	fetchHTML htmlfetcher.HTMLFetcher) Enricher {
	return &ebidatEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (se *ebidatEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errChan)

		countriesCounter := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if countriesCounter >= len(countriesSources) {
					return
				}

				countrySource := countriesSources[countriesCounter]

				hasMorePages := true
				linkToCrawl := countrySource.sourceURL
				for hasMorePages {
					htmlWithCastlesToCollect, err := se.fetchHTML(ctx, linkToCrawl, se.httpClient)
					if err != nil {
						errChan <- err
						break
					}

					castles, err := se.collectCastleNameAndLinks(htmlWithCastlesToCollect, countrySource.country)
					if err != nil {
						errChan <- err
						break
					}

					for _, c := range castles {
						castlesToEnrichChan <- c
					}

					hasMorePages, linkToCrawl = se.checkForNextPage(htmlWithCastlesToCollect)
				}

				countriesCounter++
			}
		}
	}()

	return castlesToEnrichChan, errChan
}

func (se *ebidatEnricher) collectCastleNameAndLinks(htmlContent []byte, country castle.Country) ([]castle.Model, error) {
	var castles []castle.Model
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %v", err)
	}

	doc.Find(".mainContent .burgenanzeige .burgenanreisser").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").First()
		name := link.Text()
		href, _ := link.Attr("href")
		if !strings.HasPrefix(href, "http") {
			href = ebidatHost + href
		}

		if !strings.Contains(href, "https://") {
			href = "https://" + href
		}

		castle := castle.Model{
			Name:                    name,
			CurrentEnrichmentLink:   href,
			Country:                 country,
			Sources:                 []string{href},
			CurrentEnrichmentSource: EDBIDAT.String(),
		}
		castles = append(castles, castle)
	})

	return castles, nil
}

/*
https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=..%2Fr30%2Fvc_shop%2F&var_datei_selektionen=20240613%2F212718770666ad15f4ef739.dat&var_anzahl_angezeigte_saetze=10

host: https://www.ebidat.de
resource: cgi-bin/r30msvcshop_anzeige.pl
var_hauptpfad: ..%2Fr30%2Fvc_shop%2F -->(../r30/vc_shop/)
var_datei_selektionen: 20240613%2F212718770666ad15f4ef739.dat
var_anzahl_angezeigte_saetze: 10

https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=../r30/vc_shop/&var_datei_selektionen=20240614%2F212718770666b6f64427d52.dat&var_anzahl_angezeigte_saetze=10
*/
func (se *ebidatEnricher) checkForNextPage(htmlContent []byte) (bool, string) {
	currentPage, err := se.getCurrentPage(htmlContent)
	if err != nil {
		return false, ""
	}
	nextPage := currentPage + 1
	formNameToExtractNonce := fmt.Sprintf("formseite%d", nextPage)
	found, nonce := se.getNonce(htmlContent, formNameToExtractNonce)
	if !found {
		return false, ""
	}
	return true, fmt.Sprintf("https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=../r30/vc_shop/&var_datei_selektionen=%s&var_anzahl_angezeigte_saetze=%s", nonce, se.parsePageNumber(nextPage))
}

func (se *ebidatEnricher) parsePageNumber(page int) string {
	if page == 1 {
		return "00"
	}
	return fmt.Sprintf("%d", (page-1)*10)
}

func (se *ebidatEnricher) getCurrentPage(htmlContent []byte) (int, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return 0, err
	}

	var currentPage int

	var traverseErgebnis func(*html.Node)
	traverseErgebnis = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "b" && n.FirstChild != nil {
			pageNumber, err := strconv.Atoi(n.FirstChild.Data)
			if err == nil {
				currentPage = pageNumber
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverseErgebnis(c)
		}
	}

	var parse func(*html.Node)
	parse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "section" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && attr.Val == "ergebnis" {
					traverseErgebnis(n)
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}

	parse(doc)
	return currentPage, nil
}

func (se *ebidatEnricher) getNonce(htmlContent []byte, formName string) (bool, string) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return false, ""
	}

	var findFormValue func(*html.Node, string) (string, bool)
	findFormValue = func(n *html.Node, formName string) (string, bool) {
		if n.Type == html.ElementNode && n.Data == "form" {
			var nameAttr string
			for _, attr := range n.Attr {
				if attr.Key == "name" {
					nameAttr = attr.Val
				}
			}
			if nameAttr == formName {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "input" {
						for _, inputAttr := range c.Attr {
							if inputAttr.Key == "name" && inputAttr.Val == "var_datei_selektionen" {
								for _, inputAttr := range c.Attr {
									if inputAttr.Key == "value" {
										return inputAttr.Val, true
									}
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if val, found := findFormValue(c, formName); found {
				return val, found
			}
		}
		return "", false
	}

	val, found := findFormValue(doc, formName)
	return found, val
}

func (se *ebidatEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	enrichmentURL := fmt.Sprintf("%s&m=h", c.CurrentEnrichmentLink)
	dataHTML, err := se.fetchHTML(ctx, enrichmentURL, se.httpClient)
	if err != nil {
		return castle.Model{}, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(dataHTML))
	if err != nil {
		return castle.Model{}, fmt.Errorf("error loading HTML from [%s]: %v", enrichmentURL, err)
	}

	c1 := &c
	c1.PropertyCondition = se.getPropertyConditions(doc)
	c1.PictureURL = se.collectImage(doc)
	c1.Coordinates = se.collectCoordinates(doc)
	c1.FoundationPeriod = se.collectPeriod(doc)

	c1.State = se.collectState(doc, c.Country)
	c1.City = se.collectCity(doc)
	c1.District = se.collectDistrict(doc)

	c1.CleanFields()
	return *c1, nil
}

func (se ebidatEnricher) collectCity(doc *goquery.Document) string {
	var city string
	doc.Find("li.daten").Each(func(i int, s *goquery.Selection) {
		label := s.Find(".gruppe").Text()
		if strings.Contains(label, "Stadt / Gemeinde:") {
			city = s.Find(".gruppenergebnis").Text()
		}
	})
	city = html.UnescapeString(city)
	city = strings.TrimSpace(city)
	return city
}

func (se ebidatEnricher) collectDistrict(doc *goquery.Document) string {
	var city string
	doc.Find("li.daten").Each(func(i int, s *goquery.Selection) {
		label := s.Find(".gruppe").Text()
		if strings.Contains(label, "Gemarkung / Ortsteil:") {
			city = s.Find(".gruppenergebnis").Text()
		}
	})
	city = html.UnescapeString(city)
	city = strings.TrimSpace(city)
	return city
}

func (se ebidatEnricher) collectState(doc *goquery.Document, country castle.Country) string {
	return stateCollector[country](doc)
}

func (se ebidatEnricher) collectPeriod(doc *goquery.Document) string {
	var periodText string
	// using "Datierung-Beginn:" as reference
	doc.Find("li.daten").Each(func(i int, s *goquery.Selection) {
		label := s.Find(".gruppe").Text()
		if strings.Contains(label, "Datierung-Beginn:") {
			periodText = s.Find(".gruppenergebnis").Text()
		}
	})
	periodText = strings.TrimSpace(periodText)
	re := regexp.MustCompile(`(\d{1,2})\.Jh\.`)
	matches := re.FindStringSubmatch(periodText)
	if len(matches) > 1 {
		return matches[1] + "th"
	}
	return ""
}

func (se ebidatEnricher) getPropertyConditions(doc *goquery.Document) castle.PropertyCondition {
	propertyCondition := castle.Unknown

	doc.Find("div.mainContent section article.beschreibung ul li.daten").Each(func(i int, s *goquery.Selection) {
		gruppe := s.Find("div.gruppe").Text()
		if strings.Contains(gruppe, "Erhaltung - Heutiger Zustand:") {
			collectedCondition := strings.ToLower(s.Find("div.gruppenergebnis").Text())
			switch collectedCondition {
			case "weitgehend erhalten": //largely preserved
				propertyCondition = castle.Intact
			case "stark historisierend überformt": //  heavily historicized
				propertyCondition = castle.Intact
			case "überbaut": //built over
				propertyCondition = castle.Intact
			case "geringe reste": // sall residues
				propertyCondition = castle.Ruins
			case "bedeutende reste": //significant remains
				propertyCondition = castle.Ruins
			case "fundamente": //foundations
				propertyCondition = castle.Ruins
			default:
				propertyCondition = castle.Unknown
			}
			return
		}
	})

	return propertyCondition
}

func (se ebidatEnricher) collectImage(doc *goquery.Document) string {
	var imageSrc string
	doc.Find("div.galerie img").EachWithBreak(func(i int, s *goquery.Selection) bool {
		imageSrc, _ = s.Attr("src")
		return false
	})
	collectedImageLink := fmt.Sprintf("%s%s", ebidatHost, strings.ReplaceAll(imageSrc, "..", ""))
	if !strings.Contains(collectedImageLink, "https://") {
		return fmt.Sprintf("https://%s", collectedImageLink)
	}
	return collectedImageLink
}

func (se ebidatEnricher) collectCoordinates(doc *goquery.Document) string {
	var coordinates string
	doc.Find("#verlinkungen .informationen_link a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Google Maps" {
			coordinates, _ = s.Attr("href")
			return
		}
	})
	parts := strings.Split(coordinates, "q=")
	if len(parts) > 0 {
		return parts[1]
	}
	return ""
}
