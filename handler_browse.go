package main

import (
	"fmt"
	"context"
	"strconv"
	"example.com/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	// check the args, if no limit set it 2
	limit := 2
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}
	limit32 := int32(limit)

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: limit32,
	})
	if err != nil {
		return fmt.Errorf("unable to get posts for user: %s with error: %w", user.ID, err)
	}

	for _, post := range posts {
		fmt.Printf("Title: %s\n", post.Title)
		if post.Description.Valid {
			fmt.Printf("Description: %s\n", post.Description.String)
		}
		fmt.Printf("URL: %s\n", post.Url)
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %s\n", post.PublishedAt.Time.Format("Jan 2, 2006 at 3:04 PM"))
		}
		fmt.Println("--------")
	}

	return nil
}

