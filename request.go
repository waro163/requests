package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type HTTPOptions struct {
	Headers http.Header
	Params  *url.Values
	Body    interface{}
}

func (cli *Client) ResolveURL(path string, params *url.Values) *url.URL {
	rel := &url.URL{Path: path}
	if params != nil {
		rel.RawQuery = params.Encode()
	}
	if cli.Options.BaseURL != nil {
		return cli.Options.BaseURL.ResolveReference(rel)
	}
	return rel
}

func (cli *Client) BuildRequestHeader(h http.Header) http.Header {
	header := cli.Options.DefaultHeaders.Clone()
	if header == nil {
		header = make(http.Header)
	}
	if len(h) == 0 {
		return header
	}
	for key, values := range h {
		for _, value := range values {
			header.Add(key, value)
		}
	}
	return header
}

func (cli *Client) NewRequest(ctx context.Context, method, path string, httpOptions *HTTPOptions) (*http.Request, error) {
	// build url
	url := cli.ResolveURL(path, httpOptions.Params)
	// build body reader
	reader, err := BuildReaderFromBody(httpOptions.Body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url.String(), reader)
	if err != nil {
		return nil, err
	}
	// build request header
	headers := cli.BuildRequestHeader(httpOptions.Headers)
	req.Header = headers
	return req, nil
}

func BuildReaderFromBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	var byteData []byte
	switch data := body.(type) {
	case []byte:
		byteData = data
	default:
		var err error
		byteData, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	reader := bytes.NewReader(byteData)
	return reader, nil
}
