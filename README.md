# url2mdlink

## ざっくり
Brave Browser がトップに表示しているページに関するマークダウン用のリンクをクリップボードにコピーする.

## 使い方
1. Alfred の検索窓に `url2mdlink` と入力する.
2. Brave Browser がトップに表示しているページのタイトルが表示される.
3. エンターキーを押すと, 表示されたタイトルとページの URL を含むマークダウン用リンク (e.g., `[Google](https://www.google.com)` がクリップボードにコピーされる.

## インストール手順

## workflow を構成するファイル
- urlmd: バイナリファイル. この workflow の中心的な役割を持つ. `GOOS=darwin go build` によって生成する.
  - main.go
  - entity.go
  - go.mod
- get_url_from_browser.scpt: Apple Script. Brave Browser から現在開いているページの URL を取得し, 標準出力に吐き出す.
- info.plist: XML ファイル. Alfred Workflow の全体構成を記述したもの. Alfred によって自動的に生成される.