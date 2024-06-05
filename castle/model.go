package castle

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

var (
	replacements = map[string]string{
		"(": "",
		")": "",
		"ç": "c",
		"á": "a",
		"é": "e",
		"í": "i",
		"ó": "o",
		"ú": "u",
		"â": "a",
		"ê": "e",
		"î": "i",
		"ô": "o",
		"û": "u",
		"ã": "a",
		"õ": "o",
		"à": "a",
		"è": "e",
		"ì": "i",
		"ò": "o",
		"ù": "u",
		"ä": "a",
		"ë": "e",
		"ï": "i",
		"ö": "o",
		"ü": "u",
		"ñ": "n",
		"/": "_",
	}
)

type Country string

func (c Country) String() string {
	return string(c)
}

const (
	Portugal Country = "Portugal"
	UK       Country = "UK"
	Ireland  Country = "Ireland"
)

type Model struct {
	Name             string  `json:"name"`
	Link             string  `json:"link"`
	Country          Country `json:"country"`
	State            string  `json:"state"`
	City             string  `json:"city"`
	District         string  `json:"district"`
	YearOfFoundation string  `json:"yearOfFoundation"`
	FlagLink         string  `json:"flagLink"`
}

func (m Model) NormalizeName() string {
	normalized := strings.ToLower(strings.TrimSpace(m.Name))
	normalized = strings.ReplaceAll(normalized, " ", "_")
	for old, new := range replacements {
		normalized = strings.ReplaceAll(normalized, old, new)
	}
	return normalized
}

func (m Model) StringID() string {
	normalizedName := m.NormalizeName()
	normalizedCountry := strings.ToLower(strings.TrimSpace(m.Country.String()))
	return fmt.Sprintf("%s-%s", normalizedCountry, normalizedName)
}

func (m Model) BinaryID() []byte {
	id := m.StringID()
	hash := sha256.Sum256([]byte(id))
	return hash[:]
}
