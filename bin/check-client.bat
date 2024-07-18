@ECHO OFF

rem Runs eslint for the TypeScript project

cd %~dp0\..\client

echo "=== linting client ==="
@ECHO ON
eslint .
