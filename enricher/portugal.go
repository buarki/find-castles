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

const (
	castlesSource = "https://www.castelosdeportugal.pt"
)

func enrichCastleFromPortugal(
	ctx context.Context,
	httpClient *http.Client,
	c castle.Model,
) (castle.Model, error) {
	castlePage, err := getCastleHTMLPage(ctx, c, httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := extractCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	return enrichedCastled, nil
}

func getCastleHTMLPage(ctx context.Context, c castle.Model, httpClient *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Link, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get home of castle [%s], got %v", c.Name, err)
	}
	req = req.WithContext(ctx)
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET at [%s] for castle [%s], got %v", c.Link, c.Name, err)
	}
	defer res.Body.Close()
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body content of castle [%s], got %v", c.Name, err)
	}
	return rawBody, nil
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

	fmt.Println("Table Data:", c.Name, tableData)
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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
