package mocks

import (
	"crypto/tls"
	"github.com/streadway/amqp"
	"net"
)

type AmqpConnectionMock struct {
	ChannelError         error
	DidClose             bool
	TlsConnectionState   tls.ConnectionState
	ChannelMock          *AmqpChannelMock
	ChannelCallCount     int
	BlockedNotifications []chan amqp.Blocking
	CloseNotifications   []chan *amqp.Error
}

func NewAmqpConnectionMock(channelMock *AmqpChannelMock) *AmqpConnectionMock {
	return &AmqpConnectionMock{
		ChannelError:         nil,
		ChannelMock:          channelMock,
		DidClose:             false,
		TlsConnectionState:   tls.ConnectionState{},
		BlockedNotifications: make([]chan amqp.Blocking, 0),
	}
}

func (m *AmqpConnectionMock) Channel() (AmqpChannelInterface, error) {
	m.ChannelCallCount += 1
	if m.ChannelError != nil {
		return nil, m.ChannelError
	} else {
		return m.ChannelMock, nil
	}
}

func (m *AmqpConnectionMock) Close() error {
	m.DidClose = true
	return nil
}

func (m *AmqpConnectionMock) ConnectionState() tls.ConnectionState {
	return m.TlsConnectionState
}

func (m *AmqpConnectionMock) IsClosed() bool {
	return m.DidClose
}

func (m *AmqpConnectionMock) LocalAddr() net.Addr {
	return &net.IPAddr{
		IP: net.IP{127, 0, 0, 1},
	}
}

func (m *AmqpConnectionMock) NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking {
	m.BlockedNotifications = append(m.BlockedNotifications, receiver)
	return receiver
}

func (m *AmqpConnectionMock) NotifyClose(receiver chan *amqp.Error) chan *amqp.Error {
	m.CloseNotifications = append(m.CloseNotifications, receiver)
	return receiver
}
