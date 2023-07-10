@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Visualizes space usage for the release binary

cd %~dpnx0\..\..

make build-release
go tool nm -size build\release\projectforge
