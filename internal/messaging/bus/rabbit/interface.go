package rabbit

import (
	"github.com/streadway/amqp"
)

type rabbitConnection interface {
	Close() error
	IsClosed() bool
	Channel() (rabbitChannel, error)
}

// This Rabbit Channel interface is mimicking functions from the library, https://github.com/streadway/amqp/blob/master/channel.go
type rabbitChannel interface {
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	Confirm(bool) error
	NotifyPublish(chan amqp.Confirmation) chan amqp.Confirmation
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type rabbitConn amqp.Connection

func (r *rabbitConn) Close() error {
	return (*amqp.Connection)(r).Close()
}

func (r *rabbitConn) IsClosed() bool {
	return (*amqp.Connection)(r).IsClosed()
}

// This get Channel in the amqp library
// only those functions in the amqp Channel can be called which are mimicked in rabbitChannel interface
func (r *rabbitConn) Channel() (rabbitChannel, error) {
	return (*amqp.Connection)(r).Channel()
}
