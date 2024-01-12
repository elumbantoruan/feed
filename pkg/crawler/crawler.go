package crawler

import (
	"encoding/xml"
	"fmt"
	"github/elumbantoruan/feed/pkg/feed"
	"io"
	"net/http"
	"time"
)

func CrawlerFactory(site feed.Feed) Crawler {
	if site.Type == "atom" {
		return NewAtomCrawler(site.RSS)
	}
	return NewRssCrawler(site.RSS)
}

type Crawler interface {
	Download(url string) (*feed.Feed, error)
}

type Content interface {
	feed.Atom | feed.Rss
}

func download[T Content](url string) (*T, error) {

	data, err := downloader(url)
	if err != nil {
		return nil, err
	}

	var t T
	err = xml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func downloader(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "PostmanRuntime/7.32.3")

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s - %v", "download", url, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func parseDateTime(dateTime string) (time.Time, error) {
	layouts := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC3339Nano,
	}
	var (
		t   time.Time
		err error
	)

	for _, layout := range layouts {
		t, err = time.Parse(layout, dateTime)
		if err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("parseDateTime: %w", err)
}
