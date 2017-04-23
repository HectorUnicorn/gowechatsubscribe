package uploader

import (
	"net/http"
	"os"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"fmt"
	"github.com/going/toolkit/log"
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
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

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func DoUploadFile2Wechat(source string, accessToken string) string {
	extraParams := map[string]string{
		"access_token": accessToken,
		"type": "image",
	}
	path, _ := os.Getwd()
	fmt.Println("current path:", path)
	request, err := newfileUploadRequest("https://api.weixin.qq.com/cgi-bin/media/upload", extraParams, "media", source)
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
		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body)
		return fmt.Sprintf("%v", body)
	}
	return "no response";
}
