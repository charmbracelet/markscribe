package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// RSSEntry represents a single RSS entry.
type RSSEntry struct {
	Title       string
	Author      string
	Description string
	URL         string
	PublishedAt time.Time
}

func rssFeed(url string, count int) []RSSEntry {
	var r []RSSEntry

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		panic(err)
	}

	for _, v := range feed.Items {
		// fmt.Printf("%+v\n", v)

		r = append(r, RSSEntry{
			Title:       v.Title,
			Author:      v.Author.Name,
			Description: v.Description,
			URL:         v.Link,
			PublishedAt: *v.PublishedParsed,
		})
		if len(r) == count {
			break
		}
	}

	return r
}
