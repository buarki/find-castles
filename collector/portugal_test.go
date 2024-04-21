package collector

import (
	"testing"

	"github.com/buarki/find-castles/fileloader"
)

const (
	castlesHomePageHTMLPath = "../data/portugal.html"
	expectedCastlesAsJSON   = "../data/portugal.json"
)

func TestCollectCastleNameAndLinks(t *testing.T) {
	htmlToParse, err := fileloader.LoadHTMLFile(castlesHomePageHTMLPath)
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}
	castles, err := collectCastleNameAndLinks(htmlToParse)
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}
	expectedCastles, err := fileloader.LoadCastlesAsJSONList(expectedCastlesAsJSON)
	if err != nil {
		t.Errorf("expected to have err nil, got %v", err)
	}
	if !slicesWithSameContent(castles, expectedCastles) {
		t.Errorf("parsed castles do not match expected castles")
	}
}
