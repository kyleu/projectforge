@ECHO OFF

rem Formatting code from all projects

cd %~dp0\..

echo "=== formatting ==="
@ECHO ON
gofumpt -w $(find . -type f -name "*.go" | grep -v \\./data | grep -v \\./module | grep -v \\./testproject | grep -v \\./assets/module | grep -v .html.go | grep -v .sql.go)
