module {{{ .Package }}}/tools/desktop

go 1.22

require (
	{{{ .Package }}} v0.0.0
	github.com/webview/webview v0.0.0-20210330151455-f540d88dde4e
)

replace {{{ .Package }}} => ../..
