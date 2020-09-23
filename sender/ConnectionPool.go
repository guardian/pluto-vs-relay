package sender

import (
	"github.com/google/uuid"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
	"log"
	"sync"
	"time"
)

type AmqpConnectionPool interface {
	Send(exchange string, routingKey string, content *[]byte) error
}

type connectionPoolMessage struct {
	exchange   string
	routingKey string
	content    *[]byte
	msgId      uuid.UUID
}

type connectionPoolResult struct {
	msgId     uuid.UUID
	isSuccess bool
}

type connectionPoolMap map[uuid.UUID]chan bool

type AmqpConnectionPoolImpl struct {
	connection  mocks.AmqpConnectionInterface
	inputQueue  chan connectionPoolMessage
	outputQueue chan connectionPoolResult
	amqpWrapper ConnectionWrapper
	mutex       sync.Mutex
}

func NewAmqpConnectionPool(conn mocks.AmqpConnectionInterface) AmqpConnectionPool {
	return &AmqpConnectionPoolImpl{
		connection:  conn,
		inputQueue:  make(chan connectionPoolMessage),
		outputQueue: make(chan connectionPoolResult),
		amqpWrapper: nil,
		mutex:       sync.Mutex{},
	}
}

func (p *AmqpConnectionPoolImpl) setupWrapper() error {
	p.mutex.Lock()
	newWrapper, createErr := NewConnnectionWrapper(p.connection)
	if createErr != nil {
		log.Print("Could not create a connection wrapper: ", createErr)
		return createErr
	}
	p.amqpWrapper = newWrapper
	p.mutex.Unlock()
	return nil
}

func (p *AmqpConnectionPoolImpl) Send(exchange string, routingKey string, content *[]byte) error {
	if p.amqpWrapper == nil {
		setupErr := p.setupWrapper()
		if setupErr != nil {
			return setupErr
		}
	}

	msg := connectionPoolMessage{
		exchange:   exchange,
		routingKey: routingKey,
		content:    content,
		msgId:      uuid.New(),
	}

	for {
		sendErr := p.amqpWrapper.Send(msg)
		if sendErr != nil {
			log.Print("WARNING ConnectionPool.Send could not send, re-initialising")
			setupErr := p.setupWrapper()
			if setupErr != nil {
				log.Print("FATAL can't set up another wrapper, exiting")
				panic("Can't set up wrapper")
			}
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return nil
}
