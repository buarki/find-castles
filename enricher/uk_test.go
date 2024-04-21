package enricher

import (
	"testing"

	"github.com/buarki/find-castles/castle"
	"github.com/buarki/find-castles/fileloader"
)

const (
	htmlWithEnglandCastlesPath    = "../data/uk-england-list.html"
	jsonWithEnglandCastlesPath    = "../data/uk-england-list.json"
	jsonWithEnglandCastlePagePath = "../data/uk-england-castle.html"
)

func TestExtractUkState(t *testing.T) {
	castleHTMLPage, err := fileloader.LoadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	expectedCaste := castle.Model{
		State: "Berkshire, GreaterLondon",
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
	castleHTMLPage, err := fileloader.LoadHTMLFile(jsonWithEnglandCastlePagePath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	expectedCaste := castle.Model{
		State: "Berkshire, GreaterLondon",
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
