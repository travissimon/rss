package rss

// THIS IS A GENERATED FILE, EDITS WILL BE OVERWRITTEN
// EDIT THE .haml FILE INSTEAD

import (
	"fmt"
	"html/template"
	"net/http"
)

func NewFeedWriter() (*FeedWriter) {
	wr := &FeedWriter{}
	
	for idx, pattern := range FeedTemplatePatterns {
		tmpl, err := template.New("FeedTemplates" + string(idx)).Parse(pattern)
		if err != nil {
			fmt.Errorf("Could not parse template: %d", idx)
			panic(err)
		}
		FeedTemplates = append(FeedTemplates, tmpl)
	}
	return wr
}

type FeedWriter struct {
	data *EntryList
}

func (wr *FeedWriter) SetData(data interface{}) {
	wr.data = data.(*EntryList)
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
				<p>
					`,
					`
				</p>
				`,
				`
			</div>
		</div>
	</body>
</html>
`,
}

var FeedTemplatePatterns = []string{
	`{{.feed.Title}}`,
	`{{.feed.Title}}`,
	`{{.feed.Subtitle}}`,
	`{{.Title}}`,
	`{{.Subtitle}}`,
	`{{.Summary}}`,
	`<a href='{{.Link}}'>Full article</a>`,
	`<a href='{{.Url}}'>Download</a>`,
}

var FeedTemplates = make([]*template.Template, 0, len(FeedTemplatePatterns))

func (wr FeedWriter) Execute(w http.ResponseWriter, r *http.Request) {
	wr.ExecuteData(w, r, wr.data)
}

func (wr *FeedWriter) ExecuteData(w http.ResponseWriter, r *http.Request, data *EntryList) {
	var err error = nil
	fmt.Fprint(w, FeedHtml[0])
	err = FeedTemplates[0].Execute(w, data)
	handleFeedError(err)
	fmt.Fprint(w, FeedHtml[1])
	err = FeedTemplates[1].Execute(w, data)
	handleFeedError(err)
	fmt.Fprint(w, FeedHtml[2])
	err = FeedTemplates[2].Execute(w, data)
	handleFeedError(err)
	fmt.Fprint(w, FeedHtml[3])
	for _, entry := range data.entries {
		fmt.Fprint(w, FeedHtml[4])
		err = FeedTemplates[3].Execute(w, entry)
		handleFeedError(err)
		fmt.Fprint(w, FeedHtml[5])
		err = FeedTemplates[4].Execute(w, entry)
		handleFeedError(err)
		fmt.Fprint(w, FeedHtml[6])
		err = FeedTemplates[5].Execute(w, entry)
		handleFeedError(err)
		fmt.Fprint(w, FeedHtml[7])
		err = FeedTemplates[6].Execute(w, data)
		handleFeedError(err)
		if entry.Url != "" {
			fmt.Fprint(w, FeedHtml[8])
			err = FeedTemplates[7].Execute(w, data)
			handleFeedError(err)
		}
		fmt.Fprint(w, FeedHtml[9])
	}
	fmt.Fprint(w, FeedHtml[10])
if err != nil {err = nil}}

func handleFeedError(err error) {
	if err != nil {fmt.Println(err)}}