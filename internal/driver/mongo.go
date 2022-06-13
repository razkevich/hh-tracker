package mongo

import (
	"context"
	"fmt"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/health"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConfig Mongo Config
type MongoConfig struct {
	DSN          string
	DatabaseName string
	Timeout      int
}

// Client provides a Client to interact with mongoDB
type Client struct {
	health.IHealther
	Client *mongo.Client
	Config *MongoConfig
	Logger zerolog.Logger
}

// ProvideClient factory method for wire
func ProvideClient(config *MongoConfig, logger zerolog.Logger, health health.IHealth) (*Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.DSN))
	if err != nil {
		return nil, fmt.Errorf("could not create mongo Client: %w", err)
	}
	c := Client{
		Client: client,
		Config: config,
		Logger: logger,
	}
	health.Register(&c)
	return &c, nil
}

// ProvideDatabase factory method for wire
func ProvideDatabase(c *Client) *mongo.Database {
	database := c.Client.Database(c.Config.DatabaseName)
	return database
}

// Name returns the name of the technology used
func (c *Client) Name() string {
	return "Mongo"
}

// Close close database connection
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Millisecond)
	defer cancel()
	return c.Client.Disconnect(ctx)
}

// Connect method to connect to a database
func (c *Client) Connect() error {
	log.Info().Msg("Connecting to mongo")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Millisecond)
	defer cancel()

	if err := c.Client.Connect(ctx); err != nil {
		log.Error().Err(err).Msg("mongo connect failed")
		return err
	}

	if err := c.ping(); err != nil {
		c.Close()
		log.Error().Err(err).Msg("mongo is unhealthy")
		return err
	}

	return nil
}

// Health method to connect to a database
func (c *Client) Health() bool {
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	return c.Client.Ping(timeout, nil) == nil
}

// OK method to report on health of mongo connection
func (c *Client) OK() bool {
	// TODO: seek more info on context use.
	return c.ping() == nil
}

func (c *Client) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	return c.Client.Ping(ctx, nil)
}
