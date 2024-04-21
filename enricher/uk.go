package enricher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
)

func getHTMLPageOfUKCastle(ctx context.Context, castle castle.Model, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", castle.Link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GET request for castle [%s], got %v", castle.Name, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get the HTML of castle [%s] from [%s], got %v", castle.Name, castle.Link, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body of GET at [%s] for castle [%s], got %v", castle.Link, castle.Name, err)
	}
	return rawBody, nil
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
func extractUkState(rawHTML []byte, c castle.Model) (string, error) {
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
func extractUkCity(rawHTML []byte, c castle.Model) (string, error) {
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

func extractDataOfUKCastle(rawHTML []byte, c castle.Model) (castle.Model, error) {
	state, err := extractUkState(rawHTML, c)
	if err != nil {
		return castle.Model{}, err
	}
	city, err := extractUkCity(rawHTML, c)
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

func enrichCastleFromUK(
	ctx context.Context,
	httpClient *http.Client,
	c castle.Model,
) (castle.Model, error) {
	castlePage, err := getHTMLPageOfUKCastle(ctx, c, httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := extractDataOfUKCastle(castlePage, c)
	if err != nil {
		return castle.Model{}, err
	}
	fmt.Println("finished castle", enrichedCastled)
	return enrichedCastled, nil
}

func EnrichCastleFromUK(
	ctx context.Context,
	httpClient *http.Client,
	c castle.Model,
) (castle.Model, error) {
	castlePage, err := getHTMLPageOfUKCastle(ctx, c, httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := extractDataOfUKCastle(castlePage, c)
	if err != nil {
		return castle.Model{}, err
	}
	fmt.Println("finished castle", enrichedCastled)
	return enrichedCastled, nil
}
