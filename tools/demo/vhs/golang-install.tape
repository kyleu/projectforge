# VHS for installation of Project Forge
#
# https://github.com/charmbracelet/vhs

Output ../assets/golang-install.gif
Output ../assets/golang-install.mp4
# Output ../assets/golang-install.webm
# Output ../assets/golang-install/frames/

Set TypingSpeed 50ms
Set FontSize 18
Set Width 1000
Set Height 140
Set Padding 20

Type "go install projectforge.dev/projectforge@latest"

Sleep 500ms

Enter

Sleep 5s

Type "projectforge version"

Sleep 500ms

Enter

Sleep 2s

Type "Project Forge"

Sleep 1s

Type " - A long-lived code generator and project manager for Go web apps"

Sleep 3s
