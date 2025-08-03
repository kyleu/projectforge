@ECHO OFF

rem Formatting code from the TypeScript projects

cd %~dp0\..
cd client

echo "=== formatting ==="
@ECHO ON
npm run format
