@ECHO OFF

rem Starts the app. It doesn't reload on Windows

cd %~dp0\..

set {{{ .CleanKey }}}_encryption_key=TEMP_SECRET_KEY
set GOEXPERIMENT=jsonv2

echo "Windows doesn't allow reloading... sorry"
@ECHO ON
go.exe build -gcflags "all=-N -l" -o build/debug/{{{ .Exec }}}.exe .
build\debug\{{{ .Exec }}}.exe -v --addr=0.0.0.0 {{{ if .HasModule "marketing" }}}all{{{ else }}}server{{{ end }}}
