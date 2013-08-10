package rss

import (
	"fmt"
	"time"
)

type Feed struct {
	Id          int64
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

// GetEntriesForFeed gets all entries for a feed.
func (rss *RssEngine) GetEntriesForFeed(feedId int64) (entries []*Entry, err error) {
	entries, err = rss.db.getEntriesByFeedId(feedId)

	if err != nil {
		fmt.Println(err)
	}

	return
}
