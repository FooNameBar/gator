package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/FooNameBar/gator/internal/data"
	"github.com/FooNameBar/gator/internal/database"
	"github.com/google/uuid"
)

func UserFollowFeed(s *data.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Follow command requires a URL")
	}

	url := cmd.Args[0]
	feed, err := s.DB.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("DB.GetFeed: url: %s: %v\n", url, err)
	}

	result, err := s.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("DB.CreateFeedFollow: %v\n", err)
	}

	fmt.Printf("Feed: %s\nUser: %s\n", result.FeedName, result.UserName)
	return nil
}

func UserFollowing(s *data.State, cmd Command, user database.User) error {
	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("DB.GetFeedFollowsForUser: user.ID: %d: %v\n", user.ID, err)
	}

	var feednames strings.Builder
	for _, f := range follows {
		fmt.Fprintf(&feednames, "%v\n", f.Feedname)
	}

	fmt.Println(feednames.String())

	return nil
}

func UserUnfollowFeed(s *data.State, cmd Command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Unfollow command requires a URL")
	}

	url := cmd.Args[0]
	feed, err := s.DB.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("DB.GetFeed: url: %s: %v\n", url, err)
	}

	result, err := s.DB.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("DB.DeleteFeedFollow: %v\n", err)
	}

	fmt.Printf("%s unfollowed %s\n", result.UserName, result.FeedName)
	return nil
}

