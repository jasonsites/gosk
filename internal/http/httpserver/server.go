package httpserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk-api/internal/core/logger"
	"github.com/jasonsites/gosk-api/internal/core/validation"
	"github.com/jasonsites/gosk-api/internal/domain"
)

// ServerConfig defines the input to NewServer
type ServerConfig struct {
	BaseURL   string         `validate:"required"`
	Domain    *domain.Domain `validate:"required"`
	Logger    *logger.Logger `validate:"required"`
	Mode      string         `validate:"required"`
	Namespace string         `validate:"required"`
	Port      uint           `validate:"required"`
}

// Server defines a server for handling HTTP API requests
type Server struct {
	Logger *logger.Logger
	Port   uint
	Server *http.Server
}

// NewServer returns a new Server instance
func NewServer(c *ServerConfig) (*Server, error) {
	if err := validation.Validate.Struct(c); err != nil {
		return nil, err
	}

	log := c.Logger.Log.With().Str("tags", "http").Logger()
	logger := &logger.Logger{
		Enabled: c.Logger.Enabled,
		Level:   c.Logger.Level,
		Log:     &log,
	}

	controllers := registerControllers(c.Domain.Services, logger)
	router := chi.NewRouter()
	configureMiddleware(router, c.Namespace, logger)
	registerRoutes(router, controllers, c.Namespace)

	addr := fmt.Sprintf(":%s", strconv.FormatUint(uint64(c.Port), 10))

	s := &Server{
		Logger: logger,
		Port:   c.Port,
		Server: &http.Server{Addr: addr, Handler: router},
	}

	return s, nil
}

// Serve starts the HTTP server on the configured address
func (s *Server) Serve() error {
	s.Logger.Log.Info().Msg(fmt.Sprintf("server listening on port :%d", s.Port))
	return s.Server.ListenAndServe()
}
