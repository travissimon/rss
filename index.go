package main

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"net/http"
)

func NewIndexWriter(data []*Feed) (*IndexWriter) {
	wr := &IndexWriter {
		data: data,
	}
	
	return wr
}

type IndexWriter struct {
	data []*Feed
}

var IndexHtml = [...]string{
`<html>
	<head>
		<title>Feed listing</title>
	</head>
	<body>
		<style>
			 body { font-family: Helvetica, Arial, sans-serif; background: #ddd } #content { width: 80%; background:
			#fff; border-color: #333; margin: 20px; padding: 10px; -webkit-border-radius: 10px; -moz-border-radius:
			10px; border-radius: 10px; } pre, code { font-family: Menlo, monospace; font-size: 14px; } pre { line-height:
			18px; }
		</style>
		<div id="content">
			<h1>Feeds</h1>
			<ul>
				`,
				`
				<li>
					`,
					`
				</li>
				`,
				`
			</ul>
		</div>
	</body>
</html>
`,
}

func (wr IndexWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *IndexWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data []*Feed) {
	fmt.Fprint(w, IndexHtml[0])
	for _, feed := range data {
		fmt.Fprint(w, IndexHtml[1])
		fmt.Fprint(w, "<a href='/feed/", feed.Id, "'>", feed.Title, "</a>")
		fmt.Fprint(w, IndexHtml[2])
	}
}
