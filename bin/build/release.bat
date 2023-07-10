@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Runs goreleaser

cd %~dpnx0\..\..

@ECHO ON
goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
