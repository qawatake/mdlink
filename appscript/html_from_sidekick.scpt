tell application "Sidekick"
	set HTML to execute front window's active tab javascript "document.documentElement.outerHTML"
end tell
return HTML