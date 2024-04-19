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
	castlesSource = "https://www.castelosdeportugal.pt"
)

func collectHomePageHTML(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get home, got %v", err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET at [%s], got %v", link, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content, got %v", err)
	}
	return rawBody, nil
}

func collectCastleNameAndLinks(rawHTML []byte) ([]castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return nil, fmt.Errorf("error loading HTML: %v", err)
	}
	var castles []castle.Model
	doc.Find("#div-list-alfa-cast p a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		name := s.Text()
		castles = append(castles, castle.Model{Name: name, Link: link})
	})
	return castles, nil
}

func getCastleHTMLPage(ctx context.Context, c castle.Model, link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get home of castle [%s], got %v", c.Name, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET at [%s] for castle [%s], got %v", link, c.Name, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content of castle [%s], got %v", c.Name, err)
	}
	return rawBody, nil
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func extractCastleInfo(c castle.Model, rawHTMLPage []byte) (castle.Model, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTMLPage))
	if err != nil {
		return castle.Model{}, fmt.Errorf("failed to load page, got %v", err)
	}

	var tableData = make(map[string]string)

	rowsToExtract := []string{"Distrito", "Concelho", "Freguesia", "Construção"}

	doc.Find("#info-table tbody tr").Each(func(i int, s *goquery.Selection) {
		key := strings.TrimSpace(s.Find("td:nth-child(1)").Text())
		if contains(rowsToExtract, key) {
			value := strings.TrimSpace(s.Find("td:nth-child(2)").Text())
			tableData[key] = value
		}
	})

	fmt.Println("Table Data:", tableData)
	return castle.Model{
		Name:             c.Name,
		Country:          "Portugal",
		Link:             fmt.Sprintf("%s/castelos/%s", castlesSource, strings.ReplaceAll(c.Link, "../", "")),
		City:             tableData["Concelho"],
		State:            tableData["Distrito"],
		District:         tableData["Freguesia"],
		YearOfFoundation: tableData["Construção"],
		FlagLink:         "/pt-flag.webp",
	}, nil
}

func collectCastleInfo(ctx context.Context, castle castle.Model, collectedCastles chan CollectResult, httpClient *http.Client) {
	fmt.Println("Processing castle", castle.Name)
	castlePageLink := fmt.Sprintf("%s/castelos/%s", castlesSource, strings.ReplaceAll(castle.Link, "../", ""))
	fmt.Println("castlePageLink", castlePageLink)
	castlePage, err := getCastleHTMLPage(ctx, castle, castlePageLink, httpClient)
	if err != nil {
		collectedCastles <- CollectResult{
			Castle: castle,
			Err:    err,
		}
	} else {
		enrichedCastled, err := extractCastleInfo(castle, castlePage)
		if err != nil {
			collectedCastles <- CollectResult{
				Castle: castle,
				Err:    err,
			}
		} else {
			collectedCastles <- CollectResult{
				Castle: enrichedCastled,
			}
		}
		fmt.Println("finished castle", castle.Name)
	}
}

func CollectForPotugal(ctx context.Context, httpClient *http.Client, appWg *sync.WaitGroup, results chan CollectResult) {
	defer appWg.Done()

	homePage, err := collectHomePageHTML(ctx, castlesSource, httpClient)
	if err != nil {
		results <- CollectResult{
			Err: err,
		}
	}

	castles, err := collectCastleNameAndLinks(homePage)
	if err != nil {
		results <- CollectResult{
			Err: err,
		}
	}
	fmt.Println("collected castles", len(castles))

	availableCPUS := runtime.NumCPU()
	fmt.Println("availableCPUS", availableCPUS)

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, availableCPUS)

	for i := 0; i < len(castles); i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Request canceled")
			return
		default:
			wg.Add(1)
			semaphore <- struct{}{}
			go func(c castle.Model) {
				collectCastleInfo(ctx, c, results, httpClient)
				time.Sleep(1 * time.Second)
				<-semaphore
				wg.Done()
			}(castles[i])
		}
	}
	wg.Wait()
}
