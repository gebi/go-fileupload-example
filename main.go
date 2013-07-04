package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Streams upload directly from file -> mime/multipart -> pipe -> http-request
func asyncUploadFile(params map[string]string, paramName, path string, w *io.PipeWriter, file *os.File) {
	defer file.Close()
	defer w.Close()
	writer := multipart.NewWriter(w)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	go asyncUploadFile(params, paramName, path, w, file)
	return http.NewRequest("POST", uri, r)
}

func main() {
	path, _ := os.Getwd()
	path += "/test.pdf"
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest("https://google.com/upload", extraParams, "file", "/tmp/doc.pdf")
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
	}
}
