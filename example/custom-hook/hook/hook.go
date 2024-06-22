package hook

import (
	"context"
	"net/http"

	"github.com/waro163/requests/hooks"
)

type HeaderHook struct {
	hooks.NoOps
	Key   string
	Value string
}

func (h *HeaderHook) PrepareRequest(ctx context.Context, req *http.Request) error {
	req.Header.Add(h.Key, h.Value)
	return nil
}
