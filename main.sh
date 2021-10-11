#!/bin/bash
front_app=$(osascript get_frontmost_appname.scpt)

if [ "$front_app" = "Brave Browser.app" ]; then
	url=$(osascript url_from_brave.scpt)
	html=$(osascript html_from_brave.scpt)
fi
if [ "$front_app" = "Google Chrome.app" ]; then
	url=$(osascript url_from_chrome.scpt)
	html=$(osascript html_from_chrome.scpt)
fi
if [ "$front_app" = "Safari.app" ]; then
	url=$(osascript url_from_safari.scpt)
	html=$(osascript html_from_safari.scpt)
fi
echo $html | ./url2mdlink $url