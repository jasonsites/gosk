package example

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk/internal/app"
	cerror "github.com/jasonsites/gosk/internal/cerror"
	"github.com/jasonsites/gosk/internal/http/jsonapi"
	"github.com/jasonsites/gosk/internal/http/jsonio"
	"github.com/jasonsites/gosk/internal/http/trace"
	"github.com/jasonsites/gosk/internal/logger"
)

// ExampleService
type ExampleService interface {
	Create(context.Context, any) (*ModelContainer, error)
	Delete(context.Context, uuid.UUID) error
	Detail(context.Context, uuid.UUID) (*ModelContainer, error)
	List(context.Context, ExampleQueryData) (*ModelContainer, error)
	Update(context.Context, any, uuid.UUID) (*ModelContainer, error)
}

// ControllerConfig defines the input to NewController
type ControllerConfig struct {
	Logger  *logger.CustomLogger `validate:"required"`
	Query   *ExampleQueryHandler `validate:"required"`
	Service ExampleService       `validate:"required"`
}

// exampleController
type exampleController struct {
	logger  *logger.CustomLogger
	query   *ExampleQueryHandler
	service ExampleService
}

// NewController returns a new Controller instance
func NewController(c *ControllerConfig) (*exampleController, error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	ctrl := &exampleController{
		logger:  c.Logger,
		query:   c.Query,
		service: c.Service,
	}

	return ctrl, nil
}

// Create
func (c *exampleController) Create(f func() *jsonapi.RequestBody) http.HandlerFunc {
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
func (c *exampleController) Delete() http.HandlerFunc {
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
func (c *exampleController) Detail() http.HandlerFunc {
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
func (c *exampleController) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		qs := []byte(r.URL.RawQuery)
		query := c.query.ParseQuery(qs)

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
func (c *exampleController) Update(f func() *jsonapi.RequestBody) http.HandlerFunc {
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
