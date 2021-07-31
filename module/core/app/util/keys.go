package util

const (
	AppKey     = "{{{ .Key }}}"
	AppName    = "{{{ .Name }}}"
	AppSummary = "{{{ .Summary }}}"
	AppPort    = {{{ .Port }}}
	AppContact = "{{{ .AuthorName }}} <{{{ .AuthorEmail }}}>"
	AppURL     = "{{{ .Homepage }}}"
	AppSource  = "{{{ .Sourcecode }}}"
	AppLegal   = `Built by <a href="mailto:{{{ .AuthorEmail }}}">{{{ .AuthorName }}}</a>, all rights reserved`
)
