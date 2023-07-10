@ECHO OFF

rem Builds all the templates using quicktemplate

cd %~dpnx0\..

{{{ if .HasModule "database" }}}qtc -ext sql -dir "queries"
{{{ end }}}qtc -ext html -dir "views"
