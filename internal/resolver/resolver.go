package resolver

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk/config"
	app "github.com/jasonsites/gosk/internal/app"
	"github.com/jasonsites/gosk/internal/http/httpserver"
	query "github.com/jasonsites/gosk/internal/modules/common/models/query"
	"github.com/jasonsites/gosk/internal/modules/example"
)

type ResolverEntry string

const (
	Unset ResolverEntry = ""
	HTTP  ResolverEntry = "http"
	gRPC  ResolverEntry = "grpc"
)

// Config defines the input to NewResolver
type Config struct {
	Config            *config.Configuration
	ExampleController example.ExampleController
	ExampleRepo       example.ExampleRepository
	ExampleService    example.ExampleService
	HTTPServer        *httpserver.Server
	Log               *slog.Logger
	Metadata          *app.Metadata
	PostgreSQLClient  *pgxpool.Pool
	QueryHandler      *query.QueryHandler
}

// Resolver provides a configurable app component graph
type Resolver struct {
	appContext        context.Context
	config            *config.Configuration
	exampleController example.ExampleController
	exampleRepo       example.ExampleRepository
	exampleService    example.ExampleService
	httpServer        *httpserver.Server
	log               *slog.Logger
	metadata          *app.Metadata
	postgreSQLClient  *pgxpool.Pool
	queryHandler      *query.QueryHandler
}

// NewResolver returns a new Resolver instance
func NewResolver(ctx context.Context, c *Config) *Resolver {
	if c == nil {
		c = &Config{}
	}

	r := &Resolver{
		appContext:        ctx,
		config:            c.Config,
		exampleController: c.ExampleController,
		exampleRepo:       c.ExampleRepo,
		exampleService:    c.ExampleService,
		httpServer:        c.HTTPServer,
		log:               c.Log,
		metadata:          c.Metadata,
		postgreSQLClient:  c.PostgreSQLClient,
		queryHandler:      c.QueryHandler,
	}

	return r
}

// Load resolves app components starting from the given entry node of the component graph
func (r *Resolver) Load(entry ResolverEntry) {
	switch entry {
	case "http":
		r.HTTPServer()
	default:
		panic(fmt.Errorf("invalid resolver load entry point '%s'", entry))
	}
}
