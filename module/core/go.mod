// $PF_IGNORE$
module {{{ .Package }}}

go 1.17

require (
	github.com/alecthomas/chroma v0.9.2
	github.com/dustin/go-humanize v1.0.0
	github.com/fasthttp/router v1.4.2
	github.com/gertd/go-pluralize v0.1.7
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/markbates/going v1.0.3
	github.com/markbates/goth v1.68.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cobra v1.2.1
	github.com/valyala/fasthttp v1.30.0
	github.com/valyala/quicktemplate v1.6.3
	go.opentelemetry.io/otel v1.0.0-RC3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.0.0-RC3
	go.opentelemetry.io/otel/sdk v1.0.0-RC3
	go.opentelemetry.io/otel/trace v1.0.0-RC3
	go.uber.org/zap v1.19.0{{{ if .HasModule "mobile" }}}
	golang.org/x/mobile v0.0.0-20210902104108-5d9a33257ab5{{{ end }}}
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f
	gopkg.in/yaml.v2 v2.4.0
)
