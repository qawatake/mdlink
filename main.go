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

// 第1引数はコピーする URL
// 標準入力で HTML のソースを受け取る
// 標準出力で Script Filter の JSON を出力 (https://www.alfredapp.com/help/workflows/inputs/script-filter/json/)
func main() {
	if len(os.Args) != 6 {
		log.Fatal("引数の数が不正です")
		return
	}

	browserURL := strings.Trim(os.Args[1], "\n")
	browser := os.Args[2]
	var runned float32
	if i, err := strconv.ParseFloat(os.Args[3], 32); err != nil {
		log.Fatal(err)
	} else {
		runned = float32(i)
	}
	clipboard := os.Args[4]
	htmlSource := strings.NewReader(os.Args[5])

	output := &ScriptFilterOutput{}

	// ブラウザを開いている場合, ブラウザがトップページに開いている HTML を使用する
	if browser != "" {
		doc, err := html.Parse(htmlSource)
		if err != nil {
			log.Fatal(err)
		}
		f := new(TitleFinderImpl)
		f.FindPageTitle(doc)
		item := NewScriptFilterItem(f.title, fmt.Sprintf("from %s", browser), fmt.Sprintf("[%v](%v)", f.title, browserURL), true)
		output.addItem(item)
	}

	// ブラウザを開いている & 初回の場合, クリップボード上の URL を使ってインターネットから HTML を取得しない
	if browser != "" && runned == 0 {
		// URL を使って html を取得
		output.Rerun = 0.1
		output.Variables.Runned = 1
	}

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
