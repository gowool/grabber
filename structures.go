package grabber

import (
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Page struct {
	URL         string    `json:"url,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Keywords    string    `json:"keywords,omitempty"`
	Authors     []string  `json:"authors,omitempty"`
	Favicons    []string  `json:"favicons,omitempty"`
	OpenGraph   OpenGraph `json:"open_graph,omitempty"`
}

type OpenGraph struct {
	Title       string    `json:"title,omitempty"`
	Type        string    `json:"type,omitempty"`
	URL         string    `json:"url,omitempty"`
	Description string    `json:"description,omitempty"`
	Locale      string    `json:"locale,omitempty"`
	SiteName    string    `json:"site_name,omitempty"`
	UpdatedTime time.Time `json:"updated_time,omitzero"`
	Article     Article   `json:"article,omitempty,omitzero"`
	Images      []Image   `json:"images,omitempty,omitzero"`
	Videos      []Video   `json:"videos,omitempty,omitzero"`
	Audio       []Audio   `json:"audio,omitempty,omitzero"`
}

// Article represents a structure of "article".
type Article struct {
	PublishedTime  time.Time `json:"published_time,omitzero"`
	ModifiedTime   time.Time `json:"modified_time,omitzero"`
	ExpirationTime time.Time `json:"expiration_time,omitzero"`
	Authors        []string  `json:"authors,omitempty"`
	Sections       []string  `json:"sections,omitempty"`
	Tags           []string  `json:"tags,omitempty"`
}

// Image represents a structure of "og:image".
// "og:image" might have following properties:
//   - og:image:url
//   - og:image:secure_url
//   - og:image:type
//   - og:image:width
//   - og:image:height
//   - og:image:alt
type Image struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
	Type      string `json:"type,omitempty"` // Content-Type
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

// Video represents a structure of "og:video".
// "og:video" might have following properties:
//   - og:video:url
//   - og:video:secure_url
//   - og:video:type
//   - og:video:width
//   - og:video:height
//   - og:video:tag
type Video struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
	Type      string `json:"type,omitempty"` // Content-Type
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	// Duration in seconds
	Duration int      `json:"duration,omitempty"`
	Tag      []string `json:"tag,omitempty"`
}

// Audio represents a structure of "og:audio".
// "og:audio" might have following properties:
//   - og:audio:url
//   - og:audio:secure_url
//   - og:audio:type
type Audio struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secure_url,omitempty"`
	Type      string `json:"type,omitempty"` // Content-Type
}

func NewPage(url string) *Page {
	return &Page{URL: url}
}

// Parse parses http.Response.Body and construct Page informations.
// Caller should close body after it gets parsed.
func (p *Page) Parse(body io.Reader) error {
	tokens := html.NewTokenizer(body)
	isTitle := false

	for {
		switch tokens.Next() {
		case html.ErrorToken:
			if err := tokens.Err(); err != io.EOF {
				return err
			}
			return nil
		case html.EndTagToken:
			if tokens.Token().Data == "head" {
				return nil
			}
		case html.StartTagToken:
			t := tokens.Token()
			isTitle = false

			switch t.Data {
			case "title":
				isTitle = true
			case "meta":
				if err := MetaTag(t.Attr).Contribute(p); err != nil {
					return err
				}
			case "link":
				LinkTag(t.Attr).Contribute(p)
			}
		case html.TextToken:
			if isTitle && p.Title == "" {
				isTitle = false
				p.Title = tokens.Token().Data
			}
		case html.SelfClosingTagToken:
			t := tokens.Token()
			switch t.Data {
			case "meta":
				if err := MetaTag(t.Attr).Contribute(p); err != nil {
					return err
				}
			case "link":
				LinkTag(t.Attr).Contribute(p)
			}
		default:
		}
	}
}

// ToAbs makes all relative URLs to absolute URLs
func (p *Page) ToAbs() error {
	raw := p.URL
	base, err := url.Parse(raw)
	if err != nil {
		return err
	}
	// For og:image.
	for i, img := range p.OpenGraph.Images {
		p.OpenGraph.Images[i].URL = joinToAbsolute(base, img.URL)
	}
	// For og:audio
	for i, audio := range p.OpenGraph.Audio {
		p.OpenGraph.Audio[i].URL = joinToAbsolute(base, audio.URL)
	}
	// For og:video
	for i, video := range p.OpenGraph.Videos {
		p.OpenGraph.Videos[i].URL = joinToAbsolute(base, video.URL)
	}

	if len(p.Favicons) == 0 {
		p.Favicons = []string{"/favicon.ico"}
	}

	for i, favicon := range p.Favicons {
		p.Favicons[i] = joinToAbsolute(base, favicon)
	}

	return nil
}

func joinToAbsolute(base *url.URL, relpath string) string {
	src, err := url.Parse(relpath)
	if err == nil && src.IsAbs() {
		return src.String()
	}
	if strings.HasPrefix(relpath, "//") {
		return fmt.Sprintf("%s:%s", base.Scheme, relpath)
	}
	if strings.HasPrefix(relpath, "/") {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, relpath)
	}
	return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, path.Join(base.Path, relpath))
}
