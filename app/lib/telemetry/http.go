package telemetry

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func WrapHTTPClient(cl *http.Client) *http.Client {
	cl.Transport = otelhttp.NewTransport(cl.Transport)
	return cl
}

func HTTPClient() *http.Client {
	return WrapHTTPClient(http.DefaultClient)
}
