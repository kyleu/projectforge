package telemetry

import (
	"context"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"{{{ .Package }}}/app/util"
)

var (
	initialized    = false
	tracerProvider *sdktrace.TracerProvider
)

func InitializeIfNeeded(enabled bool, logger *zap.SugaredLogger) bool {
	if initialized {
		return false
	}
	if enabled {
		Initialize(logger)
	}
	return true
}

func Initialize(logger *zap.SugaredLogger) {
	if initialized {
		logger.Warn("double telemetry initialization")
		return
	}
	otel.SetErrorHandler(&ErrHandler{logger: logger})
	initialized = true
	logger.Debug("initializing tracing telemetry")

	endpoint := "localhost:55681"
	if env := os.Getenv(util.AppKey + "_telemetry_endpoint"); env != "" {
		endpoint = env
	}
	tp, err := buildTP(endpoint)
	if err != nil {
		logger.Errorf("unable to create tracing provider: %+v", err)
		return
	}
	tracerProvider = tp
	p := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(p)
}

func buildTP(endpoint string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint(endpoint), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(util.AppKey))),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func Close() error {
	return tracerProvider.Shutdown(context.Background())
}

func StartSpan(ctx context.Context, tracerKey string, spanName string) (context.Context, trace.Span) {
	tr := otel.GetTracerProvider().Tracer(tracerKey)
	return tr.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindServer))
}

type ErrHandler struct {
	logger     *zap.SugaredLogger
	hasPrinted bool
}

func (e *ErrHandler) Handle(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "Post \"") {
		if e.hasPrinted {
			return
		}
		if idx := strings.Index(msg, "\":"); idx > -1 {
			msg = strings.TrimSpace(msg[idx + 2:])
		}
		e.logger.Warn("telemetry seems to be unavailable: [" + msg + "] (this message will appear only once)")
		e.hasPrinted = true
		return
	}
	e.logger.Warnf("telemetry error: %+v", err)
}
