#!/bin/bash
BRAVE="Brave Browser"
CHROME="Google Chrome"
SAFARI="Safari"

# 2回め以降の動作
# ブラウザが最前面にある場合には, 初回実行時にはクリップボードの内容を処理しない
if [ $runned -ge 1 ]; then
	clipboard=$(pbpaste)
	./mdlink "$browserurl" "$browser" "$runned" "$clipboard" "$browsertitle"
	exit 0
fi

# 初回起動時の動作
front_app=$(osascript appscript/frontmost_appname.scpt)

if [ "$front_app" = "$BRAVE"".app" ]; then
  browser=$BRAVE
	browserurl=$(osascript appscript/url_from_brave.scpt)
	html=$(osascript appscript/html_from_brave.scpt)
fi
if [ "$front_app" = "$CHROME"".app" ]; then
  browser=$CHROME
	url=$(osascript appscript/url_from_chrome.scpt)
	html=$(osascript appscript/html_from_chrome.scpt)
fi
if [ "$front_app" = "$SAFARI"".app" ]; then
  browser=$SAFARI
	browserurl=$(osascript appscript/url_from_safari.scpt)
	html=$(osascript appscript/html_from_safari.scpt)
fi

clipboard=$(pbpaste)
echo $html | ./mdlink "$browserurl" "$browser" "$runned" "$clipboard" ""