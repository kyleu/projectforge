@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Formatting code from all projects

cd %~dp0\..

echo "=== formatting ==="
@ECHO ON
gofumpt -w $(find . -type f -name "*.go" | grep -v \\./data | grep -v \\./assets/module | grep -v \\./module | grep -v \\./testproject | grep -v .html.go | grep -v .sql.go)
