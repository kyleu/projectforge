@ECHO OFF

rem Runs goreleaser

cd %~dp0\..\..

@ECHO ON
goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
