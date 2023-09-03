package resolver

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk-api/config"
	"github.com/jasonsites/gosk-api/internal/domain"
	"github.com/jasonsites/gosk-api/internal/httpapi"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/rs/zerolog"
)

// Config defines the input to NewResolver
type Config struct {
	Config           *config.Configuration
	Domain           *domain.Domain
	ExampleRepo      types.ExampleRepository
	ExampleService   types.Service
	HTTPServer       *httpapi.Server
	Log              *zerolog.Logger
	Metadata         *Metadata
	PostgreSQLClient *pgxpool.Pool
}

// Application metadata
type Metadata struct {
	Name    string
	Version string
}

// Resolver provides singleton instances of app components
type Resolver struct {
	appContext       context.Context
	config           *config.Configuration
	domain           *domain.Domain
	exampleRepo      types.ExampleRepository
	exampleService   types.Service
	httpServer       *httpapi.Server
	log              *zerolog.Logger
	metadata         *Metadata
	postgreSQLClient *pgxpool.Pool
}

// NewResolver returns a new Resolver instance
func NewResolver(ctx context.Context, c *Config) *Resolver {
	if c == nil {
		c = &Config{}
	}

	r := &Resolver{
		appContext:       ctx,
		config:           c.Config,
		domain:           c.Domain,
		exampleRepo:      c.ExampleRepo,
		exampleService:   c.ExampleService,
		httpServer:       c.HTTPServer,
		log:              c.Log,
		metadata:         c.Metadata,
		postgreSQLClient: c.PostgreSQLClient,
	}

	return r
}

// initialize bootstraps the application in dependency order
func (r *Resolver) Initialize() error {
	if _, err := r.Config(); err != nil {
		return err
	}
	if _, err := r.Metadata(); err != nil {
		return err
	}
	if _, err := r.Log(); err != nil {
		return err
	}
	if _, err := r.PostgreSQLClient(); err != nil {
		return err
	}
	if _, err := r.ExampleRepository(); err != nil {
		return err
	}
	if _, err := r.ExampleService(); err != nil {
		return err
	}
	if _, err := r.Domain(); err != nil {
		return err
	}
	if _, err := r.HTTPServer(); err != nil {
		return err
	}

	return nil
}
