@ECHO OFF
rem Content managed by Project Forge, see [projectforge.md] for details.

rem Downloads and installs the Go libraries and tools needed in other scripts

@ECHO ON
go install github.com/cosmtrek/air@latest
go install github.com/valyala/quicktemplate/qtc@latest
go install gotest.tools/gotestsum@latest
go install mvdan.cc/gofumpt@latest
