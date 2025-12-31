# $PF_HAS_MODULE(homebrew)$
# VHS for installation of Project Forge
#
# https://github.com/charmbracelet/vhs

Output ../assets/homebrew-install.light.gif
# Output ../assets/homebrew-install.light.mp4
# Output ../assets/homebrew-install.light.webm
# Output ../assets/homebrew-install/frames.light/

Set TypingSpeed 50ms
Set FontSize 18
Set Width 1000
Set Height 230
Set Padding 20
Set Theme "AtomOneLight"

Type "brew install {{{ .Info.Org }}}/{{{ .Info.Org }}}/{{{ .Key }}}"

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
