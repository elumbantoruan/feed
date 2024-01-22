package crawler

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/stretchr/testify/assert"
)

func TestRssCrawler_Download(t *testing.T) {
	rssContent := createRSSMock()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := xml.Marshal(rssContent)
		assert.NoError(t, err)
		w.Write(data)
	}))
	rssCrawler := NewRssCrawler()
	feedData, err := rssCrawler.Download(srv.URL)
	assert.NoError(t, err)

	assert.NotNil(t, feedData)
	assert.Len(t, feedData.Articles, 2)
	for i, article := range feedData.Articles {
		assert.Equal(t, article.ID, rssContent.Channel.Item[i].Guid.Text)
		pubDate, _ := ParseDateTime(rssContent.Channel.Item[i].PubDate)

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
