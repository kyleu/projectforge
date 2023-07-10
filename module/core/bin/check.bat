@ECHO OFF

rem Runs code statistics, checks for outdated dependencies, then runs linters

cd %~dpnx0\..

echo "=== linting ==="
golangci-lint run --fix --max-issues-per-linter=0 --sort-results{{{ .IgnoredSetting }}} ./...
