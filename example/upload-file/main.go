package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/waro163/requests"
	"github.com/waro163/requests/example/upload-file/model"
)

func main() {
	// read a file or get file from remote
	content, err := os.ReadFile("demo.jpg")
	if err != nil {
		fmt.Println("read file error: ", err)
		return
	}
	fileBuf := bytes.NewBuffer(content)
	fields := []model.MultipartFormField{
		{FieldName: "name", FieldValue: []byte("my-file")},
		{FieldName: "file", IsFile: true, FileName: "demo.jpeg", FileReader: fileBuf},
	}

	// build requst body
	bodyBuf := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, err := bodyWriter.CreateFormFile(field.FieldName, field.FileName)
			if err != nil {
				return
			}
			if _, err = io.Copy(fileWriter, field.FileReader); err != nil {
				return
			}
		} else {
			// bodyWriter.WriteField(field.FieldName, string(field.FieldValue))
			partWriter, err := bodyWriter.CreateFormField(field.FieldName)
			if err != nil {
				return
			}
			if _, err := partWriter.Write(field.FieldValue); err != nil {
				return
			}
		}
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	header := http.Header{"Content-Type": []string{contentType}}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost:8080/api/upload/demo", &bodyBuf)
	if err != nil {
		return
	}
	req.Header = header

	// new client
	cli := requests.NewClient()
	res, err := cli.DoRequest(context.Background(), req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	fmt.Println(res.StatusCode)
	body, err := io.ReadAll(res.Body)
	fmt.Println(string(body), err)
}
