package crawler

import (
	"context"
	"fmt"

	"github.com/elumbantoruan/feed/pkg/feed"
	"go.opentelemetry.io/otel/trace"
)

type AtomCrawler struct {
}

func NewAtomCrawler() *AtomCrawler {
	return &AtomCrawler{}
}

func (ac *AtomCrawler) Download(ctx context.Context, url string) (*feed.FeedSite[int64], error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Download Atom feed " + url)

	atom, err := download[feed.Atom](url)
	if err != nil {
		return nil, err
	}
	return ac.ConvertToFeed(atom)
}

func (ac *AtomCrawler) ConvertToFeed(a *feed.Atom) (*feed.FeedSite[int64], error) {
	updated, err := ParseDateTime(a.Updated)
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
		published, err := ParseDateTime(entry.Published)
		if err != nil {
			fmt.Println(err)
		}
		updated, err := ParseDateTime(entry.Updated)
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
