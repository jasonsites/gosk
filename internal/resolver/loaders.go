package resolver

import (
	"encoding/json"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk-api/config"
	"github.com/jasonsites/gosk-api/internal/domain"
	"github.com/jasonsites/gosk-api/internal/httpapi"
	repo "github.com/jasonsites/gosk-api/internal/repository"
	"github.com/jasonsites/gosk-api/internal/types"
	"github.com/jasonsites/gosk-api/internal/validation"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config provides a singleton config.Configuration instance
func (r *Resolver) Config() (*config.Configuration, error) {
	if r.config == nil {
		c, err := config.LoadConfiguration()
		if err != nil {
			err = errors.Errorf("error resolving config: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}
		r.config = c
	}

	return r.config, nil
}

// Domain provides a singleton domain.Domain instance
func (r *Resolver) Domain() (*domain.Domain, error) {
	if r.domain == nil {
		services := &domain.Services{
			Example: r.exampleService,
		}

		app, err := domain.NewDomain(services)
		if err != nil {
			err = errors.Errorf("error resolving domain: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		r.domain = app
	}

	return r.domain, nil
}

// ExampleRepository provides a singleton repo.exampleRepository instance
func (r *Resolver) ExampleRepository() (types.ExampleRepository, error) {
	if r.exampleRepo == nil {
		repo, err := repo.NewExampleRepository(&repo.ExampleRepoConfig{
			DBClient: r.postgreSQLClient,
			Logger: &types.Logger{
				Enabled: r.config.Logger.Repo.Enabled,
				Level:   r.config.Logger.Repo.Level,
				Log:     r.log,
			},
		})
		if err != nil {
			err = errors.Errorf("error resolving example respository: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		r.exampleRepo = repo
	}

	return r.exampleRepo, nil
}

// ExampleService provides a singleton domain.exampleService instance
func (r *Resolver) ExampleService() (types.Service, error) {
	if r.exampleService == nil {
		svc, err := domain.NewExampleService(&domain.ExampleServiceConfig{
			Logger: &types.Logger{
				Enabled: r.config.Logger.Domain.Enabled,
				Level:   r.config.Logger.Domain.Level,
				Log:     r.log,
			},
			Repo: r.exampleRepo,
		})
		if err != nil {
			err = errors.Errorf("error resolving example service: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		r.exampleService = svc
	}

	return r.exampleService, nil
}

// HTTPServer provides a singleton httpapi.HTTPServer instance
func (r *Resolver) HTTPServer() (*httpapi.HTTPServer, error) {
	if r.httpServer == nil {
		server, err := httpapi.NewServer(&httpapi.HTTPServerConfig{
			BaseURL: r.config.HttpAPI.BaseURL,
			Domain:  r.domain,
			Logger: &types.Logger{
				Enabled: r.config.Logger.Http.Enabled,
				Level:   r.config.Logger.Http.Level,
				Log:     r.log,
			},
			Mode:      r.config.HttpAPI.Mode,
			Namespace: r.config.HttpAPI.Namespace,
			Port:      r.config.HttpAPI.Port,
		})
		if err != nil {
			err = errors.Errorf("error resolving http server: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}
		r.httpServer = server
	}

	return r.httpServer, nil
}

// Log provides a singleton zerolog.Logger instance
func (r *Resolver) Log() (*zerolog.Logger, error) {
	if r.log == nil {
		logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().
			Int("pid", os.Getpid()).
			Str("name", r.metadata.Name).
			Str("version", r.metadata.Version).
			Timestamp().Logger()

		r.log = &logger
	}

	return r.log, nil
}

// Metadata provides a singleton application Metadata instance
func (r *Resolver) Metadata() (*Metadata, error) {
	if r.metadata == nil {
		var metadata *Metadata

		jsondata, err := os.ReadFile(r.config.Metadata.Path)
		if err != nil {
			err = errors.Errorf("error reading package.json: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		if err := json.Unmarshal(jsondata, &metadata); err != nil {
			err = errors.Errorf("error unmarshalling package.json: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		r.metadata = metadata
	}

	return r.metadata, nil
}

// PostgreSQLClient provides a singleton postgres pgxpool.Pool instance
func (r *Resolver) PostgreSQLClient() (*pgxpool.Pool, error) {
	if r.postgreSQLClient == nil {
		if err := validation.Validate.StructPartial(r.config, "Postgres"); err != nil {
			err = errors.Errorf("invalid postgres config: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		client, err := pgxpool.New(r.appContext, postgresDSN(r.config.Postgres))
		if err != nil {
			err = errors.Errorf("error resolving postgres client: %+v", err)
			log.Error().Err(err).Send()
			return nil, err
		}

		r.postgreSQLClient = client
	}

	return r.postgreSQLClient, nil
}
