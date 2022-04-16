package httpmetrics

import (
	"context"

	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

var _ propagation.TextMapCarrier = (*headerCarrier)(nil)

type headerCarrier struct {
	h *fasthttp.RequestHeader
}

func (hc headerCarrier) Get(key string) string {
	return string(hc.h.Peek(key))
}

func (hc headerCarrier) Set(key string, value string) {
	hc.h.Set(key, value)
}

func (hc headerCarrier) Keys() []string {
	var keys []string
	hc.h.VisitAll(func(key []byte, _ []byte) {
		keys = append(keys, string(key))
	})
	return keys
}

func ExtractHeaders(rc *fasthttp.RequestCtx, logger *zap.SugaredLogger) context.Context {
	return otel.GetTextMapPropagator().Extract(rc, headerCarrier{h: &rc.Request.Header})
}
