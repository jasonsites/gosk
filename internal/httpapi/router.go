package httpapi

import (
	"compress/gzip"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goddtriffin/helmet"
	"github.com/jasonsites/gosk-api/internal/domain"
	ctrl "github.com/jasonsites/gosk-api/internal/httpapi/controllers"
	mw "github.com/jasonsites/gosk-api/internal/httpapi/middleware"
	"github.com/jasonsites/gosk-api/internal/httpapi/routes"
	"github.com/jasonsites/gosk-api/internal/types"
)

type controllers struct {
	ExampleController *ctrl.Controller
}

// configureMiddleware
func configureMiddleware(r *chi.Mux, ns string, logger *types.Logger) {

	skipHealth := func(r *http.Request) bool {
		return r.URL.Path == fmt.Sprintf("/%s/health", ns)
	}

	r.Use(middleware.Compress(gzip.DefaultCompression))
	r.Use(mw.Correlation(&mw.CorrelationConfig{Next: skipHealth}))
	r.Use(mw.ResponseLogger(&mw.ResponseLoggerConfig{Logger: logger, Next: skipHealth}))
	r.Use(helmet.Default().Secure)
	// r.Use(chimw.RealIP)
	r.Use(mw.RequestLogger(&mw.RequestLoggerConfig{Logger: logger, Next: skipHealth}))
	r.Use(middleware.Recoverer)
	// r.Use(cors.Handler(cors.Options{
	// 	// AllowedOrigins:   []string{"http://*", "https://*"},
	// 	// AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
	// 	// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	// ExposedHeaders:   []string{"Link"},
	// 	// AllowCredentials: true,
	// 	MaxAge: 300,
	// }))
}

// registerControllers
func registerControllers(services *domain.Services, logger *types.Logger) *controllers {
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

// ------------------------------------------------------------------------------------------------

// // errorHandler provides custom error handling (end of chain middleware) for all routes
// func errorHandler(ctx *fiber.Ctx, err error) error {
// 	composeErrorResponse := func(code int, title, detail string) types.ErrorResponse {
// 		return types.ErrorResponse{
// 			Errors: []types.ErrorData{{
// 				Status: code,
// 				Title:  title,
// 				Detail: detail,
// 			}},
// 		}
// 	}

// 	var (
// 		code     = http.StatusInternalServerError // default error status code (500)
// 		detail   = "internal server error"
// 		title    string
// 		fiberErr *fiber.Error
// 		response types.ErrorResponse
// 	)

// 	cerr, ok := err.(*types.CustomError)
// 	// custom (controlled) errors
// 	if ok {
// 		code = types.HTTPStatusCodeMap[cerr.Type]
// 		detail = cerr.Message
// 		title = cerr.Type

// 		// fiber errors
// 	} else if errors.As(err, &fiberErr) {
// 		code = fiberErr.Code
// 		detail = fiberErr.Message

// 		// unknown errors
// 	} else {
// 		title = types.ErrorType.InternalServer
// 	}

// 	response = composeErrorResponse(code, title, detail)

// 	ctx.Status(code)
// 	if err := ctx.JSON(response); err != nil {
// 		return ctx.SendString(detail)
// 	}

// 	return nil
// }
