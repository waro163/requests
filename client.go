package requests

import (
	"context"
	"net/http"
)

type Client struct {
	Options   Options
	RawClient *http.Client

	hooks []Hook
}

func NewClient(options ...Option) *Client {
	config := &config{}
	for _, option := range options {
		option(config)
	}
	if config.timeout == nil {
		config.timeout = &defaultRequestTimeout
	}
	if config.transport == nil {
		config.transport = defaultTransport
	}
	if config.Options.DefaultHeaders == nil {
		config.Options.DefaultHeaders = defaultRequestHeaders
	}

	return &Client{
		Options: config.Options,
		RawClient: &http.Client{
			Timeout:   *config.timeout,
			Transport: config.transport,
		},
	}
}

func (cli *Client) Use(hook Hook) {
	cli.hooks = append(cli.hooks, hook)
}

func (cli *Client) prepareRequest(ctx context.Context, req *http.Request) error {
	for _, globalHook := range globalHooks {
		if err := globalHook.PrepareRequest(ctx, req); err != nil {
			if err := globalHook.OnRequestError(ctx, req, err); err != nil {
				return err
			}
		}
	}
	for _, hook := range cli.hooks {
		if err := hook.PrepareRequest(ctx, req); err != nil {
			if err := hook.OnRequestError(ctx, req, err); err != nil {
				return err
			}
		}
	}
	for _, globalHook := range lastGlobalHooks {
		if err := globalHook.PrepareRequest(ctx, req); err != nil {
			if err := globalHook.OnRequestError(ctx, req, err); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cli *Client) processResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	for i := len(lastGlobalHooks) - 1; i >= 0; i-- {
		if err := lastGlobalHooks[i].ProcessResponse(ctx, req, resp); err != nil {
			if err := lastGlobalHooks[i].OnResponseError(ctx, req, resp, err); err != nil {
				return err
			}
		}
	}
	for i := len(cli.hooks) - 1; i >= 0; i-- {
		if err := cli.hooks[i].ProcessResponse(ctx, req, resp); err != nil {
			if err := cli.hooks[i].OnResponseError(ctx, req, resp, err); err != nil {
				return err
			}
		}
	}
	for i := len(globalHooks) - 1; i >= 0; i-- {
		if err := globalHooks[i].ProcessResponse(ctx, req, resp); err != nil {
			if err := globalHooks[i].OnResponseError(ctx, req, resp, err); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cli *Client) Do(req *http.Request) (resp *http.Response, err error) {
	return cli.DoRequest(req.Context(), req)
}

func (cli *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	if err := cli.prepareRequest(ctx, req); err != nil {
		return nil, err
	}
	resp, err := cli.RawClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := cli.processResponse(ctx, req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (cli *Client) Get(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodGet, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}

func (cli *Client) Post(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodPost, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}

func (cli *Client) Put(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodPut, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}

func (cli *Client) Patch(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodPatch, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}

func (cli *Client) Delete(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodDelete, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}

func (cli *Client) Head(ctx context.Context, path string, opts *HTTPOptions) (*http.Response, error) {
	req, err := cli.NewRequest(ctx, http.MethodHead, path, opts)
	if err != nil {
		return nil, err
	}
	return cli.DoRequest(ctx, req)
}
