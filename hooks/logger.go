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
	reader, _ := req.GetBody()
	reqBody, _ := io.ReadAll(reader)
	log.Ctx(ctx).Info().
		Str(strStartTime, time.Now().String()).
		Str(strHost, req.Host).
		Str(strUrl, req.RequestURI).
		RawJSON(strHeader, byteHeaders).
		RawJSON(strReqBody, reqBody).
		Msg(defaultLogName)
	return nil
}

func (l *Logger) OnRequestError(context.Context, *http.Request, error) error {
	return nil
}

func (l *Logger) ProcessResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	byteBody, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(byteBody))
	log.Ctx(ctx).Info().
		Str(strEndTime, time.Now().String()).
		Str(strHost, req.Host).
		Str(strUrl, req.RequestURI).
		RawJSON(strRespBody, byteBody).
		Int(strStatusCode, resp.StatusCode).
		Msg(defaultLogName)
	return nil
}

func (l *Logger) OnResponseError(context.Context, *http.Request, *http.Response, error) error {
	return nil
}
