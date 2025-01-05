package httpmetrics

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"projectforge.dev/projectforge/app/util"
)

var _ propagation.TextMapCarrier = (*headerCarrier)(nil)

type headerCarrier struct {
	h http.Header
}

func (hc headerCarrier) Get(key string) string {
	return hc.h.Get(key)
}

func (hc headerCarrier) Set(key string, value string) {
	hc.h.Set(key, value)
}

func (hc headerCarrier) Keys() []string {
	return util.MapKeys(hc.h)
}

func ExtractHeaders(r *http.Request, logger util.Logger) (context.Context, util.Logger) {
	return otel.GetTextMapPropagator().Extract(r.Context(), headerCarrier{h: r.Header}), logger
}
