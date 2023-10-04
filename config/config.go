package config

import (
	"fmt"
	"log/slog"

	"github.com/spf13/viper"
)

// Configuration defines app configuration on startup
type Configuration struct {
	External External `validate:"required"`
	HTTP     HTTP     `validate:"required"`
	Logger   Logger   `validate:"required"`
	Metadata Metadata `validate:"required"`
	Postgres Postgres `validate:"required"`
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
			DefaultLimit  uint `validate:"required"`
			DefaultOffset uint `validate:"required"`
		} `validate:"required"`
		Sorting struct {
			DefaultAttr  string `validate:"required"`
			DefaultOrder string `validate:"required"`
		} `validate:"required"`
	} `validate:"required"`
	Server struct {
		Host string
		Mode string `validate:"required,oneof=debug release test"`
		Port uint   `validate:"required,max=65535"`
	} `validate:"required"`
}

// Logger defines the primary logger configuration
type Logger struct {
	Enabled bool   `validate:"required,oneof=false true"`
	Level   string `validate:"required,oneof=debug info warn error"`
	Verbose bool   `validate:"required,oneof=false true"`
}

// Metadata defines application metadata
type Metadata struct {
	Environment string `validate:"required,oneof=development production"`
	Name        string `validate:"required"`
	Version     string `validate:"required"`
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
	var config Configuration

	viper.SetConfigName("config")

	viper.AddConfigPath("/app/config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.AllowEmptyEnv(true)

	// http server
	if err := viper.BindEnv("http.server.host", "HTTP_SERVER_HOST"); err != nil {
		err := fmt.Errorf("error binding env var `HTTP_SERVER_HOST`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("http.server.mode", "HTTP_SERVER_MODE"); err != nil {
		err := fmt.Errorf("error binding env var `HTTP_SERVER_MODE`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("http.server.port", "HTTP_SERVER_PORT"); err != nil {
		err := fmt.Errorf("error binding env var `HTTP_SERVER_PORT`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// logger
	if err := viper.BindEnv("logger.enabled", "LOG_ENABLED"); err != nil {
		err := fmt.Errorf("error binding env var `LOG_ENABLED`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("logger.level", "LOG_LEVEL"); err != nil {
		err := fmt.Errorf("error binding env var `LOG_LEVEL`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("logger.verbose", "LOG_VERBOSE"); err != nil {
		err := fmt.Errorf("error binding env var `LOG_VERBOSE`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// metadata
	if err := viper.BindEnv("metadata.environment", "APP_ENV"); err != nil {
		err := fmt.Errorf("error binding env var `APP_ENV`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("metadata.version", "APP_VERSION"); err != nil {
		err := fmt.Errorf("error binding env var `APP_VERSION`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// postgres
	if err := viper.BindEnv("postgres.database", "POSTGRES_DB"); err != nil {
		err := fmt.Errorf("error binding env var `POSTGRES_DB`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("postgres.host", "POSTGRES_HOST"); err != nil {
		err := fmt.Errorf("error binding env var `POSTGRES_HOST`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("postgres.password", "POSTGRES_PASSWORD"); err != nil {
		err := fmt.Errorf("error binding env var `POSTGRES_PASSWORD`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("postgres.port", "POSTGRES_PORT"); err != nil {
		err := fmt.Errorf("error binding env var `POSTGRES_PORT`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.BindEnv("postgres.user", "POSTGRES_USER"); err != nil {
		err := fmt.Errorf("error binding env var `POSTGRES_USER`: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// read and unmarshal config
	if err := viper.ReadInConfig(); err != nil {
		err := fmt.Errorf("error reading config file: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		err := fmt.Errorf("error unmarshalling configuration: %w", err)
		slog.Error(err.Error())
		panic(err)
	}

	// fmt.Printf("\n%#v\n\n", config)

	return &config, nil
}
