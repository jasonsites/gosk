package httpserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	app "github.com/jasonsites/gosk/internal/app"
	"github.com/jasonsites/gosk/internal/logger"
)

// ServerConfig defines the input to NewServer
type ServerConfig struct {
	Controllers  *ControllerRegistry  `validate:"required"`
	Host         string               `validate:"required"`
	Logger       *logger.CustomLogger `validate:"required"`
	Port         uint                 `validate:"required"`
	RouterConfig *RouterConfig        `validate:"required"`
}

// Server defines a server for handling HTTP API requests
type Server struct {
	Logger *logger.CustomLogger
	Port   uint
	Server *http.Server
}

// NewServer returns a new Server instance
func NewServer(c *ServerConfig) (*Server, error) {
	if err := app.Validator.Validate.Struct(c); err != nil {
		return nil, err
	}

	mux := chi.NewRouter()
	configureMiddleware(c.RouterConfig, mux, c.Logger)
	registerRoutes(c.RouterConfig, mux, c.Controllers)

	addr := fmt.Sprintf(":%s", strconv.FormatUint(uint64(c.Port), 10))
	s := &Server{
		Logger: c.Logger,
		Port:   c.Port,
		Server: &http.Server{Addr: addr, Handler: mux},
	}

	return s, nil
}

// Serve starts the HTTP server on the configured address
func (s *Server) Serve() error {
	s.Logger.Log.Info(fmt.Sprintf("server listening on port :%d", s.Port))
	return s.Server.ListenAndServe()
}
