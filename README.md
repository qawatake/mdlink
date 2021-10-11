# mdlink

## ざっくり
ブラウザがトップに表示しているページに関するマークダウン用のリンクをクリップボードにコピーする.

## 使い方
1. Alfred の検索窓に `mdlink` と入力する.
2. ブラウザがトップに表示しているページのタイトルが表示される.
3. エンターキーを押すと, 表示されたタイトルとページの URL を含むマークダウン用リンク (e.g., `[Google](https://www.google.com)` がクリップボードにコピーされる.

## インストール手順

## workflow を構成するファイル
- main.sh: シュルスクリプト. AppleScript と Go の橋渡し的な処理を行う.
- mdlink: バイナリファイル. この workflow の中心的な役割を持つ. `GOOS=darwin go build` によって生成する.
  - main.go
  - entity.go
  - go.mod
- appscript/url_from_*.scpt: AppleScript. ブラウザから現在開いているページの URL を取得し, 標準出力に吐き出す.
- appscript/html_from_*.scpt: AppleScript. ブラウザから現在開いているページの HTML のソースを取得し, 標準出力に吐き出す.
- appscript/frontmost_appname.scpt: AppleScript. 最前面にあるアプリ名を取得し, 標準出力に吐き出す.
- info.plist: XML ファイル. Alfred Workflow の全体構成を記述したもの. Alfred によって自動的に生成される.