package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main_old() {
	if len(os.Args) < 2 {
		fmt.Println("You must include a file to parse")
		return
	}
	filepath := os.Args[1]

	fmt.Printf("%s\n", filepath)

	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	db := NewRssDatabase()

	parser := NewParser(filepath, string(fileContents))
	go parser.Parse()

	entries := make([]Entry, 0, 20)
	feedId := -1
	feedOpen, entryOpen := true, true
parseLoop:
	for {
		if !feedOpen && !entryOpen {
			break parseLoop
		}
		select {
		case feed, feedOk := <-parser.feed:
			if feedOk {
				fmt.Printf("feed: %v\n", feed.Title)
				feedId = db.insertFeed(&feed)
			} else {
				feedOpen = false
			}
		case entry, entryOk := <-parser.entry:
			if entryOk {
				fmt.Printf("entry: %v\n", entry.Title)
				entries = append(entries, entry)
			} else {
				entryOpen = false
			}
		}
	}
	// save all buffered entries
	for _, entry := range entries {
		entry.FeedId = feedId
		db.insertEntry(&entry)
	}
}
