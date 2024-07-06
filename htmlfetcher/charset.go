package htmlfetcher

import (
	"strings"

	"golang.org/x/net/html"
)

func getCharset(responseBody []byte) (string, error) {
	doc, err := html.Parse(strings.NewReader(string(responseBody)))
	if err != nil {
		return "", err
	}

	var charset string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, attr := range n.Attr {
				if attr.Key == "charset" {
					charset = attr.Val
					return
				} else if attr.Key == "http-equiv" && strings.ToLower(attr.Val) == "content-type" {
					for _, subAttr := range n.Attr {
						if subAttr.Key == "content" {
							content := subAttr.Val
							idx := strings.Index(content, "charset=")
							if idx != -1 {
								charset = strings.ToUpper(strings.TrimSpace(content[idx+len("charset="):]))
							}
							return
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	if charset == "" {
		charset = "UTF-8"
	}

	return strings.ToUpper(charset), nil
}
