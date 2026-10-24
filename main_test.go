package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func TestTitleFinderImpl(t *testing.T) {
	// Serve fixed HTML instead of fetching live pages so the test does not
	// depend on network access or on upstream sites changing their titles.
	mux := http.NewServeMux()
	mux.HandleFunc("/content-management/taxonomies/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Taxonomies | Hugo</title></head>
<body>
<h1>Taxonomies</h1>
<h2 id="add-taxonomies-to-content">Add Taxonomies to Content</h2>
<p>...</p>
</body>
</html>`))
	})
	mux.HandleFunc("/net/url", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>url package - net/url - pkg.go.dev</title></head>
<body>
<div id="URL"><h4>type URL</h4></div>
</body>
</html>`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	cases := []struct {
		rawurl            string
		wantPageTitle     string
		wantFragmentTitle string
	}{
		{
			rawurl:            srv.URL + "/content-management/taxonomies/#add-taxonomies-to-content",
			wantPageTitle:     "Taxonomies | Hugo",
			wantFragmentTitle: "Add Taxonomies to Content",
		},
		{
			rawurl:            srv.URL + "/net/url#URL",
			wantPageTitle:     "url package - net/url - pkg.go.dev",
			wantFragmentTitle: "",
		},
	}

	for _, tt := range cases {
		f, err := NewTitleFinderImpl(tt.rawurl)
		if err != nil {
			t.Fatalf("NewTitleFinderImpl failed: %v", err)
		}

		resp, err := http.Get(tt.rawurl)
		if err != nil {
			t.Fatalf("HTTP GET request failed: %v", err)
		}

		defer resp.Body.Close()
		doc, err := html.Parse(resp.Body)
		if err != nil {
			t.Fatalf("html.Parse failed: %v", err)
		}
		f.FindPageTitle(doc)
		f.FindFragmentTitle(doc)

		if gotPageTitle := f.title; gotPageTitle != tt.wantPageTitle {
			t.Errorf("[ERROR] got: %q, want: %q", gotPageTitle, tt.wantPageTitle)
		}
		if gotFragmentTitle := f.fragmentTitle; gotFragmentTitle != tt.wantFragmentTitle {
			t.Errorf("[ERROR] got: %q, want: %q", gotFragmentTitle, tt.wantFragmentTitle)
		}
	}
}
