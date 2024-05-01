package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/waro163/requests"
	"github.com/waro163/requests/hooks"
)

func main() {
	requests.UseGlobalHook(&hooks.Logger{})
	// url := "http://localhost:8081/test/input?a=12&c=ab"
	url := "https://google.com"
	cli := requests.NewClient()
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(`{"name":"davis","app_id":19}`))
	if err != nil {
		fmt.Println("new request error, ", err)
		return
	}
	req.Header.Add("x-app", "requests-demo")
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
