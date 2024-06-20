package enricher

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
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

	britishCollector := NewMedievalBritainEnricher(httpclient.New(), htmlFetcher)

	castlesChan, errChan := britishCollector.CollectCastlesToEnrich(context.Background())
	if len(errChan) > 0 {
		t.Fatalf("expected no errors, got %d", len(errChan))
	}
	var foundCastles []castle.Model
	for c := range castlesChan {
		foundCastles = append(foundCastles, c)
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
		Country:     castle.UK,
		City:        "windsor sl4 1nj",
		State:       "berkshire, greaterlondon",
		PictureLink: "https://medievalbritain.com/wp-content/uploads/2021/05/medieval-castles-england_windsor.jpg",
		Coordinates: "51°29'0\"N,00°36'15\"W",
	}

	britishCollector := NewMedievalBritainEnricher(httpclient.New(), htmlFetcher)

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
	if receivedCastle.PictureLink != expectedCastle.PictureLink {
		t.Errorf("expected PictureLink to be [%s], got [%s]", expectedCastle.PictureLink, receivedCastle.PictureLink)
	}
	if receivedCastle.Coordinates != expectedCastle.Coordinates {
		t.Errorf("expected Coordinates to be [%s], got [%s]", expectedCastle.Coordinates, receivedCastle.Coordinates)
	}
}

func TestExtractPictureOfMedievalBritain(t *testing.T) {
	content, err := fileloader.LoadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}
	expectedImageLink := `https://medievalbritain.com/wp-content/uploads/2021/05/medieval-castles-england_windsor.jpg`
	e := medievalbritainEnricher{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}

	collectedImageLink := e.collectImage(doc)

	if collectedImageLink != expectedImageLink {
		t.Errorf("expected to find link [%s], got [%s]", expectedImageLink, collectedImageLink)
	}
}

func TestCollectingLocalizationCoordinatesOfMedievalBritain(t *testing.T) {
	testCases := []struct {
		name                string
		htmlChunk           []byte
		expectedCoordinates string
	}{
		{
			name: "latitude and longitude together",
			htmlChunk: []byte(`
			<span class="geo-default">
				<span
					class="geo-dec"
					title="Maps, aerial photos, and other data for this location">54.9904°N 2.0000°W</span>
			</span>
			`),
			expectedCoordinates: "54.9904°N 2.0000°W",
		},
		{
			name: "latitude and longitude separated",
			htmlChunk: []byte(`
			<span
				class="geo-default">
				<span
					class="geo-dms"
					title="Maps, aerial photos, and other data for this location">
					<span class="latitude">51°29′0″N</span> <span
						class="longitude">00°36′15″W</span></span></span>
			`),
			expectedCoordinates: "51°29'0\"N,00°36'15\"W",
		},
	}
	e := medievalbritainEnricher{}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(currentTT.htmlChunk))
			if err != nil {
				t.Errorf("expected to have err nil, got [%v]", err)
			}

			receivedCoordinates := e.collectCoordinates(doc)

			if receivedCoordinates != currentTT.expectedCoordinates {
				t.Errorf("expected to find [%s], got [%s]", currentTT.expectedCoordinates, receivedCoordinates)
			}
		})
	}
}
