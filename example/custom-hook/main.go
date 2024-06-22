package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/waro163/requests"
	"github.com/waro163/requests/example/custom-hook/hook"
	"github.com/waro163/requests/hooks"
)

func main() {
	// add global hook for client
	// requests.UseGlobalHook(&hooks.Logger{})
	requests.UseLastGlobalHook(&hooks.Logger{})
	url := "http://localhost:8081/ping?a=12&c=ab"
	// url := "https://baidu.com"
	cli := requests.NewClient()

	// add custom hook for client
	h := hook.HeaderHook{Key: "x-app", Value: "requests-demo"}
	cli.Use(&h)

	// new reuqest
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(`{"name":"davis","app_id":19}`))
	if err != nil {
		fmt.Println("new request error, ", err)
		return
	}
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
