package http_util

import (
	"bytes"
	"net/http"
	"time"
)

func PostAsJSON(url string, data []byte, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	setHeader(headers, req)

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetAsJSON(url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	setHeader(headers, req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func setHeader(headers map[string]string, req *http.Request) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
