package collector

import (
	"reflect"
	"testing"

	"github.com/buarki/find-castles/castle"
)

const (
	htmlWithEnglandCastlesPath    = "../data/uk-england-list.html"
	jsonWithEnglandCastlesPath    = "../data/uk-england-list.json"
	jsonWithEnglandCastlePagePath = "../data/uk-england-castle.html"
)

func TestExtractTheListOfCastlesInEngland(t *testing.T) {
	htmlWithEnglandCastles, err := loadHTMLFile(htmlWithEnglandCastlesPath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	collectedCastles, err := extractTheListOfCastlesInEngland(htmlWithEnglandCastles)
	if err != nil {
		t.Errorf("expected to have err nil when collecting the list of castles in england, got %v", err)
	}

	expected, err := loadJSONToCompare(jsonWithEnglandCastlesPath)
	if err != nil {
		t.Errorf("expected to have err nil when loading castles of england to compare, got %v", err)
	}

	if !reflect.DeepEqual(collectedCastles, expected) {
		t.Errorf("parsed castles do not match expected castles")
	}
}

func TestExtractUkState(t *testing.T) {
	castleHTMLPage, err := loadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	expectedCaste := castle.Model{
		State: "Berkshire",
		City:  "",
	}
	receivedState, err := extractUkState(castleHTMLPage, expectedCaste)
	if err != nil {
		t.Errorf("expected err nil when getting data of UK castle, got %v", err)
	}
	if receivedState != expectedCaste.State {
		t.Errorf("expected state to be [%s], got [%s]", expectedCaste.State, receivedState)
	}
}

func TestExtractUkCity(t *testing.T) {
	castleHTMLPage, err := loadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	expectedCaste := castle.Model{
		State: "Berkshire",
		City:  "Windsor SL4 1NJ",
	}
	receivedCity, err := extractUkCity(castleHTMLPage, expectedCaste)
	if err != nil {
		t.Errorf("expected err nil when getting data of UK castle, got %v", err)
	}
	if receivedCity != expectedCaste.City {
		t.Errorf("expected city to be [%s], got [%s]", expectedCaste.City, receivedCity)
	}
}
