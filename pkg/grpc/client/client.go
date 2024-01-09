package client

import (
	"context"
	"github/elumbantoruan/feed/pkg/feed"
	pb "github/elumbantoruan/feed/pkg/feedproto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type grpcClient struct {
	serverAddr    string
	serviceClient pb.FeedServiceClient
}

func NewGRPCClient(serverAddr string) grpcClient {

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewFeedServiceClient(conn)

	return grpcClient{serviceClient: client}
}

func (g grpcClient) AddSiteFeed(ctx context.Context, site feed.Feed) error {
	pbFeed := pb.Feed{
		Site: site.Site,
		Rss:  site.RSS,
		Type: site.Type,
	}
	_, err := g.serviceClient.AddSiteFeed(ctx, &pbFeed)
	return err
}

func (g grpcClient) GetSitesFeed(ctx context.Context) ([]feed.Feed, error) {
	var feeds []feed.Feed
	pbfeeds, err := g.serviceClient.GetSitesFeed(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	for _, pbfeed := range pbfeeds.Feeds {
		ts := pbfeed.Updated.AsTime()
		feed := feed.Feed{
			ID:      pbfeed.Id,
			Site:    pbfeed.Site,
			RSS:     pbfeed.Rss,
			Type:    pbfeed.Type,
			Updated: &ts,
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (g grpcClient) UpdateSiteFeed(ctx context.Context, feed feed.Feed) error {
	pbFeed := pb.Feed{
		Id:      feed.ID,
		Site:    feed.Site,
		Icon:    feed.Icon,
		Link:    feed.Link,
		Rss:     feed.RSS,
		Type:    feed.Type,
		Updated: timestamppb.New(*feed.Updated),
	}
	_, err := g.serviceClient.UpdateSiteFeed(ctx, &pbFeed)
	return err
}

func (g grpcClient) AddArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {
	var authors []string
	for _, author := range article.Authors {
		authors = append(authors, author)
	}
	pbArticleSite := &pb.ArticleSite{
		Siteid: siteID,
		Article: &pb.Article{
			Id:          article.ID,
			Title:       article.Title,
			Link:        article.Link,
			Published:   timestamppb.New(article.Published),
			Description: article.Description,
			Content:     article.Content,
			Authors:     authors,
		},
	}
	aid, err := g.serviceClient.AddArticle(ctx, pbArticleSite)
	if err != nil {
		return -1, err
	}
	return aid.Identifier, nil
}

func (g grpcClient) GetArticles(ctx context.Context) ([]feed.ArticleSite[int64], error) {
	var articles []feed.ArticleSite[int64]
	pbArticles, err := g.serviceClient.GetArticles(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	for _, pbArticle := range pbArticles.ArticlesSite {
		var articleSite feed.ArticleSite[int64]
		articleSite.SiteID = pbArticle.Siteid
		articleSite.Article = feed.Article{
			ID:          pbArticle.Article.Id,
			Title:       pbArticle.Article.Title,
			Link:        pbArticle.Article.Link,
			Published:   pbArticle.Article.Published.AsTime(),
			Description: pbArticle.Article.Description,
			Content:     pbArticle.Article.Content,
		}
		articles = append(articles, articleSite)
	}

	return articles, nil
}
