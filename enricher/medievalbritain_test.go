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
	"github.com/google/go-cmp/cmp"
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
		PictureURL:  "https://medievalbritain.com/wp-content/uploads/2021/05/medieval-castles-england_windsor.jpg",
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
	if receivedCastle.PictureURL != expectedCastle.PictureURL {
		t.Errorf("expected PictureURL to be [%s], got [%s]", expectedCastle.PictureURL, receivedCastle.PictureURL)
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

func TestExtractContactOfMedievalBritain(t *testing.T) {
	testCases := []struct {
		name            string
		htmlChunk       []byte
		expectedContact *castle.Contact
	}{
		{
			name: "schema with number between p",
			htmlChunk: []byte(`
			<div class="elementor-text-editor elementor-clearfix">
				<p><span class="w8qArf"><strong>Phone</strong></span>
				</p>
				<p>+44 (0)303 123 7304</p>
			</div>
			`),
			expectedContact: &castle.Contact{
				Phone: "+44 (0)303 123 7304",
			},
		},
		// {
		// 	name: "schema with number between p",
		// 	htmlChunk: []byte(`
		// 	<div class="elementor-text-editor elementor-clearfix">
		// 		<p><span class="w8qArf"><strong>Phone</strong></span></p>
		// 		<div class="Z0LcW">
		// 			<span data-dtype="d3ifr" data-local-attribute="d3ph">0370 333 1181</span>
		// 		</div>
		// 		</div>
		// 		</div>
		// 	`),
		// 	expectedContact: &castle.Contact{
		// 		Phone: "0370 333 1181",
		// 	},
		// },
	}
	e := medievalbritainEnricher{}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(currentTT.htmlChunk))
			if err != nil {
				t.Errorf("expected to have err nil, got [%v]", err)
			}

			collectedContact := e.collectContactInfo(doc)

			if collectedContact == nil {
				t.Errorf("expected contact to not be null")
			}

			if collectedContact != nil && collectedContact.Phone != currentTT.expectedContact.Phone {
				t.Errorf("expected Phone to be [%s], got [%s]", currentTT.expectedContact.Phone, collectedContact.Email)
			}
		})
	}
}

func TestExtracVisitingInfotOfMedievalBritain(t *testing.T) {
	content, err := fileloader.LoadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
	if err != nil {
		t.Errorf("expected to have err nil, got [%v]", err)
	}
	expectedVisitingInfo := &castle.VisitingInfo{
		WorkingHours: "Summer: 10:00 - 16:00,Winter: 10:00 - 15:00",
		Facilities: &castle.Facilities{
			AssistanceDogsAllowed: true,
			Giftshops:             true,
			WheelchairSupport:     true,
			Restrooms:             true,
			PinicArea:             true,
			Exhibitions:           true,
			Cafe:                  true,
			Parking:               true,
		},
	}
	e := medievalbritainEnricher{}

	collectedVisitingInfo := e.collectVisitingInfo(doc)

	if collectedVisitingInfo == nil {
		t.Errorf("expected contact to not be null")
	}

	diff := cmp.Diff(expectedVisitingInfo, collectedVisitingInfo)
	if diff != "" {
		t.Errorf("diff: %v", diff)
	}
}

func TestCollectWorkingHoursOfMedievalBritain(t *testing.T) {
	testCases := []struct {
		name                 string
		htmlChunk            []byte
		expectedWorkingHours string
	}{
		{
			name: "default working hours scheme",
			htmlChunk: []byte(`
			<div class="elementor-widget-container">
				<div
					class="elementor-text-editor elementor-clearfix">
					<p><strong>Hours</strong></p>
					<p>Summer: 10:00 &#8211; 16:00</p>
					<p>Winter: 10:00 &#8211; 15:00</p>
				</div>
			</div>
			`),
			expectedWorkingHours: "Summer: 10:00 - 16:00,Winter: 10:00 - 15:00",
		},
		{
			name: "horking hours is pure text",
			htmlChunk: []byte(`
			<div class="elementor-widget-container">
					<div class="elementor-text-editor elementor-clearfix">
						<p><strong>Hours</strong></p>
						<div class="soft bg-brand-inverlochy">
							<p class="beta white">Castle Sween is open year-round.</p>
						</div>
					</div>
				</div>
			`),
			expectedWorkingHours: "Castle Sween is open year-round.",
		},
		{
			name: "horking hours is pure text (1)",
			htmlChunk: []byte(`
			<div class="elementor-widget-container">
				<div class="elementor-text-editor elementor-clearfix">
					<p><strong>Hours</strong></p>
					<p>Open 24 hours, year-round.</p>
				</div>
			</div>
			`),
			expectedWorkingHours: "Open 24 hours, year-round.",
		},
		{
			name: "when working hours is given on a single line table",
			htmlChunk: []byte(`
			<div class="elementor-widget-container">
				<div class="elementor-text-editor elementor-clearfix">
					<p><strong>Hours</strong></p>
					<table class="WgFkxc">
						<tbody>
							<tr class="K7Ltle">
								<td class="SKNSIb">Summer</td>
								<td>10am–4pm</td>
							</tr>
							<tr>
								<td class="SKNSIb">Winter</td>
								<td>10am–5pm</td>
							</tr>
						</tbody>
					</table>
				</div>
			</div>
			`),
			expectedWorkingHours: "Summer: 10am-4pm, Winter: 10am-5pm",
		},
		{
			name: "when working hours is given on a multi line table",
			htmlChunk: []byte(`
			<div class="elementor-text-editor elementor-clearfix">
				<p><strong>Hours</strong></p>
				<table class="WgFkxc">
					<tbody>
						<tr class="K7Ltle">
							<td class="SKNSIb">Tuesday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Wednesday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Thursday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Friday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Saturday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Sunday</td>
							<td>11am–4pm</td>
						</tr>
						<tr>
							<td class="SKNSIb">Monday</td>
							<td>Closed</td>
						</tr>
					</tbody>
				</table>
			</div>
			`),
			expectedWorkingHours: "Tuesday: 11am-4pm, Wednesday: 11am-4pm, Thursday: 11am-4pm, Friday: 11am-4pm, Saturday: 11am-4pm, Sunday: 11am-4pm, Monday: Closed",
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

			collectedContact := e.collectHorkingHours(doc)

			if collectedContact != currentTT.expectedWorkingHours {
				t.Errorf("expected working hours to be [%s], got [%s]", currentTT.expectedWorkingHours, collectedContact)
			}
		})
	}
}
