package mocks

import (
	"errors"
	"github.com/streadway/amqp"
)

type AmqpChannelMock struct {
	PublishedMessages       []amqp.Publishing
	PublishNotifications    []chan amqp.Confirmation
	ConfirmPublishImmediate bool
	DidClose                bool
	PublishCallCount        int
	LastPublishExchange     string
	LastPublishKey          string
	ConfirmationsMode       bool
}

func NewAmqpChannelMock() *AmqpChannelMock {
	return &AmqpChannelMock{
		PublishedMessages:    make([]amqp.Publishing, 0),
		PublishNotifications: make([]chan amqp.Confirmation, 0),
	}
}

func (m *AmqpChannelMock) Ack(tag uint64, multiple bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Cancel(consumer string, noWait bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Close() error {
	m.DidClose = true
	return nil
}

func (m *AmqpChannelMock) Confirm(noWait bool) error {
	m.ConfirmationsMode = true
	return nil
}

func (m *AmqpChannelMock) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, errors.New("Not implemented")
}
func (m *AmqpChannelMock) ExchangeBind(destination, key, source string, noWait bool, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) ExchangeDeclarePassive(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) ExchangeDelete(name string, ifUnused, noWait bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) ExchangeUnbind(destination, key, source string, noWait bool, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Flow(active bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Get(queue string, autoAck bool) (msg amqp.Delivery, ok bool, err error) {
	return amqp.Delivery{}, false, errors.New("Not implemented")
}
func (m *AmqpChannelMock) Nack(tag uint64, multiple bool, requeue bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) NotifyCancel(c chan string) chan string {
	return nil
}
func (m *AmqpChannelMock) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return nil
}
func (m *AmqpChannelMock) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return nil, nil
}
func (m *AmqpChannelMock) NotifyFlow(c chan bool) chan bool {
	return nil
}
func (m *AmqpChannelMock) NotifyPublish(confirm chan amqp.Confirmation) chan amqp.Confirmation {
	m.PublishNotifications = append(m.PublishNotifications, confirm)
	return confirm
}
func (m *AmqpChannelMock) NotifyReturn(c chan amqp.Return) chan amqp.Return {
	return nil
}

func (m *AmqpChannelMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	m.LastPublishExchange = exchange
	m.LastPublishKey = key
	m.PublishCallCount += 1
	m.PublishedMessages = append(m.PublishedMessages, msg)
	if m.ConfirmPublishImmediate {
		go func() {
			for _, ch := range m.PublishNotifications {
				ch <- amqp.Confirmation{
					DeliveryTag: 1234567,
					Ack:         true,
				}
			}
		}()
	}
	return nil
}

func (m *AmqpChannelMock) Qos(prefetchCount, prefetchSize int, global bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	return 0, errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueInspect(name string) (amqp.Queue, error) {
	return amqp.Queue{}, errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueuePurge(name string, noWait bool) (int, error) {
	return 0, errors.New("Not implemented")
}
func (m *AmqpChannelMock) QueueUnbind(name, key, exchange string, args amqp.Table) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Recover(requeue bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Reject(tag uint64, requeue bool) error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) Tx() error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) TxCommit() error {
	return errors.New("Not implemented")
}
func (m *AmqpChannelMock) TxRollback() error {
	return errors.New("Not implemented")
}
