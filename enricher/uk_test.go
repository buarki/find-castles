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
	htmlWithEnglandCastlesPath    = "../data/uk-england-list.html"
	jsonWithEnglandCastlesPath    = "../data/uk-england-list.json"
	jsonWithEnglandCastlePagePath = "../data/uk-england-castle.html"
)

func TestCollectBritishCastlesToEnrich(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(htmlWithEnglandCastlesPath)
	}
	expectedCastles, err := fileloader.LoadCastlesAsJSONList(jsonWithEnglandCastlesPath)
	if err != nil {
		t.Fatalf("expected to have err nil, got %v", err)
	}

	britishCollector := NewBritishEnricher(httpclient.New(), htmlFetcher)

	foundCastles, err := britishCollector.CollectCastlesToEnrich(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !slicesWithSameContent(foundCastles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}

func TestExtractBritishCastleInfo(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(jsonWithEnglandCastlePagePath)
	}

	expectedCastle := castle.Model{
		Country: castle.UK,
		City:    "Windsor SL4 1NJ",
		State:   "Berkshire, GreaterLondon",
	}

	britishCollector := NewBritishEnricher(httpclient.New(), htmlFetcher)

	receivedCastle, err := britishCollector.EnrichCastle(context.Background(), expectedCastle)
	if err != nil {
		t.Errorf("expected err nil, got %v", err)
	}

	if receivedCastle.City != expectedCastle.City {
		t.Errorf("expected city to be [%s], got [%s]", expectedCastle.City, receivedCastle.City)
	}
	if receivedCastle.State != expectedCastle.State {
		t.Errorf("expected State to be [%s], got [%s]", expectedCastle.State, receivedCastle.State)
	}
}
