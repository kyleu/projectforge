package util

const (
	AppKey     = "{{{ .Key }}}"
	AppName    = "{{{ .Name }}}"
	AppSummary = "{{{ .Info.Summary }}}"
	AppPort    = {{{ .Port }}}
	AppContact = "{{{ .Info.AuthorName }}} <{{{ .Info.AuthorEmail }}}>"
	AppURL     = "{{{ .Info.Homepage }}}"
	AppSource  = "{{{ .Info.Sourcecode }}}"
	AppLegal   = `Built by <a href="mailto:{{{ .Info.AuthorEmail }}}">{{{ .Info.AuthorName }}}</a>, all rights reserved`

	BoolTrue   = "true"
	BoolFalse  = "false"
	KeyError   = "error"
	KeyRandom  = "random"
	KeyUnknown = "unknown"

	ExtJSON     = ".json"
	ExtMarkdown = ".md"

	// $PF_SECTION_START(keys)$
	// $PF_SECTION_END(keys)$.
)
