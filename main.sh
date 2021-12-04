#!/bin/bash
BRAVE="Brave Browser"
CHROME="Google Chrome"
SAFARI="Safari"
VIVALDI="Vivaldi"

# 2回め以降の動作
# ブラウザが最前面にある場合には, 初回実行時にはクリップボードの内容を処理しない
if [ $runned -ge 1 ]; then
	clipboard=$(pbpaste)
	./mdlink "$browserurl" "$browser" "$runned" "$clipboard" "$browsertitle"
	exit 0
fi

# 初回起動時の動作
front_app=$(osascript appscript/frontmost_appname.scpt)
clipboard=$(pbpaste)

case "$front_app" in
	"$BRAVE"".app")
		browserurl=$(osascript appscript/url_from_brave.scpt)
		osascript appscript/html_from_brave.scpt | ./mdlink "$browserurl" "$BRAVE" "$runned" "$clipboard" ""
		;;
	"$CHROME"".app")
		browserurl=$(osascript appscript/url_from_chrome.scpt)
		osascript appscript/html_from_chrome.scpt | ./mdlink "$browserurl" "$CHROME" "$runned" "$clipboard" ""
		;;
	"$SAFARI"".app")
		browserurl=$(osascript appscript/url_from_safari.scpt)
		osascript appscript/html_from_safari.scpt | ./mdlink "$browserurl" "$SAFARI" "$runned" "$clipboard" ""
		;;
	"$VIVALDI"".app")
		browserurl=$(osascript appscript/url_from_vivaldi.scpt)
		osascript appscript/html_from_vivaldi.scpt | ./mdlink "$browserurl" "$VIVALDI" "$runned" "$clipboard" ""
		;;
	*)
		osascript appscript/html_from_safari.scpt | ./mdlink "" "" "$runned" "$clipboard" ""
		;;
esac