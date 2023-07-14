@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs eslint for the TypeScript project

cd %~dp0\..\client

echo "=== linting client ==="
@ECHO ON
eslint --ext .js,.ts,.tsx .
