package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/markdown"
)

const _date = "Jan 2006" // preferred date format for display

var (
	_md          = markdown.New(markdown.HTML(true)) // allow raw HTML in Markdown
	_post        = template.Must(template.ParseFiles("page.html"))
	_created     = flag.String("created", "", "date created")
	_updated     = flag.String("updated", "", "date updated")
	_title       = flag.String("title", "", "title (as Markdown H1)") // e.g., "# A Cool Title"
	_hideDates   = flag.Bool("nodates", false, "hide publish dates")
	_hideHome    = flag.Bool("nohome", false, "hide homepage link")
	_hideLicense = flag.Bool("nolicense", false, "hide license link")
)

type Site struct {
	BaseURL                string
	Title                  string
	Author                 string
	Description            string
	GoogleSiteVerification string
	LastChanged            string
}

type Page struct {
	Site        Site
	TitlePlain  string
	Title       template.HTML
	Created     string
	Updated     string
	Permalink   string
	HideDates   bool
	HideHome    bool
	HideLicense bool
	Content     template.HTML
}

func post(w io.Writer, site Site, post string) {
	p := Page{
		Site:        site,
		TitlePlain:  strings.TrimSpace(strings.TrimPrefix(*_title, "#")),
		Title:       template.HTML(_md.RenderToString([]byte(*_title))),
		Permalink:   site.BaseURL + "/" + path.Clean(strings.TrimSuffix(post, ".md")) + "/",
		HideDates:   *_hideDates,
		HideHome:    *_hideHome,
		HideLicense: *_hideLicense,
	}
	if !p.HideDates {
		// If we can't parse these flags, assume the file is uncommitted: it was
		// created now and hasn't been updated after publishing.
		created, err := time.Parse(time.RFC3339, *_created)
		if err != nil {
			created = time.Now()
		}
		p.Created = created.Format(_date)
		if updated, err := time.Parse(time.RFC3339, *_updated); err == nil {
			p.Updated = updated.Format(_date)
		}
	}
	content, err := ioutil.ReadFile(post)
	if err != nil {
		log.Fatalf("read Markdown file: %v", err)
	}
	content = bytes.TrimPrefix(content, []byte(*_title))
	p.Content = template.HTML(_md.RenderToString(content))
	if err := _post.Execute(w, p); err != nil {
		log.Fatalf("format HTML for %q: %v", post, err)
	}
}

func main() {
	site := Site{
		BaseURL:                "http://www.akshayshah.org",
		Title:                  "Akshay Shah",
		Author:                 "Akshay Shah",
		Description:            "Thoughts on code and human factors from a physician-turned-engineer.",
		GoogleSiteVerification: "TODO",
		LastChanged:            fmt.Sprint(time.Now().Year()),
	}
	flag.Parse()
	var paths []string
	for _, p := range flag.Args() {
		if p != "" {
			paths = append(paths, p)
		}
	}
	if len(paths) > 1 {
		log.Fatalf("one post at a time: got %v", paths)
	}
	post(os.Stdout, site, paths[0])
}
