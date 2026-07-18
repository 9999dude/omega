package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testSite(t *testing.T) *site {
	t.Helper()
	docs := fstest.MapFS{
		"systems/intro.md":       {Data: []byte("# Systems Thinking\n\nA guide to seeing the whole system.\n\n## Feedback loops\n\nUse feedback well.\n\n```mermaid\nflowchart LR\n  A[Signal] --> B[Response]\n```\n\n```go\npackage main\n```\n\n```yaml\nenabled: true\n```\n\n```json\n{\"enabled\": true}\n```\n\n| Signal | Meaning |\n|---|---|\n| Errors | Reliability |")},
		"systems/deep/queues.md": {Data: []byte("# Queues\n\nQueues absorb bursts. Read the [introduction](../intro.md#feedback-loops).")},
		"labs/README.md":         {Data: []byte("# Hands-on labs\n\nPractice the ideas safely.")},
	}
	s, err := newSite(docs)
	if err != nil {
		t.Fatalf("newSite() error = %v", err)
	}
	return s
}

func TestLoadArticles(t *testing.T) {
	s := testSite(t)
	if got, want := len(s.articles), 3; got != want {
		t.Fatalf("article count = %d, want %d", got, want)
	}
	if _, ok := s.bySlug["labs"]; !ok {
		t.Error("README.md should be available at its directory slug")
	}
	queues := s.bySlug["systems/deep/queues"]
	if !strings.Contains(string(queues.Content), `href="/docs/systems/intro#feedback-loops"`) {
		t.Errorf("local Markdown link was not rewritten: %s", queues.Content)
	}
	intro := s.bySlug["systems/intro"]
	if strings.Contains(string(intro.Content), "<h1") {
		t.Errorf("rendered content repeats the page title: %s", intro.Content)
	}
	if !strings.Contains(string(intro.Content), `class="language-mermaid"`) {
		t.Errorf("Mermaid fence does not retain its language class: %s", intro.Content)
	}
	for _, language := range []string{"go", "yaml", "json"} {
		if !strings.Contains(string(intro.Content), `class="language-`+language+`"`) {
			t.Errorf("%s fence does not retain its language class: %s", language, intro.Content)
		}
	}
	if len(intro.TOC) != 1 || intro.TOC[0].ID != "feedback-loops" {
		t.Errorf("TOC = %#v", intro.TOC)
	}
}

func TestRoutes(t *testing.T) {
	s := testSite(t)
	tests := []struct {
		path        string
		status      int
		contains    string
		contentType string
	}{
		{path: "/", status: http.StatusOK, contains: "highlight.min.js?v=11.11.1", contentType: "text/html"},
		{path: "/docs/systems/intro", status: http.StatusOK, contains: "Feedback loops", contentType: "text/html"},
		{path: "/docs/missing", status: http.StatusNotFound, contains: "404 page not found", contentType: "text/plain"},
		{path: "/healthz", status: http.StatusOK, contains: "ok", contentType: "text/plain"},
		{path: "/assets/app.css", status: http.StatusOK, contains: "--accent", contentType: "text/css"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			s.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if recorder.Code != tt.status {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.status)
			}
			if !strings.Contains(recorder.Body.String(), tt.contains) {
				t.Errorf("body does not contain %q", tt.contains)
			}
			if !strings.Contains(recorder.Header().Get("Content-Type"), tt.contentType) {
				t.Errorf("content type = %q, want %q", recorder.Header().Get("Content-Type"), tt.contentType)
			}
			if recorder.Header().Get("Content-Security-Policy") == "" {
				t.Error("missing Content-Security-Policy")
			}
		})
	}
}

func TestSearch(t *testing.T) {
	s := testSite(t)
	recorder := httptest.NewRecorder()
	s.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/search?q=whole", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	var results []searchResult
	if err := json.Unmarshal(recorder.Body.Bytes(), &results); err != nil {
		t.Fatalf("decode search results: %v", err)
	}
	if len(results) != 1 || results[0].Title != "Systems Thinking" {
		t.Errorf("results = %#v", results)
	}
}

func TestEmptyDocsFails(t *testing.T) {
	_, err := newSite(fstest.MapFS{})
	if err == nil {
		t.Fatal("newSite() should fail when no Markdown exists")
	}
}
