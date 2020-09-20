package mocks

/**
This file provides interfaces for the AMQP objects that we use, so that they can be mocked out for tests
*/
import (
	"crypto/tls"
	"github.com/streadway/amqp"
	"net"
)

type AmqpChannelInterface interface {
	Ack(tag uint64, multiple bool) error
	Cancel(consumer string, noWait bool) error
	Close() error
	Confirm(noWait bool) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	ExchangeDelete(name string, ifUnused, noWait bool) error
	ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error
	Flow(active bool) error
	Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error)
	Nack(tag uint64, multiple bool, requeue bool) error
	NotifyCancel(c chan string) chan string
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64)
	NotifyFlow(c chan bool) chan bool
	NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation
	NotifyReturn(c chan amqp.Return) chan amqp.Return
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueuePurge(name string, noWait bool) (int, error)
	QueueUnbind(name, key, exchange string, args amqp.Table) error
	Recover(requeue bool) error
	Reject(tag uint64, requeue bool) error
	Tx() error
	TxCommit() error
	TxRollback() error
}

//thin layer to allow us to call down to the real library through our interface
type AmqpConnectionInterface interface {
	Channel() (AmqpChannelInterface, error)
	Close() error
	ConnectionState() tls.ConnectionState
	IsClosed() bool
	LocalAddr() net.Addr
	NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
}

type AmqpConnectionShim struct {
	Connection *amqp.Connection
}

func (c *AmqpConnectionShim) Channel() (AmqpChannelInterface, error) {
	return c.Connection.Channel()
}

func (c *AmqpConnectionShim) Close() error {
	return c.Connection.Close()
}
func (c *AmqpConnectionShim) ConnectionState() tls.ConnectionState {
	return c.Connection.ConnectionState()
}
func (c *AmqpConnectionShim) IsClosed() bool {
	return c.Connection.IsClosed()
}
func (c *AmqpConnectionShim) LocalAddr() net.Addr {
	return c.Connection.LocalAddr()
}
func (c *AmqpConnectionShim) NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking {
	return c.Connection.NotifyBlocked(receiver)
}
func (c *AmqpConnectionShim) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	return c.Connection.NotifyClose(receiver)
}
