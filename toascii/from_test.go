package toascii_test

import (
	"testing"

	"github.com/buarki/find-castles/toascii"
)

func TestFrom(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "Guimarães",
			output: "Guimaraes",
		},
		{
			input:  "hradište",
			output: "hradiste",
		},
		{
			input:  "halicský",
			output: "halicsky",
		},
		{
			input:  "bojnický",
			output: "bojnicky",
		},
		{
			input:  "zámok dechtice-hradišco",
			output: "zamok dechtice-hradisco",
		},
		{
			input:  "Loulé",
			output: "Loule",
		},
	}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.input, func(t *testing.T) {
			t.Helper()

			received, err := toascii.From(currentTT.input)

			if err != nil {
				t.Errorf("expected to have err nil, got %v", err)
			}

			if received != currentTT.output {
				t.Errorf("expected to get [%s], got [%s]", currentTT.output, received)
			}
		})
	}
}
