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

func (w *WebStorage) GetArticles(ctx context.Context) (feed.Feeds, error) {
	var feeds feed.Feeds
	sites, err := w.GRCPClient.GetSitesFeed(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetArticles.GetSitesFeed error: %w", err)
	}
	for _, site := range sites {
		articles, err := w.GRCPClient.GetArticlesWithSite(ctx, site.ID, 10)
		if err != nil {
			return nil, fmt.Errorf("GetArticles.GetArticlesWithSite error: %w", err)
		}
		site.Articles = articles
		feeds = append(feeds, site)
	}
	return feeds, nil
}
