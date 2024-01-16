package crawler

import (
	"fmt"

	"github.com/elumbantoruan/feed/pkg/feed"
)

type AtomCrawler struct {
	URL string
}

func NewAtomCrawler(url string) *AtomCrawler {
	return &AtomCrawler{
		URL: url,
	}
}

func (ac *AtomCrawler) Download(url string) (*feed.FeedSite[int64], error) {
	atom, err := download[feed.Atom](url)
	if err != nil {
		return nil, err
	}
	return ac.ConvertToFeed(atom)
}

func (ac *AtomCrawler) ConvertToFeed(a *feed.Atom) (*feed.FeedSite[int64], error) {
	updated, err := parseDateTime(a.Updated)
	if err != nil {
		return nil, err
	}
	f := &feed.FeedSite[int64]{
		Site: feed.Site[int64]{
			Site:    a.Title,
			Updated: &updated,
			Link:    a.Link.Href,
			Icon:    a.Icon,
			RSS:     a.ID,
		},
	}

	for _, entry := range a.Entry {
		published, err := parseDateTime(entry.Published)
		if err != nil {
			fmt.Println(err)
		}
		updated, err := parseDateTime(entry.Updated)
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
			ID:          entry.ID,
			Published:   published,
			Updated:     updated,
			Title:       entry.Title,
			Authors:     authors(entry.Author.Name),
			Link:        entry.Link.Href,
			Description: entry.Title,
			Content:     entry.Content.Text,
		}
		f.Articles = append(f.Articles, article)
	}
	return f, nil
}
