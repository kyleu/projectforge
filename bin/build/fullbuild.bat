@ECHO OFF

rem Builds the app (or just use make build)

cd %~dp0\..\..

echo "Building release build, set GOOS/GOARCH to change target..."
@ECHO ON

cd client
call npm install
cd ..
call ./bin/build/client.bat

call ./bin/templates.bat
go mod tidy
call ./bin/build/build.bat
