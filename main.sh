#!/bin/bash
BRAVE="Brave Browser"
CHROME="Google Chrome"
SAFARI="SAFARI"
front_app=$(osascript appscript/get_frontmost_appname.scpt)

# if [ "$front_app" = "Brave Browser.app" ]; then
if [ "$front_app" = "$BRAVE"".app" ]; then
  source=$BRAVE
	url=$(osascript appscript/url_from_brave.scpt)
	html=$(osascript appscripthtml_from_brave.scpt)
fi
if [ "$front_app" = "$CHROME"".app" ]; then
  source=$CHROME
	url=$(osascript appscript/url_from_chrome.scpt)
	html=$(osascript appscript/html_from_chrome.scpt)
fi
if [ "$front_app" = "$SAFARI"".app" ]; then
  source=$SAFARI
	url=$(osascript appscript/url_from_safari.scpt)
	html=$(osascript appscript/html_from_safari.scpt)
fi

clipboard=$(pbpaste)
echo $html | ./mdlink "$url" "$source" "$runned" "$clipboard" "$html"