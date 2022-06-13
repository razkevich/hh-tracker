//+build wireinject

package main

import (
	"github.com/google/wire"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/config"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/health"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/bus"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/bus/rabbit"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/middleware"
	product "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/repository"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/service"
)

var providers = wire.NewSet(
	service.ProvideLogEntriesService,
	controller.ProvideLogEntriesController,
	controller.ProvideHealthController,
	rabbit.ProvideMessageBus,
	bus.Providers,
	consumer.ProvideConsumer,
	middleware.ProvideMiddleware,
	mongo.ProvideDatabase,
	mongo.ProvideClient,
	product.ProvideRepository,
)

var appProviders = wire.NewSet(
	config.Providers,
	health.ProvideHealth,
	internal.ProvideLogger,
	internal.ProvideServer,
	internal.ProvideApp,
	providers,
)

func setupApp() (*internal.App, error) {
	panic(wire.Build(
		appProviders,
	))
}

// TODO run go generate ./...
