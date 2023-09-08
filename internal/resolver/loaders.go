package resolver

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk-api/config"
	"github.com/jasonsites/gosk-api/internal/domain"
	"github.com/jasonsites/gosk-api/internal/httpapi"
	repo "github.com/jasonsites/gosk-api/internal/repository"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config provides a singleton config.Configuration instance
func (r *Resolver) Config() *config.Configuration {
	if r.config == nil {
		c, err := config.LoadConfiguration()
		if err != nil {
			err = fmt.Errorf("error resolving config: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}
		r.config = c
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
			err = fmt.Errorf("error resolving domain: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}

		r.domain = app
	}

	return r.domain
}

// ExampleRepository provides a singleton repo.exampleRepository instance
func (r *Resolver) ExampleRepository() types.ExampleRepository {
	if r.exampleRepo == nil {
		repo, err := repo.NewExampleRepository(&repo.ExampleRepoConfig{
			DBClient: r.PostgreSQLClient(),
			Logger: &types.Logger{
				Enabled: r.Config().Logger.Repo.Enabled,
				Level:   r.Config().Logger.Repo.Level,
				Log:     r.Log(),
			},
		})
		if err != nil {
			err = fmt.Errorf("error resolving example respository: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}

		r.exampleRepo = repo
	}

	return r.exampleRepo
}

// ExampleService provides a singleton domain.exampleService instance
func (r *Resolver) ExampleService() types.Service {
	if r.exampleService == nil {
		svc, err := domain.NewExampleService(&domain.ExampleServiceConfig{
			Logger: &types.Logger{
				Enabled: r.Config().Logger.Domain.Enabled,
				Level:   r.Config().Logger.Domain.Level,
				Log:     r.Log(),
			},
			Repo: r.ExampleRepository(),
		})
		if err != nil {
			err = fmt.Errorf("error resolving example service: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}

		r.exampleService = svc
	}

	return r.exampleService
}

// HTTPServer provides a singleton httpapi.HTTPServer instance
func (r *Resolver) HTTPServer() *httpapi.HTTPServer {
	if r.httpServer == nil {
		c := r.Config()
		server, err := httpapi.NewServer(&httpapi.HTTPServerConfig{
			BaseURL: c.HttpAPI.BaseURL,
			Domain:  r.Domain(),
			Logger: &types.Logger{
				Enabled: c.Logger.Http.Enabled,
				Level:   c.Logger.Http.Level,
				Log:     r.Log(),
			},
			Mode:      c.HttpAPI.Mode,
			Namespace: c.HttpAPI.Namespace,
			Port:      c.HttpAPI.Port,
		})
		if err != nil {
			err = fmt.Errorf("error resolving http server: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}
		r.httpServer = server
	}

	return r.httpServer
}

// Log provides a singleton zerolog.Logger instance
func (r *Resolver) Log() *zerolog.Logger {
	if r.log == nil {
		logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().
			Int("pid", os.Getpid()).
			Str("name", r.Metadata().Name).
			Str("version", r.Metadata().Version).
			Timestamp().Logger()

		r.log = &logger
	}

	return r.log
}

// Metadata provides a singleton application Metadata instance
func (r *Resolver) Metadata() *Metadata {
	if r.metadata == nil {
		var metadata *Metadata

		jsondata, err := os.ReadFile(r.config.Metadata.Path)
		if err != nil {
			err = fmt.Errorf("error reading package.json: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}

		if err := json.Unmarshal(jsondata, &metadata); err != nil {
			err = fmt.Errorf("error unmarshalling package.json: %w", err)
			log.Error().Err(err).Send()
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
			log.Error().Err(err).Send()
			panic(err)
		}

		client, err := pgxpool.New(r.appContext, postgresDSN(r.config.Postgres))
		if err != nil {
			err = fmt.Errorf("error resolving postgres client: %w", err)
			log.Error().Err(err).Send()
			panic(err)
		}

		r.postgreSQLClient = client
	}

	return r.postgreSQLClient
}
