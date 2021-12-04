tell application "Vivaldi"
	set HTML to execute front window's active tab javascript "document.documentElement.outerHTML"
end tell
return HTML