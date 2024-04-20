package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
)

const (
	UK                     = "UK"
	EnglandCastlesListPage = "https://medievalbritain.com/medieval-castles-of-england"
)

func getHTMLHavingTheListOfCastlesInEngland(ctx context.Context, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", EnglandCastlesListPage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the GET request for [%s] to get the list of castles in england, got %v", EnglandCastlesListPage, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get the HTML with the list of castles in england by doing GET at [%s], got %v", EnglandCastlesListPage, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body of GET at [%s], got %v", EnglandCastlesListPage, err)
	}
	return rawBody, nil
}

func extractTheListOfCastlesInEngland(rawHTML []byte) ([]castle.Model, error) {
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
			Country:  UK,
			FlagLink: "/uk-flag.webp",
		})
	})
	return castles, nil
}

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

func enrichEnglandCastleData(ctx context.Context, castle castle.Model, httpClient *http.Client, results chan CollectResult) {
	castleHTMLPage, err := getHTMLPageOfUKCastle(ctx, castle, httpClient)
	if err != nil {
		results <- CollectResult{
			Err:    err,
			Castle: castle,
		}
	}
	castleData, err := extractDataOfUKCastle(castleHTMLPage, castle)
	if err != nil {
		results <- CollectResult{
			Err:    err,
			Castle: castle,
		}
	} else {
		results <- CollectResult{
			Castle: castleData,
		}
	}
	fmt.Println("finished castle", castle.Name)
}

func CollectForUk(ctx context.Context, httpClient *http.Client, appWg *sync.WaitGroup, results chan CollectResult) {
	defer appWg.Done()

	htmlOfCastlesInEngland, err := getHTMLHavingTheListOfCastlesInEngland(ctx, httpClient)
	if err != nil {
		fmt.Println("error", err)
		results <- CollectResult{
			Err: err,
		}
	}

	englandCastles, err := extractTheListOfCastlesInEngland(htmlOfCastlesInEngland)
	if err != nil {
		results <- CollectResult{
			Err: err,
		}
	}

	availableCPUs := runtime.NumCPU()
	semaphore := make(chan struct{}, availableCPUs)

	var wg sync.WaitGroup

	for _, foundCastle := range englandCastles {
		select {
		case <-ctx.Done():
			fmt.Println("BYE UK")
			return
		default:
			wg.Add(1)
			semaphore <- struct{}{}

			go func(c castle.Model) {
				enrichEnglandCastleData(ctx, c, httpClient, results)
				time.Sleep(1 * time.Second)
				<-semaphore
				wg.Done()
			}(foundCastle)
		}
	}

	wg.Wait()
}
