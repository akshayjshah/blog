package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

const (
	_markdownDir = "posts"
	_recipesDir  = "recipes"
	_buildDir    = "dist"
	_style       = "style.css"
	_baseURL     = "https://akshayshah.org"
)

var _indexTemplate = template.Must(template.ParseFiles("index.md"))

func main() {
	log.SetOutput(os.Stdout)

	minCSS, err := readCSS(_style)
	if err != nil {
		log.Fatalf("read CSS: %v", err)
	}

	md := newMD()
	year := time.Now().Year()
	posts, err := readPosts(_markdownDir, md, year, minCSS)
	if err != nil {
		log.Fatalf("read posts: %v", err)
	}

	index, err := readIndex(posts, md, year, minCSS)
	if err != nil {
		log.Fatalf("read index: %v", err)
	}
	posts = append(posts, index)

	root, err := os.OpenRoot(_buildDir)
	if err != nil {
		log.Fatalf("open %q: %v", _buildDir, err)
	}
	for _, post := range posts {
		if err := post.RenderTo(root); err != nil {
			log.Fatalf("render %q in %q: %v", post.Slug, _buildDir, err)
		}
	}
}

func newMD() goldmark.Markdown {
	// By default, NewTypographer does the usual smartypants conversions (with
	// LaTeX-style dashes).
	typographyOverrides := make(map[extension.TypographicPunctuation]string)
	typographyOverrides[extension.Ellipsis] = "..."       // breaks monospace
	typographyOverrides[extension.LeftAngleQuote] = "<<"  // not English
	typographyOverrides[extension.RightAngleQuote] = ">>" // not English

	return goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewTypographer(
				extension.WithTypographicSubstitutions(typographyOverrides),
			),
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromahtml.WithClasses(true), // no inline styles
					chromahtml.WithLineNumbers(true),
				),
			),
			&frontmatter.Extender{},
		),
	)
}

func readCSS(fname string) (template.CSS, error) {
	rawCSS, err := os.ReadFile(fname)
	if err != nil {
		return "", fmt.Errorf("read %q: %w", fname, err)
	}
	minifier := minify.New()
	minifier.AddFunc("text/css", css.Minify)
	minCSS, err := minifier.Bytes("text/css", rawCSS)
	if err != nil {
		return "", fmt.Errorf("minify CSS: %w", err)
	}
	return template.CSS(minCSS), nil
}

func readPosts(dir string, md goldmark.Markdown, copyright int, css template.CSS) ([]page, error) {
	root, err := os.OpenRoot(dir)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", err)
	}
	var pages []page
	err = fs.WalkDir(root.FS(), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		raw, err := root.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %q: %w", path, err)
		}
		p, err := NewPage(path, raw, md, copyright, css)
		if err != nil {
			return err
		}
		pages = append(pages, p)
		return nil
	})
	slices.SortStableFunc(pages, comparePages)
	return pages, err
}

func readIndex(posts []page, md goldmark.Markdown, copyright int, css template.CSS) (page, error) {
	type item struct {
		Created string
		Title   template.HTML
		Link    string
		Via     string
	}
	var index struct {
		Posts   []item
		Recipes []item
	}
	for _, post := range posts {
		if post.Meta.Hidden {
			continue
		}
		link := post.Meta.Link
		if link == "" {
			link = fmt.Sprintf("/%s", post.Slug)
		}
		item := item{
			Created: post.Created(),
			Title:   post.Title,
			Link:    link,
			Via:     post.Meta.Via,
		}
		switch filepath.Dir(post.Slug) {
		case _recipesDir:
			index.Recipes = append(index.Recipes, item)
		default:
			index.Posts = append(index.Posts, item)
		}
	}
	var expanded bytes.Buffer
	if err := _indexTemplate.Execute(&expanded, index); err != nil {
		return page{}, fmt.Errorf("generate index Markdown: %w", err)
	}
	return NewPage("index.md", expanded.Bytes(), md, copyright, css)
}
