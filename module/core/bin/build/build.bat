@ECHO OFF

rem Builds the app (or just use make build)

cd %~dp0\..\..

echo "Building release build, set GOOS/GOARCH to change target..."
@ECHO ON
make build-release
