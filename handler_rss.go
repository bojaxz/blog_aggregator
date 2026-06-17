package main

import (
    "fmt"
    "context"
    "time"
		"strings"
		"log"
		"database/sql"
    "example.com/internal/database"
    "github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
        // get the rssfeed
        if len(cmd.Args) != 1 {
                return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
        }

				duration := cmd.Args[0]

				timeParse, err := time.ParseDuration(duration)
				if err != nil {
					return fmt.Errorf("unable to parse duration: %s with error %w", duration, err)
				}

				fmt.Printf("Collecting feeds every %v\n", timeParse)

				ticker := time.NewTicker(timeParse)
				for ; ; <-ticker.C {
					err = scrapeFeeds(s)
					if err != nil {
						fmt.Printf("error scraping feed: %v\n", err)
					}
				}

        return nil
}

func scrapeFeeds(s *state) error {
	// get the next feed to fetch from the db
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get next feed with error: %w", err)
	}

	// mark it as fetched
	markedFeed, err := s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("unable to mark feed: %v as fetched with error: %w", markedFeed, err)
	}

	// fetch the feed using the URL
	feedXML, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("unable to fetch feed at url: %s with error: %w", feed.Url, err)
	}

	// iterate over the items in the feed and print their titles to the console
	for _, feedItem := range feedXML.Channel.Item {
		fmt.Println(feedItem.Title)

		// parse the feedItem.PubDate to standard go format with time.Parse()
		var publishedAt time.Time
		var err error

		// try standard RFC1123 with a numeric timezone offset first
		publishedAt, err = time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			// if that fails, try standard RFC1123 with a named timezone
			publishedAt, err = time.Parse(time.RFC1123, feedItem.PubDate)
			if err != nil {
				// handle the error if the time doesnt match either layout
				return fmt.Errorf("unable to parse PubDate to RFC1123Z or RFC1123 with error: %w", err)
			}
		}

		// handle nullable description string (sql.NullString is different than an empty go string)
		description := sql.NullString{}
		if feedItem.Description != "" {
			description.String = feedItem.Description
			description.Valid = true
		}

		createdFeed, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: feedItem.Title,
			Url: feedItem.Link,
			Description: description,
			PublishedAt: sql.NullTime{
				Time: publishedAt,
				Valid: true,
			},
			FeedID: feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
				continue // ingore this duplicate and keep processing other posts!
			}
			log.Printf("unable to create feed %s with error: %v", feed.ID, err)
		}

		// print createdFeed
		fmt.Printf("new post created: %s", createdFeed)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
        // handle the logic for adding a feed
        if len(cmd.Args) != 2 {
                return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
        }

        name := cmd.Args[0]
        url := cmd.Args[1]

        feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
                ID: uuid.New(),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                Name: name,
                Url: url,
                UserID: user.ID,
							})
        if err != nil {
                return fmt.Errorf("couldn't add new feed to user: %v: %w", user.Name, err)
        }

        followingFeed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
                ID: uuid.New(),
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
                UserID: user.ID,
                FeedID: feed.ID,
        })
        if err != nil {
                return fmt.Errorf("unable to follow new feed: %s with error: %w", followingFeed, err)
        }

        printFeed(feed)

        return nil
}

func printFeed(feed database.Feed) {
        fmt.Printf(" * ID:             %s\n", feed.ID)
        fmt.Printf(" * Created:        %s\n", feed.CreatedAt)
        fmt.Printf(" * Updated:        %s\n", feed.UpdatedAt)
        fmt.Printf(" * Name:           %s\n", feed.Name)
        fmt.Printf(" * URL:            %s\n", feed.Url)
        fmt.Printf(" * UserID:         %s\n", feed.UserID)
}

func handlerFeeds(s *state, cmd command) error {
        // print all of the feeds to the terminal
        if len(cmd.Args) > 0 {
                return fmt.Errorf("usage: %s", cmd.Name)
        }

        feeds, err := s.db.GetFeeds(context.Background())
        if err != nil {
                return fmt.Errorf("couldn't get feeds: %w", err)
        }

        fmt.Println(feeds)

        for _, feed := range feeds {
                fmt.Printf(" - Feed Name: %s\n - Feed URL: %s\n - User Name: %s\n", feed.Name, feed.Url, feed.Name_2)
        }

        return nil
}

