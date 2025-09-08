package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

var _pageTemplate = template.Must(template.ParseFiles("page.html"))

type page struct {
	Slug      string
	Meta      metadata
	Title     template.HTML
	Content   template.HTML
	Copyright int
	CSS       template.CSS
}

func NewPage(path string, raw []byte, md goldmark.Markdown, copyright int, css template.CSS) (page, error) {
	ctx := parser.NewContext()
	var rendered bytes.Buffer
	if err := md.Convert(raw, &rendered, parser.WithContext(ctx)); err != nil {
		return page{}, fmt.Errorf("render %q: %w", path, err)
	}
	var meta metadata
	if err := frontmatter.Get(ctx).Decode(&meta); err != nil {
		return page{}, fmt.Errorf("read frontmatter in %q: %w", path, err)
	}
	if err := meta.Validate(); err != nil {
		return page{}, fmt.Errorf("invalid frontmatter in %q: %w", path, err)
	}

	// We want a version of the title without any HTML tags, but with punctuation
	// and special characters converted to HTML entities.
	var titleParagraph bytes.Buffer
	if err := md.Convert([]byte(meta.Title), &titleParagraph); err != nil {
		return page{}, fmt.Errorf("render title in %q: %w", path, err)
	}
	title := strings.TrimSpace(titleParagraph.String())
	title = strings.TrimPrefix(title, "<p>")
	title = strings.TrimSuffix(title, "</p>")
	title = strings.TrimSpace(title)

	return page{
		Slug:      filepath.Clean(strings.TrimSuffix(path, ".md")),
		Meta:      meta,
		Title:     template.HTML(title),
		Content:   template.HTML(rendered.String()),
		Copyright: copyright,
		CSS:       css,
	}, nil
}

func (p page) RenderTo(root *os.Root) error {
	if dir := filepath.Dir(p.Slug); dir != "" {
		if err := root.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("create directory %q: %w", dir, err)
		}
	}
	path := p.Slug + ".html"
	w, err := root.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer w.Close()
	if err := _pageTemplate.Execute(w, p); err != nil {
		return fmt.Errorf("render %q: %w", path, err)
	}
	return nil
}

func (p page) Permalink() string {
	u := _baseURL + "/" + strings.TrimSuffix(p.Slug, "index")
	return strings.TrimSuffix(u, "/")
}

func (p page) Created() string {
	return p.Meta.Created()
}

func (p page) Updated() string {
	return p.Meta.Updated()
}

func (p page) LastModified() string {
	// For sitemap.xml.
	if p.Meta.RawUpdated != "" {
		return p.Meta.RawUpdated
	}
	return p.Meta.RawCreated
}

func (p page) TitlePlainText() string {
	return p.Meta.Title
}

func (p page) Description() string {
	return p.Meta.Description
}

func (p page) HideHome() bool {
	return p.Meta.HideHome
}

func (p page) HideLicense() bool {
	return p.Meta.HideLicense
}

func (p page) HideDates() bool {
	return p.Meta.HideDates
}

func (p page) Priority() float64 {
	pri := p.Meta.Priority
	if pri == 0 {
		return 0.5
	}
	return pri
}

func (p page) IsExternalLink() bool {
	return p.Meta.Via != ""
}

func comparePages(l, r page) int {
	return l.Meta.Compare(r.Meta)
}
