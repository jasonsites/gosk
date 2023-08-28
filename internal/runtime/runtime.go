package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/jasonsites/gosk-api/internal/resolver"
)

type Runtime struct {
	Config *resolver.Config
}

func NewRuntime(c *resolver.Config) *Runtime {
	return &Runtime{Config: c}
}

type RunConfig struct {
	HTTPServer bool
}

// run creates a new resolver with associated context group,
// then runs goroutines for serving http requests and graceful app shutdown
func (rt *Runtime) Run(conf *RunConfig) *resolver.Resolver {
	c, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(c)
	r := resolver.NewResolver(ctx, rt.Config)

	// initialize the app resolver and start the http server
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

	// gracefully shut down the http server and close the db connection pool
	g.Go(func() error {
		<-ctx.Done()

		log.Info().Msg("shutdown initiated")

		// shutdown server
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
