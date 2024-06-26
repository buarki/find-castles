package castle

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	castlePattern = regexp.MustCompile(".*castle.*")
)

func TestGenerateMatchingTags(t *testing.T) {
	c := Model{
		Name:     "braga",
		Country:  Portugal,
		State:    "braga",
		City:     "braga",
		District: "fafe",
	}

	itemsSet := map[string]struct{}{
		strings.ToLower(c.Name):             {},
		strings.ToLower(c.Country.String()): {},
		strings.ToLower(c.State):            {},
		strings.ToLower(c.City):             {},
		strings.ToLower(c.District):         {},
	}

	generatedMatchingTags := c.GetMatchingTags()

	if len(generatedMatchingTags) == 0 {
		t.Errorf("expected len > 0")
	}

	for _, mt := range generatedMatchingTags {
		if len(mt) == 0 {
			t.Errorf("expected to not have empty strings on matching tags")
		}
		if castlePattern.Match([]byte(mt)) {
			t.Error("word [castle] MUST NOT BE PRESENT")
		}
		if _, has := itemsSet[mt]; !has {
			t.Errorf("expected to find tag [%s]", mt)
		}
	}
}

func TestIsProbably(t *testing.T) {
	testCases := []struct {
		c1      Model
		c2      Model
		matches bool
	}{
		{
			c1: Model{
				Country: UK,
			},
			c2: Model{
				Country: Ireland,
			},
			matches: false,
		},
		{
			c1: Model{
				Country: UK,
			},
			c2: Model{
				Country: UK,
			},
			matches: false,
		},
		{
			c1: Model{
				Country: UK,
				Name:    "Windsor",
			},
			c2: Model{
				Country: UK,
				Name:    "Windsor",
			},
			matches: true,
		},
		{
			c1: Model{
				Country: UK,
				Name:    "Windsor castle",
			},
			c2: Model{
				Country: UK,
				Name:    "Windsor",
			},
			matches: true,
		},
		{
			c1: Model{
				Country: UK,
				Name:    "Kirby Castle",
			},
			c2: Model{
				Country: UK,
				Name:    "Kirby Muxloe Castle",
			},
			matches: true,
		},
		{
			c1: Model{
				Country: UK,
				Name:    "Kirby Muxloe Castle",
			},
			c2: Model{
				Country: UK,
				Name:    "Kirby Castle",
			},
			matches: true,
		},
		{
			c1: Model{
				Country: Ireland,
				Name:    "St Kirby Muxloe Castle",
			},
			c2: Model{
				Country: Ireland,
				Name:    "Kirby Castle",
			},
			matches: true,
		},
		{
			c1: Model{
				Country:  Portugal,
				Name:     "Castelo do Mau Vizinho(Évora)",
				State:    "Évora",
				City:     "Évora",
				District: "Igrejinha",
			},
			c2: Model{
				Country: Portugal,
				Name:    "Castelo de Évora",
				State:   "Évora",
				City:    "Arraiolos",
			},
			matches: false,
		},
	}

	for _, tt := range testCases {
		receivedResult := tt.c1.IsProbably(tt.c2)

		if receivedResult != tt.matches {
			t.Errorf("expected to match")
		}
	}
}

func TestReconcileWith(t *testing.T) {
	testCases := []struct {
		name         string
		c1           Model
		c2           Model
		resultCastle Model
		err          error
	}{
		{
			name: "when castles are not from same country",
			c1: Model{
				Country: Portugal,
			},
			c2: Model{
				Country: UK,
			},
			resultCastle: Model{},
			err:          ErrCastlesShouldProbablyBeTheSameToReconcile,
		},
		{
			name: "when one castle has more data than other",
			c1: Model{
				Country: Portugal,
				Name:    "guimaraes",
			},
			c2: Model{
				Country:          Portugal,
				Name:             "guimaraes",
				FoundationPeriod: "XX",
			},
			resultCastle: Model{
				Country:          Portugal,
				Name:             "guimaraes",
				FoundationPeriod: "XX",
			},
			err: nil,
		},
		{
			name: "when castles are equal",
			c1: Model{
				Country:          Portugal,
				Name:             "guimaraes",
				State:            "Braga",
				City:             "guimaraes",
				District:         "Castle",
				FoundationPeriod: "XX",
			},
			c2: Model{
				Country:          Portugal,
				Name:             "guimaraes",
				State:            "Braga",
				City:             "guimaraes",
				District:         "Castle",
				FoundationPeriod: "XX",
			},
			resultCastle: Model{
				Country:          Portugal,
				Name:             "guimaraes",
				State:            "Braga",
				City:             "guimaraes",
				District:         "Castle",
				FoundationPeriod: "XX",
			},
			err: nil,
		},
		{
			name: "when castles names are slightly different",
			c1: Model{
				Country: UK,
				Name:    "kirby muxloe castle",
			},
			c2: Model{
				Country: UK,
				Name:    "kirby castle",
			},
			resultCastle: Model{
				Country: UK,
				Name:    "kirby castle",
			},
			err: nil,
		},
		{
			name: "when castles names are slightly different in reverse order",
			c1: Model{
				Country: UK,
				Name:    "kirby castle",
			},
			c2: Model{
				Country: UK,
				Name:    "kirby muxloe castle",
			},
			resultCastle: Model{
				Country: UK,
				Name:    "kirby castle",
			},
			err: nil,
		},
		{
			name: "when castles states are slightly different",
			c1: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "castelo de guimaraes",
			},
			c2: Model{
				Country: Portugal,
				State:   "Braga",
				Name:    "guimaraes",
			},
			resultCastle: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "guimaraes",
			},
			err: nil,
		},
		{
			name: "when castles states are slightly different in reverse order",
			c1: Model{
				Country: Portugal,
				State:   "Braga",
				Name:    "guimaraes",
			},
			c2: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "castelo de guimaraes",
			},
			resultCastle: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "guimaraes",
			},
			err: nil,
		},
		{
			name: "when one castle has state while other have",
			c1: Model{
				Country: Portugal,
				Name:    "guimaraes",
			},
			c2: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "castelo de guimaraes",
			},
			resultCastle: Model{
				Country: Portugal,
				State:   "Distrito de Braga",
				Name:    "guimaraes",
			},
			err: nil,
		},
		{
			name: "when one castle has coordinates while other don't",
			c1: Model{
				Country: Portugal,
				Name:    "guimaraes",
			},
			c2: Model{
				Country:     Portugal,
				Name:        "castelo de guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			resultCastle: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			err: nil,
		},
		{
			name: "when one castle has coordinates while other don't (REVERSE)",
			c1: Model{
				Country:     Portugal,
				Name:        "castelo de guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			c2: Model{
				Country: Portugal,
				Name:    "guimaraes",
			},
			resultCastle: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			err: nil,
		},
		{
			name: "when BOTH castles have coordinates, biggest one stays",
			c1: Model{
				Country:     Portugal,
				Name:        "castelo de guimaraes",
				Coordinates: "47.921271, 18.642998",
			},
			c2: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			resultCastle: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271, 18.642998",
			},
			err: nil,
		},
		{
			name: "when BOTH castles have coordinates, biggest one stays (REVERSE)",
			c1: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271,18.642998",
			},
			c2: Model{
				Country:     Portugal,
				Name:        "castelo de guimaraes",
				Coordinates: "47.921271, 18.642998",
			},
			resultCastle: Model{
				Country:     Portugal,
				Name:        "guimaraes",
				Coordinates: "47.921271, 18.642998",
			},
			err: nil,
		},
		{
			name: "when one castle has contacts while other have",
			c1: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "8383883383",
					Email: "nuno@portugal.pt",
				},
			},
			c2: Model{
				Country: Portugal,
				Name:    "castelo de guimaraes",
			},
			resultCastle: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "8383883383",
					Email: "nuno@portugal.pt",
				},
			},
			err: nil,
		},
		{
			name: "when one castle has contacts while other have (REVERSE)",
			c1: Model{
				Country: Portugal,
				Name:    "castelo de guimaraes",
			},
			c2: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "8383883383",
					Email: "nuno@portugal.pt",
				},
			},
			resultCastle: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "8383883383",
					Email: "nuno@portugal.pt",
				},
			},
			err: nil,
		},
		{
			name: "when both castles have contact the bigger data stays",
			c1: Model{
				Country: Portugal,
				Name:    "castelo de guimaraes",
				Contact: &Contact{
					Phone: "1234567",
					Email: "nuno@portugal.pt",
				},
			},
			c2: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "12345",
					Email: "nuno-manel@portugal.pt",
				},
			},
			resultCastle: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Contact: &Contact{
					Phone: "1234567",
					Email: "nuno-manel@portugal.pt",
				},
			},
			err: nil,
		},
		{
			name: "when castles have different sources",
			c1: Model{
				Country: Portugal,
				Name:    "castelo de guimaraes",
				Sources: []string{"https://castle.pt"},
			},
			c2: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Sources: []string{"https://castelos.pt"},
			},
			resultCastle: Model{
				Country: Portugal,
				Name:    "guimaraes",
				Sources: []string{"https://castle.pt", "https://castelos.pt"},
			},
			err: nil,
		},
	}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.name, func(t *testing.T) {
			t.Helper()

			received, err := currentTT.c1.ReconcileWith(currentTT.c2)

			if currentTT.err != tt.err {
				t.Errorf("expected to have err [%v], got [%v]", currentTT.err, err)
			}

			if !reflect.DeepEqual(received, currentTT.resultCastle) {
				t.Errorf("expected to have received castle [%+v], got [%+v]", currentTT.resultCastle, received)
				diff := cmp.Diff(currentTT.resultCastle, received)
				t.Errorf("diff: %v", diff)
			}
		})
	}
}

func TestWebName(t *testing.T) {
	testCases := []struct {
		c       Model
		webName string
	}{
		{
			c:       Model{Country: Slovakia, Name: "zámok dechtice-hradišco"},
			webName: "zamok-dechtice-hradisco-sk",
		},
		{
			c:       Model{Country: Ireland, Name: "Roscrea, Gardens And Damer House/black Mills"},
			webName: "roscrea-gardens-and-damer-house-black-mills-ie",
		},
		{
			c:       Model{Country: Ireland, Name: "Parke’s"},
			webName: "parke-s-ie",
		},
		{
			c:       Model{Country: Slovakia, Name: "Bojnický Zámok/schloss Bojnice"},
			webName: "bojnicky-zamok-schloss-bojnice-sk",
		},
		{
			c:       Model{Country: Slovakia, Name: "Nitrianska Blatnica - jurko a púst"},
			webName: "nitrianska-blatnica-jurko-a-pust-sk",
		},
	}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.webName, func(t *testing.T) {
			received, err := currentTT.c.WebName()

			if err != nil {
				t.Errorf("expected err nil, got %v", err)
			}

			if received != currentTT.webName {
				t.Errorf("expected to have [%s], got [%s]", currentTT.webName, received)
			}
		})
	}
}
