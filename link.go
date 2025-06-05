package grabber

import (
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Rel  string
	Ref  string
	Href string
}

// LinkTag constructs LinkTag.
func LinkTag(attrs []html.Attribute) *Link {
	link := new(Link)
	for _, attr := range attrs {
		switch attr.Key {
		case "rel":
			link.Rel = attr.Val
		case "ref":
			link.Ref = attr.Val
		case "href":
			link.Href = attr.Val
		}
	}
	return link
}

func (link *Link) Contribute(p *Page) {
	switch {
	case link.IsIcon():
		if len(p.Favicons) == 0 || p.Favicons[len(p.Favicons)-1] != link.Href {
			p.Favicons = append(p.Favicons, link.Href)
		}
	}
}

func (link *Link) IsIcon() bool {
	s := link.Rel + link.Ref
	return (strings.Contains(s, "shortcut") || strings.Contains(s, "icon")) && link.Href != ""
}
