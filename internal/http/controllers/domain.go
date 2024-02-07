package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/app"
	"github.com/jasonsites/gosk/internal/core/cerror"
	"github.com/jasonsites/gosk/internal/core/interfaces"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/trace"
	"github.com/jasonsites/gosk/internal/http/jsonio"
)

// Config defines the input to NewController
type Config struct {
	Logger      *logger.CustomLogger      `validate:"required"`
	QueryConfig *QueryConfig              `validate:"required"`
	Service     interfaces.ExampleService `validate:"required"`
}

// Controller
type Controller struct {
	logger  *logger.CustomLogger
	query   *queryHandler
	service interfaces.ExampleService
}

// NewController returns a new Controller instance
func NewController(c *Config) (*Controller, error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	queryHandler, err := NewQueryHandler(c.QueryConfig)
	if err != nil {
		return nil, err
	}
	ctrl := &Controller{
		logger:  c.Logger,
		query:   queryHandler,
		service: c.Service,
	}

	return ctrl, nil
}

// Create
func (c *Controller) Create(f func() *jsonapi.RequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		resource := f()
		if err := jsonio.DecodeRequest(w, r, resource); err != nil {
			err = cerror.NewValidationError(err, "request body decode error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		data := resource.Data.Attributes
		model, err := c.service.Create(ctx, data)
		if err != nil {
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError(err, "model format response error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusCreated, response)
	}
}

// Delete
func (c *Controller) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			err = cerror.NewInternalServerError(err, "error parsing resource id")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		if err := c.service.Delete(ctx, uuid); err != nil {
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Detail
func (c *Controller) Detail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			err = cerror.NewInternalServerError(err, "error parsing resource id")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		model, err := c.service.Detail(ctx, uuid)
		if err != nil {
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError(err, "error formatting response from model")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, response)
	}
}

// List
func (c *Controller) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		qs := []byte(r.URL.RawQuery)
		query := c.query.parseQuery(qs)

		model, err := c.service.List(ctx, *query)
		if err != nil {
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError(err, "error formatting response from model")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, response)
	}
}

// Update
func (c *Controller) Update(f func() *jsonapi.RequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			err = cerror.NewValidationError(err, "resource id parse error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		resource := f()
		if err := jsonio.DecodeRequest(w, r, resource); err != nil {
			err = cerror.NewValidationError(err, "request body decode error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		data := resource.Data.Attributes // TODO: problem here with ID
		model, err := c.service.Update(ctx, data, uuid)
		if err != nil {
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError(err, "model format response error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, response)
	}
}
