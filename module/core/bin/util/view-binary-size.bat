@ECHO OFF

rem Visualizes space usage for the release binary

cd %~dp0\..\..

@ECHO ON
make build-release
go tool nm -size build\release\{{{ .Exec }}}
