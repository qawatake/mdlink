package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	if len(os.Args) != 6 {
		log.Fatal("引数の数が不正です")
		return
	}

	browserURL := strings.Trim(os.Args[1], "\n")
	browser := os.Args[2]
	var runned int
	if i, err := strconv.ParseInt(os.Args[3], 10, 32); err != nil {
		log.Fatal(err)
	} else {
		runned = int(i)
	}
	clipboard := os.Args[4]
	browserTitle := os.Args[5]

	output := &ScriptFilterOutput{}

	// ブラウザの内容を処理
	if browser != "" {
		var item *ScriptFilterItem
		if runned >= 1 {
			item = NewScriptFilterItem(browserTitle, fmt.Sprintf("from %s", browser), fmt.Sprintf("[%v](%v)", browserTitle, browserURL), true)
		} else { // 初回起動 && ブラウザを開いている場合, ブラウザがトップページに開いている HTML を使用する
			doc, err := html.Parse(os.Stdin)
			if err != nil {
				log.Fatal(err)
			}
			f := new(TitleFinderImpl)
			f.FindPageTitle(doc)
			item = NewScriptFilterItem(f.title, fmt.Sprintf("from %s", browser), fmt.Sprintf("[%v](%v)", f.title, browserURL), true)
			browserTitle = f.title
		}
		output.addItem(item)
	}

	// クリップボードの内容を処理
	// ブラウザを開いていない OR 2回目以降のコード実行の場合, クリップボード上の URL を使ってインターネットから HTML
	if (browser == "" || runned >= 1) && isURL(clipboard) {
		resp, err := http.Get(clipboard)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		doc, err := html.Parse(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		f := new(TitleFinderImpl)
		f.FindPageTitle(doc)
		item := NewScriptFilterItem(f.title, "from Clipboard", fmt.Sprintf("[%v](%v)", f.title, clipboard), true)
		output.addItem(item)
	}

	// 再実行するかどうかを判定
	// ブラウザを開いている & 初回の場合, クリップボード上の URL を使ってインターネットから HTML を取得しない
	if browser != "" && runned == 0 {
		output.Rerun = 0.1
	}

	output.Variables.Runned = int(runned) + 1
	output.Variables.Browser = browser
	output.Variables.BrowserURL = browserURL
	output.Variables.Title = browserTitle

	ec := json.NewEncoder(os.Stdout)
	if err := ec.Encode(output); err != nil {
		log.Fatal(err)
	}
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return false
	}
	matched, _ := regexp.Match(`(http)|(https)(:\d+)?`, []byte(u.Scheme))
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return false
	}
	return matched
}
