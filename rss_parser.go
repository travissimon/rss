package rss

import (
	"fmt"
	"html"
	"strings"
	"time"
)

type RssParser struct {
	lexer         *lexer
	feed          *Feed
	entries       []*Entry
	feedHandlers  map[string]feedHandler
	entryHandlers map[string]entryHandler
}

func NewParser(name, input string) *RssParser {
	return &RssParser{
		lexer:   lex(name, input),
		entries: make([]*Entry, 0, 20),
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
			"summary":     handleEntrySummary,
			"description": handleEntrySummary,
			"encoded":     handleEntryEncoded,
			"content":     handleEntryContent,
			"source":      handleEntrySource,
			"comments":    handleEntryComments,
			"enclosure":   handleEntryEnclosure,
		},
	}
}

func (r *RssParser) Parse() (feed *Feed, entries []*Entry, err error) {
	// skip everything before the feed as unnecessary
	r.skipUntilFeedTag()
	r.populateFeed()
	r.populateEntries()

	return r.feed, r.entries, nil
}

type feedHandler func(l *lexer, feed *Feed)
type entryHandler func(l *lexer, entry *Entry)

func skipUntilTagClose(l *lexer) {
	for lexeme := l.nextItem(); lexeme.typ != itemCloseTag && lexeme.typ != itemSelfClosingTag; lexeme = l.nextItem() {
	}
}

func extractTextAndSkip(l *lexer) *lexeme {
	var lexeme lexeme
	for lexeme = l.nextItem(); lexeme.typ != itemText && lexeme.typ != itemHtml && lexeme.typ != itemCloseTag && lexeme.typ != itemSelfClosingTag; lexeme = l.nextItem() {
	}
	if lexeme.typ == itemText || lexeme.typ == itemHtml {
		skipUntilTagClose(l)
		return &lexeme
	}
	return nil
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
	r.feed = new(Feed)
FeedLoop:
	for {
		lexeme := r.lexer.nextItem()

		switch lexeme.typ {
		case itemEOF:
			break FeedLoop
		case itemOpenTag:
			if lexeme.val == "item" {
				return
			}
		}

		handler := r.feedHandlers[lexeme.val]
		if handler == nil {
			continue
		}

		handler(r.lexer, r.feed)
	}
}

// handleFeedTitle handles title tags for the feed secion
func handleFeedTitle(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Title = lexeme.val
}

// handleFeedTitle handles link tags for the feed secion
func handleFeedLink(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Link = lexeme.val
}

func handleFeedSubtitle(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Subtitle = lexeme.val
}

func handleFeedCopyright(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Copyright = lexeme.val
}

func handleFeedAuthor(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Author = lexeme.val
}

func handleFeedPubDate(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	var err error
	feed.PublishDate, err = parseDate(lexeme.val)
	if err != nil {
		fmt.Println(err)
	}
}

func handleFeedCategory(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Category = lexeme.val
}

func handleFeedGenerator(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Generator = lexeme.val
}

func handleFeedLogo(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Logo = lexeme.val
}

func handleFeedIcon(l *lexer, feed *Feed) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	feed.Icon = lexeme.val
}

// Entry handlers
func (r *RssParser) populateEntries() {
DocumentLoop:
	for {
		entry := new(Entry)
	EntryLoop:
		for {
			lexeme := r.lexer.nextItem()

			switch lexeme.typ {
			case itemEOF:
				r.entries = append(r.entries, entry)
				break DocumentLoop
			case itemOpenTag:
				if lexeme.val == "item" {
					r.entries = append(r.entries, entry)
					break EntryLoop
				}
			}

			handler := r.entryHandlers[lexeme.val]
			if handler == nil {
				continue
			}

			handler(r.lexer, entry)
		}
	}
}

func handleEntryTitle(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Title = lexeme.val
}

func handleEntryLink(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Link = lexeme.val
}

func handleEntrySubtitle(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Subtitle = lexeme.val
}

func handleEntryGuid(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Guid = lexeme.val
}

func handleEntryUpdatedDate(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	var err error
	entry.UpdatedDate, err = parseDate(lexeme.val)
	if err != nil {
		fmt.Println(err)
	}
}

func handleEntrySummary(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	if lexeme.typ == itemText {
		entry.Summary = html.UnescapeString(lexeme.val)
	} else if lexeme.typ == itemHtml {
		entry.Summary = lexeme.val
	}
}

func handleEntryEncoded(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	if lexeme.typ == itemText {
		entry.Encoded = html.UnescapeString(lexeme.val)
	} else if lexeme.typ == itemHtml {
		entry.Encoded = lexeme.val
	}
}

func handleEntryContent(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	if lexeme.typ == itemText {
		entry.Content = html.UnescapeString(lexeme.val)
	} else if lexeme.typ == itemHtml {
		entry.Content = lexeme.val
	}
}

func handleEntrySource(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Source = lexeme.val
}

func handleEntryComments(l *lexer, entry *Entry) {
	lexeme := extractTextAndSkip(l)
	if lexeme == nil {
		return
	}
	entry.Comments = lexeme.val
}

func handleEntryThumbnail(l *lexer, entry *Entry) {
	lexeme := l.nextItem() // should be url attribute name
	lexeme = l.nextItem()  // attribute val
	entry.Thumbnail = lexeme.val
	skipUntilTagClose(l)
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
