package requests

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	OtelTransport = otelhttp.NewTransport(defaultTransport)
)
