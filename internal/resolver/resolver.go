package resolver

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk-api/config"
	"github.com/jasonsites/gosk-api/internal/domain"
	"github.com/jasonsites/gosk-api/internal/httpapi"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/rs/zerolog"
)

// Application metadata
type Metadata struct {
	Name    string
	Version string
}

// Config defines the input to NewResolver
type Config struct {
	Config           *config.Configuration
	Domain           *domain.Domain
	ExampleRepo      types.ExampleRepository
	ExampleService   types.Service
	HTTPServer       *httpapi.HTTPServer
	Log              *zerolog.Logger
	Metadata         *Metadata
	PostgreSQLClient *pgxpool.Pool
}

// Resolver provides a configurable app component graph
type Resolver struct {
	appContext       context.Context
	config           *config.Configuration
	domain           *domain.Domain
	exampleRepo      types.ExampleRepository
	exampleService   types.Service
	httpServer       *httpapi.HTTPServer
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

// LoadEntries provides option strings for loading the resolver from various entry nodes
// in the app component graph (cli, grpc, http)
var LoadEntries = struct{ HTTPServer string }{
	HTTPServer: "http",
}

// Load resolves app components starting from the given entry node of the component graph
func (r *Resolver) Load(entry string) {
	switch entry {
	case LoadEntries.HTTPServer:
		r.HTTPServer()
	default:
		panic(fmt.Errorf("invalid resolver load entry point '%s'", entry))
	}
}
