@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs eslint for the TypeScript project

cd %~dpnx0\..\client

@ECHO ON
echo "=== linting client ==="
eslint --ext .js,.ts,.tsx .
