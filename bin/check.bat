@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs code statistics, checks for outdated dependencies, then runs linters

cd %~dp0\..

echo "=== linting ==="
@ECHO ON
golangci-lint run --fix --max-issues-per-linter=0 --sort-results --skip-dirs "/data|/module|/testproject" ./...
