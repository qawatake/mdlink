tell application "Google Chrome"
	set HTML to execute front window's active tab javascript "document.documentElement.outerHTML"
end tell
return HTML