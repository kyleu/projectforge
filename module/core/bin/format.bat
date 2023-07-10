@ECHO OFF

rem Formatting code from all projects

cd %~dpnx0\..

@ECHO ON
echo "=== formatting ==="
gofumpt -w $(find . -type f -name "*.go"{{{ .IgnoreGrep}}} | grep -v .html.go | grep -v .sql.go)
