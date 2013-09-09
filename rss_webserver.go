package rss

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		filepath := os.Args[1]
		fmt.Printf("importing %q\n", filepath)
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

	//db := NewRssDatabase("rss", "travis", "")

	parser := NewParser(filepath, string(fileContents))
	feed, entries, err := parser.Parse()

	fmt.Printf("Err: %v\nFeed: %v\nEntries: %v\n", err, feed, entries)

	//feedId, err := db.insertFeed(feed)
	// save all buffered entries
	//for _, entry := range entries {
	//entry.FeedId = feedId
	//db.insertEntry(entry)
	//}
}
