package testutils

import (
	"fmt"
	"io"
	"net/http"
)

// Expected
type Expected struct {
	Code int
}

// RequestData
type RequestData struct {
	Body    io.Reader
	Headers map[string]string
	Method  string
	Route   string
}

// RequestOptions
type RequestOptions struct {
	JSON bool
}

// SetRequestData creates a new HTTP Request instance from the given data
func (r *RequestData) SetRequestData(opts *RequestOptions) (*http.Request, error) {
	if r == nil {
		return nil, fmt.Errorf("request data must be non-nil")
	}
	if r.Body == nil && (r.Method == http.MethodPost || r.Method == http.MethodPut) {
		return nil, fmt.Errorf("request body must be non-nil for %q method", r.Method)
	}
	req, err := http.NewRequest(r.Method, r.Route, r.Body)
	if err != nil {
		return nil, err
	}
	if opts == nil && r.Body != nil {
		opts = &RequestOptions{JSON: true}
	}
	req = r.SetRequestHeaders(req, r.Headers, opts)
	return req, nil
}

// SetRequestHeaders set all headers on the given request
func (r *RequestData) SetRequestHeaders(req *http.Request, headers map[string]string, opts *RequestOptions) *http.Request {
	if opts != nil && opts.JSON {
		req.Header.Add("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}
