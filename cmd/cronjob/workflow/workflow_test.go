package workflow

import (
	"context"
	"encoding/xml"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elumbantoruan/feed/cmd/cronjob/mock"
	"github.com/elumbantoruan/feed/pkg/crawler"
	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/stretchr/testify/assert"
)

func TestWorkflow_Run(t *testing.T) {
	rssContent := createRSSMock()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := xml.Marshal(rssContent)
		assert.NoError(t, err)
		w.Write(data)
	}))

	// for unit test, embed the URL from httptest into the context
	// so the URL can be passed into gRPC Client in RSS field
	// which will be used in the Crawler.Download(url)
	ctx := context.WithValue(context.Background(), "url", srv.URL)

	var storage []feed.Article
	client := &mock.MockGRPCClient{Storage: storage}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	work := New(client, logger)
	res, err := work.Run(ctx)
	assert.NoError(t, err)

	// GetSites only called once
	assert.Equal(t, 1, client.GetSitesCount)
	// UpdateSiteFeedCount called twice
	assert.Equal(t, 1, client.UpdateSiteFeedCount)
	// AddArticleCount called twice
	assert.Equal(t, 2, client.AddArticleCount)

	// res is result which reflects the number of site
	assert.Equal(t, 1, len(res))
	// there are two articles, reflected from data provided in createRSSMock
	assert.Equal(t, 2, res[0].Metric.NewArticles)

	// asserting from mock storage
	for i, article := range storage {
		assert.Equal(t, article.ID, rssContent.Channel.Item[i].Guid.Text)
		pubDate, _ := crawler.ParseDateTime(rssContent.Channel.Item[i].PubDate)

		assert.Equal(t, article.Published, pubDate)
		assert.Equal(t, article.Updated, pubDate)
		assert.Equal(t, article.Title, rssContent.Channel.Item[i].Title)
		assert.Equal(t, article.Link, rssContent.Channel.Item[i].Link)
		assert.Equal(t, article.Description, rssContent.Channel.Item[i].Description)
		assert.Equal(t, article.Content, rssContent.Channel.Item[i].Content)
	}

}

func createRSSMock() *feed.Rss {
	pubDate := "Sun, 21 Jan 2024 21:48:54 +0000"
	items := []feed.Item{
		{
			Text:        "Text1",
			Title:       "Title1",
			Link:        "Link1",
			PubDate:     pubDate,
			Content:     "Content1",
			Description: "Description1",
		},
		{
			Text:        "Text2",
			Title:       "Title2",
			Link:        "Link2",
			PubDate:     pubDate,
			Content:     "Content2",
			Description: "Description2",
		},
	}
	items[0].Guid.Text = "guid1"
	items[1].Guid.Text = "guid2"

	return &feed.Rss{
		Channel: &feed.Channel{
			Item:          items,
			LastBuildDate: pubDate,
		},
	}
}
