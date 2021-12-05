package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/net/html"
)

// 引数1: ブラウザから取得した URL
// 引数2: ブラウザ名
// 引数3: プログラムが実行された回数
// 引数4: クリップボードの内容
// 引数5: ブラウザから取得した ページタイトル
// 標準入力: ブラウザから取得した HTML のソース
// 標準出力で Script Filter の JSON を出力 (https://www.alfredapp.com/help/workflows/inputs/script-filter/json/)
func main() {
	if err := run(os.Args, os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func run(args []string, htmlfrom io.Reader, jsonto io.Writer) error {
	if len(args) != 6 {
		return fmt.Errorf("wrong number of args | want: %d, got: %d", 5, len(args)-1)
	}

	browserURL := strings.Trim(args[1], "\n")
	browser := args[2]
	var runned int
	if i, err := strconv.ParseInt(args[3], 10, 32); err != nil {
		return errors.Wrap(err, "failed to parse int")
	} else {
		runned = int(i)
	}
	clipboard := args[4]
	browserTitle := args[5]
	var htmlSource []byte
	if browser != "" && runned < 1 {
		var err error
		htmlSource, err = io.ReadAll(htmlfrom)
		if err != nil {
			return errors.Wrap(err, "failed to read html source")
		}
	}

	output := &ScriptFilterOutput{}

	// 再実行するかどうかを判定
	// ブラウザを開いている & 初回の場合, クリップボード上の URL を使ってインターネットから HTML を取得しない
	if browser != "" && runned == 0 {
		output.Rerun = 0.1
	}

	// Script Filter の variables field を設定
	output.Variables.Runned = int(runned) + 1
	output.Variables.Browser = browser
	output.Variables.BrowserURL = browserURL
	if browser != "" && runned < 1 {
		pageTitle, fragmentTitle, err := getTitles(browserURL, bytes.NewReader(htmlSource))
		if err != nil {
			return errors.Wrap(err, "getTitles failed")
		}
		output.Variables.Title = buildTitle(pageTitle, fragmentTitle)
	}

	// <==== Script Filter の items field を設定
	// ブラウザで処理済みの2回め
	if browser != "" && runned >= 1 {
		item := NewScriptFilterItem(browserTitle, fmt.Sprintf("from %s", browser), fmt.Sprintf("[%v](%v)", browserTitle, browserURL), true)
		output.addItem(item)
	}

	// ブラウザの内容を処理
	if browser != "" && runned < 1 {
		item, err := handleBrowser(htmlSource, browserURL, browser)
		if err != nil {
			return errors.Wrap(err, "handleBrowser failed")
		}
		if item != nil {
			output.addItem(item)
		}
	}

	// クリップボードの内容を処理
	// ブラウザを開いていない OR 2回目以降のコード実行の場合, クリップボード上の URL を使ってインターネットから HTML
	if (browser == "" || runned >= 1) && isURL(clipboard) {
		item, err := handleClipboard(clipboard)
		if err != nil {
			return errors.Wrap(err, "handleClipboard failed")
		}
		if item != nil {
			output.addItem(item)
		}
	}
	// <====

	ec := json.NewEncoder(jsonto)
	if err := ec.Encode(output); err != nil {
		return errors.Wrap(err, "failed to encode json")
	}
	return nil
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func handleBrowser(htmlSource []byte, browserURL string, browser string) (*ScriptFilterItem, error) {
	pageTitle, fragmentTitle, err := getTitles(browserURL, bytes.NewReader(htmlSource))
	if err != nil {
		return nil, errors.Wrap(err, "getTitles failed")
	}
	title := buildTitle(pageTitle, fragmentTitle)
	return NewScriptFilterItem(title, fmt.Sprintf("from %s", browser), fmt.Sprintf("[%v](%v)", title, browserURL), true), nil
}

func handleClipboard(clipboard string) (*ScriptFilterItem, error) {
	resp, err := http.Get(clipboard)
	if err != nil {
		return nil, errors.Wrap(err, "http get failed")
	}
	defer resp.Body.Close()
	pageTitle, fragmentTitle, err := getTitles(clipboard, resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "getTitle failed")
	}
	title := buildTitle(pageTitle, fragmentTitle)
	item := NewScriptFilterItem(title, "from Clipboard", fmt.Sprintf("[%v](%v)", title, clipboard), true)
	return item, nil
}

func getTitles(urltxt string, htmlfrom io.Reader) (pageTitle, fragmentTitle string, err error) {
	doc, err := html.Parse(htmlfrom)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to parse html")
	}

	f, err := NewTitleFinderImpl(urltxt)
	if err != nil {
		return "", "", errors.Wrap(err, "NewTitleFinderImpl failed")
	}
	f.FindFragmentTitle(doc)
	f.FindPageTitle(doc)
	return f.title, f.fragmentTitle, nil
}

func buildTitle(pageTitle, fragmentTitle string) (title string) {
	if fragmentTitle != "" {
		return fmt.Sprintf("%s - %s", fragmentTitle, pageTitle)
	}
	return pageTitle
}
