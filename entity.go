package main

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type TitleFinder interface {
	FindPageTitle(n *html.Node)
}

type TitleFinderImpl struct {
	title string
}

type ScriptFilterOutput struct {
	Rerun     float32 `json:"rerun"` // 再実行への待機時間. 再実行する場合に設定.
	Variables struct {
		Runned int `json:"runned"` // すでに実行した回数
	} `json:"variables"`
	Items []*ScriptFilterItem `json:"items"`
}

type ScriptFilterItem struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Arg      string `json:"arg"`
	Valid    bool   `json:"valid,omitempty"`
}

func (f *TitleFinderImpl) FindPageTitle(n *html.Node) {
	if f.title != "" {
		return
	}

	// title タグを検索
	if n.Type == html.ElementNode && n.DataAtom == atom.Title && n.FirstChild.Type == html.TextNode {
		f.title = n.FirstChild.Data
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

func NewScriptFilterItem(title string, subtitle string, arg string, valid bool) *ScriptFilterItem {
	output := ScriptFilterItem{
		Title:    title,
		Subtitle: subtitle,
		Arg:      arg,
		Valid:    valid,
	}
	return &output
}

func (output *ScriptFilterOutput) addItem(item *ScriptFilterItem) {
	output.Items = append(output.Items, item)
}
