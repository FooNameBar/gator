package commands

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
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

		for _, item := range rssFeed.Channel.Item {
			date, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				// TODO: Add a fallback instead
				return fmt.Errorf("time.Parse: %v\n", err)
			}

			fmt.Printf("adding %s to DB\n", item.Link)
			_, err = s.DB.CreatePost(context.Background(), database.CreatePostParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title:     item.Title,
				Url:       item.Link,
				Description: sql.NullString{
					String: item.Description,
					Valid:  true,
				},
				PublishedAt: date,
				FeedID:      feed.ID,
			})
			if err != nil && !strings.Contains(err.Error(), "duplicate") {
				return fmt.Errorf("DB.CreatePost: %v\n", err)
			}
		}
	}

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

func ExploreFeedPosts(s *data.State, cmd Command, user database.User) error {
	limit := int32(2)
	if len(cmd.Args) == 1 {
		val, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt: %s is invalid: %v\n", cmd.Args[0], err)
		}
		limit = int32(val)
	}

	posts, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("DB.GetPostsForUser: %v\n", err)
	}

	var postBldr strings.Builder
	for _, p := range posts {
		postObj := rss.RSSItem{
			Title:       p.Title,
			Link:        p.Url,
			Description: p.Description.String,
			PubDate:     p.PublishedAt.String(),
		}
		fmt.Fprintf(&postBldr, "%s\n", postObj.String())
	}
	postStr := postBldr.String()
	if postStr == "" {
		postStr = "No posts to browse. Run 'agg' first\n"
	}
	fmt.Print(postStr)
	return nil
}
