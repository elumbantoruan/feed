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

func (m *MockGRPCClient) AddSite(ctx context.Context, site feed.Site[int64]) error {
	return nil
}

func (m *MockGRPCClient) GetSites(ctx context.Context) ([]feed.Site[int64], error) {
	m.GetSitesCount++
	sites := []feed.Site[int64]{createTestSite(1), createTestSite(2)}
	return sites, nil
}

func (m *MockGRPCClient) UpdateSite(ctx context.Context, site feed.Site[int64]) error {
	m.UpdateSiteFeedCount++
	return nil
}

func (m *MockGRPCClient) UpsertArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {
	m.AddArticleCount++
	return int64(m.AddArticleCount), nil
}

func (m *MockGRPCClient) GetArticles(ctx context.Context) ([]feed.ArticleSite[int64], error) {
	return nil, nil
}

func (m *MockGRPCClient) GetArticlesWithSite(ctx context.Context, siteID int64, limit int32) ([]feed.Article, error) {
	return nil, nil
}

func createTestSite(id int64) feed.Site[int64] {
	ts := time.Time{}
	return feed.Site[int64]{
		ID:      id,
		Site:    fmt.Sprintf("TestSite%d", id),
		Link:    fmt.Sprintf("http://testsite%d", id),
		RSS:     fmt.Sprintf("http://testsite%d", id),
		Type:    "test",
		Updated: &ts,
	}
}

func createTestFeedData(id int64) feed.FeedSite[int64] {
	ts := time.Time{}
	return feed.FeedSite[int64]{
		Site: feed.Site[int64]{
			ID:      id,
			Site:    fmt.Sprintf("TestSite%d", id),
			Link:    fmt.Sprintf("http://testsite%d", id),
			RSS:     fmt.Sprintf("http://testsite%d", id),
			Type:    "test",
			Updated: &ts,
		},
	}
}
