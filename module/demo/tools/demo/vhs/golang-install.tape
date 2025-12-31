# VHS for installation of Project Forge
#
# https://github.com/charmbracelet/vhs

Output ../assets/golang-install.gif
# Output ../assets/golang-install.mp4
# Output ../assets/golang-install.webm
# Output ../assets/golang-install/frames/

Set TypingSpeed 50ms
Set FontSize 18
Set Width 1000
Set Height 140
Set Padding 20
Set Theme "Atom"

Type "go install {{{ .Package }}}@latest"

Sleep 500ms

Enter

Wait

Sleep 500ms

Type "{{{ .Exec }}} version"

Sleep 500ms

Enter

Sleep 1s

Type "{{{ .Name }}}"

Sleep 1s

Type " - {{{ .Info.Summary }}}"

Sleep 3s
