package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
)

type MockGRPCClient struct {
	AddArticleCount     int
	GetSitesCount       int
	UpdateSiteFeedCount int
}

func (m *MockGRPCClient) AddSiteFeed(ctx context.Context, site feed.Feed) error {
	return nil
}

func (m *MockGRPCClient) GetSitesFeed(ctx context.Context) ([]feed.Feed, error) {
	m.GetSitesCount++
	feeds := []feed.Feed{createTestFeedData(1), createTestFeedData(2)}
	return feeds, nil
}

func (m *MockGRPCClient) UpdateSiteFeed(ctx context.Context, feed feed.Feed) error {
	m.UpdateSiteFeedCount++
	return nil
}

func (m *MockGRPCClient) AddArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {
	m.AddArticleCount++
	return int64(m.AddArticleCount), nil
}

func (m *MockGRPCClient) GetArticles(ctx context.Context) ([]feed.ArticleSite[int64], error) {
	return nil, nil
}

func createTestFeedData(id int64) feed.Feed {
	ts := time.Time{}
	return feed.Feed{
		ID:      id,
		Site:    fmt.Sprintf("TestSite%d", id),
		Link:    fmt.Sprintf("http://testsite%d", id),
		RSS:     fmt.Sprintf("http://testsite%d", id),
		Type:    "test",
		Updated: &ts,
	}
}
