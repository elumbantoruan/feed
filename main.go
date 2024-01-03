package main

import (
	"fmt"
	"github/elumbantoruan/feed/pkg/crawler"
	"github/elumbantoruan/feed/pkg/feed"
)

func main() {

	urlTypes := []feed.URLType{
		{URL: "https://www.theverge.com/rss/index.xml", Type: "atom"},
		{URL: "https://www.wired.com/feed/tag/ai/latest/rss", Type: "rss"},
		{URL: "https://mashable.com/feeds/rss/all", Type: "rss"},
		{URL: "https://gizmodo.com/rss", Type: "rss"},
		{URL: "https://www.engadget.com/rss.xml", Type: "rss"},
		{URL: "https://readwrite.com/feed/", Type: "rss"},
	}

	for _, ut := range urlTypes {
		cr := crawler.CrawlerFactory(ut)
		f, err := cr.Download(ut.URL)
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println(f.RSS, f.Updated)
	}
}
