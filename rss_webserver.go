package rss

import (
	"fmt"
	"io/ioutil"
	"os"
)

var database = NewRssDatabase()

func main() {
	if len(os.Args) == 2 {
		filepath := os.Args[1]
		parseFile(filepath)
		return
	}
}

func parseFile(filepath string) {
	fmt.Printf("Importing file: %s\n", filepath)

	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	db := NewRssDatabase()

	parser := NewParser(filepath, string(fileContents))
	go parser.Parse()

	entries := make([]Entry, 0, 20)
	feedId := int64(-1)
	feedOpen, entryOpen := true, true
parseLoop:
	for {
		if !feedOpen && !entryOpen {
			break parseLoop
		}
		select {
		case feed, feedOk := <-parser.feed:
			if feedOk {
				feedId, err = db.insertFeed(&feed)
				if err != nil {
					panic(err)
				}
			} else {
				feedOpen = false
			}
		case entry, entryOk := <-parser.entry:
			if entryOk {
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
