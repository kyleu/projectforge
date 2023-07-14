@ECHO OFF

rem Builds all the templates using quicktemplate

cd %~dp0\..

@ECHO ON
{{{ if .HasModule "database" }}}qtc -ext sql -dir "queries"
{{{ end }}}qtc -ext html -dir "views"
