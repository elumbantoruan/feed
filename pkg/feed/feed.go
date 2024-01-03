package feed

import (
	"time"
)

type Feed struct {
	Site     string    `json:"site"`
	Icon     string    `json:"icon"`
	Link     string    `json:"link"`
	RSS      string    `json:"rss"`
	Updated  time.Time `json:"updated"`
	Articles []Article `json:"articless"`
}

type Article struct {
	ID        string    `json:"id"`
	Published time.Time `json:"published"`
	Updated   time.Time `json:"updated"`
	Title     string    `json:"title"`
	Authors   []string  `json:"author"`
	Link      string    `json:"link"`
	Content   string    `json:"content"`
	Blob      *string   `json:"blob"`
}

type URLType struct {
	URL  string
	Type string
}
