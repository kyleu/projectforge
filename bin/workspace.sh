#!/usr/bin/osascript

set scriptPath to POSIX path of ((path to me as text) & "::")
set scriptDir to do shell script "dirname " & quoted form of scriptPath

tell application "iTerm2"
	set lwin to current session of current tab of current window
	set lwidth to columns of lwin
	tell lwin
    write text "cd " & quoted form of scriptDir
    write text "bin/dev.sh"
    split vertically with default profile
		set columns of lwin to lwidth - 37
	end tell
	set rwin to second session of current tab of current window
	tell rwin
    write text "cd " & quoted form of (scriptDir & "/client")
    write text "clear"
    write text "../bin/build/client-watch.sh"
    set columns of rwin to 36
	end tell
end tell
