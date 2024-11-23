package commands

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rss_feed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &rss_feed)
	if err != nil {
		return nil, err
	}

	unescape(&rss_feed)

	return &rss_feed, err
}

func unescape(rss_feed *RSSFeed) {
	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i := range rss_feed.Channel.Item {
		rss_feed.Channel.Item[i].Title = html.UnescapeString(rss_feed.Channel.Item[i].Title)
		rss_feed.Channel.Item[i].Description = html.UnescapeString(rss_feed.Channel.Item[i].Description)
	}
}
