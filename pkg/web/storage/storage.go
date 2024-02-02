package storage

import (
	"context"
	"fmt"
	"runtime"
	"sort"

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

// GetArticles returns articles for all feedsites
// It starts with getting the list of the sites
// Then concurrently get the articles for each site
func (w *WebStorage) GetArticles(ctx context.Context) (feed.FeedSites[int64], error) {
	var feeds feed.FeedSites[int64]
	sites, err := w.GRCPClient.GetSites(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetArticles.GetSitesFeed error: %w", err)
	}

	var (
		sitesStream     = make(chan feed.Site[int64], len(sites))
		workers         = runtime.NumCPU()
		feedSitesStream = make(chan FeedSitesResult[int64], len(sites))
	)

	for i := 1; i <= workers; i++ {
		go func(wid int) {
			w.workerGetArticles(ctx, wid, sitesStream, feedSitesStream)
		}(i)
	}
	for _, site := range sites {
		sitesStream <- site
	}
	for i := 0; i < len(sites); i++ {
		result := <-feedSitesStream
		if result.Error != nil {
			return nil, fmt.Errorf("GetArticles.GetArticlesWithSite error: %w", err)
		}
		feedSite := result.FeedSite
		feeds = append(feeds, feedSite)
	}

	sort.SliceStable(feeds, func(i, j int) bool {
		return feeds[i].Site.ID < feeds[j].Site.ID
	})

	return feeds, nil
}

func (w *WebStorage) workerGetArticles(ctx context.Context, wID int, siteStream <-chan feed.Site[int64], feedSiteResult chan<- FeedSitesResult[int64]) {
	site := <-siteStream
	articles, err := w.GRCPClient.GetArticlesWithSite(ctx, site.ID, 10)
	feedSiteResult <- FeedSitesResult[int64]{
		FeedSite: feed.FeedSite[int64]{
			Site:     site,
			Articles: articles,
		},
		Error: err,
	}
}

type FeedSitesResult[T any] struct {
	FeedSite feed.FeedSite[T]
	Error    error
}
