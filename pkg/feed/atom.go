package feed

import (
	"encoding/xml"
	"fmt"
	"time"
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

type AtomCrawler struct {
	URL string
}

func NewAtomCrawler(url string) *AtomCrawler {
	return &AtomCrawler{
		URL: url,
	}
}

func (ac *AtomCrawler) Download(url string) (*Feed, error) {
	atom, err := download[Atom](url)
	if err != nil {
		return nil, err
	}
	return ac.ConvertToFeed(atom)
}

func (ac *AtomCrawler) ConvertToFeed(a *Atom) (*Feed, error) {
	layout := "2006-01-02T15:04:05Z04:00"
	layout = time.RFC3339
	updated, err := time.Parse(layout, a.Updated)
	if err != nil {
		return nil, err
	}
	f := &Feed{
		Site:    a.Title,
		Updated: updated,
		Link:    a.Link.Href,
		Icon:    a.Icon,
		RSS:     a.ID,
	}

	for _, entry := range a.Entry {
		published, err := time.Parse(layout, entry.Published)
		if err != nil {
			fmt.Println(err)
		}
		updated, err := time.Parse(layout, entry.Updated)
		if err != nil {
			fmt.Println(err)
		}
		authors := func(vals []string) []string {
			var auths []string
			for _, val := range vals {
				auths = append(auths, val)
			}
			return auths
		}
		article := Article{
			ID:        entry.ID,
			Published: published,
			Updated:   updated,
			Title:     entry.Title,
			Authors:   authors(entry.Author.Name),
			Link:      entry.Link.Href,
			Content:   entry.Content.Text,
		}
		f.Articles = append(f.Articles, article)
	}
	return f, nil
}
