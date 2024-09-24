package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const userAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

type FetchConfig struct {
	Method     string
	Url        string
	Data       interface{}
	RetryTimes int
	Headers    map[string]string
	Cookies    map[string]string
}

type FetchResponse struct {
	StatusCode int
	Data       []byte
	Header     http.Header
	Cookie     []*http.Cookie
}

func Fetch(cfg FetchConfig) (*FetchResponse, error) {
	var req *http.Request
	var resp *http.Response
	var backoff time.Duration = 1
	var reqBody io.Reader = nil
	var err error

	if cfg.RetryTimes == 0 {
		cfg.RetryTimes = 3
	}
	if cfg.Method == "" {
		cfg.Method = "GET"
	}

	if cfg.Data != nil && (cfg.Method == "POST" || cfg.Method == "PUT") {
		if cfg.Headers == nil {
			cfg.Headers = map[string]string{}
		}
		if _, exist := cfg.Headers["Content-Type"]; !exist {
			cfg.Headers["Content-Type"] = "application/json"
		}
		v := cfg.Headers["Content-Type"]
		if v == "application/x-www-form-urlencoded" {
			switch data := cfg.Data.(type) {
			case map[string]string:
				params := url.Values{}
				for k, v := range data {
					params.Set(k, v)
				}
				reqBody = strings.NewReader(params.Encode())
			case string:
				reqBody = strings.NewReader(data)
			case url.Values:
				reqBody = strings.NewReader(data.Encode())
			default:
				return nil, errors.New("unsupported data type")
			}
		} else if v == "application/json" {
			var jsonData []byte
			jsonData, err = json.Marshal(cfg.Data)
			if err != nil {
				return nil, err
			}
			reqBody = bytes.NewReader(jsonData)
		} else {
			reqBody = strings.NewReader(cfg.Data.(string))
		}
	}

	for retryTime := 0; retryTime <= cfg.RetryTimes; retryTime++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err = http.NewRequestWithContext(ctx, cfg.Method, cfg.Url, reqBody)
		if err != nil {
			return nil, err
		}
		if cfg.Method == "POST" || cfg.Method == "PUT" {
			req.Header.Set("Content-Type", "application/json")
		}
		if v, exist := cfg.Headers["User-Agent"]; exist {
			if v != "" {
				req.Header.Set("User-Agent", v)
			}
		} else {
			req.Header.Set("User-Agent", userAgent)
		}
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}
		for k, v := range cfg.Cookies {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			if isRetryableError(err) {
				err = errors.New("request error: " + err.Error())
				time.Sleep(backoff * time.Second)
				backoff *= 2
				continue
			}
		}

		if resp == nil {
			return nil, errors.New("response is nil")
		}

		if isRetryableStatusCode(resp.StatusCode) {
			err = errors.New("response status code: " + resp.Status)
			time.Sleep(backoff * time.Second)
			backoff *= 2
			continue
		}

		contentType := resp.Header.Get("Content-Type")
		var reader io.Reader
		if strings.Contains(contentType, "charset=") {
			reader, err = charset.NewReader(resp.Body, contentType)
		} else {
			reader = resp.Body
		}
		if err != nil {
			return nil, err
		}
		dataBytes, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		res := &FetchResponse{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Cookie:     resp.Cookies(),
			Data:       dataBytes,
		}

		return res, nil
	}
	return nil, err
}

func isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

func isRetryableError(err error) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}
	}
	return false
}
