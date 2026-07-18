package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

//go:embed web/templates/*.html web/static/*
var webFiles embed.FS

var (
	headingPattern = regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`)
	linkPattern    = regexp.MustCompile(`!?\[([^\]]+)\]\([^)]+\)`)
	markupPattern  = regexp.MustCompile("[`*_~\\[\\]<>]")
)

type site struct {
	articles   []*article
	bySlug     map[string]*article
	sections   []navSection
	home       *template.Template
	article    *template.Template
	markdown   goldmark.Markdown
	docsFS     fs.FS
	pathToSlug map[string]string
}

type article struct {
	Slug        string
	SourcePath  string
	Title       string
	Description string
	Section     string
	SectionName string
	Trail       string
	ReadMinutes int
	WordCount   int
	SearchText  string
	Content     template.HTML
	TOC         []tocItem
	Previous    *navItem
	Next        *navItem
}

type tocItem struct {
	Level int
	Title string
	ID    string
}

type navSection struct {
	Slug  string
	Name  string
	Count int
	Docs  []navItem
}

type navItem struct {
	Slug    string
	Title   string
	Trail   string
	Section string
}

type viewData struct {
	PageTitle    string
	Description  string
	IsHome       bool
	ActiveSlug   string
	Article      *article
	Articles     []*article
	Sections     []navSection
	ArticleCount int
	WordCount    int
}

type searchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Section     string `json:"section"`
	URL         string `json:"url"`
	Score       int    `json:"-"`
}

func newSite(docsFS fs.FS) (*site, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithXHTML()),
	)

	funcs := template.FuncMap{
		"initial": func(value string) string {
			for _, r := range value {
				return strings.ToUpper(string(r))
			}
			return "O"
		},
	}
	homeTemplate, err := template.New("base.html").Funcs(funcs).ParseFS(webFiles, "web/templates/base.html", "web/templates/home.html")
	if err != nil {
		return nil, fmt.Errorf("parse home template: %w", err)
	}
	articleTemplate, err := template.New("base.html").Funcs(funcs).ParseFS(webFiles, "web/templates/base.html", "web/templates/article.html")
	if err != nil {
		return nil, fmt.Errorf("parse article template: %w", err)
	}

	s := &site{
		bySlug:     make(map[string]*article),
		home:       homeTemplate,
		article:    articleTemplate,
		markdown:   md,
		docsFS:     docsFS,
		pathToSlug: make(map[string]string),
	}
	if err := s.loadArticles(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *site) loadArticles() error {
	var paths []string
	err := fs.WalkDir(s.docsFS, ".", func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() && strings.EqualFold(path.Ext(filePath), ".md") {
			paths = append(paths, filePath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk Markdown files: %w", err)
	}
	if len(paths) == 0 {
		return errors.New("no Markdown files found")
	}
	sort.Strings(paths)

	for _, filePath := range paths {
		slug := strings.TrimSuffix(filePath, path.Ext(filePath))
		if strings.EqualFold(path.Base(slug), "README") {
			slug = path.Dir(slug)
		}
		slug = strings.TrimPrefix(slug, "./")
		s.pathToSlug[filePath] = slug
	}

	for _, filePath := range paths {
		source, err := fs.ReadFile(s.docsFS, filePath)
		if err != nil {
			return fmt.Errorf("read %s: %w", filePath, err)
		}
		slug := s.pathToSlug[filePath]
		parts := strings.Split(slug, "/")
		section := parts[0]
		trail := ""
		if len(parts) > 2 {
			trail = strings.Join(parts[1:len(parts)-1], " / ")
		} else if len(parts) == 2 && strings.EqualFold(path.Base(filePath), "README.md") {
			trail = strings.Join(parts[1:], " / ")
		}
		title := extractTitle(source, filePath)
		description := extractDescription(source, title)
		words := len(strings.Fields(string(source)))
		a := &article{
			Slug:        slug,
			SourcePath:  filePath,
			Title:       title,
			Description: description,
			Section:     section,
			SectionName: displayName(section),
			Trail:       displayTrail(trail),
			WordCount:   words,
			ReadMinutes: maxInt(1, (words+219)/220),
			SearchText:  strings.ToLower(title + " " + description + " " + string(source)),
		}
		s.articles = append(s.articles, a)
		s.bySlug[slug] = a
	}

	for i, a := range s.articles {
		content, toc, err := s.renderMarkdown(a.SourcePath)
		if err != nil {
			return fmt.Errorf("render %s: %w", a.SourcePath, err)
		}
		a.Content = template.HTML(content) // Goldmark escapes raw HTML by default.
		a.TOC = toc
		if i > 0 {
			a.Previous = itemFor(s.articles[i-1])
		}
		if i+1 < len(s.articles) {
			a.Next = itemFor(s.articles[i+1])
		}
	}
	s.sections = buildSections(s.articles)
	return nil
}

func (s *site) renderMarkdown(filePath string) (string, []tocItem, error) {
	source, err := fs.ReadFile(s.docsFS, filePath)
	if err != nil {
		return "", nil, err
	}
	doc := s.markdown.Parser().Parse(text.NewReader(source))
	// The page header already presents the document title. Remove only the first
	// H1 from the rendered tree so the source remains untouched and the title is
	// not repeated in the article body.
	var titleHeading ast.Node
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if heading, ok := node.(*ast.Heading); ok && heading.Level == 1 {
				titleHeading = node
				return ast.WalkStop, nil
			}
		}
		return ast.WalkContinue, nil
	})
	if titleHeading != nil {
		parent := titleHeading.Parent()
		parent.RemoveChild(parent, titleHeading)
	}
	ids := make(map[string]int)
	var toc []tocItem
	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := node.(type) {
		case *ast.Heading:
			title := strings.TrimSpace(string(n.Text(source)))
			id := uniqueID(slugify(title), ids)
			n.SetAttributeString("id", []byte(id))
			if n.Level >= 2 && n.Level <= 3 && len(toc) < 24 {
				toc = append(toc, tocItem{Level: n.Level, Title: title, ID: id})
			}
		case *ast.Link:
			destination := string(n.Destination)
			if replacement := s.localLink(filePath, destination); replacement != "" {
				n.Destination = []byte(replacement)
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return "", nil, err
	}

	var output bytes.Buffer
	if err := s.markdown.Renderer().Render(&output, source, doc); err != nil {
		return "", nil, err
	}
	return output.String(), toc, nil
}

func (s *site) localLink(from, destination string) string {
	if destination == "" || strings.HasPrefix(destination, "#") || strings.Contains(destination, "://") || strings.HasPrefix(destination, "mailto:") {
		return ""
	}
	filePart, fragment, _ := strings.Cut(destination, "#")
	if !strings.EqualFold(path.Ext(filePart), ".md") {
		return ""
	}
	resolved := path.Clean(path.Join(path.Dir(from), filePart))
	slug, ok := s.pathToSlug[resolved]
	if !ok {
		return ""
	}
	if fragment != "" {
		fragment = "#" + fragment
	}
	return "/docs/" + slug + fragment
}

func (s *site) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleHome)
	mux.HandleFunc("GET /docs/{slug...}", s.handleArticle)
	mux.HandleFunc("GET /api/search", s.handleSearch)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	assets, _ := fs.Sub(webFiles, "web/static")
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", cacheAssets(http.FileServer(http.FS(assets)))))
	return securityHeaders(mux)
}

func (s *site) handleHome(w http.ResponseWriter, _ *http.Request) {
	data := s.baseData()
	data.PageTitle = "Omega Learn — a growing technical field guide"
	data.Description = "Explore practical notes on data structures, Kubernetes, SRE, and more."
	data.IsHome = true
	data.Articles = s.articles
	if err := s.home.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "render page", http.StatusInternalServerError)
	}
}

func (s *site) handleArticle(w http.ResponseWriter, r *http.Request) {
	slug := strings.Trim(strings.TrimSpace(r.PathValue("slug")), "/")
	a, ok := s.bySlug[slug]
	if !ok {
		http.NotFound(w, r)
		return
	}
	data := s.baseData()
	data.PageTitle = a.Title + " — Omega Learn"
	data.Description = a.Description
	data.ActiveSlug = a.Slug
	data.Article = a
	if err := s.article.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, "render page", http.StatusInternalServerError)
	}
}

func (s *site) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q")))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if query == "" {
		_ = json.NewEncoder(w).Encode([]searchResult{})
		return
	}
	terms := strings.Fields(query)
	results := make([]searchResult, 0, 12)
	for _, a := range s.articles {
		score := 0
		title := strings.ToLower(a.Title)
		for _, term := range terms {
			if !strings.Contains(a.SearchText, term) {
				score = 0
				break
			}
			score += strings.Count(a.SearchText, term)
			if strings.Contains(title, term) {
				score += 20
			}
			if strings.Contains(a.Slug, term) {
				score += 8
			}
		}
		if score > 0 {
			results = append(results, searchResult{Title: a.Title, Description: a.Description, Section: a.SectionName, URL: "/docs/" + a.Slug, Score: score})
		}
	}
	sort.SliceStable(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	if len(results) > 12 {
		results = results[:12]
	}
	_ = json.NewEncoder(w).Encode(results)
}

func (s *site) baseData() viewData {
	data := viewData{Sections: s.sections, ArticleCount: len(s.articles)}
	for _, a := range s.articles {
		data.WordCount += a.WordCount
	}
	return data
}

func buildSections(articles []*article) []navSection {
	order := make([]string, 0)
	groups := make(map[string]*navSection)
	for _, a := range articles {
		group, ok := groups[a.Section]
		if !ok {
			group = &navSection{Slug: a.Section, Name: a.SectionName}
			groups[a.Section] = group
			order = append(order, a.Section)
		}
		group.Docs = append(group.Docs, *itemFor(a))
		group.Count++
	}
	sections := make([]navSection, 0, len(order))
	for _, key := range order {
		sections = append(sections, *groups[key])
	}
	return sections
}

func itemFor(a *article) *navItem {
	return &navItem{Slug: a.Slug, Title: a.Title, Trail: a.Trail, Section: a.SectionName}
}

func extractTitle(source []byte, filePath string) string {
	if match := headingPattern.FindSubmatch(source); len(match) == 2 {
		return cleanInline(string(match[1]))
	}
	return displayName(strings.TrimSuffix(path.Base(filePath), path.Ext(filePath)))
}

func extractDescription(source []byte, title string) string {
	lines := strings.Split(string(source), "\n")
	var paragraph []string
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```") || strings.HasPrefix(line, "|") || line == "---" {
			if len(paragraph) > 0 {
				break
			}
			continue
		}
		paragraph = append(paragraph, strings.TrimLeft(line, ">-+* "))
	}
	description := cleanInline(strings.Join(paragraph, " "))
	if description == "" || strings.EqualFold(description, title) {
		return "A practical guide in the Omega learning library."
	}
	runes := []rune(description)
	if len(runes) > 180 {
		description = strings.TrimSpace(string(runes[:177])) + "…"
	}
	return description
}

func cleanInline(value string) string {
	value = linkPattern.ReplaceAllString(value, "$1")
	value = markupPattern.ReplaceAllString(value, "")
	value = strings.ReplaceAll(value, "&amp;", "&")
	return strings.TrimSpace(value)
}

func displayName(value string) string {
	value = strings.ReplaceAll(value, "-", " ")
	value = strings.ReplaceAll(value, "_", " ")
	words := strings.Fields(value)
	for i, word := range words {
		switch strings.ToLower(word) {
		case "sre", "k8s", "api", "dsa":
			words[i] = strings.ToUpper(word)
		default:
			runes := []rune(word)
			if len(runes) > 0 {
				runes[0] = unicode.ToUpper(runes[0])
				words[i] = string(runes)
			}
		}
	}
	return strings.Join(words, " ")
}

func displayTrail(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, " / ")
	for i := range parts {
		parts[i] = displayName(parts[i])
	}
	return strings.Join(parts, " / ")
}

func slugify(value string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
		} else if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func uniqueID(base string, seen map[string]int) string {
	if base == "" {
		base = "section"
	}
	seen[base]++
	if seen[base] == 1 {
		return base
	}
	return base + "-" + strconv.Itoa(seen[base])
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'; img-src 'self' data:; base-uri 'none'; frame-ancestors 'none'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

func cacheAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Revalidate assets so restarting the local server immediately exposes UI
		// changes instead of leaving a browser on an hour-old stylesheet.
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
