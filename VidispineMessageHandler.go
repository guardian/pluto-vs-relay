package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type VidispineMessageHandler struct {
	Connection     mocks.AmqpConnectionInterface
	ChannelTimeout time.Duration
	ExchangeName   string
}

/**
tries to establish a channel to the broker. Retries up to the given timeout.
*/
func (h VidispineMessageHandler) EstablishChannel() (mocks.AmqpChannelInterface, error) {
	doneChan := make(chan mocks.AmqpChannelInterface)
	abortChan := make(chan interface{})

	go func() {
		for {
			rmqChnl, chnlErr := h.Connection.Channel()
			if chnlErr == nil {
				doneChan <- rmqChnl
				return
			} else {
				log.Print("ERROR VidispineMessageHandler.EstablishConnection - could not establish connection to broker: ", chnlErr)
				time.Sleep(5 * time.Second)
			}

			select {
			case <-abortChan:
				log.Print("ERROR VidispineMessageHandler.EstablishConnection - received timeout, stopping")
				return
			default:

			}
		}
	}()

	timer := time.NewTimer(h.ChannelTimeout)
	select {
	case <-timer.C:
		log.Print("ERROR VidispineMessageHandler.EstablishConnection - Connection attempt timed out, aborting")
		abortChan <- false
		return nil, errors.New("connection attempt timed out")
	case rmqChnl := <-doneChan:
		timer.Stop()
		return rmqChnl, nil
	}
}

/**
intended to be run as a goroutine, this function continues trying to send the message until successful.
it then waits for a confirmation to come back from the rabbitmq server to acknowledge that the message arrived.
sendCompletedChan - outgoing channel that will carry a single boolean value, false if we completed with error and true if we
completed ok
abortSignalChan - incoming channel that will cause the retries to stop if it has a value sent. if nothing is sent then
retries will continue indefinitely
routing_jey     - routing key for the message
bodyContent     - pointer to a byte array of the content to send
ackTimeout      - maximum time to wait for a confirmtion message
*/
func (h VidispineMessageHandler) backgroundRetryUntilSent(sendCompletedChan chan bool,
	abortSignalChan chan bool,
	routing_key string,
	bodyContent *[]byte,
	ackTimeout time.Duration) {
	rmqChan, chanErr := h.EstablishChannel()
	if chanErr != nil {
		log.Printf("ERROR VidispineMessageHandler.ServeHTTP could not establish channel to rabbitmq: %s", chanErr)
		sendCompletedChan <- false
		return
	}

	defer rmqChan.Close()
	confirmChan := make(chan amqp.Confirmation)

	for {
		//non-blocking check, if we have been sent abort signal
		select {
		case <-abortSignalChan:
			log.Print("WARNING VidispineMessageHandler.backgroundRetryUntilSent aborting retries")
			sendCompletedChan <- false
			return
		default:
		}
		//tell the client library to send confirmation messages to us
		rmqChan.NotifyPublish(confirmChan)

		setConfirmationErr := rmqChan.Confirm(false)
		if setConfirmationErr != nil {
			log.Print("Could not request channel confirmations: ", setConfirmationErr)
			time.Sleep(5 * time.Second)
			continue
		}
		msg := amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
			Body:            *bodyContent,
		}
		sendErr := rmqChan.Publish(h.ExchangeName, routing_key, true, false, msg)
		if sendErr != nil {
			log.Print("Could not send message to channel: ", sendErr)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	timer := time.NewTimer(ackTimeout)
	select {
	case confirmMsg := <-confirmChan:
		log.Print("INFO VidispineMessageHandler.backgroundRetryUntilSent Server confirmed, delivery tag is ", confirmMsg.DeliveryTag)
		sendCompletedChan <- true
	case <-timer.C:
		log.Print("ERROR VidispineMessageHandler.backgroundRetryUntilSent Server did not respond")
		sendCompletedChan <- false
	}
}

/**
receive the message from VS and pass it on
*/
func (h VidispineMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.Header().Add("Content-Type", "application/json")
		resp := GenericResponse{Status: "error", Detail: "Expected a POST"}
		responseBytes, _ := json.Marshal(resp)

		w.WriteHeader(405)
		io.Copy(w, bytes.NewReader(responseBytes))
		return
	}

	bodyContent, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		log.Print("ERROR VidispineMessageHandler could not read content sent by server: ", readErr)
		w.WriteHeader(400)
		return
	}

	var notification VSNotificationDocument
	parseErr := json.Unmarshal(bodyContent, &notification)
	if parseErr != nil {
		log.Printf("ERROR Could not parse content from server (expecting JSON). Offending content was: ")
		log.Printf(string(bodyContent))
		w.WriteHeader(400)
		return
	}

	notificationPtr := &notification
	routing_key := fmt.Sprintf("vidispine.job.%s", strings.ToLower(notificationPtr.GetAction()))

	log.Printf("DEBUG VidispineMessageHandler.ServeHTTP I will send the content %s to the routing key %s", string(bodyContent), routing_key)

	sendCompletedChan := make(chan bool)
	abortChan := make(chan bool)
	//set up the send loop asynchronously so we can manage a timeout
	go h.backgroundRetryUntilSent(sendCompletedChan, abortChan, routing_key, &bodyContent, 30*time.Second)

	timer := time.NewTimer(45 * time.Second)

	select {
	case successFlag := <-sendCompletedChan:
		if successFlag {
			w.WriteHeader(200)
			return
		} else {
			log.Print("ERROR VidispineMessageHandler.ServeHTTP could not send message")
			w.WriteHeader(500) //VS should re-send in this instance
			return
		}
	case <-timer.C:
		log.Print("ERROR VidispineMessageHandler timed out trying to send message")
		abortChan <- true
		w.WriteHeader(500) //VS should re-send in this instance
	}

}
