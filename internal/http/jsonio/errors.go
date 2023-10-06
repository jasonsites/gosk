package jsonio

import (
	"net/http"

	"github.com/jasonsites/gosk/internal/core/cerror"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
)

// HTTPStatusCodeMap maps custom error types to relevant HTTP status codes
var HTTPStatusCodeMap = map[string]int{
	cerror.ErrorType.Conflict:       http.StatusConflict,
	cerror.ErrorType.Forbidden:      http.StatusForbidden,
	cerror.ErrorType.InternalServer: http.StatusInternalServerError,
	cerror.ErrorType.NotFound:       http.StatusNotFound,
	cerror.ErrorType.Unauthorized:   http.StatusUnauthorized,
	cerror.ErrorType.Validation:     http.StatusBadRequest,
}

func EncodeError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		code     = http.StatusInternalServerError
		detail   = "internal server error"
		title    string
		response jsonapi.ErrorResponse
	)

	switch e := err.(type) {
	case *cerror.CustomError:
		code = HTTPStatusCodeMap[e.Type]
		title = e.Type
		if e.Type != cerror.ErrorType.InternalServer {
			detail = e.Message
		}
	default:
		title = cerror.ErrorType.InternalServer
	}

	response = composeErrorResponse(code, title, detail)
	EncodeResponse(w, r, code, response)
}

func composeErrorResponse(code int, title, detail string) jsonapi.ErrorResponse {
	return jsonapi.ErrorResponse{
		Errors: []jsonapi.ErrorData{{
			Status: code,
			Title:  title,
			Detail: detail,
		}},
	}
}
