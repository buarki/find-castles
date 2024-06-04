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

func enrichCastleFromIreland(
	ctx context.Context,
	httpClient *http.Client,
	c castle.Model,
) (castle.Model, error) {
	castlePage, err := getIrelandCastleHTMLPage(ctx, c, httpClient)
	if err != nil {
		return castle.Model{}, nil
	}
	enrichedCastled, err := extractIrelandCastleInfo(c, castlePage)
	if err != nil {
		return castle.Model{}, err
	}
	return enrichedCastled, nil
}

func getIrelandCastleHTMLPage(ctx context.Context, c castle.Model, httpClient *http.Client) ([]byte, error) {
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

func extractContact(rawHTML []byte) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(rawHTML))
	if err != nil {
		return "", fmt.Errorf("error loading HTML: %v", err)
	}

	address, err := doc.Find("#place--contact div p.address").First().Html()
	if err != nil {
		return "", err
	}

	return address, nil
}

func get(raw string) (string, string, string) {
	parts := strings.Split(raw, "<br/>")
	for i := range parts {
		parts[i] = strings.TrimSuffix(strings.ReplaceAll(parts[i], ",", ""), " ")
	}
	districtAndCityAreEqual := len(parts) < 4
	if districtAndCityAreEqual {
		return parts[0], parts[0], parts[len(parts)-2]
	}
	return parts[0], parts[1], parts[len(parts)-2]
}

func extractIrelandCastleInfo(c castle.Model, rawHTMLPage []byte) (castle.Model, error) {
	rawAddress, err := extractContact(rawHTMLPage)
	if err != nil {
		return castle.Model{}, err
	}

	district, city, state := get(rawAddress)

	return castle.Model{
		Name:     c.Name,
		Country:  castle.Ireland,
		Link:     c.Link,
		City:     city,
		State:    state,
		District: district,
		FlagLink: c.FlagLink,
	}, nil
}
