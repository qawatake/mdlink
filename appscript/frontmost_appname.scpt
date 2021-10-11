tell application "System Events"
  set front_app to name of (path to frontmost application)
end tell
return front_app