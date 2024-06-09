package hooks

import (
	"context"
	"net/http"

	"github.com/waro163/requests"
)

type NoOps struct{}

var _ requests.Hook = (*NoOps)(nil)

func (h *NoOps) PrepareRequest(context.Context, *http.Request) error {
	return nil
}

func (h *NoOps) OnRequestError(context.Context, *http.Request, error) error {
	return nil
}

func (h *NoOps) ProcessResponse(context.Context, *http.Request, *http.Response) error {
	return nil
}

func (h *NoOps) OnResponseError(context.Context, *http.Request, *http.Response, error) error {
	return nil
}
