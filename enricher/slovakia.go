package enricher

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/htmlfetcher"
	"golang.org/x/net/html"
)

const (
	sloavkHost   = "www.ebidat.de"
	slovakSource = "https://" + sloavkHost + "/cgi-bin/ebidat.pl?a=a&te53=6"
)

type extractLocation struct {
	state    string
	city     string
	district string
}

type slovakEnricher struct {
	httpClient *http.Client
	fetchHTML  htmlfetcher.HTMLFetcher
}

func NewSlovakEnricher(
	httpClient *http.Client,
	fetchHTML htmlfetcher.HTMLFetcher) Enricher {
	return &slovakEnricher{
		httpClient: httpClient,
		fetchHTML:  fetchHTML,
	}
}

func (se *slovakEnricher) CollectCastlesToEnrich(ctx context.Context) (chan castle.Model, chan error) {
	castlesToEnrichChan := make(chan castle.Model)
	errChan := make(chan error)

	go func() {
		defer close(castlesToEnrichChan)
		defer close(errChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				hasMorePages := true
				linkToCrawl := slovakSource

				for hasMorePages {
					htmlWithCastlesToCollect, err := se.fetchHTML(ctx, linkToCrawl, se.httpClient)
					if err != nil {
						errChan <- err
						return
					}
					castles, err := se.collectCastleNameAndLinks(htmlWithCastlesToCollect)
					if err != nil {
						errChan <- err
						return
					}
					for _, c := range castles {
						castlesToEnrichChan <- c
					}
					hasMorePages, linkToCrawl = se.checkForNextPage(htmlWithCastlesToCollect)
					if !hasMorePages {
						return
					}
				}
			}
		}

	}()

	return castlesToEnrichChan, errChan
}

func (se *slovakEnricher) collectCastleNameAndLinks(htmlContent []byte) ([]castle.Model, error) {
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
			href = sloavkHost + href
		}

		data := se.extractDistrictCityAndState(s.Text())
		castle := castle.Model{
			Name:     name,
			Link:     href,
			Country:  castle.Slovakia,
			City:     data.city,
			State:    data.state,
			District: data.district,
			FlagLink: "/sk.png",
		}
		castles = append(castles, castle)
	})

	return castles, nil
}

func (se *slovakEnricher) extractDistrictCityAndState(rawLocation string) extractLocation {
	mainPart := rawLocation[:strings.Index(rawLocation, "if")]
	splited := strings.Split(mainPart, "\n")
	var noEmptySpaces []string
	for _, s := range splited {
		s1 := strings.ReplaceAll(s, "\t", "")
		if len(s1) > 0 {
			noEmptySpaces = append(noEmptySpaces, s1)
		}
	}
	if len(noEmptySpaces) == 0 {
		return extractLocation{}
	}
	for i, j := 0, len(noEmptySpaces)-1; i < j; i, j = i+1, j-1 {
		noEmptySpaces[i], noEmptySpaces[j] = noEmptySpaces[j], noEmptySpaces[i]
	}
	if len(noEmptySpaces) == 4 {
		return extractLocation{
			state:    noEmptySpaces[0],
			city:     noEmptySpaces[1],
			district: noEmptySpaces[2],
		}
	}
	if len(noEmptySpaces) == 3 {
		return extractLocation{
			state:    noEmptySpaces[0],
			city:     noEmptySpaces[1],
			district: noEmptySpaces[2],
		}
	}
	return extractLocation{
		state: noEmptySpaces[0],
	}
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
func (se *slovakEnricher) checkForNextPage(htmlContent []byte) (bool, string) {
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

func (se *slovakEnricher) parsePageNumber(page int) string {
	if page == 1 {
		return "00"
	}
	return fmt.Sprintf("%d", (page-1)*10)
}

func (se *slovakEnricher) getCurrentPage(htmlContent []byte) (int, error) {
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

func (se *slovakEnricher) getNonce(htmlContent []byte, formName string) (bool, string) {
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

func (se *slovakEnricher) EnrichCastle(ctx context.Context, c castle.Model) (castle.Model, error) {
	return c, nil
}
