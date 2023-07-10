@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Builds all the templates using quicktemplate

cd %~dpnx0\..

@ECHO ON
qtc -ext html -dir "views"
