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
	Storage             []feed.Article
}

func (m *MockGRPCClient) AddSite(ctx context.Context, site feed.Site[int64]) error {
	return nil
}

func (m *MockGRPCClient) GetSites(ctx context.Context) ([]feed.Site[int64], error) {
	m.GetSitesCount++
	url := fmt.Sprintf("%v", ctx.Value("url"))
	sites := []feed.Site[int64]{createTestSite(url)}
	return sites, nil
}

func (m *MockGRPCClient) UpdateSite(ctx context.Context, site feed.Site[int64]) error {
	m.UpdateSiteFeedCount++
	return nil
}

func (m *MockGRPCClient) UpsertArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {
	m.AddArticleCount++
	m.Storage = append(m.Storage, article)
	return int64(m.AddArticleCount), nil
}

func (m *MockGRPCClient) GetArticles(ctx context.Context) ([]feed.ArticleSite[int64], error) {
	return nil, nil
}

func (m *MockGRPCClient) GetArticlesWithSite(ctx context.Context, siteID int64, limit int32) ([]feed.Article, error) {
	return nil, nil
}

func createTestSite(url string) feed.Site[int64] {
	ts := time.Time{}
	id := int64(1)
	return feed.Site[int64]{
		ID:      id,
		Site:    fmt.Sprintf("TestSite%d", id),
		Link:    url,
		RSS:     url,
		Type:    "rss",
		Updated: &ts,
	}
}
