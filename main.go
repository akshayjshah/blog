package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

const (
	_humanDate   = "Jan 2006" // preferred date format for display
	_machineDate = time.DateOnly
	_configFile  = "config.yaml"
	_markdownDir = "posts"
	_recipesDir  = "recipes"
	_buildDir    = "dist"
)

var (
	_post = template.Must(template.ParseFiles("page.html"))
	_home = template.Must(template.ParseFiles("index.md"))
)

type config struct {
	BaseURL     string       `yaml:"base_url"`
	Author      string       `yaml:"author"`
	Description string       `yaml:"description"`
	CSS         string       `yaml:"css"`
	Pages       []pageConfig `yaml:"pages"`
}

func (c *config) Sort() {
	slices.SortStableFunc(c.Pages, func(l, r pageConfig) bool {
		// Reverse chronological sort.
		return l.Date() > r.Date()
	})
}

type pageConfig struct {
	Created string `yaml:"created"`
	Hidden  bool   `yaml:"hidden"`

	// for normal posts
	Markdown    string `yaml:"markdown"`
	Description string `yaml:"description"`
	Updated     string `yaml:"updated"`
	HideHome    bool   `yaml:"hide_home"`
	HideLicense bool   `yaml:"hide_license"`

	// for external posts
	Link  string `yaml:"link"`
	Title string `yaml:"title"`
	Via   string `yaml:"via"`
}

func (p pageConfig) Date() string {
	if p.Updated != "" {
		return p.Updated
	}
	return p.Created
}

func (p pageConfig) Validate() error {
	if p.Markdown == "" && p.Link == "" {
		return errors.New("require either markdown file or external link")
	}
	if p.Markdown != "" && p.Link != "" {
		return errors.New("can't have both markdown file and external link")
	}
	if p.Link != "" && p.Via == "" {
		return errors.New("external links must have via text")
	}
	if !p.Hidden && p.Created == "" {
		return errors.New("require created date")
	}
	if p.Updated != "" && p.Updated <= p.Created {
		return errors.New("updated date is before created date")
	}
	return nil
}

type page struct {
	Author      string
	Description string
	Created     string
	Updated     string
	Permalink   string
	Markdown    []byte
	HideHome    bool
	HideLicense bool
	Out         string
}

func (p page) Render(endCopyright int, css template.CSS) (template.HTML, error) {
	data := struct {
		page

		Created      string
		Updated      string
		Title        template.HTML
		TitlePlain   string
		Content      template.HTML
		EndCopyright int
		CSS          template.CSS
	}{
		page:         p,
		EndCopyright: endCopyright,
		CSS:          css,
	}

	first, err := head(bytes.NewReader(p.Markdown))
	if err != nil {
		return "", fmt.Errorf("read first line of in-memory Markdown: %w", err)
	}
	data.TitlePlain = strings.TrimSpace(strings.TrimPrefix(first, "#"))
	data.Title = renderHTML([]byte(first))

	created, err := machineToHuman(p.Created)
	if err != nil {
		return "", fmt.Errorf("format created date: %w", err)
	}
	data.Created = created
	updated, err := machineToHuman(p.Updated)
	if err != nil {
		return "", fmt.Errorf("format updated date: %w", err)
	}
	data.Updated = updated

	content := bytes.TrimPrefix(p.Markdown, []byte(first))
	data.Content = renderHTML(content)

	if err := os.MkdirAll(filepath.Dir(p.Out), 0750); err != nil && !os.IsExist(err) {
		return "", fmt.Errorf("create directories for %q: %v", p.Out, err)
	}
	w, err := os.Create(p.Out)
	if err != nil {
		return "", fmt.Errorf("create %q: %w", p.Out, err)
	}
	defer w.Close()
	log.Printf("Rendering %s", p.Out)
	return template.HTML(data.TitlePlain), _post.Execute(w, data)
}

func renderHTML(markdown []byte) template.HTML {
	raw := blackfriday.Run(
		markdown,
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs),
	)
	return template.HTML(raw)
}

func head(r io.Reader) (string, error) {
	scanner := bufio.NewScanner(r)
	var first string
	for scanner.Scan() {
		first = scanner.Text()
		break
	}
	return first, scanner.Err()
}

func machineToHuman(date string) (string, error) {
	if date == "" {
		return "", nil
	}
	d, err := time.Parse(_machineDate, date)
	if err != nil {
		return "", fmt.Errorf("parse date %q: %w", date, err)
	}
	return d.Format(_humanDate), nil
}

func homepage(cfg config, titles []template.HTML) (page, error) {
	type item struct {
		Published string
		Title     template.HTML
		Link      string
		Via       string
	}
	var index struct {
		Posts   []item
		Recipes []item
	}
	for i, pc := range cfg.Pages {
		if pc.Hidden {
			continue
		}
		published, err := machineToHuman(pc.Created)
		if err != nil {
			return page{}, fmt.Errorf("format created date for %q", pc.Markdown)
		}
		link := pc.Link
		if link == "" {
			link = fmt.Sprintf("/%s", strings.TrimSuffix(pc.Markdown, ".md"))
		}
		item := item{
			Published: published,
			Title:     titles[i],
			Link:      link,
			Via:       pc.Via,
		}
		switch filepath.Dir(pc.Markdown) {
		case _recipesDir:
			index.Recipes = append(index.Recipes, item)
		default:
			index.Posts = append(index.Posts, item)
		}
	}
	var md bytes.Buffer
	if err := _home.Execute(&md, index); err != nil {
		return page{}, fmt.Errorf("generate homepage Markdown: %w", err)
	}
	return page{
		Author:      cfg.Author,
		Description: cfg.Description,
		Permalink:   cfg.BaseURL,
		Markdown:    md.Bytes(),
		Out:         filepath.Join(_buildDir, "index.html"),
	}, nil
}

func main() {
	log.SetOutput(os.Stdout)
	year := time.Now().Year()
	configYAML, err := os.ReadFile(_configFile)
	if err != nil {
		log.Fatalf("read %q: %v", _configFile, err)
	}
	var cfg config
	if err := yaml.Unmarshal(configYAML, &cfg); err != nil {
		log.Fatalf("unmarshal config: %v", err)
	}
	cfg.Sort()

	rawCSS, err := os.ReadFile(cfg.CSS)
	if err != nil {
		log.Fatalf("read %q: %v", cfg.CSS, err)
	}
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	styles, err := minifier.Bytes("text/css", rawCSS)
	if err != nil {
		log.Fatalf("minify CSS: %v", err)
	}

	titles := make([]template.HTML, len(cfg.Pages))
	for i, pc := range cfg.Pages {
		if err := pc.Validate(); err != nil {
			log.Fatalf("invalid post config %v: %v", pc, err)
		}
		if pc.Markdown == "" {
			titles[i] = template.HTML(pc.Title)
			continue // external link, no page to render
		}
		slug := path.Clean(strings.TrimSuffix(pc.Markdown, ".md"))
		md, err := os.ReadFile(filepath.Join(_markdownDir, pc.Markdown))
		if err != nil {
			log.Fatalf("read %q: %v", pc.Markdown, err)
		}
		page := page{
			Author:      cfg.Author,
			Description: pc.Description,
			Created:     pc.Created,
			Updated:     pc.Updated,
			Permalink:   cfg.BaseURL + "/" + slug,
			Markdown:    md,
			Out:         filepath.Join(_buildDir, slug+".html"),
			HideHome:    pc.HideHome,
			HideLicense: pc.HideLicense,
		}
		title, err := page.Render(year, template.CSS(styles))
		if err != nil {
			log.Fatalf("render %s: %v", slug, err)
		}
		titles[i] = title
	}

	home, err := homepage(cfg, titles)
	if err != nil {
		log.Fatalf("render homepage Markdown: %v", err)
	}
	if _, err := home.Render(year, template.CSS(styles)); err != nil {
		log.Fatalf("render homepage to HTML: %v", err)
	}
}
