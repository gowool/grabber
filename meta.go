package grabber

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type Meta struct {
	Name     string
	Property string
	Content  string
}

// MetaTag constructs MetaTag.
func MetaTag(attrs []html.Attribute) *Meta {
	meta := new(Meta)
	for _, attr := range attrs {
		switch attr.Key {
		case "property":
			meta.Property = attr.Val
		case "content":
			meta.Content = attr.Val
		case "name":
			meta.Name = attr.Val
		}
	}
	return meta
}

func (meta *Meta) Contribute(p *Page) (err error) {
	switch {
	case meta.IsDescription():
		p.Description = meta.Content
	case meta.IsKeywords():
		p.Keywords = meta.Content
	case meta.IsAuthor():
		p.Authors = append(p.Authors, meta.Content)
	case meta.IsOGTitle():
		p.OpenGraph.Title = meta.Content
	case meta.IsOGDescription():
		p.OpenGraph.Description = meta.Content
	case meta.IsOGSiteName():
		p.OpenGraph.SiteName = meta.Content
	case meta.IsOGImage():
		if len(p.OpenGraph.Images) == 0 || p.OpenGraph.Images[len(p.OpenGraph.Images)-1].URL != meta.Content {
			p.OpenGraph.Images = append(p.OpenGraph.Images, Image{URL: meta.Content})
		}
	case meta.IsPropertyOf("og:image"):
		if len(p.OpenGraph.Images) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:image:type":
			p.OpenGraph.Images[len(p.OpenGraph.Images)-1].Type = meta.Content
		case "og:image:secure_url":
			p.OpenGraph.Images[len(p.OpenGraph.Images)-1].SecureURL = meta.Content
		case "og:image:alt":
			p.OpenGraph.Images[len(p.OpenGraph.Images)-1].Alt = meta.Content
		case "og:image:width":
			p.OpenGraph.Images[len(p.OpenGraph.Images)-1].Width, err = strconv.Atoi(meta.Content)
		case "og:image:height":
			p.OpenGraph.Images[len(p.OpenGraph.Images)-1].Height, err = strconv.Atoi(meta.Content)
		}
	case meta.IsOGAudio():
		if len(p.OpenGraph.Audio) == 0 || p.OpenGraph.Audio[len(p.OpenGraph.Audio)-1].URL != meta.Content {
			p.OpenGraph.Audio = append(p.OpenGraph.Audio, Audio{URL: meta.Content})
		}
	case meta.IsPropertyOf("og:audio"):
		if len(p.OpenGraph.Audio) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:audio:type":
			p.OpenGraph.Audio[len(p.OpenGraph.Audio)-1].Type = meta.Content
		case "og:audio:secure_url":
			p.OpenGraph.Audio[len(p.OpenGraph.Audio)-1].SecureURL = meta.Content
		}
	case meta.IsOGVideo():
		if len(p.OpenGraph.Videos) == 0 || p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].URL != meta.Content {
			p.OpenGraph.Videos = append(p.OpenGraph.Videos, Video{URL: meta.Content})
		}
	case meta.IsPropertyOf("og:video"):
		if len(p.OpenGraph.Videos) == 0 {
			return nil
		}
		switch meta.Property {
		case "og:video:type":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Type = meta.Content
		case "og:video:secure_url":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].SecureURL = meta.Content
		case "og:video:width":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Width, err = strconv.Atoi(meta.Content)
		case "og:video:height":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Height, err = strconv.Atoi(meta.Content)
		case "og:video:duration":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Duration, err = strconv.Atoi(meta.Content)
		case "og:video:tag":
			p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Tag = append(p.OpenGraph.Videos[len(p.OpenGraph.Videos)-1].Tag, meta.Content)
		}
	case meta.IsOGType():
		p.OpenGraph.Type = meta.Content
	case meta.IsOGURL():
		p.OpenGraph.URL = meta.Content
	case meta.IsOGLocale():
		p.OpenGraph.Locale = meta.Content
	case meta.IsOGUpdatedTime():
		p.OpenGraph.UpdatedTime = ParseDate(meta.Content)
	case meta.IsArticlePublishedTime():
		p.OpenGraph.Article.PublishedTime = ParseDate(meta.Content)
	case meta.IsArticleModifiedTime():
		p.OpenGraph.Article.ModifiedTime = ParseDate(meta.Content)
	case meta.IsArticleExpirationTime():
		p.OpenGraph.Article.ExpirationTime = ParseDate(meta.Content)
	case meta.IsArticleAuthor():
		p.OpenGraph.Article.Authors = append(p.OpenGraph.Article.Authors, meta.Content)
	case meta.IsArticleSection():
		p.OpenGraph.Article.Sections = append(p.OpenGraph.Article.Sections, meta.Content)
	case meta.IsArticleTag():
		p.OpenGraph.Article.Tags = append(p.OpenGraph.Article.Tags, meta.Content)
	}
	return err
}

// IsDescription returns if it can be "description".
func (meta *Meta) IsDescription() bool {
	return meta.Name == "description" && meta.Content != ""
}

// IsKeywords returns if it can be "keywords".
func (meta *Meta) IsKeywords() bool {
	return meta.Name == "keywords" && meta.Content != ""
}

// IsAuthor returns if it can be "author".
func (meta *Meta) IsAuthor() bool {
	return meta.Name == "author" && meta.Content != ""
}

// IsOGTitle returns if it can be "title" of OGP
func (meta *Meta) IsOGTitle() bool {
	return meta.Property == "og:title" && meta.Content != ""
}

// IsOGDescription returns if it can be "description" of OGP
func (meta *Meta) IsOGDescription() bool {
	return meta.Property == "og:description" && meta.Content != ""
}

// IsOGImage returns if it can be a root of "og:image"
func (meta *Meta) IsOGImage() bool {
	return meta.Property == "og:image" || meta.Property == "og:image:url"
}

// IsPropertyOf returns if it can be a property of specified struct
func (meta *Meta) IsPropertyOf(name string) bool {
	return strings.HasPrefix(meta.Property, name+":")
}

// IsOGAudio reeturns if it can be a root of "og:audio"
func (meta *Meta) IsOGAudio() bool {
	return meta.Property == "og:audio" || meta.Property == "og:audio:url"
}

// IsOGVideo returns if it can be a root of "og:video"
func (meta *Meta) IsOGVideo() bool {
	return meta.Property == "og:video" || meta.Property == "og:video:url"
}

// IsOGType returns if it can be "og:type"
func (meta *Meta) IsOGType() bool {
	return meta.Property == "og:type"
}

// IsOGSiteName returns if it can be "og:site_name"
func (meta *Meta) IsOGSiteName() bool {
	return meta.Property == "og:site_name"
}

// IsOGURL returns if it can be "og:url"
func (meta *Meta) IsOGURL() bool {
	return meta.Property == "og:url"
}

// IsOGLocale returns if it can be "og:locale"
func (meta *Meta) IsOGLocale() bool {
	return meta.Property == "og:locale"
}

// IsOGUpdatedTime returns if it can be "og:updated_time"
func (meta *Meta) IsOGUpdatedTime() bool {
	return meta.Property == "og:updated_time"
}

// IsArticlePublishedTime returns if it can be "article:published_time"
func (meta *Meta) IsArticlePublishedTime() bool {
	return meta.Property == "article:published_time"
}

// IsArticleModifiedTime returns if it can be "article:modified_time"
func (meta *Meta) IsArticleModifiedTime() bool {
	return meta.Property == "article:modified_time"
}

// IsArticleExpirationTime returns if it can be "article:expiration_time"
func (meta *Meta) IsArticleExpirationTime() bool {
	return meta.Property == "article:expiration_time"
}

// IsArticleAuthor returns if it can be "article:author"
func (meta *Meta) IsArticleAuthor() bool {
	return meta.Property == "article:author"
}

// IsArticleSection returns if it can be "article:section"
func (meta *Meta) IsArticleSection() bool {
	return meta.Property == "article:section"
}

// IsArticleTag returns if it can be "article:tag"
func (meta *Meta) IsArticleTag() bool {
	return meta.Property == "article:tag"
}
