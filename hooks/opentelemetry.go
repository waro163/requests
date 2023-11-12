package hooks

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	goPackageName = "github.com/waro163/requests"
)

type Opentelemetry struct {
	Name string // this should be app name
}

func (o *Opentelemetry) PrepareRequest(c context.Context, req *http.Request) error {
	tracer := otel.Tracer(goPackageName)
	ctx, _ := tracer.Start(
		req.Context(),
		fmt.Sprintf("http.%s %s", strings.ToLower(req.Method), req.Host),
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithTimestamp(time.Now()),
		trace.WithAttributes(
			semconv.ServiceName(o.Name),
			semconv.HTTPMethod(req.Method),
			semconv.HTTPURL(req.URL.String()),
		),
	)
	propagators := otel.GetTextMapPropagator()
	propagators.Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	*req = *(req.WithContext(ctx))
	return nil
}

func (o *Opentelemetry) OnRequestError(context.Context, *http.Request, error) error {
	return nil
}

func (o *Opentelemetry) ProcessResponse(c context.Context, req *http.Request, resp *http.Response) error {
	span := trace.SpanFromContext(req.Context())
	defer span.End()
	if resp != nil {
		span.SetStatus(codes.Ok, resp.Status)
		span.SetAttributes(
			semconv.HTTPStatusCode(resp.StatusCode),
			semconv.HTTPResponseContentLength(int(resp.ContentLength)),
		)
	} else {
		span.SetStatus(codes.Unset, codes.Unset.String())
	}
	return nil
}

func (o *Opentelemetry) OnResponseError(c context.Context, req *http.Request, resp *http.Response, err error) error {
	span := trace.SpanFromContext(req.Context())
	span.RecordError(err)
	span.SetStatus(codes.Error, codes.Error.String())
	return nil
}
