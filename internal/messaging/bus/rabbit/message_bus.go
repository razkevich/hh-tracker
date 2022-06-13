package rabbit

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/internalerrors"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/logging"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/bus"
	"time"
)

// Config represents a rabbit configuration
type Config struct {
	Queue string
	Hosts []string
}

// MessageBus logs messages to a rabbit outboxRabbitExchange
type MessageBus struct {
	hosts      []amqp.URI
	queue      string
	timeout    time.Duration
	connection rabbitConnection
}

// ProvideMessageBus is the wire provider for MessageBus
func ProvideMessageBus(config Config) (bus.IMessageBus, error) {
	mb := &MessageBus{
		queue:   config.Queue,
		timeout: 5 * time.Second,
	}

	hosts, err := mb.parseHosts(config.Hosts)
	if err != nil {
		return nil, err
	}
	mb.hosts = hosts
	return mb, nil
}

func (mb *MessageBus) err(msg string, err error) error {
	log.Error().Err(err).Msg(msg)
	return fmt.Errorf("%s: %w", msg, err)
}

func (mb *MessageBus) doConnect(dial func(string) (rabbitConnection, error)) error {
	for _, uri := range mb.hosts {
		host := fmt.Sprintf("%s:%d", uri.Host, uri.Port)
		connection, err := dial(uri.String())
		if err != nil {
			log.Error().Err(err).Str("host", host).Msg("failed to connect")
			continue
		}

		channel, err := connection.Channel()
		if err != nil {
			log.Error().Err(err).Str("host", host).Msg(err.Error())
			return err
		}

		_, err = channel.QueueDeclare(mb.queue, true, false, false, false, nil)
		if err != nil {
			log.Error().Err(err).Str("host", host).Msg(err.Error())
			panic(err)
		}

		mb.connection = connection
		log.Info().Str("host", host).Msg("connected to RabbitMQ host")
		return nil
	}

	return internalerrors.ErrNoValidRabbitMQHosts
}

// ReconnectIfDisconnected reconnects the connection if it was disconnected
func (mb *MessageBus) ReconnectIfDisconnected(ctx context.Context) error {
	if mb.connection == nil || mb.connection.IsClosed() {
		log.Info().Msg("Rabbit is not connected, attempting connection")
		err := mb.Connect()
		if err != nil {
			logging.GetLoggerFromCtx(ctx).Error().Err(err).Msg("failed rabbit reconnect")
			return internalerrors.ErrCannotConnectRabbit
		}
	}
	return nil
}

// Connect connects to the rabbit service
func (mb *MessageBus) Connect() error {
	return mb.doConnect(func(host string) (rabbitConnection, error) {
		conn, err := amqp.Dial(host)
		return (*rabbitConn)(conn), err
	})
}

// Disconnect disconnects from rabbit
func (mb *MessageBus) Disconnect() error {
	return mb.connection.Close()
}

func (mb *MessageBus) parseHosts(hosts []string) ([]amqp.URI, error) {
	parsedHosts := make([]amqp.URI, 0, len(hosts))
	for _, host := range hosts {
		uri, err := amqp.ParseURI(host)
		if err != nil {
			log.Warn().Err(err).Str("host", host).Msg("could not parse RabbitMQ host")
			continue
		}
		parsedHosts = append(parsedHosts, uri)
	}

	if len(parsedHosts) == 0 {
		return nil, internalerrors.ErrNoValidRabbitMQHosts
	}

	return parsedHosts, nil
}
