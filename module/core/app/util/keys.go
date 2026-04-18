package util

const (
	AppKey     = "{{{ .Key }}}"
	AppName    = "{{{ .Name }}}"
	AppSummary = "{{{ .Info.Summary }}}"
	AppPort    = {{{ .Port }}}
	AppCmd     = {{{ .ExecSafe }}}
	AppContact = "{{{ .Info.AuthorName }}} <{{{ .Info.AuthorEmail }}}>"
	AppURL     = "{{{ .Info.Homepage }}}"
	AppSource  = "{{{ .Info.Sourcecode }}}"
	AppLegal   = `Built by <a href="mailto:{{{ .Info.AuthorEmail }}}">{{{ .Info.AuthorName }}}</a>, all rights reserved`

	KeyOK       = "ok"
	KeySuccess  = "success"
	KeyError    = "error"
	KeyDebug    = "debug"
	KeyStart    = "start"
	KeyRandom   = "random"
	KeyUnknown  = "unknown"
	KeyEllipsis = "…"

	KeyCSV  = "csv"
	KeyJSON = "json"
	KeyTOML = "toml"
	KeyXML  = "xml"
	KeyYAML = "yaml"

	ExtJSON     = ".json"
	ExtMarkdown = ".md"{{{ .DatabaseKeys }}}
	// $PF_SECTION_START(keys)$
	// $PF_SECTION_END(keys)$.
)
