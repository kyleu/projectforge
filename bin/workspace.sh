#!/usr/bin/osascript
# Content managed by Project Forge, see [projectforge.md] for details.
tell application "iTerm2"
	set lwin to current session of current tab of current window
	set lwidth to columns of lwin
	tell lwin
    write text "cd ~/go/src/github.com/kyleu/projectforge"
    write text "clear"
    write text "bin/dev.sh"
    split vertically with default profile
		set columns of lwin to lwidth - 52
	end tell
	set rwin to second session of current tab of current window
	tell rwin
    write text "cd ~/go/src/github.com/kyleu/projectforge/client"
    write text "clear"
    write text "../bin/build/client-watch.sh"
    set columns of rwin to 51
	end tell
end tell
