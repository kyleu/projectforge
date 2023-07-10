@ECHO OFF

rem Builds all the templates using quicktemplate

cd %~dpnx0\..

@ECHO ON
{{{ if .HasModule "database" }}}qtc -ext sql -dir "queries"
{{{ end }}}qtc -ext html -dir "views"
