package agg

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/samuelea/gator/internal/config"
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

func AggHandler(state *config.State, command config.Command) error {

	feedUrl := "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), feedUrl)

	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	request.Header.Set("User-Agent", "gator")

	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}

	response, err := httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(body, &feed)

	if err != nil {
		return nil, err
	}

	cleanUpRSS(&feed)

	return &feed, nil
}

func cleanUpRSS(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, rssItem := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(rssItem.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(rssItem.Description)
	}
}