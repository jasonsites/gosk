package httpserver

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/goddtriffin/helmet"
	"github.com/jasonsites/gosk-api/internal/core/logger"
	"github.com/jasonsites/gosk-api/internal/domain"
	ctrl "github.com/jasonsites/gosk-api/internal/http/controllers"
	mw "github.com/jasonsites/gosk-api/internal/http/middleware"
	"github.com/jasonsites/gosk-api/internal/http/routes"
)

type controllers struct {
	ExampleController *ctrl.Controller
}

// configureMiddleware
func configureMiddleware(r *chi.Mux, ns string, logger *logger.Logger) {

	skipHealth := func(r *http.Request) bool {
		return r.URL.Path == fmt.Sprintf("/%s/health", ns)
	}

	r.Use(middleware.Compress(gzip.DefaultCompression))
	r.Use(mw.Correlation(&mw.CorrelationConfig{Next: skipHealth}))
	r.Use(mw.ResponseLogger(&mw.ResponseLoggerConfig{Logger: logger, Next: skipHealth}))
	r.Use(helmet.Default().Secure)
	// r.Use(middleware.RealIP)
	r.Use(mw.RequestLogger(&mw.RequestLoggerConfig{Logger: logger, Next: skipHealth}))
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

// registerControllers
func registerControllers(services *domain.Services, logger *logger.Logger) *controllers {
	return &controllers{
		ExampleController: ctrl.NewController(&ctrl.Config{
			Service: services.Example,
			Logger:  logger,
		}),
	}
}

// registerRoutes
func registerRoutes(r *chi.Mux, c *controllers, ns string) {
	routes.BaseRouter(r, nil, ns)
	routes.HealthRouter(r, nil, ns)
	routes.ExampleRouter(r, c.ExampleController, ns)
}
