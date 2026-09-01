package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
	"github.com/FooNameBar/gator/internal/rss"
	"github.com/google/uuid"
)

func GetFeeds(s *data.State, cmd Command) error {

	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("rss.FetchFeed: %v\n", err)
	}

	fmt.Print(feed.String())
	return nil
}

func AddFeed(s *data.State, cmd Command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Need to specify a name and url for a feed")
	}

	feedname, url := cmd.Args[0], cmd.Args[1]
	username := s.Config.CurrentUserName
	user, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("DB.GetUser: user %s not found: %v\n", feedname, err)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedname,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("DB.CreateFeed: %v\n", err)
	}

	fmt.Printf("%v\n", feed)

	return nil
}

func ListFeeds(s *data.State, cmd Command) error {
	feeds, err := s.DB.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("DB.ListFeeds: %v\n", err)
	}

	var out strings.Builder
	for _, f := range feeds {
		fmt.Fprintf(&out, "Name: %s\nURL: %s\nUser: %s\n", f.Feedname, f.Url, f.Username)
	}

	fmt.Print(out.String())
	return nil
}
