package hooks

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Opentelemertry struct{}

func (o *Opentelemertry) PrepareRequest(ctx context.Context, req *http.Request) error {
	propagators := otel.GetTextMapPropagator()
	propagators.Inject(req.Context(), propagation.HeaderCarrier(req.Header))
	return nil
}

func (o *Opentelemertry) OnRequestError(context.Context, *http.Request, error) error {
	return nil
}

func (o *Opentelemertry) ProcessResponse(context.Context, *http.Request, *http.Response) error {
	return nil
}

func (o *Opentelemertry) OnResponseError(context.Context, *http.Request, *http.Response, error) error {
	return nil
}
