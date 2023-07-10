@ECHO OFF

rem Runs goreleaser

cd %~dpnx0\..\..

[[ -f "$HOME/bin/oauth" ]] && . $HOME/bin/oauth

goreleaser -f ./tools/release/.goreleaser.yml release --timeout 240m --clean
