package main

import (
// "fmt"
)

type RssParser struct {
	lexer         *lexer
	feed          chan Feed
	entry         chan Entry
	feedHandlers  map[string]feedHandler
	entryHandlers map[string]entryHandler
}

func NewParser(name, input string) *RssParser {
	return &RssParser{
		lexer: lex(name, input),
		feed:  make(chan Feed, 1),
		entry: make(chan Entry, 1),
		feedHandlers: map[string]feedHandler{
			"title":          handleFeedTitle,
			"link":           handleFeedLink,
			"description":    handleFeedSubtitle,
			"subtitle":       handleFeedSubtitle,
			"copyright":      handleFeedCopyright,
			"author":         handleFeedAuthor,
			"managingEditor": handleFeedAuthor,
			"pubDate":        handleFeedPubDate,
			"category":       handleFeedCategory,
			"generator":      handleFeedGenerator,
			"logo":           handleFeedLogo,
			"icon":           handleFeedIcon,
		},
		entryHandlers: map[string]entryHandler{
			"title":       handleEntryTitle,
			"link":        handleEntryLink,
			"subtitle":    handleEntrySubtitle,
			"id":          handleEntryGuid,
			"guid":        handleEntryGuid,
			"pubDate":     handleEntryUpdatedDate,
			"updatedDate": handleEntryUpdatedDate,
			"Summary":     handleEntrySummary,
			"description": handleEntrySummary,
			"content":     handleEntryContent,
			"source":      handleEntrySource,
			"comments":    handleEntryComments,
			"enclosure":   handleEntryEnclosure,
		},
	}
}

func (r *RssParser) Parse() {
	// skip everything before the feed as unnecessary
	r.skipUntilFeedTag()
	r.populateFeed()
	r.populateEntries()
	close(r.feed)
	close(r.entry)
}

type feedHandler func(l *lexer, feed *Feed)
type entryHandler func(l *lexer, entry *Entry)

func skipUntilType(l *lexer, typ lexItemType) {
	for lexeme := l.nextItem(); lexeme.typ != typ; lexeme = l.nextItem() {
	}
}

// Ignore everything (xml declarations, etc) before the openning feed tag
func (r *RssParser) skipUntilFeedTag() {
	for lexeme := r.lexer.nextItem(); lexeme.val != "channel" && lexeme.val != "feed"; lexeme = r.lexer.nextItem() {
	}
}

// Feed handlers
// These funcs assume that the lexer has just returned the open tag.
// For example, handleFeedTitle assumes lexer has just returned
// {itemOpenTaq, "Title"}

func (r *RssParser) populateFeed() {
	feed := Feed{}
FeedLoop:
	for {
		lexeme := r.lexer.nextItem()

		switch lexeme.typ {
		case itemEOF:
			break FeedLoop
		case itemOpenTag:
			if lexeme.val == "item" {
				r.feed <- feed
				break FeedLoop
			}
		}

		handler := r.feedHandlers[lexeme.val]
		if handler == nil {
			continue
		}

		handler(r.lexer, &feed)
	}
}

// handleFeedTitle handles title tags for the feed secion
func handleFeedTitle(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Title = lexeme.val
	skipUntilType(l, itemCloseTag)
}

// handleFeedTitle handles link tags for the feed secion
func handleFeedLink(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Link = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedSubtitle(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Subtitle = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedCopyright(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Copyright = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedAuthor(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Author = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedPubDate(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.PublishDate = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedCategory(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Category = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedGenerator(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Generator = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedLogo(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Logo = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleFeedIcon(l *lexer, feed *Feed) {
	lexeme := l.nextItem()
	feed.Icon = lexeme.val
	skipUntilType(l, itemCloseTag)
}

// Entry handlers
func (r *RssParser) populateEntries() {
DocumentLoop:
	for {
		entry := Entry{}
	EntryLoop:
		for {
			lexeme := r.lexer.nextItem()

			switch lexeme.typ {
			case itemEOF:
				r.entry <- entry
				break DocumentLoop
			case itemOpenTag:
				if lexeme.val == "item" {
					r.entry <- entry
					break EntryLoop
				}
			}

			handler := r.entryHandlers[lexeme.val]
			if handler == nil {
				continue
			}

			handler(r.lexer, &entry)
		}
	}
}

func handleEntryTitle(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Title = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryLink(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Link = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntrySubtitle(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Subtitle = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryGuid(l *lexer, entry *Entry) {
	for lexeme := l.nextItem(); lexeme.typ != itemCloseTag; lexeme = l.nextItem() {
		if lexeme.typ == itemText {
			entry.Guid = lexeme.val
		}
	}
}

func handleEntryUpdatedDate(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.UpdatedDate = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntrySummary(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Summary = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryContent(l *lexer, entry *Entry) {
	for lexeme := l.nextItem(); lexeme.typ != itemCloseTag; lexeme = l.nextItem() {
		if lexeme.typ == itemText {
			entry.Content = lexeme.val
		}
	}
}

func handleEntrySource(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Source = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryComments(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Comments = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryThumbnail(l *lexer, entry *Entry) {
	lexeme := l.nextItem() // should be url attribute name
	lexeme = l.nextItem()  // attribute val
	entry.Thumbnail = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryEnclosure(l *lexer, entry *Entry) {
EnclosureLoop:
	for {
		lexeme := l.nextItem()
		switch lexeme.val {
		case "length":
			lexeme := l.nextItem()
			entry.Length = lexeme.val
		case "type":
			lexeme := l.nextItem()
			entry.Type = lexeme.val
		case "url":
			lexeme := l.nextItem()
			entry.Url = lexeme.val
		}
		if lexeme.typ == itemSelfClosingTag || lexeme.typ == itemCloseTag && lexeme.val == "enclosure" {
			break EnclosureLoop
		}
	}
}
