package config

import (
	"log"

	"github.com/spf13/viper"
)

// Configuration defines app configuration on startup
type Configuration struct {
	Application Application `validate:"required"`
	External    External    `validate:"required"`
	HTTP        HTTP        `validate:"required"`
	Logger      Logger      `validate:"required"`
	Metadata    Metadata    `validate:"required"`
	Postgres    Postgres    `validate:"required"`
}

type Application struct {
	Environment string `validate:"required"`
}

type External struct {
	Example struct {
		Host    string
		Timeout uint
	}
}

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

type Logger struct {
	App    LoggerConfig `validate:"required"`
	Domain LoggerConfig `validate:"required"`
	HTTP   LoggerConfig `validate:"required"`
	Repo   LoggerConfig `validate:"required"`
}

type LoggerConfig struct {
	Enabled bool   `validate:"required,oneof=false true"`
	Level   string `validate:"required,oneof=trace debug info warn error fatal panic"`
}

type Metadata struct {
	Path string `validate:"required"`
}

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

	// application
	if err := viper.BindEnv("application.environment", "APP_ENV"); err != nil {
		log.Fatalf("error binding env var `APP_ENV`: %v", err)
	}

	// http server
	if err := viper.BindEnv("http.server.host", "HTTP_SERVER_HOST"); err != nil {
		log.Fatalf("error binding env var `HTTP_SERVER_HOST`: %v", err)
	}
	if err := viper.BindEnv("http.server.mode", "HTTP_SERVER_MODE"); err != nil {
		log.Fatalf("error binding env var `HTTP_SERVER_MODE`: %v", err)
	}
	if err := viper.BindEnv("http.server.port", "HTTP_SERVER_PORT"); err != nil {
		log.Fatalf("error binding env var `HTTP_SERVER_PORT`: %v", err)
	}

	// logger - app
	if err := viper.BindEnv("logger.app.enabled", "LOGGER_APP_ENABLED"); err != nil {
		log.Fatalf("error binding env var `LOGGER_APP_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.app.level", "LOGGER_APP_LEVEL"); err != nil {
		log.Fatalf("error binding env var `LOGGER_APP_LEVEL`: %v", err)
	}

	// logger - domain
	if err := viper.BindEnv("logger.domain.enabled", "LOGGER_DOMAIN_ENABLED"); err != nil {
		log.Fatalf("error binding env var `LOGGER_DOMAIN_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.domain.level", "LOGGER_DOMAIN_LEVEL"); err != nil {
		log.Fatalf("error binding env var `LOGGER_DOMAIN_LEVEL`: %v", err)
	}

	// logger - http
	if err := viper.BindEnv("logger.http.enabled", "LOGGER_HTTP_ENABLED"); err != nil {
		log.Fatalf("error binding env var `LOGGER_HTTP_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.http.level", "LOGGER_HTTP_LEVEL"); err != nil {
		log.Fatalf("error binding env var `LOGGER_HTTP_LEVEL`: %v", err)
	}

	// logger - repos
	if err := viper.BindEnv("logger.repos.enabled", "LOGGER_REPOS_ENABLED"); err != nil {
		log.Fatalf("error binding env var `LOGGER_REPOS_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.repos.level", "LOGGER_REPOS_LEVEL"); err != nil {
		log.Fatalf("error binding env var `LOGGER_REPOS_LEVEL`: %v", err)
	}

	// metadata
	if err := viper.BindEnv("metadata.path", "METADATA_PATH"); err != nil {
		log.Fatalf("error binding env var `METADATA_PATH`: %v", err)
	}

	// postgres
	if err := viper.BindEnv("postgres.database", "POSTGRES_DB"); err != nil {
		log.Fatalf("error binding env var `POSTGRES_DB`: %v", err)
	}
	if err := viper.BindEnv("postgres.host", "POSTGRES_HOST"); err != nil {
		log.Fatalf("error binding env var `POSTGRES_HOST`: %v", err)
	}
	if err := viper.BindEnv("postgres.password", "POSTGRES_PASSWORD"); err != nil {
		log.Fatalf("error binding env var `POSTGRES_PASSWORD`: %v", err)
	}
	if err := viper.BindEnv("postgres.port", "POSTGRES_PORT"); err != nil {
		log.Fatalf("error binding env var `POSTGRES_PORT`: %v", err)
	}
	if err := viper.BindEnv("postgres.user", "POSTGRES_USER"); err != nil {
		log.Fatalf("error binding env var `POSTGRES_USER`: %v", err)
	}

	// external service - example
	if err := viper.BindEnv("external.services.example.baseURL", "EXTSVC_EXAMPLE_BASEURL"); err != nil {
		log.Fatalf("error binding env var `EXTSVC_EXAMPLE_BASEURL`: %v", err)
	}
	if err := viper.BindEnv("external.services.example.timeout", "EXTSVC_EXAMPLE_TIMEOUT"); err != nil {
		log.Fatalf("error binding env var `EXTSVC_EXAMPLE_TIMEOUT`: %v", err)
	}

	// read and unmarshal config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("error unmarshalling configuration: %v", err)
	}

	// fmt.Printf("\n%#v\n", config)

	return &config, nil
}
