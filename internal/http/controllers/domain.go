package controllers

import (
	"fmt"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/core/cerror"
	"github.com/jasonsites/gosk-api/internal/core/interfaces"
	"github.com/jasonsites/gosk-api/internal/core/logger"
	"github.com/jasonsites/gosk-api/internal/core/query"
	"github.com/jasonsites/gosk-api/internal/core/trace"
)

func init() {
	// registers a directive named "path" to retrieve values from chi.URLParam
	httpin.UseGochiURLParam("path", chi.URLParam)
}

// Config defines the input to NewController
type Config struct {
	Service interfaces.Service
	Logger  *logger.Logger
}

// Controller
type Controller struct {
	service interfaces.Service
	logger  *logger.Logger
}

// NewController returns a new Controller instance
func NewController(c *Config) *Controller {
	return &Controller{
		service: c.Service,
		logger:  c.Logger,
	}
}

// Create
func (c *Controller) Create(f func() *RequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		resource := f()
		if err := c.JSONDecode(w, r, resource); err != nil {
			err = cerror.NewInternalServerError("error parsing request body")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		// TODO: validation errors bypass default error handler (rethink this)
		if response := validateBody(resource, log); response != nil {
			err := fmt.Errorf("validation error(s) %+v", response)
			log.Error().Err(err).Send()
			c.JSONEncode(w, r, http.StatusBadRequest, response)
			return
		}

		data := resource.Data.Attributes
		model, err := c.service.Create(ctx, data)
		if err != nil {
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError("error formatting response from model")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		c.JSONEncode(w, r, http.StatusCreated, response)
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
			err = cerror.NewInternalServerError("error parsing resource id")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		if err := c.service.Delete(ctx, uuid); err != nil {
			log.Error().Err(err).Send()
			c.Error(w, r, err)
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
			err = cerror.NewInternalServerError("error parsing resource id")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		model, err := c.service.Detail(ctx, uuid)
		if err != nil {
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError("error formatting response from model")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		c.JSONEncode(w, r, http.StatusOK, response)
	}
}

// List
func (c *Controller) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		// TEMP ----------------------------------------------
		// qs := ctx.Request().URI().QueryString()
		// query := parseQuery(qs)
		var (
			defaultLimit  = 20           // TODO: move to config
			defaultOffset = 0            // TODO: move to config
			defaultOrder  = "desc"       // TODO: move to config
			defaultProp   = "created_on" // TODO: move to config
		)

		model, err := c.service.List(ctx, query.QueryData{
			Paging: query.QueryPaging{
				Limit:  &defaultLimit,
				Offset: &defaultOffset,
			},
			Sorting: query.QuerySorting{
				Order: &defaultOrder,
				Prop:  &defaultProp,
			},
		}) // *query
		// END TEMP -------------------------------------------
		if err != nil {
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError("error formatting response from model")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		c.JSONEncode(w, r, http.StatusOK, response)
	}
}

// Update
func (c *Controller) Update(f func() *RequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := trace.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			err = cerror.NewInternalServerError("error parsing resource id")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		resource := f()
		if err := c.JSONDecode(w, r, resource); err != nil {
			err = cerror.NewInternalServerError("error parsing request body")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		// TODO: validation errors bypass default error handler
		if response := validateBody(resource, log); response != nil {
			err := fmt.Errorf("validation error(s) %+v", response)
			log.Error().Err(err).Send()
			c.JSONEncode(w, r, http.StatusBadRequest, response)
			return
		}

		data := resource.Data.Attributes // TODO: problem here with ID
		model, err := c.service.Update(ctx, data, uuid)
		if err != nil {
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		response, err := model.FormatResponse()
		if err != nil {
			err = cerror.NewInternalServerError("error formatting response from model")
			log.Error().Err(err).Send()
			c.Error(w, r, err)
			return
		}

		c.JSONEncode(w, r, http.StatusCreated, response)
	}
}
