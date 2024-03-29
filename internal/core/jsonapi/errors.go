package jsonapi

// ErrorResponse
type ErrorResponse struct {
	Errors []ErrorData `json:"errors"`
}

// ErrorData
type ErrorData struct {
	Status int          `json:"status"`
	Source *ErrorSource `json:"source,omitempty"`
	Title  string       `json:"title"`
	Detail string       `json:"detail"`
}

// ErrorSource
type ErrorSource struct {
	// TODO: add header and query options
	Pointer string `json:"pointer,omitempty"`
}
