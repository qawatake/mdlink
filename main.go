package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// 第1引数はコピーする URL
// 標準入力で HTML のソースを受け取る
// 標準出力で Script Filter の JSON を出力 (https://www.alfredapp.com/help/workflows/inputs/script-filter/)
func main() {
	if len(os.Args) != 2 {
		log.Fatal("引数の数が不正です")
		return
	}

	url := strings.Trim(os.Args[1], "\n")

	f := new(TitleFinderImpl)
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	f.FindPageTitle(doc)
	output := &ScriptFilterOutput{}
	item := NewScriptFilterItem(f.title, "", fmt.Sprintf("[%v](%v)", f.title, url), true)
	output.addItem(item)

	ec := json.NewEncoder(os.Stdout)
	if err := ec.Encode(output); err != nil {
		log.Fatal(err)
	}
}
