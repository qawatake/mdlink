package main

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TitleFinder interface {
	FindPageTitle(n *html.Node)
	FindFragmentTitle(n *html.Node)
}

type TitleFinderImpl struct {
	title         string
	fragmentId    string
	fragmentTitle string
}

func NewTitleFinderImpl(rawurl string) (*TitleFinderImpl, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("url.Parse failed: %w", err)
	}

	f := new(TitleFinderImpl)
	f.fragmentId = u.Fragment
	return f, nil
}

func (f *TitleFinderImpl) FindFragmentTitle(n *html.Node) {
	if f.fragmentTitle != "" {
		return
	}

	// 指定した id を含むタグを検索
	if n.Type == html.ElementNode && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == f.fragmentId {
				f.fragmentTitle = strings.TrimFunc(n.FirstChild.Data, func(r rune) bool {
					return unicode.IsSpace(r)
				})
				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f.FindFragmentTitle(c)
	}
}

func (f *TitleFinderImpl) FindPageTitle(n *html.Node) {
	if f.title != "" {
		return
	}

	// title タグを検索
	if n.Type == html.ElementNode && n.DataAtom == atom.Title && n.FirstChild.Type == html.TextNode {
		f.title = strings.TrimFunc(n.FirstChild.Data, func(r rune) bool {
			return unicode.IsSpace(r)
		})
		return
	}

	// <meta property="og:title" ...> を検索
	if n.Type == html.ElementNode && n.DataAtom == atom.Meta {
		noTitle := false
		for _, a := range n.Attr {
			if a.Key == "property" && a.Val == "og:title" {
				noTitle = true
			}
		}
		if !noTitle {
			return
		}
		for _, a := range n.Attr {
			if a.Key == "content" {
				f.title = a.Val
				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		f.FindPageTitle(c)
	}
}
