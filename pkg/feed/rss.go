package feed

import (
	"encoding/xml"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	URL     string
	Text    string   `xml:",chardata"`
	Atom    string   `xml:"atom,attr"`
	Dc      string   `xml:"dc,attr"`
	Media   string   `xml:"media,attr"`
	Version string   `xml:"version,attr"`
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
		Rel  string `xml:"rel,attr"`
		Type string `xml:"type,attr"`
	} `xml:"link"`
	Copyright     string `xml:"copyright"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Item          []Item `xml:"item"`
}

type Item struct {
	Text  string `xml:",chardata"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Guid  struct {
		Text        string `xml:",chardata"`
		IsPermaLink string `xml:"isPermaLink,attr"`
	} `xml:"guid"`
	PubDate     string   `xml:"pubDate"`
	Content     string   `xml:"content"`
	Description string   `xml:"description"`
	Category    []string `xml:"category"`
	Keywords    string   `xml:"keywords"`
	Creator     string   `xml:"creator"`
	Publisher   string   `xml:"publisher"`
	Subject     string   `xml:"subject"`
	Thumbnail   struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Width  string `xml:"width,attr"`
		Height string `xml:"height,attr"`
	} `xml:"thumbnail"`
}
