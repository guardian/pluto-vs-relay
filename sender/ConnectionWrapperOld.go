package sender

//
//import (
//	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
//	"github.com/streadway/amqp"
//	"log"
//	"time"
//)
//
//type ConnectionWrapper interface {
//	//initiates the connection. returns a channel that indicates when it shut down.
//	Initiate(conn mocks.AmqpConnectionInterface,
//		inputChannel chan connectionPoolMessage,
//		outputChannel chan connectionPoolResult,
//		requestTerminationChannel chan bool) chan bool
//	Run()
//}
//
//type ConnectionWrapperImpl struct {
//	amqpChannel mocks.AmqpChannelInterface
//	terminationChannel chan bool
//	inputChannel chan connectionPoolMessage
//	outputChannel chan connectionPoolResult
//	requestTerminationChannel chan bool
//	amqpConfirmChanel chan amqp.Confirmation
//}
//
//func (c *ConnectionWrapperImpl) Initiate(conn mocks.AmqpConnectionInterface,
//	inputChannel chan connectionPoolMessage,
//	outputChannel chan connectionPoolResult,
//	requestTerminationChannel chan bool,
//) chan bool {
//	c.terminationChannel = make(chan bool)
//	amqpChannel, setupErr := conn.Channel()
//	if setupErr != nil {
//		log.Printf("ERROR ConnectionWrapper.Initiate could not set up new connection: ", setupErr)
//		c.terminationChannel <- false
//		return c.terminationChannel
//	}
//
//	amqpChannel.NotifyPublish(c.amqpConfirmChanel)
//	setConfirmationErr := c.amqpChannel.Confirm(false)
//	if setConfirmationErr != nil {
//		log.Print("Could not request channel confirmations: ", setConfirmationErr)
//		c.terminationChannel <- false
//		return c.terminationChannel
//	}
//
//	c.amqpChannel = amqpChannel
//	c.inputChannel = inputChannel
//	c.outputChannel = outputChannel
//	c.requestTerminationChannel = requestTerminationChannel
//
//	return c.terminationChannel
//}
//
///**
//intended to be executed in a goroutine, this will attempt to send messages on the contained channel until it fails
//or is asked to be shut down
// */
//func (c *ConnectionWrapperImpl) Run() {
//	for {
//		select {
//		case <- c.requestTerminationChannel:
//			log.Print("INFO ConnectionWrapper.Run requested termination, shutting down")
//			return
//		case msg := <- c.inputChannel:
//			for {
//				if c.processMessage(msg) {
//					break
//				} else {
//
//				}
//			}
//		}
//	}
//}
//
//func (c *ConnectionWrapperImpl) processMessage(msg connectionPoolMessage) bool {
//	sendErr := c.amqpChannel.Publish(msg.exchange, msg.routingKey, false, false, amqp.Publishing{
//		ContentType:     "application/json",
//		ContentEncoding: "utf-8",
//		MessageId:       msg.msgId.String(),
//		Timestamp:       time.Time{},
//		Body:            *msg.content,
//	})
//	if sendErr != nil {
//		log.Print("ERROR ConnectionWrapper.processMessage could not send: ", sendErr)
//		return false
//	}
//}
