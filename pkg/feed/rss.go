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

type RssCrawler struct {
	URL string
}

func NewRssCrawler(url string) *RssCrawler {
	return &RssCrawler{
		URL: url,
	}
}

func (rc *RssCrawler) Download(url string) (*Feed, error) {
	rss, err := download[Rss](url)
	if err != nil {
		return nil, err
	}
	rss.URL = url
	return rc.ConvertToFeed(rss)
}

func (rc *RssCrawler) ConvertToFeed(r *Rss) (*Feed, error) {
	recentItemPubDate := r.Channel.LastBuildDate
	if recentItemPubDate == "" && len(r.Channel.Item) > 0 {
		recentItemPubDate = r.Channel.Item[0].PubDate
	}

	pubDate, err := ParseDateTime(recentItemPubDate)
	if err != nil {
		return nil, err
	}
	f := &Feed{
		Site:    r.Channel.Title,
		Link:    r.Channel.Link.Href,
		Updated: pubDate,
		RSS:     r.URL,
	}
	for _, item := range r.Channel.Item {
		published, err := ParseDateTime(item.PubDate)

		if err != nil {
			return nil, err
		}
		article := Article{
			ID:        item.Guid.Text,
			Published: published,
			Updated:   published,
			Title:     item.Title,
			Authors:   []string{item.Creator},
			Link:      item.Link,
			Content:   item.Description,
		}
		f.Articles = append(f.Articles, article)
	}

	return f, nil
}
