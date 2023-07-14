@ECHO OFF

rem Builds the app (or just use make build)

cd %~dp0\..\..

echo "Building release build, set GOOS/GOARCH to change target..."
@ECHO ON

cd client
npm install
cd ..
./bin/build/client.bat

./bin/templates.bat
go mod tidy
./bin/build/build.bat
