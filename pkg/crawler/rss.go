package crawler

import (
	"context"
	"strings"

	"github.com/elumbantoruan/feed/pkg/feed"
	"go.opentelemetry.io/otel/trace"
)

type RssCrawler struct {
}

func NewRssCrawler() *RssCrawler {
	return &RssCrawler{}
}

func (rc *RssCrawler) Download(ctx context.Context, url string) (*feed.FeedSite[int64], error) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Download RSS feed " + url)

	rss, err := download[feed.Rss](url)
	if err != nil {
		return nil, err
	}
	rss.URL = url
	return rc.ConvertToFeed(rss)
}

func (rc *RssCrawler) ConvertToFeed(r *feed.Rss) (*feed.FeedSite[int64], error) {
	recentItemPubDate := r.Channel.LastBuildDate
	if recentItemPubDate == "" && len(r.Channel.Item) > 0 {
		recentItemPubDate = r.Channel.Item[0].PubDate
	}

	pubDate, err := ParseDateTime(recentItemPubDate)
	if err != nil {
		return nil, err
	}
	f := &feed.FeedSite[int64]{
		Site: feed.Site[int64]{
			Site:    r.Channel.Title,
			Link:    r.Channel.Link.Href,
			Updated: &pubDate,
			RSS:     r.URL,
		},
	}
	for _, item := range r.Channel.Item {
		published, err := ParseDateTime(item.PubDate)

		if err != nil {
			return nil, err
		}
		article := feed.Article{
			ID:          item.Guid.Text,
			Published:   published,
			Updated:     published,
			Title:       item.Title,
			Authors:     []string{item.Creator},
			Link:        item.Link,
			Description: item.Description,
			Content:     strings.TrimSpace(item.Content),
		}
		if article.Content == "" {
			article.Content = item.Description
		}
		f.Articles = append(f.Articles, article)
	}

	return f, nil
}
