package resolver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jasonsites/gosk/config"
	app "github.com/jasonsites/gosk/internal/app"
	"github.com/jasonsites/gosk/internal/http/httpserver"
	"github.com/jasonsites/gosk/internal/logger"
	"github.com/jasonsites/gosk/internal/modules/common/query"
)

// Config provides a singleton config.Configuration instance
func (r *Resolver) Config() *config.Configuration {
	if r.config == nil {
		conf, err := config.LoadConfiguration()
		if err != nil {
			err = fmt.Errorf("configuration load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.config = conf
	}

	return r.config
}

// HTTPServer provides a singleton httpserver.Server instance
func (r *Resolver) HTTPServer() *httpserver.Server {
	if r.httpServer == nil {
		c := r.Config()

		log := r.Log().With(slog.String("tags", "http"))
		cLogger := &logger.CustomLogger{
			Enabled: c.Logger.Enabled,
			Level:   c.Logger.Level,
			Log:     log,
		}

		controllers := &httpserver.ControllerRegistry{
			ExampleController: r.ExampleController(),
		}
		routerConfig := &httpserver.RouterConfig{Namespace: c.HTTP.Router.Namespace}
		serverConfig := &httpserver.ServerConfig{
			Controllers:  controllers,
			Host:         c.HTTP.Server.Host,
			Logger:       cLogger,
			Port:         c.HTTP.Server.Port,
			RouterConfig: routerConfig,
		}

		server, err := httpserver.NewServer(serverConfig)
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
		c := r.Config()

		var handler slog.Handler
		opts := &slog.HandlerOptions{
			Level: logLevel(c.Logger.Level),
		}
		if c.Logger.Verbose {
			opts.AddSource = true
		}

		attrs := []slog.Attr{
			slog.Int(logger.AttrKey.PID, os.Getpid()),
			slog.String(logger.AttrKey.App.Name, r.Metadata().Name),
			slog.String(logger.AttrKey.App.Version, r.Metadata().Version),
		}

		if c.Metadata.Environment == app.Env.Development {
			handler = logger.NewDevHandler(*r.Metadata(), opts).WithAttrs(attrs)
		} else {
			handler = slog.NewJSONHandler(os.Stdout, opts).WithAttrs(attrs)
		}

		logger := slog.New(handler)
		slog.SetDefault(logger)

		r.log = logger
	}

	return r.log
}

// Metadata provides a singleton application Metadata instance
func (r *Resolver) Metadata() *app.Metadata {
	if r.metadata == nil {
		var metadata app.Metadata

		jsondata, err := os.ReadFile("/app/package.json")
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

		if r.Config().Metadata.Version != "" {
			metadata.Version = r.Config().Metadata.Version
		}

		r.metadata = &metadata
	}

	return r.metadata
}

func (r *Resolver) QueryHandler() *query.QueryHandler {
	if r.queryHandler == nil {
		c := r.Config()

		limit := int(c.HTTP.Router.Paging.DefaultLimit)

		attr := c.HTTP.Router.Sorting.DefaultAttr
		order := c.HTTP.Router.Sorting.DefaultOrder

		queryConfig := &query.QueryConfig{
			Defaults: &query.QueryDefaults{
				Paging: query.QueryPaging{
					Limit: &limit,
				},
				Sorting: &query.QuerySorting{
					Attr:  &attr,
					Order: &order,
				},
			},
		}

		queryHandler, err := query.NewQueryHandler(queryConfig)
		if err != nil {
			err = fmt.Errorf("http server load error: %w", err)
			slog.Error(err.Error())
			panic(err)
		}

		r.queryHandler = queryHandler
	}

	return r.queryHandler
}

// PostgreSQLClient provides a singleton postgres pgxpool.Pool instance
func (r *Resolver) PostgreSQLClient() *pgxpool.Pool {
	if r.postgreSQLClient == nil {
		if err := app.Validator.Validate.StructPartial(r.config, "Postgres"); err != nil {
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
