package crawler

import (
	"strings"

	"github.com/elumbantoruan/feed/pkg/feed"
)

type RssCrawler struct {
	URL string
}

func NewRssCrawler(url string) *RssCrawler {
	return &RssCrawler{
		URL: url,
	}
}

func (rc *RssCrawler) Download(url string) (*feed.FeedSite[int64], error) {
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

	pubDate, err := parseDateTime(recentItemPubDate)
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
		published, err := parseDateTime(item.PubDate)

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
