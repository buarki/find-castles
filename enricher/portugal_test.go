package enricher

import (
	"context"
	"net/http"
	"testing"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fileloader"
	"github.com/buarki/find-castles/httpclient"
)

const (
	castlesHomePageHTMLPath = "../data/portugal.html"
	expectedCastlesAsJSON   = "../data/portugal.json"
	castlePageHTMLPath      = "../data/portugal-castle.html"
)

func slicesWithSameContent(x, y []castle.Model) bool {
	if len(x) != len(y) {
		return false
	}
	diff := make(map[castle.Model]int, len(x))
	for _, _x := range x {
		diff[_x]++
	}
	for _, _y := range y {
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y]--
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

func TestCollectPortugueseCastlesToEnrich(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(castlesHomePageHTMLPath)
	}
	expectedCastles, err := fileloader.LoadCastlesAsJSONList(expectedCastlesAsJSON)
	if err != nil {
		t.Fatalf("expected to have err nil, got %v", err)
	}

	portugueseCollector := NewPortugueseEnricher(httpclient.New(), htmlFetcher)

	foundCastles, err := portugueseCollector.CollectCastlesToEnrich(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !slicesWithSameContent(foundCastles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}

func TestExtractPortugueseCastleInfo(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(castlePageHTMLPath)
	}

	expectedCastle := castle.Model{
		Name:             "Porto",
		Country:          "Portugal",
		City:             "Guimarães",
		State:            "Guimarães",
		District:         "Oliveira do Castelo",
		YearOfFoundation: "(ant. a 958)",
		Link:             "https://somelink.pt",
	}

	portugueseEnricher := NewPortugueseEnricher(httpclient.New(), htmlFetcher)

	receivedCastle, err := portugueseEnricher.EnrichCastle(context.Background(), expectedCastle)
	if err != nil {
		t.Errorf("expected err nil, got %v", err)
	}

	if receivedCastle.City != expectedCastle.City {
		t.Errorf("expected city to be [%s], got [%s]", expectedCastle.City, receivedCastle.City)
	}
	if receivedCastle.District != expectedCastle.District {
		t.Errorf("expected District to be [%s], got [%s]", expectedCastle.District, receivedCastle.District)
	}
	if receivedCastle.State != expectedCastle.State {
		t.Errorf("expected State to be [%s], got [%s]", expectedCastle.State, receivedCastle.State)
	}
	if receivedCastle.Name != expectedCastle.Name {
		t.Errorf("expected Name to be [%s], got [%s]", expectedCastle.Name, receivedCastle.Name)
	}
	if receivedCastle.YearOfFoundation != expectedCastle.YearOfFoundation {
		t.Errorf("expected YearOfFoundation to be [%s], got [%s]", expectedCastle.YearOfFoundation, receivedCastle.YearOfFoundation)
	}
}
