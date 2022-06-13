package config

import (
	"github.com/caarlos0/env"
)

// Env is a collection of config variables read from the environment
type Env struct {
	Port             int      `env:"SVC_PORT" envDefault:"8000"`
	MongoDSN         string   `env:"MONGO_DSN,required"`
	Timeout          int      `env:"MONGO_TIMEOUT" envDefault:"5000"`
	DatabaseName     string   `env:"MONGO_DATABASE_NAME" envDefault:"personal-data"`
	SvcBaseURL       string   `env:"SVC_BASE_URL,required"`
	GlobalLogLevel   string   `env:"LOGGING_LEVEL" envDefault:"info"`
	DefaultPageLimit int      `env:"DEFAULT_PAGE_LIMIT" envDefault:"20"`
	RabbitHosts      []string `env:"RABBIT_HOSTS,required"`
	RabbitQueue      string   `env:"RABBIT_CONSUME_QUEUE" envDefault:"personal-data-consumer-queue"`
	EnforceLimits    bool     `env:"ENFORCE_LIMITS" envDefault:"true"`
}

// ProvideEnv is a wire provider for configuration environment variables.
func ProvideEnv() (Env, error) {
	var e Env
	err := env.Parse(&e)
	return e, err
}
