# VHS for installation of Project Forge
#
# https://github.com/charmbracelet/vhs

Output ../assets/golang-install.light.gif
# Output ../assets/golang-install.light.mp4
# Output ../assets/golang-install.light.webm
# Output ../assets/golang-install/frames.light/

Set TypingSpeed 50ms
Set FontSize 18
Set Width 1000
Set Height 140
Set Padding 20
Set Theme "AtomOneLight"

Type "go install projectforge.dev/projectforge@latest"

Sleep 500ms

Enter

Wait

Sleep 500ms

Type "projectforge version"

Sleep 500ms

Enter

Sleep 1s

Type "Project Forge"

Sleep 1s

Type " - Build and maintain feature-rich applications using Golang"

Sleep 3s
