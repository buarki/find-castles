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

func TestExtractPictureOfHeritageIreland(t *testing.T) {
	content, err := fileloader.LoadHTMLFile(irishCastlePageHTMLPath)
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}
	expectedImageLink := `https://heritageireland.ie/assets/uploads/2020/03/Adare-Castle-Aerial-View-2-640x427.jpg`
	e := heritageirelandEnricher{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}

	collectedImageLink := e.collectImage(doc)

	if collectedImageLink != expectedImageLink {
		t.Errorf("expected to find link [%s], got [%s]", expectedImageLink, collectedImageLink)
	}
}

func TestExtractContactOfHeritageIreland(t *testing.T) {
	content, err := fileloader.LoadHTMLFile(irishCastlePageHTMLPath)
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}
	expectedContact := &castle.Contact{
		Email: "reception@adareheritagecentre.ie",
		Phone: "061 396 666",
	}
	e := heritageirelandEnricher{}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}

	collectedContact := e.collectContactInfo(doc)

	if collectedContact == nil {
		t.Errorf("expected contact to not be null")
	}

	if collectedContact != nil && collectedContact.Email != expectedContact.Email {
		t.Errorf("expected email to be [%s], got [%s]", expectedContact.Email, collectedContact.Email)
	}

	if collectedContact != nil && collectedContact.Phone != expectedContact.Phone {
		t.Errorf("expected Phone to be [%s], got [%s]", expectedContact.Phone, collectedContact.Email)
	}
}

func TestCollectWorkingHoursOfHeritageIreland(t *testing.T) {
	testCases := []struct {
		name                 string
		htmlChunk            []byte
		expectedWorkingHours string
	}{
		{
			name: "simple",
			htmlChunk: []byte(`
			<section id="place--opening" class="section">
				<h2>Opening Times</h2>
				<div>
					<p><strong>01 June &#8211; 29 September 2024</strong></p>
					<p>09:30- 16:00</p>
				</div>
			</section>
			`),
			expectedWorkingHours: "01 June - 29 September 2024 - 09:30- 16:00",
		},
		{
			name: "simple",
			htmlChunk: []byte(`
			<section id="place--opening" class="section">
        <h2>Opening Times</h2>
        <div>
          <p>15 March- 3 November 2024</p>
        </div>
        <div>
          <h3>Seasonal Opening Times</h3>
          <dl class="accordion">
            <dt>15 March - 26 October<b></b></dt>
            <dd>
              <div>
                <p>Daily 10:00 &#8211; 18:00</p>
                <p>Last admission: 17:15</p>
              </div>
            </dd>
            <dt>27 October - 03 November <b></b></dt>
            <dd>
              <div>
                <p>Daily 10:00 &#8211; 17:00</p>
                <p>Last admission: 16:15</p>
              </div>
            </dd>
          </dl>
        </div>
      </section>
			`),
			expectedWorkingHours: "15 March- 3 November 2024",
		},
	}
	e := heritageirelandEnricher{}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.name, func(t *testing.T) {
			t.Helper()
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(currentTT.htmlChunk))
			if err != nil {
				t.Errorf("expected to have err nil, got [%v]", err)
			}
			collectedWorkingHours := e.collectHorkingHours(doc)

			if collectedWorkingHours != currentTT.expectedWorkingHours {
				t.Errorf("expected to have [%s], got [%s]", currentTT.expectedWorkingHours, collectedWorkingHours)
			}
		})
	}

}
