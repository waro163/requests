package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/waro163/requests"
	"github.com/waro163/requests/hooks"
)

func main() {
	// add global hook for client
	requests.UseLastGlobalHook(&hooks.Logger{})

	// init client
	baseUrl, _ := url.Parse("http://localhost:8081")
	cli := requests.NewClient(
		requests.WithName("demo"),
		requests.WithDefaultRequestHeaders(
			http.Header{
				"x-app":        []string{"demo"},
				"Content-Type": []string{"application/json"},
			},
		),
		requests.WithBaseURL(baseUrl),
		// requests.WithTimeout(10*time.Second),
	)
	// also can add options for client after new client
	cli.WithTimeout(10 * time.Second)

	// new request by client
	v := url.Values{}
	v.Add("a", "12")
	req, err := cli.NewRequest(context.Background(), http.MethodGet, "/ping", &requests.HTTPOptions{
		Headers: http.Header{"custom-header": []string{"my-header"}},
		Params:  &v,
		Body:    []byte(`{"name":"davis","app_id":19}`),
	})
	if err != nil {
		fmt.Println("new request error, ", err)
		return
	}

	// do request by client
	res, err := cli.DoRequest(context.Background(), req)
	if err != nil {
		fmt.Println("do request error, ", err)
		return
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body), err)
}
