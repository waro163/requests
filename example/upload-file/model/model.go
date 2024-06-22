package model

import "io"

type MultipartFormField struct {
	FieldName  string
	FieldValue []byte
	IsFile     bool
	FileName   string
	FileReader io.Reader
}
