package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/htttpclient"
)

const (
	castlesSource = "https://www.castelosdeportugal.pt"
)

func collectHomePageHTML(link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get home, got %v", err)
	}
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

type collectResult struct {
	castle castle.Model
	err    error
}

func getCastleHTMLPage(c castle.Model, link string, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get home of castle [%s], got %v", c.Name, err)
	}
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
		Country:          "portugal",
		Link:             c.Link,
		City:             tableData["Concelho"],
		State:            tableData["Distrito"],
		District:         tableData["Freguesia"],
		YearOfFoundation: tableData["Construção"],
	}, nil
}

func collectCastleInfo(castle castle.Model, collectedCastles chan collectResult, httpClient *http.Client) {
	fmt.Println("Processing castle", castle.Name)
	castlePageLink := fmt.Sprintf("%s/castelos/%s", castlesSource, strings.ReplaceAll(castle.Link, "../", ""))
	fmt.Println("castlePageLink", castlePageLink)
	castlePage, err := getCastleHTMLPage(castle, castlePageLink, httpClient)
	if err != nil {
		collectedCastles <- collectResult{
			castle: castle,
			err:    err,
		}
	} else {
		enrichedCastled, err := extractCastleInfo(castle, castlePage)
		if err != nil {
			collectedCastles <- collectResult{
				castle: castle,
				err:    err,
			}
		} else {
			collectedCastles <- collectResult{
				castle: enrichedCastled,
			}
		}
		fmt.Println("finished castle", castle.Name)
	}
}

func main() {
	httpClient := htttpclient.New()

	homePage, err := collectHomePageHTML(castlesSource, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	castles, err := collectCastleNameAndLinks(homePage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("collected castles", len(castles))

	availableCPUS := runtime.NumCPU()
	fmt.Println("availableCPUS", availableCPUS)

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, availableCPUS)

	collectResults := make(chan collectResult)

	var collectedCastles []castle.Model

	go func() {
		fmt.Println("waiting...")
		for result := range collectResults {
			if result.err != nil {
				fmt.Printf("failed to collect info for castle [%s], got %v\n", result.castle.Name, err)
			} else {
				fmt.Printf("%v\n", result.castle)
				collectedCastles = append(collectedCastles, result.castle)
			}
		}
	}()

	for i := 0; i < len(castles); i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(c castle.Model) {
			collectCastleInfo(c, collectResults, httpClient)
			time.Sleep(1 * time.Second)
			<-semaphore
			wg.Done()
		}(castles[i])
	}

	fmt.Println("waiting")
	wg.Wait()
	fmt.Println("waited...")
	close(collectResults)
	fmt.Println("closed...")

	b, err := json.Marshal(collectedCastles)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("pt.json", b, 0777)
}
