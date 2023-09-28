package resolver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk/config"
	"github.com/jasonsites/gosk/internal/core/interfaces"
	"github.com/jasonsites/gosk/internal/core/logger"
	"github.com/jasonsites/gosk/internal/core/query"
	"github.com/jasonsites/gosk/internal/core/validation"
	"github.com/jasonsites/gosk/internal/domain"
	"github.com/jasonsites/gosk/internal/http/controllers"
	"github.com/jasonsites/gosk/internal/http/httpserver"
	"github.com/jasonsites/gosk/internal/repos"
)

// Config provides a singleton config.Configuration instance
func (r *Resolver) Config() *config.Configuration {
	if r.config == nil {
		conf, err := config.LoadConfiguration()
		if err != nil {
			err = fmt.Errorf("config load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.config = conf
	}

	return r.config
}

// Domain provides a singleton domain.Domain instance
func (r *Resolver) Domain() *domain.Domain {
	if r.domain == nil {
		services := &domain.Services{
			Example: r.ExampleService(),
		}

		app, err := domain.NewDomain(services)
		if err != nil {
			err = fmt.Errorf("domain load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.domain = app
	}

	return r.domain
}

// ExampleRepository provides a singleton repo.exampleRepository instance
func (r *Resolver) ExampleRepository() interfaces.ExampleRepository {
	if r.exampleRepo == nil {
		repo, err := repos.NewExampleRepository(&repos.ExampleRepoConfig{
			DBClient: r.PostgreSQLClient(),
			Logger: &logger.Logger{
				Enabled: r.Config().Logger.Repo.Enabled,
				Level:   r.Config().Logger.Repo.Level,
				Log:     r.Log(),
			},
		})
		if err != nil {
			err = fmt.Errorf("example respository load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.exampleRepo = repo
	}

	return r.exampleRepo
}

// ExampleService provides a singleton domain.exampleService instance
func (r *Resolver) ExampleService() interfaces.Service {
	if r.exampleService == nil {
		svc, err := domain.NewExampleService(&domain.ExampleServiceConfig{
			Logger: &logger.Logger{
				Enabled: r.Config().Logger.Domain.Enabled,
				Level:   r.Config().Logger.Domain.Level,
				Log:     r.Log(),
			},
			Repo: r.ExampleRepository(),
		})
		if err != nil {
			err = fmt.Errorf("example service load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.exampleService = svc
	}

	return r.exampleService
}

// HTTPServer provides a singleton httpserver.Server instance
func (r *Resolver) HTTPServer() *httpserver.Server {
	if r.httpServer == nil {
		c := r.Config()

		queryConfig := func() *controllers.QueryConfig {
			limit := int(c.HTTP.Router.Paging.DefaultLimit)
			offset := int(c.HTTP.Router.Paging.DefaultOffset)

			attr := c.HTTP.Router.Sorting.DefaultAttr
			order := c.HTTP.Router.Sorting.DefaultOrder

			return &controllers.QueryConfig{
				Defaults: &controllers.QueryDefaults{
					Paging: &query.QueryPaging{
						Limit:  &limit,
						Offset: &offset,
					},
					Sorting: &query.QuerySorting{
						Attr:  &attr,
						Order: &order,
					},
				},
			}
		}()

		routerConfig := &httpserver.RouterConfig{Namespace: c.HTTP.Router.Namespace}

		server, err := httpserver.NewServer(&httpserver.ServerConfig{
			Domain: r.Domain(),
			Host:   c.HTTP.Server.Host,
			Logger: &logger.Logger{
				Enabled: c.Logger.HTTP.Enabled,
				Level:   c.Logger.HTTP.Level,
				Log:     r.Log(),
			},
			Mode:         c.HTTP.Server.Mode,
			Port:         c.HTTP.Server.Port,
			QueryConfig:  queryConfig,
			RouterConfig: routerConfig,
		})
		if err != nil {
			err = fmt.Errorf("http server load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.httpServer = server
	}

	return r.httpServer
}

// Log provides a singleton slog.Logger instance
func (r *Resolver) Log() *slog.Logger {
	if r.log == nil {
		logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}).WithAttrs([]slog.Attr{
			slog.Int("pid", os.Getpid()),
			slog.String("name", r.Metadata().Name),
			slog.String("version", r.Metadata().Version),
		})
		logger := slog.New(logHandler)
		slog.SetDefault(logger)

		r.log = logger
	}

	return r.log
}

// Metadata provides a singleton application Metadata instance
func (r *Resolver) Metadata() *Metadata {
	if r.metadata == nil {
		var metadata *Metadata

		jsondata, err := os.ReadFile(r.config.Metadata.Path)
		if err != nil {
			err = fmt.Errorf("package.json read error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		if err := json.Unmarshal(jsondata, &metadata); err != nil {
			err = fmt.Errorf("package.json unmarshall error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.metadata = metadata
	}

	return r.metadata
}

// PostgreSQLClient provides a singleton postgres pgxpool.Pool instance
func (r *Resolver) PostgreSQLClient() *pgxpool.Pool {
	if r.postgreSQLClient == nil {
		if err := validation.Validate.StructPartial(r.config, "Postgres"); err != nil {
			err = fmt.Errorf("invalid postgres config: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		client, err := pgxpool.New(r.appContext, postgresDSN(r.config.Postgres))
		if err != nil {
			err = fmt.Errorf("postgres client load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.postgreSQLClient = client
	}

	return r.postgreSQLClient
}
