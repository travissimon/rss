package main

import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	_ "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/thrsafe"
	"reflect"
	"strings"
)

type RssDatabase struct {
	db *sql.DB
}

func (rss *RssDatabase) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewRssDatabase() *RssDatabase {
	rss := new(RssDatabase)
	db, err := sql.Open("mymysql", "rss/rss/rss")

	rss.panicOnError(err)
	rss.db = db
	return rss
}

// Returns a string with escaped single quotes
func getStr(str string) string {
	return strings.Replace(str, "'", "\\'", -1)
}

func getIntValue(num interface{}) int {
	switch reflect.TypeOf(num).Kind() {
	case reflect.Float64:
		return int(num.(float64))
	case reflect.Int32:
		return num.(int)
	case reflect.Int64:
		return int(num.(int64))
	}

	return 0
}

// Feed methods

// Get all feeds
func (rss *RssDatabase) getAllFeeds() (feeds []*Feed) {
	rows, err := rss.db.Query(`
SELECT
  Id,
  Title,
  Link,
  Subtitle,
  Copyright,
  Author,
  PublishDate,
  Category,
  Logo,
  Icon
FROM rss.Feed
ORDER BY Title
`)
	rss.panicOnError(err)
	feeds = make([]*Feed, 0, 40)
	for rows.Next() {
		feed := new(Feed)
		err = rows.Scan(
			&feed.Id,
			&feed.Title,
			&feed.Link,
			&feed.Subtitle,
			&feed.Copyright,
			&feed.Author,
			&feed.PublishDate,
			&feed.Category,
			&feed.Logo,
			&feed.Icon,
		)
		feeds = append(feeds, feed)
	}
	return
}

// Inserts a feed into the database
func (rss *RssDatabase) insertFeed(feed *Feed) (id int) {
	res, err := rss.db.Exec(`
INSERT INTO rss.Feed (
  Id,
  Title,
  Link,
  Subtitle,
  Copyright,
  Author,
  PublishDate,
  Category,
  Logo,
  Icon
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
);`,
		feed.Id,
		feed.Title,
		feed.Link,
		feed.Subtitle,
		feed.Copyright,
		feed.Author,
		feed.PublishDate,
		feed.Category,
		feed.Logo,
		feed.Icon)

	rss.panicOnError(err)
	id64, _ := res.LastInsertId()
	id = getIntValue(id64)
	return
}

func (rss *RssDatabase) getFeedById(id int) (feed *Feed) {
	rows, err := rss.db.Query(`
SELECT
  Id,
  Title,
  Link,
  Subtitle,
  Copyright,
  Author,
  PublishDate,
  Category,
  Logo,
  Icon
FROM rss.Feed
WHERE Id = ?
`,
		id)
	rss.panicOnError(err)
	feed = new(Feed)
	for rows.Next() {
		err = rows.Scan(
			&feed.Id,
			&feed.Title,
			&feed.Link,
			&feed.Subtitle,
			&feed.Copyright,
			&feed.Author,
			&feed.PublishDate,
			&feed.Category,
			&feed.Logo,
			&feed.Icon,
		)
	}
	return
}

// Entry methods

func (rss *RssDatabase) insertEntry(entry *Entry) (id int) {
	fmt.Printf("Inserting id: %d, feed id: %d\n", entry.Id, entry.FeedId)
	res, err := rss.db.Exec(`
INSERT INTO rss.Entry (
  Id,
  FeedId,
  Title,
  Link,
  Subtitle,
  Guid,
  UpdatedDate,
  Summary,
  Content,
  Source,
  Comments,
  Thumbnail,
  Length,
  Type,
  URL
) VALUES (
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
)
`,
		entry.Id,
		entry.FeedId,
		entry.Title,
		entry.Link,
		entry.Subtitle,
		entry.Guid,
		entry.UpdatedDate,
		entry.Summary,
		entry.Content,
		entry.Source,
		entry.Comments,
		entry.Thumbnail,
		entry.Length,
		entry.Type,
		entry.Url,
	)
	rss.panicOnError(err)
	id64, _ := res.LastInsertId()
	id = getIntValue(id64)
	return
}

func (rss *RssDatabase) getEntriesByFeedId(feedId int) []*Entry {
	rows, err := rss.db.Query(`
SELECT
  Id,
  FeedId,
  Title,
  Link,
  Subtitle,
  Guid,
  Summary,
  Content,
  Source,
  Comments,
  Thumbnail,
  Length,
  Type,
  URL
FROM rss.Entry
WHERE FeedId = ?
`,
		feedId,
	)

	rss.panicOnError(err)

	entries := make([]*Entry, 0, 20)
	for rows.Next() {
		entry := new(Entry)
		err = rows.Scan(
			&entry.Id,
			&entry.FeedId,
			&entry.Title,
			&entry.Link,
			&entry.Subtitle,
			&entry.Guid,
			&entry.Summary,
			&entry.Content,
			&entry.Source,
			&entry.Comments,
			&entry.Thumbnail,
			&entry.Length,
			&entry.Type,
			&entry.Url,
		)
		entries = append(entries, entry)
	}

	return entries
}
