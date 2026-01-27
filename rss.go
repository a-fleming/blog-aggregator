package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"www.github.com/a-fleming/blog-aggregator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string, timeoutSec int) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feedData RSSFeed
	err = xml.Unmarshal(data, &feedData)
	if err != nil {
		return &RSSFeed{}, err
	}

	feedData.Channel.Title = html.UnescapeString(feedData.Channel.Title)
	feedData.Channel.Link = html.UnescapeString(feedData.Channel.Link)
	feedData.Channel.Description = html.UnescapeString(feedData.Channel.Description)
	for _, item := range feedData.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Link = html.UnescapeString(item.Link)
		item.Description = html.UnescapeString(item.Description)
		item.PubDate = html.UnescapeString(item.PubDate)
	}

	return &feedData, nil
}

func printFeed(feed RSSFeed) {
	fmt.Printf("RSS Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		if len(item.Title) == 0 {
			continue
		}

		fmt.Printf(" * Title: %s\n", item.Title)
		fmt.Printf(" * URL:%s\n", item.Link)
		fmt.Printf(" * PubDate:%s\n", item.PubDate)
		fmt.Printf(" * Description:%s\n", item.Description)
		fmt.Println()
	}
}

func scrapeFeeds(ctx context.Context, s *state, timeoutSec int) error {
	feedDbInfo, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}
	rssFeed, err := fetchFeed(ctx, feedDbInfo.Url, timeoutSec)
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(ctx, feedDbInfo.ID)
	if err != nil {
		return err
	}
	unescapeAndTrimFeed(*rssFeed)
	printFeed(*rssFeed)
	err = addPosts(ctx, s, *rssFeed, feedDbInfo.ID)
	if err != nil {
		return err
	}
	return nil
}

func addPosts(ctx context.Context, s *state, feed RSSFeed, feedID uuid.UUID) error {
	for _, post := range feed.Channel.Item {
		if len(post.Title) == 0 {
			continue
		}

		pubTime, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			return err
		}

		var description sql.NullString
		if len(post.Description) != 0 {
			description = sql.NullString{
				String: post.Description,
				Valid:  true,
			}
		} else {
			description = sql.NullString{
				Valid: false,
			}
		}
		createPostParams := database.CreatePostParams{
			Title:       post.Title,
			Url:         post.Link,
			Description: description,
			PublishedAt: pubTime,
			FeedID:      feedID,
		}
		postResults, err := s.db.CreatePost(ctx, createPostParams)
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				fmt.Printf("skipping duplicate post\n")
				continue
			} else {
				return err
			}
		}
		fmt.Printf("postResults: %+v\n\n", postResults)
	}
	return nil
}

func unescapeAndTrimFeed(feed RSSFeed) {
	feed.Channel.Title = strings.TrimSpace(html.UnescapeString(feed.Channel.Title))
	feed.Channel.Link = strings.TrimSpace(html.UnescapeString(feed.Channel.Link))
	feed.Channel.Description = strings.TrimSpace(html.UnescapeString(feed.Channel.Description))
	for _, post := range feed.Channel.Item {
		post.Title = strings.TrimSpace(html.UnescapeString(post.Title))
		post.Link = strings.TrimSpace(html.UnescapeString(post.Link))
		post.PubDate = strings.TrimSpace(html.UnescapeString(post.PubDate))
		post.Description = strings.TrimSpace(html.UnescapeString(post.Description))
	}
}
