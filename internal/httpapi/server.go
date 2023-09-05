package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk-api/internal/domain"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
)

// HTTPServerConfig defines the input to NewServer
type HTTPServerConfig struct {
	BaseURL   string         `validate:"required"`
	Domain    *domain.Domain `validate:"required"`
	Logger    *types.Logger  `validate:"required"`
	Mode      string         `validate:"required"`
	Namespace string         `validate:"required"`
	Port      uint           `validate:"required"`
}

// HTTPServer defines a server for handling HTTP API requests
type HTTPServer struct {
	Logger *types.Logger
	port   uint
	server *http.Server
}

// NewServer returns a new Server instance
func NewServer(c *HTTPServerConfig) (*HTTPServer, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	log := c.Logger.Log.With().Str("tags", "http").Logger()
	logger := &types.Logger{
		Enabled: c.Logger.Enabled,
		Level:   c.Logger.Level,
		Log:     &log,
	}

	controllers := registerControllers(c.Domain.Services, logger)
	router := chi.NewRouter()
	configureMiddleware(router, c.Namespace, logger)
	registerRoutes(router, controllers, c.Namespace)

	addr := fmt.Sprintf(":%s", strconv.FormatUint(uint64(c.Port), 10))

	s := &HTTPServer{
		Logger: logger,
		port:   c.Port,
		server: &http.Server{Addr: addr, Handler: router},
	}

	return s, nil
}

// Serve starts the HTTP server on the configured address
func (s *HTTPServer) Serve() error {
	s.Logger.Log.Info().Msg(fmt.Sprintf("server listening on port :%d", s.port))
	return s.server.ListenAndServe()
}
