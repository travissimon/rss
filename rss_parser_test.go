package rss

import (
	"fmt"
	"log"
	"testing"
)

func Test_RssParserHookup(t *testing.T) {
	log.Printf("Hookup Succeeded - testing RSS Parser")
}

func Test_SuttersMill(t *testing.T) {
	fmt.Printf("Sutter's Mill: %s . . .\n", suttersMillContent[:60])
}

func Test_EricLippert(t *testing.T) {
	fmt.Printf("Eric Lippert: %s . . .\n", ericLippertContent[:60])
}

func testContent(name, content string, expF *Feed, expEs []*Entry, t *testing.T) {
	actF, actEs := parseFeed(name, content, t)
	// cmpInt64("Id", expF.Id, actF.Id, t)
	cmpStr("feed Url", expF.Url, actF.Url, t)
	cmpStr("feed Feed", expF.Feed, actF.Feed, t)
	cmpStr("feed Title", expF.Title, actF.Title, t)
	cmpStr("feed Link", expF.Link, actF.Link, t)
	cmpStr("feed Subtitle", expF.Subtitle, actF.Subtitle, t)
	cmpStr("feed Copyright", expF.Copyright, actF.Copyright, t)
	cmpStr("feed Author", expF.Author, actF.Author, t)
	cmpTime("feed PublishDate", expF.PublishDate, actF.PublishDate, t)
	cmpStr("feed Category", expF.Category, actF.Category, t)
	cmpStr("feed Generator", expF.Generator, actF.Generator, t)
	cmpStr("feed Logo", expF.Logo, actF.Logo, t)
	cmpStr("feed Icon", expF.Icon, actF.Icon, t)

	if len(expEs) != len(actEs) {
		t.Error(fmt.Sprintf("Lenght of slices are not equal. Expected slice length: %d, while actual slice entries length: %d"))
	}

	for expected, i := range expEs {
		actual := actEs[i]
		cmpStr("entry Title", expected.Title, actual.Title, t)
		cmpStr("entry Link", expected.Link, actual.Link, t)
		cmpStr("entry Subtitle", expected.Subtitle, actual.Subtitle, t)
		cmpStr("entry Guid", expected.Guid, actual.Guid, t)
		cmpTime("entry UpdatedDate", expected.UpdatedDate, actual.UpdatedDate, t)
		cmpStr("entry Summary", expected.Summary, actual.Summary, t)
		cmpStr("entry Content", expected.Content, actual.Content, t)
		cmpStr("entry Source", expected.Source, actual.Source, t)
		cmpStr("entry Comments", expected.Comments, actual.Comments, t)
		cmpStr("entry Thumbnail", expected.Thumbnail, actual.Thumbnail, t)
		cmpStr("entry Length", expected.Length, actual.Length, t)
		cmpStr("entry Type", expected.Type, actual.Type, t)
		cmpStr("entry Url", expected.Url, actual.Url, t)
	}
}

func cmpStr(name, expected, actual string, t *testing.T) {
	if v1 != v2 {
		t.Error(fmt.Sprintf("Error with %s. Expected %q, received %q.", name, expected, actual))
	}
}

func cmpTime(name string, expected, actual time.Time, t *testing.T) {
	if expected != actual {
		t.Error(fmt.Sprintf("Error with %s. Expected %v, received %v.", name, expected, actual))
	}
}

func parseFeed(name, content string, t *testing.T) (feed *Feed, entries []*Entry) {
	parser := NewParser(name, content)
	go parser.Parse()

	entrySlice := make([]*Entry, 0, 20)
	feedOpen, entryOpen := true, true
parseLoop:
	for {
		if !feedOpen && !entryOpen {
			break parseLoop
		}
		select {
		case parsedFeed, feedOk := <-parser.feed:
			if feedOk {
				feed = &parsedFeed
				feed.Url = feedUrl
			} else {
				feedOpen = false
			}
		case parsedEntry, entryOk := <-parser.entry:
			if entryOk {
				entrySlice = append(entrySlice, &parsedEntry)
			} else {
				entryOpen = false
			}
		}
	}

	entries = entrySlice
	return
}

var suttersMillContent string = `
<rss xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:wfw="http://wellformedweb.org/CommentAPI/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:sy="http://purl.org/rss/1.0/modules/syndication/" xmlns:slash="http://purl.org/rss/1.0/modules/slash/" xmlns:georss="http://www.georss.org/georss" xmlns:geo="http://www.w3.org/2003/01/geo/wgs84_pos#" xmlns:media="http://search.yahoo.com/mrss/" version="2.0">
<channel>
<title>Sutter’s Mill</title>
<atom:link href="http://herbsutter.com/feed/" rel="self" type="application/rss+xml"/>
<link>http://herbsutter.com</link>
<description>Herb Sutter on software, hardware, and concurrency</description>
<lastBuildDate>Wed, 21 Aug 2013 00:57:11 +0000</lastBuildDate>
<language>en</language>
<sy:updatePeriod>hourly</sy:updatePeriod>
<sy:updateFrequency>1</sy:updateFrequency>
<generator>http://wordpress.com/</generator>
<image>
<url>
http://0.gravatar.com/blavatar/4554b8d24c7f200dc5e2e1b18db1893f?s=96&d=http%3A%2F%2Fs2.wp.com%2Fi%2Fbuttonw-com.png
</url>
<title>Sutter’s Mill</title>
<link>http://herbsutter.com</link>
</image>
<item>
<title>
GotW #7b: Minimizing Compile-Time Dependencies, Part 2
</title>
<link>
http://herbsutter.com/2013/08/19/gotw-7b-minimizing-compile-time-dependencies-part-2/
</link>
<comments>
http://herbsutter.com/2013/08/19/gotw-7b-minimizing-compile-time-dependencies-part-2/#comments
</comments>
<pubDate>Mon, 19 Aug 2013 10:33:12 +0000</pubDate>
<dc:creator>Herb Sutter</dc:creator>
<category>
<![CDATA[ GotW ]]>
</category>
<guid isPermaLink="false">http://herbsutter.com/?p=2294</guid>
<description>
<![CDATA[
Now that the unnecessary headers have been removed, it&#8217;s time for Phase 2: How can you limit dependencies on the internals of a class? Problem JG Questions 1. What does private mean for a class member in C++? 2. Why does changing the private members of a type cause a recompilation? Guru Question 3. Below [&#8230;]<img alt="" border="0" src="http://stats.wordpress.com/b.gif?host=herbsutter.com&#038;blog=3379246&#038;post=2294&#038;subd=herbsutter&#038;ref=&#038;feed=1" width="1" height="1" />
]]>
</description>
<content:encoded>
<![CDATA[
<p><span style="color:#5a5a5a;"><em>Now that the unnecessary headers have been removed, it&#8217;s time for Phase 2: How can you limit dependencies on the internals of a class?</em></span> </p> <h1>Problem<br /> </h1> <h2>JG Questions<br /> </h2> <p>1. What does <span style="color:#2e74b5;">private</span> mean for a class member in C++? </p> <p>2. Why does changing the private members of a type cause a recompilation? </p> <h2>Guru Question<br /> </h2> <p>3. Below is how the header from the previous Item looks after the initial cleanup pass. What further <span style="color:#2e74b5;">#include</span>s could be removed if we made some suitable changes, and how? </p> ...
]]>
</content:encoded>
<wfw:commentRss>
http://herbsutter.com/2013/08/19/gotw-7b-minimizing-compile-time-dependencies-part-2/feed/
</wfw:commentRss>
<slash:comments>8</slash:comments>
<media:content url="http://0.gravatar.com/avatar/c0ba56bfd231f8f04feb057728975181?s=96&d=identicon&r=G" medium="image">
<media:title type="html">Herb Sutter</media:title>
</media:content>
</item>
</channel>
</rss>`

var ericLippertContent = `
<?xml version="1.0" encoding="UTF-8" ?>
<?xml-stylesheet type="text/xsl" href="http://blogs.msdn.com/utility/FeedStylesheets/rss.xsl" media="screen"?><rss version="2.0" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:slash="http://purl.org/rss/1.0/modules/slash/" xmlns:wfw="http://wellformedweb.org/CommentAPI/">
<channel>
<title>Fabulous Adventures In Coding</title>
<link>http://blogs.msdn.com/b/ericlippert/</link>
<description>Eric Lippert&amp;#39;s Erstwhile Blog</description>
<dc:language>en-US</dc:language>
<item>
<title>A new fabulous adventure</title>
<link>http://blogs.msdn.com/b/ericlippert/archive/2012/11/29/a-new-fabulous-adventure.aspx</link>
<pubDate>Thu, 29 Nov 2012 15:00:00 GMT</pubDate>
<guid isPermaLink="false">91d46819-8472-40ad-a661-2c78acb4018c:10369420</guid>
<dc:creator>Eric Lippert</dc:creator>
<slash:comments>48</slash:comments>
<comments>http://blogs.msdn.com/b/ericlippert/archive/2012/11/29/a-new-fabulous-adventure.aspx#comments</comments>
<description>&lt;div class="mine"&gt;
&lt;p&gt;Tomorrow, the 30th of November, 2012, is the first day of my fifth decade here on Earth, and my last day at Microsoft. (*)&lt;/p&gt;
&lt;p&gt;(*) That timing is not coincidental.&lt;/p&gt;
</description>
</item>
</channel>
</rss>`
