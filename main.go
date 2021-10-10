package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("引数の数が不正です")
		return
	}

	url := strings.Trim(os.Args[1], "\n")
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	f := new(TitleFinderImpl)
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	f.findPageTitle(doc)
	output := &ScriptFilterOutput{}
	item := NewScriptFilterItem(f.title, "", fmt.Sprintf("[%v](%v)", f.title, url), true)
	output.addItem(item)

	ec := json.NewEncoder(os.Stdout)
	if err := ec.Encode(output); err != nil {
		log.Fatal(err)
	}
}
