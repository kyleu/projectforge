#!/usr/bin/osascript
tell application "iTerm"
	set lwin to current session of current tab of current window
	tell lwin
    write text "cd ~/go/src/{{{ .SourceTrimmed }}}"
    write text "clear"
    write text "bin/dev.sh"
    split vertically with default profile
		set columns of lwin to 147
	end tell
	set rwin to second session of current tab of current window
	tell rwin
    write text "cd ~/go/src/{{{ .SourceTrimmed }}}/client"
    write text "clear"
    write text "../bin/build/client-watch.sh"
    set columns of rwin to 51
	end tell
end tell
