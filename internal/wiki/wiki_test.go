package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HistoryLabs/events-api/internal/wiki"
)

func setupServer(t *testing.T, handler http.HandlerFunc, statusCode int, body string) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})
	return ts
}

func TestFetch_SetsUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})

	wiki.Fetch(context.Background(), "March_2", wiki.FetchOpts{Section: 1})

	if gotUA == "" {
		t.Fatal("expected User-Agent to be set, got empty string")
	}
}

func TestFetch_QueryParams_NoRedirects(t *testing.T) {
	var gotQuery map[string]string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		gotQuery = map[string]string{
			"action":    q.Get("action"),
			"format":    q.Get("format"),
			"page":      q.Get("page"),
			"section":   q.Get("section"),
			"redirects": q.Get("redirects"),
		}
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})

	wiki.Fetch(context.Background(), "March_15", wiki.FetchOpts{Section: 1})

	checks := map[string]string{
		"action":  "parse",
		"format":  "json",
		"page":    "March_15",
		"section": "1",
	}
	for k, want := range checks {
		if got := gotQuery[k]; got != want {
			t.Errorf("query param %q: got %q, want %q", k, got, want)
		}
	}
	if got := gotQuery["redirects"]; got != "" {
		t.Errorf("expected no redirects param, got %q", got)
	}
}

func TestFetch_QueryParams_WithRedirects(t *testing.T) {
	var gotRedirects string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRedirects = r.URL.Query().Get("redirects")
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})

	wiki.Fetch(context.Background(), "AD_1776", wiki.FetchOpts{Section: 1, Redirects: true})

	if gotRedirects != "true" {
		t.Errorf("expected redirects=true, got %q", gotRedirects)
	}
}

func TestFetch_ReturnsBodyVerbatim(t *testing.T) {
	want := `{"parse":{"title":"test"}}`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(want))
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})

	got, err := wiki.Fetch(context.Background(), "Test", wiki.FetchOpts{Section: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFetch_NonOK_ReturnsErrUpstreamStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	origURL := wiki.DefaultClient.BaseURL
	wiki.DefaultClient.BaseURL = ts.URL
	t.Cleanup(func() {
		wiki.DefaultClient.BaseURL = origURL
		ts.Close()
	})

	_, err := wiki.Fetch(context.Background(), "Test", wiki.FetchOpts{Section: 1})
	if !errors.Is(err, wiki.ErrUpstreamStatus) {
		t.Errorf("expected ErrUpstreamStatus, got %v", err)
	}
}
