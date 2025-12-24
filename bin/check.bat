@ECHO OFF

rem Runs code statistics, checks for outdated dependencies, then runs linters

cd %~dp0\..

echo "=== linting ==="
@ECHO ON
golangci-lint-v2 run --max-issues-per-linter=0 ./...
