package sender

import (
	"fmt"
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
	timestamp  time.Time
}

func (m *connectionPoolMessage) String() string {
	return fmt.Sprintf("%s to %s on %s", m.msgId.String(), m.routingKey, m.exchange)
}

type connectionPoolResult struct {
	msgId     uuid.UUID
	isSuccess bool
}

type connectionPoolMap map[uuid.UUID]chan bool

type AmqpConnectionPoolImpl struct {
	connection  mocks.AmqpConnectionInterface
	amqpWrapper ConnectionWrapper
	mutex       sync.Mutex
}

func NewAmqpConnectionPool(conn mocks.AmqpConnectionInterface) AmqpConnectionPool {
	return &AmqpConnectionPoolImpl{
		connection:  conn,
		amqpWrapper: nil,
		mutex:       sync.Mutex{},
	}
}

func (p *AmqpConnectionPoolImpl) setupWrapper(pendingMessages []connectionPoolMessage) error {
	p.mutex.Lock()
	newWrapper, createErr := NewConnnectionWrapper(p.connection, pendingMessages)
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
		setupErr := p.setupWrapper([]connectionPoolMessage{})
		if setupErr != nil {
			return setupErr
		}
	}

	msg := connectionPoolMessage{
		exchange:   exchange,
		routingKey: routingKey,
		content:    content,
		timestamp:  time.Now(),
		msgId:      uuid.New(),
	}

	for {
		sendErr := p.amqpWrapper.Send(msg)
		if sendErr != nil {
			log.Print("WARNING ConnectionPool.Send could not send, re-initialising")
			pendingMessages := p.amqpWrapper.Finish()
			log.Printf("WARNING ConnectionPool.Send re-sending %d pending messages", len(pendingMessages))
			setupErr := p.setupWrapper(pendingMessages)
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
