@ECHO OFF

rem Runs eslint for the TypeScript project

cd %~dpnx0\..\client

@ECHO ON
echo "=== linting client ==="
eslint --ext .js,.ts,.tsx .
