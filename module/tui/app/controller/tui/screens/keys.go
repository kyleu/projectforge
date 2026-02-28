// $PF_GENERATE_ONCE$
package screens

const (
	KeyMainMenu    = "mainmenu"{{{ if .HasModule "docbrowse" }}}
	KeyDocs        = "docs"{{{ end }}}{{{ if .HasModule "filesystem" }}}
	KeyFileBrowser = "filebrowser"
	KeyFileViewer  = "fileviewer"{{{ end }}}
	KeySettings    = "settings"
	KeyAbout       = "about"

	KeyEnter     = "enter"
	KeyEsc       = "esc"
	KeyBackspace = "backspace"
	KeyLeft      = "left"
	KeyPgDown    = "pgdown"
	KeyPgUp      = "pgup"
	KeyHome      = "home"
	KeyEnd       = "end"
)
