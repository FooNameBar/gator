package rss

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func (rf *RSSFeed) String() string {
	var items strings.Builder
	for _, item := range rf.Channel.Item {
		items.WriteString(item.String())
		items.WriteString("\n")
	}
	return fmt.Sprintf(
		`Title: %s
Link: %s
Description: %s
Item:
%s`, rf.Channel.Title, rf.Channel.Link, rf.Channel.Description, items.String())
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (ri *RSSItem) String() string {
	return fmt.Sprintf(
		`- Title: %s
- Link: %s
- Description: %s
- PubDate: %s
`, ri.Title, ri.Link, ri.Description, ri.PubDate)
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	var data []byte
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("NewRequestWithContext: url %s: %v\n", feedURL, err)
	}

	req.Header.Set("User-Agent", "gator")

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do: %v\n", err)
	}
	defer res.Body.Close()

	xmlData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v\n", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(xmlData, &feed)
	if err != nil {
		return nil, fmt.Errorf("xml.Unmarshal: %v\n", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
