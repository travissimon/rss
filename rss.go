package main

import ()

type Feed struct {
	Id          int
	Feed        string
	Title       string
	Link        string
	Subtitle    string
	Copyright   string
	Author      string
	PublishDate string
	Category    string
	Generator   string
	Logo        string
	Icon        string
}

type Entry struct {
	Id          int
	FeedId      int
	Title       string
	Link        string
	Subtitle    string
	Guid        string
	UpdatedDate string
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
