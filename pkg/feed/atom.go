package feed

import (
	"encoding/xml"
)

type Atom struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	Lang    string   `xml:"lang,attr"`
	Title   string   `xml:"title"`
	Icon    string   `xml:"icon"`
	Updated string   `xml:"updated"`
	ID      string   `xml:"id"`
	Link    Link     `xml:"link"`
	Entry   []Entry  `xml:"entry"`
}

type Link struct {
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type Entry struct {
	ID        string  `xml:"id"`
	Published string  `xml:"published"`
	Updated   string  `xml:"updated"`
	Title     string  `xml:"title"`
	Content   Content `xml:"content"`
	Link      Link    `xml:"link"`
	Author    Author  `xml:"author"`
}

type Author struct {
	Text string   `xml:",chardata"`
	Name []string `xml:"name"`
}

type Content struct {
	Text string `xml:",chardata"`
	Type string `xml:"type,attr"`
}
