package rss

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewIndexWriter() *IndexWriter {
	wr := &IndexWriter{}

	for idx, pattern := range IndexTemplatePatterns {
		tmpl, err := template.New("IndexTemplates" + string(idx)).Parse(pattern)
		if err != nil {
			fmt.Errorf("Could not parse template: %d", idx)
			panic(err)
		}
		IndexTemplates = append(IndexTemplates, tmpl)
	}
	return wr
}

type IndexWriter struct {
	data []*Feed
}

func (wr *IndexWriter) SetData(data interface{}) {
	wr.data = data.([]*Feed)
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
			</ul>
		</div>
	</body>
</html>
`,
}

var IndexTemplatePatterns = []string{
	`<a href='/feed/{{.Id}}'>{{.Title}}</a>`,
}

var IndexTemplates = make([]*template.Template, 0, len(IndexTemplatePatterns))

func (wr IndexWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *IndexWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data []*Feed) {
	var err error = nil
	fmt.Fprint(w, IndexHtml[0])
	for _, feed := range data {
		fmt.Fprint(w, IndexHtml[1])
		err = IndexTemplates[0].Execute(w, feed)
		handleIndexError(err)
	}
	fmt.Fprint(w, IndexHtml[2])
	if err != nil {
		err = nil
	}
}

func handleIndexError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
