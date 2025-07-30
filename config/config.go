package config

import (
	"fmt"
	"log/slog"

	"github.com/fsnotify/fsnotify"
	"github.com/jasonsites/gosk/internal/app"
	"github.com/spf13/viper"
)

// Configuration defines application configuration
type Configuration struct {
	App      App      `validate:"required"`
	External External `validate:"required"`
	HTTP     HTTP     `validate:"required"`
	Logger   Logger   `validate:"required"`
	Postgres Postgres `validate:"required"`
}

type App struct {
	Metadata Metadata `validate:"required"`
}

// External defines external service configuration
type External struct {
	Example struct {
		Host    string
		Timeout uint
	}
}

// HTTP defines HTTP Server configuration
type HTTP struct {
	Router struct {
		Namespace string `validate:"required"`
		Paging    struct {
			DefaultLimit uint `validate:"required"`
		}
		Sorting struct {
			DefaultAttr  string `validate:"required"`
			DefaultOrder string `validate:"required"`
		} `validate:"required"`
	} `validate:"required"`
	Server struct {
		Host string
		Port uint `validate:"required,max=65535"`
	} `validate:"required"`
}

// Logger defines the primary logger configuration
type Logger struct {
	Format  string `validate:"oneof=json styled"`
	Level   string `validate:"oneof=debug info warn error"`
	Verbose bool
}

// Metadata defines application metadata
type Metadata struct {
	Environment string `validate:"oneof=development production"`
	Name        string
	Version     string
}

// Postgres defines the postgres connection parameters
type Postgres struct {
	Database string `validate:"required"`
	Host     string `validate:"required"`
	Password string `validate:"required"`
	Port     uint   `validate:"required,max=65535"`
	User     string `validate:"required"`
}

// LoadConfiguration loads config parameters on startup
func LoadConfiguration() (*Configuration, error) {
	var conf Configuration

	viper.SetConfigName("config")
	viper.AddConfigPath("/app/config")
	viper.AllowEmptyEnv(true)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
	})
	viper.WatchConfig()

	// default values
	viper.SetDefault("app.metadata.environment", "production")
	viper.SetDefault("app.metadata.name", "gosk")
	viper.SetDefault("app.metadata.version", "local")
	viper.SetDefault("external.example.baseURL", "http://www.example.com")
	viper.SetDefault("external.example.timeout", 25000)
	viper.SetDefault("http.router.namespace", "domain")
	viper.SetDefault("http.router.paging.defaultLimit", 20)
	viper.SetDefault("http.router.sorting.defaultAttr", "created_on")
	viper.SetDefault("http.router.sorting.defaultOrder", "desc")
	viper.SetDefault("http.server.host", "localhost")
	viper.SetDefault("http.server.port", 9202)
	viper.SetDefault("logger.enabled", true)
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.verbose", false)
	viper.SetDefault("postgres.database", "svcdb")
	viper.SetDefault("postgres.host", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.user", "postgres")

	// environment variables
	viper.BindEnv("app.metadata.environment", "APP_ENV")
	viper.BindEnv("app.metadata.version", "APP_VERSION")
	viper.BindEnv("http.server.host", "HTTP_SERVER_HOST")
	viper.BindEnv("http.server.port", "HTTP_SERVER_PORT")
	viper.BindEnv("logger.format", "LOGGER_FORMAT")
	viper.BindEnv("logger.level", "LOGGER_LEVEL")
	viper.BindEnv("logger.verbose", "LOGGER_VERBOSE")
	viper.BindEnv("postgres.database", "POSTGRES_DB")
	viper.BindEnv("postgres.host", "POSTGRES_HOST")
	viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("postgres.port", "POSTGRES_PORT")
	viper.BindEnv("postgres.user", "POSTGRES_USER")

	// read, unmarshal, and validate configuration
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Info("config file bypassed")
		} else {
			err := fmt.Errorf("configuration read error: %w", err)
			slog.Error(err.Error())
			return &conf, err
		}

	}
	if err := viper.Unmarshal(&conf); err != nil {
		err := fmt.Errorf("configuration unmarshal error: %w", err)
		slog.Error(err.Error())
		return &conf, err
	}
	if err := app.Validator.Validate.Struct(&conf); err != nil {
		return &conf, fmt.Errorf("invalid configuration: %v", err)
	}

	// fmt.Printf("%+v\n", conf)

	return &conf, nil
}
