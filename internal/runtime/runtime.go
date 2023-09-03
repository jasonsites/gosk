package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonsites/gosk-api/internal/resolver"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Runtime struct {
	config *resolver.Config
}

func NewRuntime(c *resolver.Config) *Runtime {
	// TODO: validate config?
	return &Runtime{config: c}
}

type RunConfig struct {
	HTTPServer bool
}

// Run creates a new resolver with associated context group, then runs goroutines for bootstrapping the application and handling graceful shutdown
func (rt *Runtime) Run(conf *RunConfig) *resolver.Resolver {
	c, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(c)
	r := resolver.NewResolver(ctx, rt.config)

	// initialize the app resolver and start the configured applications
	g.Go(func() error {
		log.Info().Msg("initializing resolver")

		if err := r.Initialize(); err != nil {
			return err
		}

		if conf.HTTPServer {
			log.Info().Msg("starting http server")
			server, err := r.HTTPServer()
			if err != nil {
				return err
			}
			if err := server.Serve(); err != nil {
				return err
			}
		}

		return nil
	})

	// gracefully shut down the configured applications and close the db connection pool
	g.Go(func() error {
		<-ctx.Done()

		log.Info().Msg("shutdown initiated")

		if conf.HTTPServer {
			server, err := r.HTTPServer()
			if err != nil {
				return err
			}
			if err := server.App.Shutdown(); err != nil {
				return err
			}
			log.Info().Msg("http server shut down")
		}

		// close db pool
		pool, err := r.PostgreSQLClient()
		if err != nil {
			return err
		}
		pool.Close()
		log.Info().Msg("db connection pool closed")

		log.Info().Msg("shutdown complete")

		return nil
	})

	if err := g.Wait(); err != nil {
		err = errors.Errorf("error running application: %+v", err)
		log.Error().Err(err).Send()
	}

	return r
}
