package commands

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
	"github.com/FooNameBar/gator/internal/rss"
	"github.com/google/uuid"
)

func GetFeeds(s *data.State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Need to specify a time between requests")
	}

	dur, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("time.ParseDuration: %v\n", err)
	}
	ticker := time.Tick(dur)

	for {
		<-ticker
		fmt.Println("Scraping feeds")
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeeds(s *data.State) error {
	feeds, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("DB.GetNextFeedToFetch: %v\n", err)
	}

	var out strings.Builder
	for _, feed := range feeds {
		_, err = s.DB.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
			LastFetchedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
			ID:        feed.ID,
		})
		if err != nil {
			return fmt.Errorf("DB.MarkFeedFetched: %v\n", err)
		}

		rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
		if err != nil {
			return fmt.Errorf("rss.FetchFeed: %v\n", err)
		}
		fmt.Fprintf(&out, "%s\n", rssFeed.String())
	}

	fmt.Print(out.String())
	return nil
}

func AddFeed(s *data.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("Need to specify a name and url for a feed")
	}

	feedname, url := cmd.Args[0], cmd.Args[1]

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

	_, err = s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("DB.CreateFeedFollow: %v\n", err)
	}

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
