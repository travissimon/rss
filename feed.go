package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"net/http"
)

func NewFeedWriter(data *EntryList) (*FeedWriter) {
	wr := &FeedWriter {
		data: data,
	}
	
	return wr
}

type FeedWriter struct {
	data *EntryList
}

var FeedHtml = [...]string{
`<html>
	<head>
		<title>
			`,
			`
		</title>
	</head>
	<body>
		<style>
			 body { font-family: Helvetica, Arial, sans-serif; background: #ddd } #content { width: 80%; background:
			#fff; border-color: #333; margin: 20px; padding: 10px; -webkit-border-radius: 10px; -moz-border-radius:
			10px; border-radius: 10px; } pre, code { font-family: Menlo, monospace; font-size: 14px; } pre { line-height:
			18px; }
		</style>
		<div id="content">
			<h1>
				`,
				`
			</h1>
			<h3>
				`,
				`
			</h3>
			<div>
				`,
				`
				<div></div>
				<h3>
					`,
					`
				</h3>
				<p>
					`,
					`
				</p>
				<p>
					`,
					`
				</p>
				<p>
					`,
					`
				</p>
				`,
				`
				<p>
					`,
					`
				</p>
				`,
				`
				`,
				`
			</div>
		</div>
	</body>
</html>
`,
}

func (wr FeedWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *FeedWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data *EntryList) {
	fmt.Fprint(w, FeedHtml[0])
	fmt.Fprint(w, data.feed.Title)
	fmt.Fprint(w, FeedHtml[1])
	fmt.Fprint(w, data.feed.Title)
	fmt.Fprint(w, FeedHtml[2])
	fmt.Fprint(w, data.feed.Subtitle)
	fmt.Fprint(w, FeedHtml[3])
	for _, entry := range data.entries {
		fmt.Fprint(w, FeedHtml[4])
		fmt.Fprint(w, entry.Title)
		fmt.Fprint(w, FeedHtml[5])
		fmt.Fprint(w, entry.Subtitle)
		fmt.Fprint(w, FeedHtml[6])
		fmt.Fprint(w, entry.Summary)
		fmt.Fprint(w, FeedHtml[7])
		fmt.Fprint(w, "<a href='", entry.Link, "'>Full article</a>")
		fmt.Fprint(w, FeedHtml[8])
		if entry.Url != "" {
			fmt.Fprint(w, FeedHtml[9])
			fmt.Fprint(w, "<a href='" + entry.Url + "'>Download</a>")
			fmt.Fprint(w, FeedHtml[10])
		}
		fmt.Fprint(w, FeedHtml[11])
	}
}
