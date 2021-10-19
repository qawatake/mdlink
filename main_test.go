package main

import (
	"net/http"
	"testing"

	"golang.org/x/net/html"
)

func TestTitleFinderImpl(t *testing.T) {
	cases := []struct {
		rawurl            string
		wantPageTitle     string
		wantFragmentTitle string
	}{
		{
			rawurl:            "https://gohugo.io/content-management/taxonomies/#add-taxonomies-to-content",
			wantPageTitle:     "Taxonomies | Hugo",
			wantFragmentTitle: "Add Taxonomies to Content",
		},
		{
			rawurl:            "https://pkg.go.dev/net/url#URL",
			wantPageTitle:     "url package - net/url - pkg.go.dev",
			wantFragmentTitle: "",
		},
	}

	for _, tt := range cases {
		f, err := NewTitleFinderImpl(tt.rawurl)
		if err != nil {
			t.Errorf("NewTitleFinderImpl failed: %v", err)
		}

		resp, err := http.Get(tt.rawurl)
		if err != nil {
			t.Errorf("HTTP GET request failed: %v", err)
		}

		defer resp.Body.Close()
		doc, err := html.Parse(resp.Body)
		if err != nil {
			t.Errorf("html.Parse failed: %v", err)
		}
		f.FindPageTitle(doc)
		f.FindFragment(doc)

		if gotPageTitle := f.title; gotPageTitle != tt.wantPageTitle {
			t.Errorf("[ERROR] got: %q, want: %q", gotPageTitle, tt.wantPageTitle)
		}
		if gotFragmentTitle := f.fragmentTitle; gotFragmentTitle != tt.wantFragmentTitle {
			t.Errorf("[ERROR] got: %q, want: %q", gotFragmentTitle, tt.wantFragmentTitle)
		}
	}
}
