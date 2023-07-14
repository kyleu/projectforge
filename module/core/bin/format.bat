@ECHO OFF

rem Formatting code from all projects

cd %~dp0\..

echo "=== formatting ==="
@ECHO ON
gofumpt -w $(find . -type f -name "*.go"{{{ .IgnoreGrep}}} | grep -v .html.go | grep -v .sql.go)
