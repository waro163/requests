package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	strStartTime   = "StartTime"
	strHost        = "Host"
	strUrl         = "Url"
	strMethod      = "Method"
	strQuery       = "Query"
	strHeader      = "Header"
	strReqBody     = "ReqBody"
	strEndTime     = "EndTime"
	strRespBody    = "RespBody"
	strStatusCode  = "StatusCode"
	defaultLogName = "OUTGOING_REQUEST_LOG"
)

type Logger struct{}

func (l *Logger) PrepareRequest(ctx context.Context, req *http.Request) error {
	byteHeaders, _ := json.Marshal(req.Header)
	buffer := bytes.Buffer{}
	if req.Body != nil && req.Body != http.NoBody && req.GetBody != nil {
		if reader, err := req.GetBody(); err == nil {
			buffer.ReadFrom(reader)
		}
	}
	log.Info().
		Str(strStartTime, time.Now().String()).
		Str(strHost, req.Host).
		Str(strMethod, req.Method).
		Str(strUrl, req.URL.Path).
		Str(strQuery, req.URL.RawQuery).
		RawJSON(strHeader, byteHeaders).
		RawJSON(strReqBody, buffer.Bytes()).
		Msg(defaultLogName)
	return nil
}

func (l *Logger) OnRequestError(context.Context, *http.Request, error) error {
	return nil
}

func (l *Logger) ProcessResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	buffer := bytes.NewBuffer(nil)
	if resp.Body != nil && resp.Body != http.NoBody {
		buffer.ReadFrom(resp.Body)
		resp.Body.Close()
		resp.Body = io.NopCloser(buffer)
	}
	log.Info().
		Str(strEndTime, time.Now().String()).
		Str(strHost, req.Host).
		Str(strMethod, req.Method).
		Str(strUrl, req.URL.Path).
		Str(strQuery, req.URL.RawQuery).
		RawJSON(strRespBody, buffer.Bytes()).
		Int(strStatusCode, resp.StatusCode).
		Msg(defaultLogName)
	return nil
}

func (l *Logger) OnResponseError(context.Context, *http.Request, *http.Response, error) error {
	return nil
}
