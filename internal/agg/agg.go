package agg

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/samuelea/gator/internal/config"
	"github.com/samuelea/gator/internal/database"
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

func AddFeedHandler(state *config.State, command config.Command, user database.User) error {
	if len(command.Args) < 2 {
		return fmt.Errorf("not enough arguments provided. name and url arguments are required")
	}

	feedName := command.Args[0]
	feedUrl := command.Args[1]

	feedEntry, err := state.DbQueries.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedUrl,
		UserID: user.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to create feed entry: %w", err)
	}

	_, err = state.DbQueries.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID: feedEntry.ID,
		UserID: user.ID,
	}) 

	if err != nil {
		return nil
	}

	fmt.Println("Added new feed entry!")
	fmt.Printf("ID: %v\n", feedEntry.ID)
	fmt.Printf("CreatedAT: %v\n", feedEntry.CreatedAt)
	fmt.Printf("UpdatedAT: %v\n", feedEntry.UpdatedAt)
	fmt.Printf("Name: %v\n", feedEntry.Name)
	fmt.Printf("Url: %v\n", feedEntry.Url)
	fmt.Printf("UserID: %v\n", feedEntry.UserID)

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

func FeedsHandler(state *config.State, command config.Command) error {
	feedsAndUsers, err := state.DbQueries.FeedsAndUsers(context.Background())

	if err != nil {
		return nil
	}
	
	for _, feed := range feedsAndUsers {
		fmt.Printf("Feed: %v\n", feed.Name)
		fmt.Printf("url: %v\n", feed.Url)
		fmt.Printf("user: %v\n", feed.Name_2)
	}

	return nil
}

func FollowHandler(state *config.State, command config.Command, user database.User) error {
	feedUrl := command.Args[0]

	feed, err := state.DbQueries.FeedFromUrl(context.Background(), feedUrl)

	if err != nil {
		return err
	}

	feedFollow, err := state.DbQueries.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID: feed.ID,
			UserID: user.ID,
		},
	)

	if err != nil {
		return nil
	}

	fmt.Printf("New feed followed: %s\n", feedFollow.FeedName)
	fmt.Printf("Followed by: %s\n", feedFollow.UserName)

	return nil
}

func FollowingHandler(state *config.State, command config.Command, user database.User) error {
	followed_feeds, err := state.DbQueries.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return nil
	}
	
	for _, feed := range followed_feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}

	return nil
}