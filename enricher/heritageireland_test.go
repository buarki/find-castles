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
	irishCastlesHomePageHTMLPath = "../data/ireland-list.html"
	expectedIrishCastlesAsJSON   = "../data/ireland-list.json"
	irishCastlePageHTMLPath      = "../data/ireland-castle.html"
)

func TestCollectIrishCastlesToEnrich(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(irishCastlesHomePageHTMLPath)
	}
	expectedCastles, err := fileloader.LoadCastlesAsJSONList(expectedIrishCastlesAsJSON)
	if err != nil {
		t.Fatalf("expected to have err nil, got %v", err)
	}

	irishCollector := NewHeritageIreland(httpclient.New(), htmlFetcher)

	castlesChan, errChan := irishCollector.CollectCastlesToEnrich(context.Background())
	if len(errChan) > 0 {
		t.Fatalf("expected no err, got %d", len(errChan))
	}
	var foundCastles []castle.Model
	for c := range castlesChan {
		foundCastles = append(foundCastles, c)
	}

	if !slicesWithSameContent(foundCastles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}

func TestEnrich(t *testing.T) {
	testCases := []struct {
		html   string
		castle castle.Model
	}{
		{
			html: `
			<div id="place--contact">
			<div>
				<h2>Contact</h2>
				<p class="address">Adare Heritage Centre<br />
					Adare<br />
					Co. Limerick<br />
					V94 DWV7</p>
				<p class="phone">061 396 666</p>
				<p class="email"><a href="mailto:reception@adareheritagecentre.ie">reception@adareheritagecentre.ie</a></p>
			</div>
			<div>
			`,
			castle: castle.Model{
				Name:     "adare",
				District: "adare heritage centre",
				City:     "adare",
				State:    "co. limerick",
			},
		},
		{
			html: `
			<div id="place--contact">
				<div>
					<h2>Contact</h2>
					<p class="address">Trim <br />
						Co Meath<br>C15 HN90</p>
					<p class="phone">046 9438619</p>
					<p class="email"><a href="mailto:trimcastle@opw.ie">trimcastle@opw.ie</a></p>
				</div>
			<div>
			`,
			castle: castle.Model{
				Name:     "trim",
				District: "trim",
				City:     "trim",
				State:    "co meath",
			},
		},
		{
			html: `
			<div id="place--contact">
				<div>
					<h2>Contact</h2>
					<p class="address">Ross Castle, <br />
						Ross Road, <br />
						Killarney, <br />
						Co. Kerry<br>V93 V304</p>
					<p class="phone">064 663 5851</p>
					<p class="email"><a href="mailto:rosscastle@opw.ie">rosscastle@opw.ie</a></p>
				</div>
			<div>`,
			castle: castle.Model{
				Name: "Ross Castle",

				District: "ross road",
				City:     "killarney",
				State:    "co. kerry",
			},
		},
	}

	for _, tt := range testCases {
		fetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
			return []byte(tt.html), nil
		}
		enricher := NewHeritageIreland(httpclient.New(), fetcher)

		castle, err := enricher.EnrichCastle(context.Background(), tt.castle)
		if err != nil {
			t.Errorf("expecte err nil, got %v", err)
		}

		if castle.City != tt.castle.City {
			t.Errorf("expected city [%s], got [%s]", tt.castle.City, castle.City)
		}
		if castle.State != tt.castle.State {
			t.Errorf("expected State [%s], got [%s]", tt.castle.State, castle.State)
		}
		if castle.District != tt.castle.District {
			t.Errorf("expected District [%s], got [%s]", tt.castle.District, castle.District)
		}
	}
}
