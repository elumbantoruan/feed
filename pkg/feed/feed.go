package feed

import (
	"time"
)

type Site[T any] struct {
	ID           T          `json:"id"`
	Site         string     `json:"site"`
	SiteURL      string     `json:"siteURL"`
	Icon         string     `json:"icon"`
	Link         string     `json:"link"`
	RSS          string     `json:"rss"`
	Type         string     `json:"type"`
	ArticlesHash string     `json:"articlesHash"`
	Updated      *time.Time `json:"updated"`
}

type FeedSite[T any] struct {
	Site     Site[T]   `json:"site"`
	Articles []Article `json:"articles"`
}

type Feed struct {
	ID       int64      `json:"id"`
	Site     string     `json:"site"`
	SiteURL  string     `json:"siteURL"`
	Icon     string     `json:"icon"`
	Link     string     `json:"link"`
	RSS      string     `json:"rss"`
	Type     string     `json:"type"`
	Updated  *time.Time `json:"updated"`
	Articles []Article  `json:"articless"`
}

type Article struct {
	ID          string    `json:"id"`
	Published   time.Time `json:"published"`
	Updated     time.Time `json:"updated"`
	Title       string    `json:"title"`
	Authors     []string  `json:"author"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Blob        *string   `json:"blob"`
}

type ArticleSite[T any] struct {
	SiteID  T       `json:"siteID"`
	Article Article `json:"article"`
}

type FeedSites[T any] []FeedSite[T]
