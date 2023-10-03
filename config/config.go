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
	Environment string `validate:"required,oneof=development production"`
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
	Enabled   bool       `validate:"required,oneof=false true"`
	Level     string     `validate:"required,oneof=debug info warn error"`
	SubLogger SubLoggers `validate:"required"`
	Verbose   bool       `validate:"required,oneof=false true"`
}

type SubLoggers struct {
	Domain SubLoggerConfig `validate:"required"`
	HTTP   SubLoggerConfig `validate:"required"`
	Repos  SubLoggerConfig `validate:"required"`
}

type SubLoggerConfig struct {
	Enabled bool   `validate:"oneof=false true"`
	Level   string `validate:"oneof=debug info warn error"`
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

	// logger
	if err := viper.BindEnv("logger.enabled", "LOG_ENABLED"); err != nil {
		log.Fatalf("error binding env var `LOG_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.level", "LOG_LEVEL"); err != nil {
		log.Fatalf("error binding env var `LOG_LEVEL`: %v", err)
	}
	if err := viper.BindEnv("logger.verbose", "LOG_VERBOSE"); err != nil {
		log.Fatalf("error binding env var `LOG_VERBOSE`: %v", err)
	}

	// sublogger - domain
	if err := viper.BindEnv("logger.sublogger.domain.enabled", "SUBLOG_DOMAIN_ENABLED"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_DOMAIN_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.sublogger.domain.level", "SUBLOG_DOMAIN_LEVEL"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_DOMAIN_LEVEL`: %v", err)
	}

	// sublogger - http
	if err := viper.BindEnv("logger.sublogger.http.enabled", "SUBLOG_HTTP_ENABLED"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_HTTP_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.sublogger.http.level", "SUBLOG_HTTP_LEVEL"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_HTTP_LEVEL`: %v", err)
	}

	// sublogger - repos
	if err := viper.BindEnv("logger.sublogger.repos.enabled", "SUBLOG_REPOS_ENABLED"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_REPOS_ENABLED`: %v", err)
	}
	if err := viper.BindEnv("logger.sublogger.repos.level", "SUBLOG_REPOS_LEVEL"); err != nil {
		log.Fatalf("error binding env var `SUBLOG_REPOS_LEVEL`: %v", err)
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
