package main

import (
	"fmt"
	"context"
	"time"
	"example.com/internal/database"
	"github.com/google/uuid"
)

func handlerFollow (s *state, cmd command, user database.User) error {
	// handle the follow command logic
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("unable to find feed at url: %s with error: %w", url, err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't add a feed follow to feed: %v for user: %v with error: %w", user.Name, feed.Name, err)
	}

	fmt.Printf("Feed: %s followed for user: %s", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// return all of the names of the feeds that the current user is following
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("unable to get feeds for user: %s with error %w", user.Name, err)
	}
	
	fmt.Printf("You are currently following these feeds:\n")
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// handle the unfollow command logic
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <feed url>", cmd.Name)
	}

	feedURL := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("unable to find feed at %s with error: %w", feedURL, err)
	}

	err = s.db.DeleteFeedFollowByUserAndFeedID(context.Background(), database.DeleteFeedFollowByUserAndFeedIDParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to unfollow %s with error: %w", feedURL, err)
	}

	fmt.Printf("unfollowed feed at %s\n", feedURL)

	return nil
}

