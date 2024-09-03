package xhttp

import (
	"bytes"
	"context"
	"errors"

	"github.com/go-resty/resty/v2"
)

var (
	// ErrPostOptionsRequired defines the error returned when post options is nil
	ErrPostOptionsRequired = errors.New("post options required")
	// ErrPostURLRequired defines the error returned when post url is empty
	ErrPostURLRequired = errors.New("post url required")
)

// GetOptions get options
type GetOptions struct {
	Context     context.Context
	QueryParams map[string]string
	Headers     map[string]string
}

// Get send get request with options
func Get(client *resty.Client, url string, options *GetOptions) (*resty.Response, error) {
	httpReq := client.R()

	if options != nil {
		if options.QueryParams != nil {
			httpReq.SetQueryParams(options.QueryParams)
		}

		if options.Headers != nil {
			httpReq.SetHeaders(options.Headers)
		}

		if options.Context != nil {
			httpReq.SetContext(options.Context)
		}
	}

	return httpReq.Get(url)
}

// PostOptions request options
type PostOptions struct {
	Context context.Context
	Url     string
	Data    []byte
	Json    bool
	Headers map[string]string
}

// Post send post request with options
func Post(client *resty.Client, options *PostOptions) (*resty.Response, error) {
	if options == nil {
		return nil, ErrPostOptionsRequired
	}

	if options.Url == "" {
		return nil, ErrPostURLRequired
	}

	httpReq := client.R()

	if options.Json {
		httpReq.SetHeader("Content-Type", "application/json")
	}

	if options.Headers != nil {
		httpReq.SetHeaders(options.Headers)
	}

	if options.Data != nil && len(options.Data) > 0 {
		httpReq.SetBody(bytes.NewReader(options.Data))
	}

	if options.Context != nil {
		httpReq.SetContext(options.Context)
	}

	return httpReq.Post(options.Url)
}
