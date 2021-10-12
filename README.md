# mdlink

## ざっくり
表示中のブラウザあるいはクリップボード上の URL から, マークダウン用のリンクを生成する.

## 対応ブラウザ
- Safari
- Google Chrome
- Brave Browser

## 使い方
1. Alfred の検索窓に `mdlink` と入力する.
2. 実行した状況により, 動作は変わる.
  - ブラウザを表示しながら実行した場合, 表示中のページのタイトルが表示される.
  - そうでない場合, クリップボードにある URL をもとにして, 参照先のページのタイトルが表示される.
4. エンターキーを押すと, 表示されたタイトルとページの URL を含むマークダウン用リンク (e.g., `[Google](https://www.google.com)` がクリップボードにコピーされる.

[![Image from Gyazo](https://i.gyazo.com/d1e9f97f64f9c365ceb1edb3562b66bd.gif)](https://gyazo.com/d1e9f97f64f9c365ceb1edb3562b66bd)

## インストール手順
1. https://github.com/qawatake/mdlink/releases/latest から `mdlink.alfredworkflow` をダウンロード.
2. ダウンロードしたファイルを開けば, Alfred が自動的に workflow を追加してくれる (はず).
3. Google Chrome あるいは Brave Browser を使用する場合, ↓が必要.
    1. ブラウザを起動.
    2. ツールバー: [View] -> [Developer] -> [Allow JavaScript from Apple Events]

## workflow を構成するファイル
- `main.sh`: シュルスクリプト. AppleScript と Go の橋渡し的な処理を行う.
- `mdlink`: バイナリファイル. この workflow の中心的な役割を持つ. `GOOS=darwin go build` によって生成する.
  - `main.go`
  - `entity.go`
  - `go.mod`
- `appscript/url_from_*.scpt`: AppleScript. ブラウザから現在開いているページの URL を取得し, 標準出力に吐き出す.
- `appscript/html_from_*.scpt`: AppleScript. ブラウザから現在開いているページの HTML のソースを取得し, 標準出力に吐き出す.
- `appscript/frontmost_appname.scpt`: AppleScript. 最前面にあるアプリ名を取得し, 標準出力に吐き出す.
- `info.plist`: XML ファイル. Alfred Workflow の全体構成を記述したもの. Alfred によって自動的に生成される.