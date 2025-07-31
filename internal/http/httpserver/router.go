package httpserver

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/goddtriffin/helmet"
	mw "github.com/jasonsites/gosk/internal/http/middleware"
	"github.com/jasonsites/gosk/internal/logger"
	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
	"github.com/jasonsites/gosk/internal/modules/example"
	"github.com/jasonsites/gosk/internal/modules/health"
)

type ControllerRegistry struct {
	ExampleController example.ExampleController
}

type RouterConfig struct {
	Namespace    string              `validate:"required"`
	QueryHandler *query.QueryHandler `validate:"required"`
}

// configureMiddleware
func configureMiddleware(conf *RouterConfig, r *chi.Mux, logger *logger.CustomLogger) {
	skipHealth := func(r *http.Request) bool {
		return r.URL.Path == fmt.Sprintf("/%s/health", conf.Namespace)
	}

	r.Use(middleware.Compress(gzip.DefaultCompression))
	r.Use(mw.Correlation(&mw.CorrelationConfig{Next: skipHealth}))
	r.Use(mw.ResponseLogger(&mw.ResponseLoggerConfig{Logger: logger, Next: skipHealth}))
	r.Use(helmet.Default().Secure)
	r.Use(mw.RequestLogger(&mw.RequestLoggerConfig{Logger: logger, QueryHandler: conf.QueryHandler, Next: skipHealth}))
	r.Use(mw.NotFound)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// registerRoutes
func registerRoutes(conf *RouterConfig, r *chi.Mux, c *ControllerRegistry) {
	ns := conf.Namespace
	BaseRouter(r, ns)
	health.HealthRouter(r, ns)
	example.ExampleRouter(r, ns, c.ExampleController)
}
