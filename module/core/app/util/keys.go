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

	KeyDebug   = "debug"
	KeyError   = "error"
	KeyRandom  = "random"
	KeyStart   = "start"
	KeyUnknown = "unknown"

	KeyCSV  = "csv"
	KeyJSON = "json"
	KeyTOML = "toml"
	KeyXML  = "xml"
	KeyYAML = "yaml"

	ExtJSON     = ".json"
	ExtMarkdown = ".md"

	OK    = "ok"
	Error = "error"{{{ if .HasDatabase }}}

{{{ if .MySQL }}}	DatabaseMySQL = "mysql"
{{{ end }}}{{{ if .PostgreSQL }}}	DatabasePostgreSQL = "postgres"
{{{ end }}}{{{ if .SQLite }}}	DatabaseSQLite = "sqlite"
{{{ end }}}{{{ if .SQLServer }}}	DatabaseSQLServer = "sqlserver"
{{{ end }}}{{{ end }}}
	// $PF_SECTION_START(keys)$
	// $PF_SECTION_END(keys)$.
)
