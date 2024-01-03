package feed

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func download[T any](url string) (*T, error) {

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

func ParseDateTime(dateTime string) (time.Time, error) {
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
			return t, nil
		}
	}
	return time.Time{}, errors.New("parseDateTime: " + err.Error())
}
