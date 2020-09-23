package sender

import (
	"github.com/streadway/amqp"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
	"log"
	"sync"
	"time"
)

type ConnectionWrapper interface {
	Send(message connectionPoolMessage) error
	//call this to shut down the connection. Returns a list of messages that did not get published.
	Finish() []connectionPoolMessage
}

type ConnectionWrapperImpl struct {
	amqpChannel        mocks.AmqpChannelInterface
	terminationChannel chan bool
	amqpConfirmChanel  chan amqp.Confirmation
	sendCounter        uint64
	//this is a map of (delivery_tag, source) that contains content which has not yet been acked
	awaitingConfirmation map[uint64]connectionPoolMessage
	mutex                sync.Mutex
}

func NewConnnectionWrapper(conn mocks.AmqpConnectionInterface) (ConnectionWrapper, error) {
	amqpChnl, chnlErr := conn.Channel()
	if chnlErr != nil {
		log.Print("ERROR NewConnectionWrapper could not create a channel: ", chnlErr)
		return nil, chnlErr
	}

	w := &ConnectionWrapperImpl{
		amqpChannel:          amqpChnl,
		terminationChannel:   make(chan bool),
		amqpConfirmChanel:    make(chan amqp.Confirmation),
		sendCounter:          0,
		awaitingConfirmation: make(map[uint64]connectionPoolMessage),
		mutex:                sync.Mutex{},
	}

	amqpChnl.NotifyPublish(w.amqpConfirmChanel)
	setNotifyErr := amqpChnl.Confirm(false)
	if setNotifyErr != nil {
		log.Print("ERROR NewConnectionWrapper could not enter confirmation mode: ", setNotifyErr)
		return nil, setNotifyErr
	}
	go w.pickUpConfirms()

	return w, nil
}

/**
this is a goroutine that listens for confirmation messages
*/
func (w *ConnectionWrapperImpl) pickUpConfirms() {
	for {
		select {
		case <-w.terminationChannel:
			log.Print("INFO ConnectionWrapperImpl.pickUpConfirms requested termination, shutting down")
			return
		case msg := <-w.amqpConfirmChanel:
			if msg.Ack {
				delete(w.awaitingConfirmation, msg.DeliveryTag)
			} else {
				log.Print("WARNING ConnectionWrapperImpl.pickUpConfirms server rejected message, resending")
				originalSource := w.awaitingConfirmation[msg.DeliveryTag]
				go func() {
					time.Sleep(3 * time.Second)
					w.Send(originalSource)
				}()
			}
		}
	}
}

/**
sends content to the channel. If this returns an error, then the ConnectionWrapper object should be disposed and
re-initialised
*/
func (w *ConnectionWrapperImpl) Send(message connectionPoolMessage) error {
	w.mutex.Lock()
	w.sendCounter += 1
	w.awaitingConfirmation[w.sendCounter] = message
	w.mutex.Unlock()

	return w.amqpChannel.Publish(message.exchange,
		message.routingKey,
		true,
		false,
		amqp.Publishing{
			ContentEncoding: "utf-8",
			ContentType:     "application/json",
			MessageId:       message.msgId.String(),
			Body:            *message.content,
		})
}

func (w *ConnectionWrapperImpl) Finish() []connectionPoolMessage {
	closeErr := w.amqpChannel.Close()
	if closeErr != nil {
		log.Print("ERROR ConnectionWrapper.Finish could not close broker channel: ", closeErr)
	}

	w.terminationChannel <- true
	pendingMessages := make([]connectionPoolMessage, len(w.awaitingConfirmation))
	i := 0
	for _, msg := range w.awaitingConfirmation {
		pendingMessages[i] = msg
	}

	return pendingMessages
}
