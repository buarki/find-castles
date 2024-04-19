package collector

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/buarki/find-castles/castle"
)

const (
	castlesHomePageHTMLPath = "../data/portugal.html"
	expectedCastlesAsJSON   = "../data/portugal.json"
	castlePageHTMLPath      = "../data/portugal-castle.html"
)

func loadHTMLHomePage() ([]byte, error) {
	b, err := os.ReadFile(castlesHomePageHTMLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load file with HTML page, got %v", err)
	}
	return b, nil
}

func loadJSONToCompare() ([]castle.Model, error) {
	var castles []castle.Model
	b, err := os.ReadFile(expectedCastlesAsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to open portuguese json, got %v", err)
	}
	err = json.Unmarshal(b, &castles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON, got %v", err)
	}
	return castles, nil
}

func loadCastlePage() ([]byte, error) {
	b, err := os.ReadFile(castlePageHTMLPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load file with HTML page for castle, got %v", err)
	}
	return b, nil
}

func TestCollectCastleNameAndLinks(t *testing.T) {
	htmlToParse, err := loadHTMLHomePage()
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}
	castles, err := collectCastleNameAndLinks(htmlToParse)
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}

	expectedCastles, err := loadJSONToCompare()
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}

	if !reflect.DeepEqual(castles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}

func TestExtractCastleInfo(t *testing.T) {
	expectedCastle := castle.Model{
		Name:             "Porto",
		Country:          "Portugal",
		City:             "Guimarães",
		State:            "Guimarães",
		District:         "Oliveira do Castelo",
		YearOfFoundation: "(ant. a 958)",
	}

	castlePage, err := loadCastlePage()
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
