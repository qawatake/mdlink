tell application "Brave Browser"
	set HTML to execute front window's active tab javascript "document.documentElement.outerHTML"
end tell
return HTML