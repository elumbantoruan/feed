package crawler

import (
	"fmt"
	"github/elumbantoruan/feed/pkg/feed"
	"time"
)

type AtomCrawler struct {
	URL string
}

func NewAtomCrawler(url string) *AtomCrawler {
	return &AtomCrawler{
		URL: url,
	}
}

func (ac *AtomCrawler) Download(url string) (*feed.Feed, error) {
	atom, err := download[feed.Atom](url)
	if err != nil {
		return nil, err
	}
	return ac.ConvertToFeed(atom)
}

func (ac *AtomCrawler) ConvertToFeed(a *feed.Atom) (*feed.Feed, error) {
	layout := "2006-01-02T15:04:05Z04:00"
	layout = time.RFC3339
	updated, err := time.Parse(layout, a.Updated)
	if err != nil {
		return nil, err
	}
	f := &feed.Feed{
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
		article := feed.Article{
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
