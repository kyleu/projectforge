#!/usr/bin/osascript
tell application "iTerm2"
    tell current session of current tab of current window
        write text "cd ~/go/src/{{{ .SourceTrimmed }}}"
        write text "clear"
        write text "bin/dev.sh"
        split vertically with default profile
    end tell
    tell second session of current tab of current window
        write text "cd ~/go/src/{{{ .SourceTrimmed }}}/client"
        write text "clear"
        write text "../bin/build/client-watch.sh"
    end tell
end tell
