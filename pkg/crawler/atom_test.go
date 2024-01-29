package crawler

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/stretchr/testify/assert"
)

func TestAtomCrawler_Download(t *testing.T) {
	atomContent := createAtomMock()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := xml.Marshal(atomContent)
		assert.NoError(t, err)
		w.Write(data)
	}))
	atomCrawler := NewAtomCrawler()
	feedData, err := atomCrawler.Download(context.Background(), srv.URL)
	assert.NoError(t, err)

	assert.NotNil(t, feedData)
	assert.Len(t, feedData.Articles, 2)
	for i, article := range feedData.Articles {
		assert.Equal(t, article.ID, atomContent.Entry[i].ID)
		pubDate, _ := ParseDateTime(atomContent.Entry[i].Published)
		updated, _ := ParseDateTime(atomContent.Entry[i].Updated)
		assert.Equal(t, article.Published, pubDate)
		assert.Equal(t, article.Updated, updated)
		assert.Equal(t, article.Title, atomContent.Entry[i].Title)
		assert.Equal(t, article.Link, atomContent.Entry[i].Link.Href)
		assert.Equal(t, article.Description, atomContent.Entry[i].Title)
		assert.Equal(t, article.Content, atomContent.Entry[i].Content.Text)
	}
}

func createAtomMock() *feed.Atom {
	ts := "2024-01-21T11:26:32-05:00"
	atom := &feed.Atom{
		Title:   "Title",
		Updated: ts,
		Link:    feed.Link{Href: "href1"},
		Icon:    "icon1",
		ID:      "id1",
		Entry: []feed.Entry{
			{
				ID:        "articleID1",
				Published: ts,
				Updated:   ts,
				Title:     "entry title1",
				Link:      feed.Link{Href: "href1"},
				Content:   feed.Content{Text: "Content1"},
			},
			{
				ID:        "articleID2",
				Published: ts,
				Updated:   ts,
				Title:     "entry title2",
				Link:      feed.Link{Href: "href2"},
				Content:   feed.Content{Text: "Content2"},
			},
		},
	}
	return atom
}
