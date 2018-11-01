package funcs

import (
	"bytes"
	"mime/multipart"
	"path/filepath"
	"net/http"
	"os"
	"io"
)

type ProgressReader struct {
	io.Reader
	Reporter func(r int64)
}

type ProgressWriter struct {
	io.Writer
	Reporter func(r int64)
}

func (pr *ProgressWriter) Write(p []byte) (n int, err error) {
	n, err = pr.Writer.Write(p)
	pr.Reporter(int64(n))
	return
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Reporter(int64(n))
	return
}

func CreateFileUploadRequest(uri string, params map[string]string, paramName, path string, report func(r int64)) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, &ProgressReader{body, report})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
