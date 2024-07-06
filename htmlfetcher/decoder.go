package htmlfetcher

import (
	"golang.org/x/text/encoding/charmap"
)

func decode(body []byte) ([]byte, error) {
	charset, err := getCharset(body)
	if err != nil {
		return nil, err
	}
	switch charset {
	case "ISO-8859-1":
		decoder := charmap.ISO8859_1.NewDecoder()
		return decoder.Bytes(body)
	default:
		return body, nil
	}
}
