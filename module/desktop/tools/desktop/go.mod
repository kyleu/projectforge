// $PF_GENERATE_ONCE$
module {{{ .Package }}}/tools/desktop

go {{{ .GoMajorVersionSafe }}}

require (
	{{{ .Package }}} v0.0.0
	github.com/webview/webview v0.0.0-20210330151455-f540d88dde4e
)

replace {{{ .Package }}} => ../..
