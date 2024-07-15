package main

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
)

var team = map[string]string{
	"Christian Rocha": "meowgorithm",
	"Bashbunni":       "bashbunni",
	"Ayman Bagabas":   "aymanbagabas",
	"Carlos Becker":   "caarlos0",
	"Maas Lalani":     "maaslalani",
	"Charm":           "charmbracelet",
}

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
		author := v.Author.Name
		if profile, ok := team[v.Author.Name]; ok {
			author = fmt.Sprintf("[%s](https://github.com/%s)", profile, profile)
		}

		r = append(r, RSSEntry{
			Title:       v.Title,
			Author:      author,
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
