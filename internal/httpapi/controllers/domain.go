package controllers

import (
	"fmt"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jasonsites/gosk-api/internal/types"
)

func init() {
	// registers a directive named "path" to retrieve values from chi.URLParam
	httpin.UseGochiURLParam("path", chi.URLParam)
}

// Config
type Config struct {
	Service types.Service
	Logger  *types.Logger
}

// Controller
type Controller struct {
	service types.Service
	logger  *types.Logger
}

// NewController
func NewController(c *Config) *Controller {
	return &Controller{
		service: c.Service,
		logger:  c.Logger,
	}
}

// Create
func (c *Controller) Create(f func() *types.JSONRequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := types.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)
		log.Info().Msg("Create Controller called")

		resource := f()
		if err := c.ReadJSON(w, r, resource); err != nil {
			fmt.Printf("JSON PARSING ERROR: %+v\n", err) // TODO
			message := "error parsing request body"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		// TODO: validation errors bypass default error handler
		if err := validateBody(resource, log); err != nil {
			message := "validation error"
			log.Error().Msg(message)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		data := resource.Data.Attributes
		model, err := c.service.Create(ctx, data)
		if err != nil {
			log.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := model.Serialize()
		if err != nil {
			message := "serialization error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		if err := c.JSONResponse(w, http.StatusCreated, result); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// Delete
func (c *Controller) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := types.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)
		log.Info().Msg("Delete Controller called")

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			message := "error parsing uuid"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		if err := c.service.Delete(ctx, uuid); err != nil {
			message := "example service error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Detail
func (c *Controller) Detail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := types.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)
		log.Info().Msg("Detail Controller called")

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			message := "error parsing uuid"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		model, err := c.service.Detail(ctx, uuid)
		if err != nil {
			message := "example service error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		result, err := model.Serialize()
		if err != nil {
			message := "serialization error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		if err := c.JSONResponse(w, http.StatusOK, result); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// List
func (c *Controller) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := types.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)
		log.Info().Msg("List Controller called")

		// TEMP ----------------------------------------------
		// qs := ctx.Request().URI().QueryString()
		// query := parseQuery(qs)
		var (
			defaultLimit  = 20           // TODO: move to config
			defaultOffset = 0            // TODO: move to config
			defaultOrder  = "desc"       // TODO: move to config
			defaultProp   = "created_on" // TODO: move to config
		)

		model, err := c.service.List(ctx, types.QueryData{
			Paging: types.QueryPaging{
				Limit:  &defaultLimit,
				Offset: &defaultOffset,
			},
			Sorting: types.QuerySorting{
				Order: &defaultOrder,
				Prop:  &defaultProp,
			},
		}) // *query
		// END TEMP -------------------------------------------
		if err != nil {
			message := "example service error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		result, err := model.Serialize()
		if err != nil {
			message := "serialization error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		if err := c.JSONResponse(w, http.StatusOK, result); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

// Update
func (c *Controller) Update(f func() *types.JSONRequestBody) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		traceID := types.GetTraceIDFromContext(ctx)
		log := c.logger.CreateContextLogger(traceID)
		log.Info().Msg("Update Controller called")

		id := chi.URLParam(r, "id")
		uuid, err := uuid.Parse(id)
		if err != nil {
			message := "error parsing uuid"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		resource := f()
		if err := c.ReadJSON(w, r, resource); err != nil {
			fmt.Printf("JSON PARSING ERROR: %+v\n", err) // TODO
			message := "error parsing request body"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		// TODO: validation errors bypass default error handler
		if err := validateBody(resource, log); err != nil {
			message := "validation error"
			log.Error().Msg(message)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		data := resource.Data.Attributes // TODO: problem here with ID
		model, err := c.service.Update(ctx, data, uuid)
		if err != nil {
			log.Error().Err(err).Send()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := model.Serialize()
		if err != nil {
			message := "serialization error"
			log.Error().Err(err).Msg(message)
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		if err := c.JSONResponse(w, http.StatusCreated, result); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}
