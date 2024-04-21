package collector

import (
	"testing"

	"github.com/buarki/find-castles/fileloader"
)

const (
	htmlWithEnglandCastlesPath = "../data/uk-england-list.html"
	jsonWithEnglandCastlesPath = "../data/uk-england-list.json"
)

func TestExtractTheListOfCastlesInEngland(t *testing.T) {
	htmlWithEnglandCastles, err := fileloader.LoadHTMLFile(htmlWithEnglandCastlesPath)
	if err != nil {
		t.Errorf("expected to have err nil when loading HTML with the list of castles in england, got %v", err)
	}
	collectedCastles, err := extractTheListOfCastlesInEngland(htmlWithEnglandCastles)
	if err != nil {
		t.Errorf("expected to have err nil when collecting the list of castles in england, got %v", err)
	}
	expectedCastles, err := fileloader.LoadCastlesAsJSONList(jsonWithEnglandCastlesPath)
	if err != nil {
		t.Errorf("expected to have err nil when loading castles of england to compare, got %v", err)
	}
	if !slicesWithSameContent(expectedCastles, collectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}
