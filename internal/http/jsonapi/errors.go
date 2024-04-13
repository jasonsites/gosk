package jsonapi

// ErrorResponse
type ErrorResponse struct {
	Errors []ErrorData `json:"errors"`
}

// ErrorData
type ErrorData struct {
	Source *ErrorSourcePointer `json:"source,omitempty"`
	Title  string              `json:"title"`
	Detail string              `json:"detail"`
}

// // ErrorSourcePointer
// type ErrorSourceHeader struct {
// 	Header string `json:"header,omitempty"`
// }

// // ErrorSourceParameter
// type ErrorSourceParameter struct {
// 	Parameter string `json:"parameter,omitempty"`
// }

// ErrorSourcePointer
type ErrorSourcePointer struct {
	Pointer string `json:"pointer,omitempty"`
}
