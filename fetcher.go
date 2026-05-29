package main
import (
	"fmt"
	"io"
	"encoding/xml"
	"net/http"	
	"context"
	"html"
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to complete request: %s", err)
	}

	req.Header.Set("User-Agent", "gator")
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain response: %s", err)
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %s", err)
	}

	var rf *RSSFeed
	err = xml.Unmarshal(data, &rf)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %s", err)
	}
	rf.Channel.Title = html.UnescapeString(rf.Channel.Title)
	rf.Channel.Description = html.UnescapeString(rf.Channel.Description)
	for i := range rf.Channel.Item {
		rf.Channel.Item[i].Title = html.UnescapeString(rf.Channel.Item[i].Title)
		rf.Channel.Item[i].Description = html.UnescapeString(rf.Channel.Item[i].Description)
	}
	return rf, nil	
}