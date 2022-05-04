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
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
)

const _date = "Jan 2006" // preferred date format for display

var (
	_minifier    = minify.New()
	_post        = template.Must(template.ParseFiles("page.html"))
	_home        = template.Must(template.ParseFiles("index.md"))
	_index       = flag.Bool("index", false, "write Markdown index instead of single post")
	_hideDates   = flag.Bool("nodates", false, "hide publish dates")
	_hideHome    = flag.Bool("nohome", false, "hide homepage link")
	_hideLicense = flag.Bool("nolicense", false, "hide license link")
	_recipes     = flag.String("recipes", "", "recipes directory (with -index)")
	_style       = flag.String("style", "", "CSS file")
)

func init() {
	_minifier.AddFunc("text/html", html.Minify)
	_minifier.AddFunc("text/css", css.Minify)
}

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

type Index struct {
	Posts   []IndexEntry
	Recipes []IndexEntry
}

func (i *Index) Sort() {
	for _, s := range [][]IndexEntry{i.Posts, i.Recipes} {
		sort.Slice(s, func(i, j int) bool {
			// Reverse chronological sort.
			return s[j].Created.Before(s[i].Created)
		})
	}
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

func render(md []byte) template.HTML {
	raw := blackfriday.Run(
		md,
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs),
	)
	min, err := _minifier.Bytes("text/html", raw)
	must(err, "minify HTML")
	return template.HTML(min)
}

func style() template.CSS {
	css, err := ioutil.ReadFile(*_style)
	must(err, "read CSS from %q", *_style)
	mincss, err := _minifier.Bytes("text/css", css)
	must(err, "minify css")
	return template.CSS(mincss)
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
	permalink := site.BaseURL + "/" + path.Clean(strings.TrimSuffix(post, ".md")) + "/"
	if strings.HasSuffix(permalink, "/index/") {
		permalink = site.BaseURL + "/"
	}
	p := Page{
		Site:        site,
		TitlePlain:  strings.TrimSpace(strings.TrimPrefix(first, "#")),
		Title:       render([]byte(first)),
		Permalink:   permalink,
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
	p.Content = render(content)
	must(_post.Execute(w, p), "format HTML for %q", post)
}

func homepage(w io.Writer, files []string) {
	var (
		posts   []IndexEntry
		recipes []IndexEntry
	)
	for _, f := range files {
		title := strings.TrimSpace(strings.TrimPrefix(head(f), "#"))
		c := created(f)
		entry := IndexEntry{
			Created:   c,
			Published: c.Format(_date),
			Title:     title,
			Link:      fmt.Sprintf("/%s/", strings.TrimSuffix(f, ".md")),
		}
		switch filepath.Dir(f) {
		case *_recipes:
			recipes = append(recipes, entry)
		default:
			posts = append(posts, entry)
		}
	}
	idx := Index{
		Posts:   posts,
		Recipes: recipes,
	}
	idx.Sort()
	must(_home.Execute(w, idx), "generate Markdown homepage")
}

func main() {
	flag.Parse()
	site := Site{
		BaseURL:     "https://akshayshah.org",
		Title:       "Akshay Shah",
		Author:      "Akshay Shah",
		Description: "Code, cooking, and caffeine.",
		LastChanged: fmt.Sprint(time.Now().Year()),
		CSS:         style(),
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
