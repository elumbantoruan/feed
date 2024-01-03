package feed

func CrawlerFactory(urlType URLType) Crawler {
	if urlType.Type == "atom" {
		return NewAtomCrawler(urlType.URL)
	}
	return NewRssCrawler(urlType.URL)
}

type Crawler interface {
	Download(url string) (*Feed, error)
}
