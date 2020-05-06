package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"time"

	"gitlab.com/golang-commonmark/markdown"
)

const _date = "Jan 2006" // preferred date format for display

var (
	_md          = markdown.New(markdown.HTML(true)) // allow raw HTML in Markdown
	_post        = template.Must(template.ParseFiles("page.html"))
	_home        = template.Must(template.ParseFiles("index.md"))
	_index       = flag.Bool("index", false, "write Markdown index instead of single post")
	_hideDates   = flag.Bool("nodates", false, "hide publish dates")
	_hideHome    = flag.Bool("nohome", false, "hide homepage link")
	_hideLicense = flag.Bool("nolicense", false, "hide license link")
	_style       = flag.String("style", "", "CSS file")
)

type Site struct {
	BaseURL                string
	Title                  string
	Author                 string
	Description            string
	GoogleSiteVerification string
	LastChanged            string
	CSS                    template.CSS
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

type IndexEntry struct {
	Created   time.Time
	Published string
	Title     string
	Link      string
}

func must(err error, prefix string, args ...interface{}) {
	if err != nil {
		msg := fmt.Sprintf(prefix, args...)
		log.Fatalf("%s: %v", msg, err)
	}
}

func head(path string) string {
	f, err := os.Open(path)
	must(err, "open %q", path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var first string
	for scanner.Scan() {
		first = scanner.Text()
		break
	}
	must(scanner.Err(), "scan %q", path)
	return first
}

func created(path string) time.Time {
	git := exec.Command("git", "log", "--diff-filter=A", "--follow", "--format=%aI", path)
	tail := exec.Command("tail", "-1")

	r, w := io.Pipe()
	git.Stdout = w
	tail.Stdin = r

	var out bytes.Buffer
	tail.Stdout = &out

	must(git.Start(), "start git")
	must(tail.Start(), "start tail")
	must(git.Wait(), "wait git")
	must(w.Close(), "close pipe writer")
	must(tail.Wait(), "wait tail")
	c, err := time.Parse(time.RFC3339, strings.TrimSpace(out.String()))
	if err != nil {
		// Assume the file is uncommitted, default to now.
		return time.Now()
	}
	return c
}

func updated(path string) *time.Time {
	git := exec.Command("git", "log", "--follow", "-1", "--format=%aI", path)
	var out bytes.Buffer
	git.Stdout = &out
	must(git.Run(), "git")
	u, err := time.Parse(time.RFC3339, strings.TrimSpace(out.String()))
	if err != nil {
		// Assume file hasn't been updated.
		return nil
	}
	return &u
}

func post(w io.Writer, site Site, post string) {
	first := head(post)
	p := Page{
		Site:        site,
		TitlePlain:  strings.TrimSpace(strings.TrimPrefix(first, "#")),
		Title:       template.HTML(_md.RenderToString([]byte(first))),
		Permalink:   site.BaseURL + "/" + path.Clean(strings.TrimSuffix(post, ".md")) + "/",
		HideDates:   *_hideDates,
		HideHome:    *_hideHome,
		HideLicense: *_hideLicense,
	}
	if !p.HideDates {
		p.Created = created(post).Format(_date)
		if u := updated(post); u != nil {
			p.Updated = u.Format(_date)
		}
	}
	content, err := ioutil.ReadFile(post)
	must(err, "read Markdown file")
	content = bytes.TrimPrefix(content, []byte(first))
	p.Content = template.HTML(_md.RenderToString(content))
	must(_post.Execute(w, p), "format HTML for %q", post)
}

func homepage(w io.Writer, posts []string) {
	entries := make([]IndexEntry, len(posts))
	for i, p := range posts {
		title := strings.TrimSpace(strings.TrimPrefix(head(p), "#"))
		c := created(p)
		entries[i] = IndexEntry{
			Created:   c,
			Published: c.Format(_date),
			Title:     title,
			Link:      fmt.Sprintf("/%s/", strings.TrimSuffix(p, ".md")),
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		// Reverse chronological sort.
		return entries[j].Created.Before(entries[i].Created)
	})
	must(_home.Execute(w, entries), "generate Markdown homepage")
}

func main() {
	flag.Parse()
	css, err := ioutil.ReadFile(*_style)
	must(err, "read CSS from %q", *_style)
	site := Site{
		BaseURL:                "http://www.akshayshah.org",
		Title:                  "Akshay Shah",
		Author:                 "Akshay Shah",
		Description:            "Thoughts on code and human factors from a physician-turned-engineer.",
		GoogleSiteVerification: "TODO",
		LastChanged:            fmt.Sprint(time.Now().Year()),
		CSS:                    template.CSS(css),
	}
	var paths []string
	for _, p := range flag.Args() {
		if p != "" {
			paths = append(paths, p)
		}
	}
	if *_index {
		homepage(os.Stdout, paths)
		return
	}
	if len(paths) > 1 {
		log.Fatalf("one post at a time: got %v", paths)
	}
	post(os.Stdout, site, paths[0])
}
