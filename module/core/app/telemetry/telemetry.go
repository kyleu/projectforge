package telemetry

import (
	"context"

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
	initialized = true
	logger.Debug("initializing tracing telemetry")

	tp, err := buildTP("localhost:55681")
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
