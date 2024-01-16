package mock

import (
	"fmt"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
)

type MockCrawler[T any] struct {
	Counter int
}

func (m *MockCrawler[T]) Download(url string) (*feed.FeedSite[int64], error) {
	m.Counter++
	feed := createCrawlerData(url, m.Counter)
	return &feed, nil
}

func createCrawlerData(url string, counter int) feed.FeedSite[int64] {
	feedData := createTestFeedData(int64(counter))
	ts := time.Now()
	feedData.Site.Updated = &ts
	feedData.Articles = []feed.Article{
		{
			Content: fmt.Sprintf("Content test%d", counter),
		},
	}
	return feedData

}
