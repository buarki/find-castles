package enricher

import (
	"testing"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fileloader"
)

const (
	castlesHomePageHTMLPath = "../data/portugal.html"
	expectedCastlesAsJSON   = "../data/portugal.json"
	castlePageHTMLPath      = "../data/portugal-castle.html"
)

func TestExtractCastleInfo(t *testing.T) {
	expectedCastle := castle.Model{
		Name:             "Porto",
		Country:          "Portugal",
		City:             "Guimarães",
		State:            "Guimarães",
		District:         "Oliveira do Castelo",
		YearOfFoundation: "(ant. a 958)",
	}

	castlePage, err := fileloader.LoadHTMLFile(castlePageHTMLPath)
	if err != nil {
		t.Errorf("failed to load castle page, got %v", err)
	}

	receivedCastle, err := extractCastleInfo(expectedCastle, castlePage)
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
