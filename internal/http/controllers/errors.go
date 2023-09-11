package controllers

import (
	"net/http"

	"github.com/jasonsites/gosk-api/internal/core/cerrors"
	"github.com/jasonsites/gosk-api/internal/core/trace"
)

// HTTPStatusCodeMap maps application error types to respective HTTP status codes
var HTTPStatusCodeMap = map[string]int{
	cerrors.ErrorType.Conflict:       http.StatusConflict,
	cerrors.ErrorType.Forbidden:      http.StatusForbidden,
	cerrors.ErrorType.InternalServer: http.StatusInternalServerError,
	cerrors.ErrorType.NotFound:       http.StatusNotFound,
	cerrors.ErrorType.Unauthorized:   http.StatusUnauthorized,
	cerrors.ErrorType.Validation:     http.StatusBadRequest,
}

func (c *Controller) Error(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	traceID := trace.GetTraceIDFromContext(ctx)
	log := c.logger.CreateContextLogger(traceID)

	// span := trace.SpanFromContext(r.Context())
	// span.RecordError(err)

	switch err.(type) {
	// case StatusError:
	// 	// retrieve the status here and write out a specific HTTP status code
	// 	log.Error().Err(err).Msg("error")
	// 	http.Error(w, http.StatusText(e.Code), e.Code)
	// 	return
	default:
		// Any error types we don't specifically look out for default to serving an HTTP Internal Server Error
		log.Error().Err(err).Msg("unhandled error")
		// Str("request-id", middlewares.GetReqID(r.Context())).
		// Str("trace.id", span.SpanContext().TraceID().String()).
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
