package enricher

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fileloader"
	"github.com/buarki/find-castles/httpclient"
	"github.com/google/go-cmp/cmp"
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

	countX := make(map[string]int)
	countY := make(map[string]int)

	for _, cx := range x {
		key := fmt.Sprintf("%s:%s", cx.Country, cx.Name)
		countX[key]++
	}

	for _, cy := range y {
		key := fmt.Sprintf("%s:%s", cy.Country, cy.Name)
		countY[key]++
	}

	for key, cntX := range countX {
		cntY, ok := countY[key]
		if !ok || cntX != cntY {
			return false
		}
	}

	return true
}

func TestCollectCastleNamesAndLinks(t *testing.T) {
	fakeHTMLFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return []byte(
			`
			<div id="indice">
				<h1>Castelos de Portugal</h1>
				<div class="row">
					<h4>A</h4>
					<a rel="nofollow" href="CastelosSECXII/abrantes.html">Abrantes</a>
					<a rel="nofollow" href="CastelosSECXII/alpalhao.html">Alpalhão</a>
					<a rel="nofollow" href="Castelos(pre)SECXII/atouguiaBaleia.html">Atouguia da Baleia</a>
				</div>
				<div class="row">
					<h4>B</h4>
					<a rel="nofollow" href="CastelosSECXIII/braga.html">Braga</a>
					<a rel="nofollow" href="CastelosSECXIII/braganca.html">Bragança</a>
				</div>
				<div class="row">
					<h4>C</h4>
					<a rel="nofollow" href="CastelosSECXII/cabecoVide.html">Cabeço de Vide</a>
					<a rel="nofollow" href="CastelosSECXII/casteloRodrigo.html">Castelo Rodrigo</a>
				</div>
			</div>
		`), nil
	}

	expectedCastles := []castle.Model{
		{Country: castle.Portugal, Name: "Abrantes", Link: "https://www.castelosdeportugal.pt/castelos/CastelosSECXII/abrantes.html"},
		{Country: castle.Portugal, Name: "Alpalhão", Link: "https://www.castelosdeportugal.pt/castelos/CastelosSECXII/alpalhao.html"},
		{Country: castle.Portugal, Name: "Atouguia da Baleia", Link: "https://www.castelosdeportugal.pt/castelos/Castelos(pre)SECXII/atouguiaBaleia.html"},
		{Country: castle.Portugal, Name: "Braga", Link: "https://www.castelosdeportugal.pt/castelos/CastelosSECXIII/braga.html"},
		{Country: castle.Portugal, Name: "Bragança", Link: "https://www.castelosdeportugal.pt/castelos/CastelosSECXIII/braganca.html"},
		{Country: castle.Portugal, Name: "Cabeço de Vide", Link: "https://www.castelosdeportugal.pt/castelos/CastelosSECXII/cabecoVide.html"},
		{Country: castle.Portugal, Name: "Castelo Rodrigo", Link: "https://www.castelosdeportugal.pt/CastelosSECXII/casteloRodrigo.html"},
	}

	portugueseCollector := NewCastelosDePortugalEnricher(httpclient.New(), fakeHTMLFetcher)

	castleChan, errChan := portugueseCollector.CollectCastlesToEnrich(context.Background())
	if len(errChan) > 0 {
		t.Error("expected to have no values on err chan")
	}
	var foundCastles []castle.Model
	for c := range castleChan {
		foundCastles = append(foundCastles, c)
	}

	if !slicesWithSameContent(foundCastles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
		diff := cmp.Diff(foundCastles, expectedCastles)
		t.Errorf("diff: %v", diff)
	}
}

func TestExtractPortugueseCastleInfo(t *testing.T) {
	htmlFetcher := func(ctx context.Context, link string, httpClient *http.Client) ([]byte, error) {
		return fileloader.LoadHTMLFile(castlePageHTMLPath)
	}

	expectedCastle := castle.Model{
		Name:              "porto",
		Country:           castle.Portugal,
		City:              "guimarães",
		State:             "guimarães",
		District:          "oliveira do castelo",
		FoundationPeriod:  "(ant. a 958)",
		Link:              "https://somelink.pt",
		PropertyCondition: castle.Intact,
	}

	castelosDePortugalEnricher := NewCastelosDePortugalEnricher(httpclient.New(), htmlFetcher)

	receivedCastle, err := castelosDePortugalEnricher.EnrichCastle(context.Background(), expectedCastle)
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
	if receivedCastle.FoundationPeriod != expectedCastle.FoundationPeriod {
		t.Errorf("expected FoundationPeriod to be [%s], got [%s]", expectedCastle.FoundationPeriod, receivedCastle.FoundationPeriod)
	}
	if receivedCastle.PropertyCondition != expectedCastle.PropertyCondition {
		t.Errorf("expected PropertyCondition to be [%s], got [%s]", expectedCastle.PropertyCondition, receivedCastle.PropertyCondition)
	}
}

func TestParseCondition(t *testing.T) {
	testCases := []struct {
		rawCondition      string
		expectedCondition castle.PropertyCondition
	}{
		{rawCondition: "Boa", expectedCondition: castle.Intact},
		{rawCondition: "", expectedCondition: castle.Unknown},
		{rawCondition: "()", expectedCondition: castle.Unknown},
		{rawCondition: "Submerso", expectedCondition: castle.Ruins},
		{rawCondition: "Mau", expectedCondition: castle.Ruins},
		{rawCondition: "Razoável", expectedCondition: castle.Damaged},
	}
	cp := &castelosDePortugalEnricher{}

	for _, tt := range testCases {
		cTT := tt
		t.Run(cTT.rawCondition, func(t *testing.T) {
			received := cp.parseCondition(cTT.rawCondition)

			if received != cTT.expectedCondition {
				t.Errorf("expected to have condition [%s], got [%s]", cTT.expectedCondition.String(), received.String())
			}
		})
	}
}
