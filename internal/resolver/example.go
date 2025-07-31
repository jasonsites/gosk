package resolver

import (
	"fmt"
	"log/slog"

	"github.com/jasonsites/gosk/internal/logger"
	"github.com/jasonsites/gosk/internal/modules/example"
)

// ExampleController provides a singleton example.exampleController instance
func (r *Resolver) ExampleController() example.ExampleController {
	if r.exampleController == nil {
		c := r.Config()

		log := r.Log().With(slog.String("tags", "controller,example"))
		cLogger := &logger.CustomLogger{
			Level: c.Logger.Level,
			Log:   log,
		}

		ctrlConfig := &example.ControllerConfig{
			Logger:  cLogger,
			Query:   r.QueryHandler(),
			Service: r.ExampleService(),
		}
		ctrl, err := example.NewController(ctrlConfig)
		if err != nil {
			err = fmt.Errorf("example controller load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.exampleController = ctrl
	}

	return r.exampleController
}

// ExampleRepository provides a singleton example.exampleRepository instance
func (r *Resolver) ExampleRepository() example.ExampleRepository {
	if r.exampleRepo == nil {
		c := r.Config()

		log := r.Log().With(slog.String("tags", "repo,example"))
		cLogger := &logger.CustomLogger{
			Level: c.Logger.Level,
			Log:   log,
		}
		repoConfig := &example.ExampleRepoConfig{
			DBClient: r.PostgreSQLClient(),
			Logger:   cLogger,
		}

		repo, err := example.NewExampleRepository(repoConfig)
		if err != nil {
			err = fmt.Errorf("example respository load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.exampleRepo = repo
	}

	return r.exampleRepo
}

// ExampleService provides a singleton example.exampleService instance
func (r *Resolver) ExampleService() example.ExampleService {
	if r.exampleService == nil {
		c := r.Config()

		log := r.Log().With(slog.String("tags", "service,example"))
		cLogger := &logger.CustomLogger{
			Level: c.Logger.Level,
			Log:   log,
		}
		svcConfig := &example.ExampleServiceConfig{
			Logger: cLogger,
			Repo:   r.ExampleRepository(),
		}

		svc, err := example.NewExampleService(svcConfig)
		if err != nil {
			err = fmt.Errorf("example service load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.exampleService = svc
	}

	return r.exampleService
}
