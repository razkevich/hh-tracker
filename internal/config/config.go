// Package config is a collection of wire producers returning configuration
// from environment variables
package config

import (
	"github.com/google/wire"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/middleware"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/service"
)

// Providers is a wire set of all the config providers
var Providers = wire.NewSet(
	ProvideEnv,
	ProvideLoggingConfig,
	ProvideServerConfig,
	ProvideControllerConfig,
	ProvideServiceConfig,
	ProvideMiddlewareConfig,
	ProvideMongoConfig,
)

// ProvideLoggingConfig returns a logging Config struct populated from environment variables
func ProvideLoggingConfig(e Env) internal.LoggerConfig {
	return internal.LoggerConfig{
		GlobalLogLevel: e.GlobalLogLevel,
	}
}

// ProvideServerConfig returns a Server Config struct populated from environment variables for use by wire
func ProvideServerConfig(e Env) internal.ServerConfig {
	return internal.ServerConfig{
		Port: e.Port,
	}
}

// ProvideMiddlewareConfig returns a middleware Config struct populated from environment variables for use by wire
func ProvideMiddlewareConfig(e Env) middleware.MiddlewareConfig {
	return middleware.MiddlewareConfig{
		EnforceLimits:    e.EnforceLimits,
		DefaultPageLimit: e.DefaultPageLimit,
	}
}

// ProvideControllerConfig returns a controller Config struct populated from environment variables for use by wire
func ProvideControllerConfig(e Env) controller.ControllerConfig {
	return controller.ControllerConfig{
		SvcBaseURL: e.SvcBaseURL,
	}
}

// ProvideServiceConfig returns a Service Config struct populated from environment variables for use by wire
func ProvideServiceConfig(e Env) (*service.ServiceConfig, error) {
	return &service.ServiceConfig{}, nil
}

// ProvideMongoConfig returns a mongo Config struct populated from environment variables for use by wire
func ProvideMongoConfig(e Env) *mongo.MongoConfig {
	return &mongo.MongoConfig{
		DSN:          e.MongoDSN,
		DatabaseName: e.DatabaseName,
		Timeout:      e.Timeout,
	}
}
