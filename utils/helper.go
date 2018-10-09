package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func MultipartReq(url string, params map[string]io.Reader) (*http.Response, error) {
	client := &http.Client{}
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var err error
	for key, r := range params {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if _, ok := r.(*strings.Reader); ok {
			if fw, err = w.CreateFormField(key); err != nil {
				return nil, err
			}
		} else {
			if fw, err = w.CreateFormFile(key, key); err != nil {
				return nil, err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return nil, err
		}
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, url, &b)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return client.Do(req)
}
