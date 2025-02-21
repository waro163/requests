package hooks

import (
	"context"
	"net/http"

	"github.com/waro163/requests"
)

type customHook struct {
	preReq     func(context.Context, *http.Request) error
	onPreReq   func(context.Context, *http.Request, error) error
	procResp   func(context.Context, *http.Request, *http.Response) error
	onProcResp func(context.Context, *http.Request, *http.Response, error) error
}

func NewCustomHook(
	preReq func(context.Context, *http.Request) error,
	onPreReq func(context.Context, *http.Request, error) error,
	procResp func(context.Context, *http.Request, *http.Response) error,
	onProcResp func(context.Context, *http.Request, *http.Response, error) error,
) requests.Hook {
	return &customHook{
		preReq:     preReq,
		onPreReq:   onPreReq,
		procResp:   procResp,
		onProcResp: onProcResp,
	}
}

func (h *customHook) PrepareRequest(ctx context.Context, req *http.Request) error {
	if h.preReq != nil {
		return h.preReq(ctx, req)
	}
	return nil
}

func (h *customHook) OnRequestError(ctx context.Context, req *http.Request, err error) error {
	if h.onPreReq != nil {
		return h.onPreReq(ctx, req, err)
	}
	return nil
}

func (h *customHook) ProcessResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	if h.procResp != nil {
		return h.procResp(ctx, req, resp)
	}
	return nil
}

func (h *customHook) OnResponseError(ctx context.Context, req *http.Request, resp *http.Response, err error) error {
	if h.onProcResp != nil {
		return h.onProcResp(ctx, req, resp, err)
	}
	return nil
}
