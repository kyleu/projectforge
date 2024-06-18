@ECHO OFF

rem Builds all the templates using quicktemplate

cd %~dp0\..

@ECHO ON
qtc -ext html -dir "views"
