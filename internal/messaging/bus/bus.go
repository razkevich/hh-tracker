package bus

import (
	"context"
	cloud "github.com/cloudevents/sdk-go/v2"
	"github.com/google/wire"
)

// IAppMessageBus gathers all the interfaces that a message bus must satisfy
type IAppMessageBus interface {
	Connect() error
	Disconnect() error
}

// IMessageBusClient is the interface for receiving messages from a bus
type IMessageBusClient interface {
	// Consume takes messages from the bus and calls the process function for each, acking
	// or rejecting those messages based on the value returned from that function.
	// Consume runs until it the quit chan is closed, returning nil, or
	// until process returns an error when it returns that error.
	Consume(context context.Context, quit <-chan struct{}, process func(context context.Context, event cloud.Event) (ack bool, err error)) error

	// ReconnectIfDisconnected reconnects the connection if it was disconnected
	ReconnectIfDisconnected(ctx context.Context) error
}

// IMessageBus gathers all the interfaces that a message bus must satisfy
type IMessageBus interface {
	IAppMessageBus
	IMessageBusClient
}

// Providers is the set of wire providers offering repository interfaces
var Providers = wire.NewSet(
	ProvideAppMessageBus,
	ProvideConsumerMessageBus,
)

// ProvideAppMessageBus is a wire function for providing an app message bus
func ProvideAppMessageBus(mb IMessageBus) IAppMessageBus {
	return mb
}

// ProvideConsumerMessageBus is a wire function for providing a consumer message bus
func ProvideConsumerMessageBus(mb IMessageBus) IMessageBusClient {
	return mb
}
