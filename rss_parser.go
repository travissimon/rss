package rss

import (
	"fmt"
	"html"
	"strings"
	"time"
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
			"encoded":     handleEntryEncoded,
			// "content":     handleEntryContent,
			"source":    handleEntrySource,
			"comments":  handleEntryComments,
			"enclosure": handleEntryEnclosure,
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
	var err error
	feed.PublishDate, err = parseDate(lexeme.val)
	if err != nil {
		fmt.Println(err)
	}
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
	var err error
	entry.UpdatedDate, err = parseDate(lexeme.val)
	if err != nil {
		fmt.Println(err)
	}
	skipUntilType(l, itemCloseTag)
}

func handleEntrySummary(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Summary = lexeme.val
	skipUntilType(l, itemCloseTag)
}

func handleEntryEncoded(l *lexer, entry *Entry) {
	lexeme := l.nextItem()
	entry.Summary = html.UnescapeString(lexeme.val)
	skipUntilType(l, itemCloseTag)
}

/*func handleEntryContent(l *lexer, entry *Entry) {
	for lexeme := l.nextItem(); lexeme.typ != itemCloseTag; lexeme = l.nextItem() {
		if lexeme.typ == itemText {
			entry.Content = lexeme.val
		}
	}
}*/

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

// date parsing 'borrowed' from mjibson's wondeful goread
// https://github.com/mjibson/goread/blob/master/goapp/utils.go
var dateFormats = []string{
	"01-02-2006",
	"01/02/2006 15:04:05 MST",
	"02 Jan 2006 15:04 MST",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05 MST",
	"02 Jan 2006 15:04:05 UT",
	"02 Jan 2006",
	"02-01-2006 15:04:05 MST",
	"02.01.2006 -0700",
	"02.01.2006 15:04:05",
	"02/01/2006 15:04:05",
	"02/01/2006",
	"06-1-2 15:04",
	"06/1/2 15:04",
	"1/2/2006 15:04:05 MST",
	"1/2/2006 3:04:05 PM",
	"15:04 02.01.2006 -0700",
	"2 Jan 2006 15:04:05 MST",
	"2 Jan 2006",
	"2 January 2006 15:04:05 -0700",
	"2 January 2006",
	"2006 January 02",
	"2006-01-02 00:00:00.0 15:04:05.0 -0700",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05 -0700",
	"2006-01-02 15:04:05 MST",
	"2006-01-02 15:04:05-07:00",
	"2006-01-02 15:04:05Z",
	"2006-01-02",
	"2006-01-02T15:04-07:00",
	"2006-01-02T15:04:05 -0700",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05-0700",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05-07:00:00",
	"2006-01-02T15:04:05:-0700",
	"2006-01-02T15:04:05:00",
	"2006-01-02T15:04:05Z",
	"2006-1-02T15:04:05Z",
	"2006-1-2 15:04:05",
	"2006-1-2",
	"2006/01/02",
	"6-1-2 15:04",
	"6/1/2 15:04",
	"Jan 02 2006 03:04:05PM",
	"Jan 2, 2006 15:04:05 MST",
	"Jan 2, 2006 3:04:05 PM MST",
	"January 02, 2006 03:04 PM",
	"January 02, 2006 15:04",
	"January 02, 2006 15:04:05 MST",
	"January 02, 2006",
	"January 2, 2006 03:04 PM",
	"January 2, 2006 15:04:05 MST",
	"January 2, 2006 15:04:05",
	"January 2, 2006",
	"January 2, 2006, 3:04 p.m.",
	"Mon 02 Jan 2006 15:04:05 -0700",
	"Mon 2 Jan 2006 15:04:05 MST",
	"Mon Jan 2 15:04 2006",
	"Mon Jan 2 15:04:05 2006 MST",
	"Mon, 02 Jan 06 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 02 Jan 2006 15:04 MST",
	"Mon, 02 Jan 2006 15:04:05 --0700",
	"Mon, 02 Jan 2006 15:04:05 -07",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 -07:00",
	"Mon, 02 Jan 2006 15:04:05 00",
	"Mon, 02 Jan 2006 15:04:05 MST -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 MST-07:00",
	"Mon, 02 Jan 2006 15:04:05 UT",
	"Mon, 02 Jan 2006 15:04:05 Z",
	"Mon, 02 Jan 2006 15:04:05",
	"Mon, 02 Jan 2006 15:04:05MST",
	"Mon, 02 Jan 2006 3:04:05 PM MST",
	"Mon, 02 Jan 2006",
	"Mon, 02 January 2006",
	"Mon, 2 Jan 06 15:04:05 -0700",
	"Mon, 2 Jan 06 15:04:05 MST",
	"Mon, 2 Jan 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04",
	"Mon, 2 Jan 2006 15:04:05 -0700 MST",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05 UT",
	"Mon, 2 Jan 2006 15:04:05",
	"Mon, 2 Jan 2006 15:04:05-0700",
	"Mon, 2 Jan 2006 15:04:05MST",
	"Mon, 2 Jan 2006 15:4:5 MST",
	"Mon, 2 Jan 2006",
	"Mon, 2 Jan 2006, 15:04 -0700",
	"Mon, 2 January 2006 15:04:05 -0700",
	"Mon, 2 January 2006 15:04:05 MST",
	"Mon, 2 January 2006, 15:04 -0700",
	"Mon, 2 January 2006, 15:04:05 MST",
	"Mon, 2, Jan 2006 15:4",
	"Mon, Jan 2 2006 15:04:05 -0700",
	"Mon, Jan 2 2006 15:04:05 -700",
	"Mon, January 02, 2006, 15:04:05 MST",
	"Mon, January 2 2006 15:04:05 -0700",
	"Mon,02 Jan 2006 15:04:05 -0700",
	"Mon,02 January 2006 14:04:05 MST",
	"Monday, 02 January 2006 15:04:05 -0700",
	"Monday, 02 January 2006 15:04:05 MST",
	"Monday, 02 January 2006 15:04:05",
	"Monday, 2 Jan 2006 15:04:05 -0700",
	"Monday, 2 Jan 2006 15:04:05 MST",
	"Monday, 2 January 2006 15:04:05 -0700",
	"Monday, 2 January 2006 15:04:05 MST",
	"Monday, January 02, 2006",
	"Monday, January 2, 2006 03:04 PM",
	"Monday, January 2, 2006 15:04:05 MST",
	"Monday, January 2, 2006",
	"Updated January 2, 2006",
	"mon,2 Jan 2006 15:04:05 MST",
	time.ANSIC,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RubyDate,
	time.UnixDate,
}

func parseDate(dateStr string) (t time.Time, err error) {
	d := strings.TrimSpace(dateStr)
	if d == "" {
		err = fmt.Errorf("Empty date string")
		return
	}
	for _, f := range dateFormats {
		if t, err = time.Parse(f, d); err == nil {
			return
		}
	}
	err = fmt.Errorf("Could not parse date: %v", dateStr)
	return
}
