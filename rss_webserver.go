package main

import (
	"fmt"
	"github.com/travissimon/go-mvc"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var database = NewRssDatabase()

func main() {
	if len(os.Args) == 2 {
		filepath := os.Args[1]
		parseFile(filepath)
		return
	}

	startHttp()
}

func indexController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	feeds := database.getAllFeeds()
	writer := NewIndexWriter(feeds)
	return mvc.Haml(writer, feeds, ctx)
}

func entryController(ctx *mvc.WebContext, params url.Values) mvc.ControllerResult {
	idStr := params.Get("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		panic(err)
	}
	id := int(id64)
	feed := database.getFeedById(id)
	entries := database.getEntriesByFeedId(id)
	entryList := new(EntryList)
	entryList.feed = feed
	entryList.entries = entries

	writer := NewFeedWriter(entryList)
	return mvc.Haml(writer, entryList, ctx)
}

func startHttp() {
	url := "localhost:8080"
	fmt.Printf("Listenting on http://%s\n", url)

	handler := mvc.NewMvcHandler()
	handler.AddRoute("Homepage", "/", mvc.GET, indexController)
	handler.AddRoute("Feed", "/feed/{id}", mvc.GET, entryController)

	http.Handle("/", handler)
	http.ListenAndServe(url, nil)
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
