// Package telemetry - Content managed by Project Forge, see [projectforge.md] for details.
package telemetry

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func WrapHTTPClient(cl *http.Client) *http.Client {
	cl.Transport = otelhttp.NewTransport(cl.Transport)
	return cl
}
