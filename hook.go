package requests

import (
	"context"
	"net/http"
)

var globalHooks []Hook

type Hook interface {
	// PrepareRequest handle this request before do requesting
	PrepareRequest(context.Context, *http.Request) error
	// OnRequestError is for when PrepareRequest return error, it will handle this error, if this function also return error, all the work will stop
	OnRequestError(context.Context, *http.Request, error) error
	// ProcessResponse handle this request and response after do requesting
	ProcessResponse(context.Context, *http.Request, *http.Response) error
	// OnRequestError is for when ProcessResponse return error, it will handle this error, if this function also return error, all the work will stop
	OnResponseError(context.Context, *http.Request, *http.Response, error) error
}
