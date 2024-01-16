package storage

import (
	"context"
	"fmt"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
)

type WebStorage struct {
	GRCPClient client.GRPCFeedClient
}

func NewWebStorage(grpcClient client.GRPCFeedClient) *WebStorage {
	return &WebStorage{
		GRCPClient: grpcClient,
	}
}

func (w *WebStorage) GetArticles(ctx context.Context) (feed.FeedSites[int64], error) {
	var feeds feed.FeedSites[int64]
	sites, err := w.GRCPClient.GetSites(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetArticles.GetSitesFeed error: %w", err)
	}
	for _, site := range sites {
		articles, err := w.GRCPClient.GetArticlesWithSite(ctx, site.ID, 10)
		if err != nil {
			return nil, fmt.Errorf("GetArticles.GetArticlesWithSite error: %w", err)
		}
		feedSite := feed.FeedSite[int64]{
			Site:     site,
			Articles: articles,
		}
		feeds = append(feeds, feedSite)
	}
	return feeds, nil
}
