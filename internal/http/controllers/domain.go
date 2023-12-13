package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/core/cerror"
	"github.com/jasonsites/gosk/internal/core/interfaces"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/core/trace"
	"github.com/jasonsites/gosk/internal/core/validation"
	"github.com/jasonsites/gosk/internal/http/jsonio"
)

// Config defines the input to NewController
type Config struct {
	Logger      *logger.CustomLogger      `validate:"required"`
	QueryConfig *QueryConfig              `validate:"required"`
	Service     interfaces.ExampleService `validate:"required"`
}

// ExampleController
type ExampleController struct {
	logger  *logger.CustomLogger
	query   *queryHandler
	service interfaces.ExampleService
}

// NewExampleController returns a new ExampleController instance
func NewExampleController(c *Config) (*ExampleController, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	queryHandler, err := NewQueryHandler(c.QueryConfig)
	if err != nil {
		return nil, err
	}
	ctrl := &ExampleController{
		logger:  c.Logger,
		query:   queryHandler,
		service: c.Service,
	}

	return ctrl, nil
}

// Create
func (c *ExampleController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		body := &jsonapi.RequestBody[models.ExampleRequestData]{
			Data: &jsonapi.RequestResource[models.ExampleRequestData]{
				Attributes: &models.ExampleRequestData{},
			},
		}
		if err := jsonio.DecodeRequest(w, r, body); err != nil {
			err = cerror.NewValidationError(err, "invalid request body")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		model, err := c.service.Create(ctx, body.Data.Attributes)
		if err != nil {
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusCreated, model.FormatDetailResponse())
	}
}

// Delete
func (c *ExampleController) Delete() http.HandlerFunc {
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
			jsonio.EncodeError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Detail
func (c *ExampleController) Detail() http.HandlerFunc {
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

		jsonio.EncodeResponse(w, r, http.StatusOK, model.FormatDetailResponse())
	}
}

// List
func (c *ExampleController) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		qs := []byte(r.URL.RawQuery)
		query := c.query.parseQuery(qs)

		model, err := c.service.List(ctx, *query)
		if err != nil {
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, model.FormatListResponse())
	}
}

// Update
func (c *ExampleController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			err = cerror.NewInternalServerError(err, "resource id parse error")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		body := &jsonapi.RequestBody[models.ExampleRequestData]{
			Data: &jsonapi.RequestResource[models.ExampleRequestData]{
				Attributes: &models.ExampleRequestData{},
			},
		}
		if err := jsonio.DecodeRequest(w, r, body); err != nil {
			err = cerror.NewValidationError(err, "invalid request body")
			log.Error(err.Error())
			jsonio.EncodeError(w, r, err)
			return
		}

		data := body.Data.Attributes
		model, err := c.service.Update(ctx, data, uuid)
		if err != nil {
			jsonio.EncodeError(w, r, err)
			return
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, model.FormatDetailResponse())
	}
}
