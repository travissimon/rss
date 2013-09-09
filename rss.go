package rss

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type FeedStub struct {
	Id    int64
	Title string
}

type Feed struct {
	Id          int64
	Url         string
	Feed        string
	Title       string
	Link        string
	Subtitle    string
	Copyright   string
	Author      string
	PublishDate time.Time
	Category    string
	Generator   string
	Logo        string
	Icon        string
}

type Entry struct {
	Id          int64
	FeedId      int64
	Title       string
	Link        string
	Subtitle    string
	Guid        string
	UpdatedDate time.Time
	Summary     string
	Encoded     string
	Content     string
	Source      string
	Comments    string
	Thumbnail   string
	Length      string
	Type        string
	Url         string
}

func (e *Entry) String() string {
	return "{" + string(e.Id) + " " + e.Title + "}"
}

type EntryList struct {
	feed    *Feed
	entries []*Entry
}

type RssEngine struct {
	db *RssDatabase
}

func NewRssEngine(database, username, password string) *RssEngine {
	rss := new(RssEngine)
	rss.db = NewRssDatabase(database, username, password)
	return rss
}

func (rss *RssEngine) GetFeedsForUser(userId int64) (feeds []*Feed, err error) {
	feeds, err = rss.db.getFeedsForUser(userId)

	if err != nil {
		fmt.Println(err)
	}

	return
}

func (rss *RssEngine) GetFeedStubsForUser(userId int64) (feeds []*FeedStub, err error) {
	feeds, err = rss.db.getFeedStubsForUser(userId)

	if err != nil {
		fmt.Println(err)
	}

	return
}

// GetEntriesForFeed gets all entries for a feed.
func (rss *RssEngine) GetEntriesForFeed(feedId int64) (entries []*Entry, err error) {
	entries, err = rss.db.getEntriesByFeedId(feedId)

	if err != nil {
		fmt.Println(err)
	}

	return
}

func (rss *RssEngine) AddFeedForUser(userId int64, feedUrl string) (feed *Feed, entries []*Entry, err error) {
	feedExists, subscribed := rss.db.getFeedStatusForUser(userId, feedUrl)

	fmt.Printf("Adding %v, exists: %v, subscribed: %v\n", feedUrl, feedExists, subscribed)

	if subscribed {
		feed, err = rss.db.getFeedByUrl(feedUrl)
		entries, err = rss.db.getEntriesByFeedId(feed.Id)
		return
	}

	var feedId int64 = -1

	if !feedExists {
		// download feed
		contents, err := rss.downloadRssFile(feedUrl)
		fmt.Printf("Contents: %s\n", contents[0:120])
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return nil, nil, err
		}
		// parse feed
		feed, entries, err = rss.parseFeed(feedUrl, string(contents))
		fmt.Printf("err: %v\nfeed: %v\nEntries: %v\n", err, feed, entries)
		// store in database
		if err != nil {
			fmt.Printf("Error parsing %s: %s\n", feedUrl, err.Error())
			return nil, nil, err
		}

		feedId, err = rss.db.insertFeed(feed)
		fmt.Printf("feedid: %v, err: %v\n", feedId, err)
		if err != nil {
			return nil, nil, err
		}
		for _, entry := range entries {
			entry.FeedId = feedId
			rss.db.insertEntry(entry)
		}
		// start updater go routine

	}

	// add subscription for feed
	fmt.Printf("Adding subscription for %v to feed %v\n", userId, feedId)
	err = rss.AddSubscription(userId, feedId)
	return
}

// Currently this just adds subscription to feed, not to entries. Need to fix
func (rss *RssEngine) AddSubscription(userId, feedId int64) (err error) {
	err = rss.db.AddSubscription(userId, feedId, true)
	return
}

func (rss *RssEngine) downloadRssFile(feedUrl string) (contents string, err error) {
	resp, err := http.Get(feedUrl)
	if err != nil {
		fmt.Printf(err.Error())
		return "", err
	}
	byteArray, e := ioutil.ReadAll(resp.Body)
	contents = string(byteArray)
	err = e
	return
}

func (rss *RssEngine) parseFeed(feedUrl, rssContents string) (feed *Feed, entries []*Entry, err error) {
	parser := NewParser(feedUrl, rssContents)
	feed, entries, err = parser.Parse()
	return
}
