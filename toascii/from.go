package toascii

import (
	"fmt"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func From(s string) (string, error) {
	result, _, err := transform.String(transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn))), s)
	if err != nil {
		return "", fmt.Errorf("failed to parse [%s] to ascii, got [%v]", s, err)
	}
	return result, nil
}
