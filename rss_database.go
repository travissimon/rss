package rss

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
	db                         *sql.DB
	getAllFeedsStmt            *sql.Stmt
	getFeedsByUserIdStmt       *sql.Stmt
	getFeedStubsByUserIdStmt   *sql.Stmt
	insertSubscriptionStmt     *sql.Stmt
	insertFeedStmt             *sql.Stmt
	getFeedByIdStmt            *sql.Stmt
	insertEntryStmt            *sql.Stmt
	getEntriesByFeedIdStmt     *sql.Stmt
	isUserSubscribedToFeedStmt *sql.Stmt
	getFeedIdByUrlStmt         *sql.Stmt
	getFeedByUrlStmt           *sql.Stmt
}

func (rss *RssDatabase) panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewRssDatabase(database, username, password string) *RssDatabase {
	rss := new(RssDatabase)
	db, err := sql.Open("mymysql", fmt.Sprintf("%s/%s/%s", database, username, password))
	rss.panicOnError(err)
	rss.db = db

	getAll, err := db.Prepare(getAllFeedsSQL)
	rss.panicOnError(err)
	rss.getAllFeedsStmt = getAll

	getForUser, err := db.Prepare(getFeedsForUserIdSQL)
	rss.panicOnError(err)
	rss.getFeedsByUserIdStmt = getForUser

	getStubsForUser, err := db.Prepare(getFeedStubsForUserIdSQL)
	rss.panicOnError(err)
	rss.getFeedStubsByUserIdStmt = getStubsForUser

	insSub, err := db.Prepare(insertSubscriptionSQL)
	rss.panicOnError(err)
	rss.insertSubscriptionStmt = insSub

	insFeed, err := db.Prepare(insertFeedSQL)
	rss.panicOnError(err)
	rss.insertFeedStmt = insFeed

	feedById, err := db.Prepare(getFeedByIdSQL)
	rss.panicOnError(err)
	rss.getFeedByIdStmt = feedById

	insEntry, err := db.Prepare(insertEntrySQL)
	rss.panicOnError(err)
	rss.insertEntryStmt = insEntry

	entriesByFeed, err := db.Prepare(getEntriesByFeedIdSQL)
	rss.panicOnError(err)
	rss.getEntriesByFeedIdStmt = entriesByFeed

	feedIdByUrl, err := db.Prepare(getFeedIdByUrlSQL)
	rss.panicOnError(err)
	rss.getFeedIdByUrlStmt = feedIdByUrl

	feedByUrl, err := db.Prepare(getFeedByUrlSQL)
	rss.panicOnError(err)
	rss.getFeedByUrlStmt = feedByUrl

	isSubscribed, err := db.Prepare(isUserSubscribedToFeedSQL)
	rss.panicOnError(err)
	rss.isUserSubscribedToFeedStmt = isSubscribed

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

// GetAllFeeds returns all the feeds!
func (rss *RssDatabase) getAllFeeds() (feeds []*Feed, err error) {
	rows, err := rss.getAllFeedsStmt.Query()

	if err != nil {
		return nil, err
	}

	feeds = make([]*Feed, 0, 40)
	for rows.Next() {
		feed, err := getFeedFromRow(rows)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (rss *RssDatabase) getFeedsForUser(userId int64) (feeds []*Feed, err error) {
	rows, err := rss.getFeedsByUserIdStmt.Query(userId)

	if err != nil {
		return nil, err
	}

	feeds = make([]*Feed, 0, 40)
	for rows.Next() {
		feed, err := getFeedFromRow(rows)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	return feeds, nil
}

func (rss *RssDatabase) getFeedStubsForUser(userId int64) (feeds []*FeedStub, err error) {
	rows, err := rss.getFeedStubsByUserIdStmt.Query(userId)

	if err != nil {
		return nil, err
	}

	feeds = make([]*FeedStub, 0, 40)
	for rows.Next() {
		stub := new(FeedStub)
		err = rows.Scan(
			&stub.Id,
			&stub.Title,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, stub)
	}
	return feeds, nil
}

// InsertFeed inserts a feed into the database
func (rss *RssDatabase) insertFeed(feed *Feed) (id int64, err error) {
	res, err := rss.insertFeedStmt.Exec(
		feed.Id,
		feed.Url,
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
	id, err = res.LastInsertId()
	return
}

func (rss *RssDatabase) getFeedById(id uint64) (feed *Feed, err error) {
	rows := rss.getFeedByIdStmt.QueryRow(id)
	feed, err = getFeedFromRow(rows)
	return
}

func (rss *RssDatabase) getFeedByUrl(feedUrl string) (feed *Feed, err error) {
	rows := rss.getFeedByUrlStmt.QueryRow(feedUrl)
	feed, err = getFeedFromRow(rows)
	return
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func getFeedFromRow(row Scanner) (feed *Feed, err error) {
	f := new(Feed)
	err = row.Scan(
		&f.Id,
		&f.Url,
		&f.Title,
		&f.Link,
		&f.Subtitle,
		&f.Copyright,
		&f.Author,
		&f.PublishDate,
		&f.Category,
		&f.Logo,
		&f.Icon,
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Entry methods

func (rss *RssDatabase) insertEntry(entry *Entry) (id int64, err error) {
	res, err := rss.insertEntryStmt.Exec(
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
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (rss *RssDatabase) getEntriesByFeedId(feedId int64) (entries []*Entry, err error) {
	rows, err := rss.getEntriesByFeedIdStmt.Query(feedId)

	if err != nil {
		return nil, err
	}

	entries = make([]*Entry, 0, 20)
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
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (rss *RssDatabase) getFeedStatusForUser(userId int64, feedUrl string) (feedExists, subscriptionExists bool) {
	feedExists = false
	subscriptionExists = false

	_, err := rss.GetFeedIdByUrl(feedUrl)

	// presence of an error means there was no row
	feedExists = (err == nil)

	var rowId int64
	err = rss.isUserSubscribedToFeedStmt.QueryRow(userId, feedUrl).Scan(&rowId)
	subscriptionExists = (err == nil)

	return
}

func (rss *RssDatabase) GetFeedIdByUrl(feedUrl string) (id int64, err error) {
	err = rss.getFeedIdByUrlStmt.QueryRow(feedUrl).Scan(&id)
	return
}

func (rss *RssDatabase) AddSubscription(userId, feedId int64, isRead bool) error {
	_, err := rss.insertSubscriptionStmt.Exec(userId, feedId, isRead)
	return err
}

// Feed SQL
var getAllFeedsSQL string = `
SELECT
  Id,
  Url,
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
;`

var insertFeedSQL string = `
INSERT INTO rss.Feed (
  Id,
  Url,
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
  ?,
  ?
);`

var getFeedByIdSQL string = `
SELECT
  Id,
  Url,
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
;`

var getFeedIdByUrlSQL string = `
SELECT Id
FROM rss.Feed
WHERE Url = ?
`
var getFeedByUrlSQL string = `
SELECT
  Id,
  Url,
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
WHERE Url = ? 
`

var getFeedsForUserIdSQL string = `
SELECT
  feed.Id,
  feed.Url,
  feed.Title,
  feed.Link,
  feed.Subtitle,
  feed.Copyright,
  feed.Author,
  feed.PublishDate,
  feed.Category,
  feed.Logo,
  feed.Icon
FROM rss.Subscription sub
  INNER JOIN rss.Feed feed
    ON sub.FeedId = feed.Id
WHERE sub.UserId = ?
ORDER BY feed.Title
;`

var getFeedStubsForUserIdSQL string = `
SELECT
  feed.Id,
  feed.Title
FROM rss.Subscription sub
  INNER JOIN rss.Feed feed
    ON sub.FeedId = feed.Id
WHERE sub.UserId = ?
ORDER BY feed.Title
;`

// Entry SQL
var insertEntrySQL string = `
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
;`

var getEntriesByFeedIdSQL string = `
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
;`

// Subscription SQL

var insertSubscriptionSQL string = `
INSERT INTO rss.Subscription (
  UserId,
  FeedId,
  UnreadItems
) VALUES (
  ?,
  ?,
  ?
);`

var isUserSubscribedToFeedSQL string = `
SELECT
  feed.Id
FROM rss.Subscription sub
  INNER JOIN rss.Feed feed
    ON sub.FeedId = feed.Id
WHERE sub.UserId = ?
AND feed.Url = ?
;`
