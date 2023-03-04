// $PF_GENERATE_ONCE$
module {{{ .Package }}}

go {{{ .GoVersionSafe }}}

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/dustin/go-humanize v1.0.1
	github.com/fasthttp/router v1.4.17
	github.com/gertd/go-pluralize v0.2.1
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/json-iterator/go v1.1.12
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kirsle/configdir v0.0.0-20170128060238-e45d2f54772f
	github.com/muesli/coral v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.14.0
	github.com/valyala/fasthttp v1.44.0
	github.com/valyala/quicktemplate v1.7.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.40.0
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.14.0
	go.opentelemetry.io/otel/sdk v1.14.0
	go.opentelemetry.io/otel/trace v1.14.0
	go.uber.org/zap v1.24.0
	golang.org/x/exp v0.0.0-20230304125523-9ff063c70017
	gopkg.in/yaml.v2 v2.4.0
)
