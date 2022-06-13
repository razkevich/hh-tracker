package internal

import (
	"context"
	"github.com/rs/zerolog/log"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/bus"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// App provides the context to run the application.
type App struct {
	server      *http.Server
	consumer    *consumer.Consumer
	mongoClient *mongo.Client
	quit        chan (struct{})
	shutdown    chan (struct{})
	messageBus  bus.IMessageBus
}

// ProvideApp is a wire provider
func ProvideApp(
	server *http.Server,
	mongoClient *mongo.Client,
	consumer *consumer.Consumer,
	messageBus bus.IMessageBus,
) *App {

	return &App{
		server:      server,
		consumer:    consumer,
		mongoClient: mongoClient,
		messageBus:  messageBus,
	}
}

var notify = signal.Notify

func (app *App) supervisor() {
	signals := make(chan os.Signal, 1)
	notify(signals, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-signals:
		log.Info().Str("signal", sig.String()).Msg("received signal, shutting down")
	case <-app.quit:
	}
	if err := app.server.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("failed to shutdown HTTP server")
	}
	log.Debug().Msg("shut down HTTP server")
	if err := app.mongoClient.Close(); err != nil {
		log.Error().Err(err).Msg("failed to disconnect from database")
	}
	app.consumer.Stop()
	log.Debug().Msg("disconnected from database")
	if err := app.messageBus.Disconnect(); err != nil {
		log.Error().Err(err).Msg("failed to disconnect from message bus")
	}
	log.Debug().Msg("disconnected from message bus")
	close(app.shutdown)
}

func (app *App) retryConnect(wg *sync.WaitGroup, service string, connectFunc func() error) {
	var attempt int
	log.Debug().Msgf("connecting to %s", service)
	for attempt = 1; attempt < 7; attempt++ {
		err := connectFunc()
		if err == nil {
			log.Debug().Msgf("connected to %s on attempt %d", service, attempt)
			wg.Done()
			return
		}
		log.Error().Err(err).Msgf("failed to connect to %s on attempt %d", service, attempt)
		time.Sleep(10 * time.Second)
	}
	wg.Done()
	log.Fatal().Msgf("cannot connect to %s", service)
}

func (app *App) connect(services map[string]func() error) {
	var wg sync.WaitGroup

	wg.Add(len(services))
	for service, connectFunc := range services {
		go app.retryConnect(&wg, service, connectFunc)
	}
	wg.Wait()
}

// Start starts the application running until a shutdown signal is received
func (app *App) Start() {
	app.connect(map[string]func() error{
		"database":    app.mongoClient.Connect,
		"message bus": app.messageBus.Connect,
	})

	log.Debug().Msg("starting message relay")
	log.Debug().Msg("starting consumer queue")
	go app.consumer.Run()

	app.shutdown = make(chan struct{})
	app.quit = make(chan struct{})
	go app.supervisor()
	log.Debug().Msg("starting HTTP server")
	if err := app.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("failed to serve HTTP")
	}
	<-app.shutdown
	log.Info().Msg("Shutdown")
}

// Stop cleanly shuts down the application, waiting for it to complete
// shutting down before returning.
func (app *App) Stop() {
	close(app.quit)
	<-app.shutdown
}
