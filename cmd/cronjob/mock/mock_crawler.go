package mock

import (
	"fmt"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
)

type MockCrawler struct {
	Counter int
}

func (m *MockCrawler) Download(url string) (*feed.Feed, error) {
	m.Counter++
	feed := createCrawlerData(url, m.Counter)
	return &feed, nil
}

func createCrawlerData(url string, counter int) feed.Feed {
	feedData := createTestFeedData(int64(counter))
	ts := time.Now()
	feedData.Updated = &ts
	feedData.Articles = []feed.Article{
		{
			Content: fmt.Sprintf("Content test%d", counter),
		},
	}
	return feedData

}
