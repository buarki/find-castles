package toascii

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	replacer = strings.NewReplacer(
		`∅`, "o",
		`ø`, "o",
	)
)

func From(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFKC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return "", fmt.Errorf("failed to parse [%s] to ascii, got [%v]", s, err)
	}
	return replacer.Replace(result), nil
}
