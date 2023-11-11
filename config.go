package requests

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	defaultRequestHeaders = http.Header{
		"Content-Type": []string{"application/json"},
	}
	defaultRequestTimeout = 30 * time.Second
	defaultTransport      = &http.Transport{
		MaxIdleConns:        300,
		MaxIdleConnsPerHost: 30,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
)

// Options passed into hooks
type Options struct {
	DefaultHeaders http.Header
	BaseURL        *url.URL
	Name           string
}

type config struct {
	Options
	timeout   *time.Duration
	transport http.RoundTripper
}

// Option defines the function signature to set the optional configuration properties
type Option func(c *config)

func WithTimeout(t time.Duration) Option {
	return func(c *config) {
		c.timeout = &t
	}
}

func WithTransport(t http.RoundTripper) Option {
	return func(c *config) {
		c.transport = t
	}
}

func WithDefaultRequestHeaders(headers http.Header) Option {
	return func(c *config) {
		c.DefaultHeaders = headers
	}
}

func WithBaseURL(base *url.URL) Option {
	return func(c *config) {
		c.BaseURL = base
	}
}

func WithName(name string) Option {
	return func(c *config) {
		c.Name = name
	}
}

func (cli *Client) WithName(name string) {
	cli.Options.Name = name
}

func (cli *Client) WithBaseURL(base *url.URL) {
	cli.Options.BaseURL = base
}

func (cli *Client) WithTimeout(t time.Duration) {
	cli.RawClient.Timeout = t
}

func (cli *Client) WithDefaultRequestHeaders(headers http.Header) {
	cli.Options.DefaultHeaders = headers
}

func (cli *Client) WithTransport(t http.RoundTripper) {
	cli.RawClient.Transport = t
}
