package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/sender"
)

type VidispineMessageHandler struct {
	ConnectionPool sender.AmqpConnectionPool
	ChannelTimeout time.Duration
	ExchangeName   string
}

/*
*
receive the message from VS and pass it on
*/
func (h VidispineMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	notificationPtr, bodyContentPtr := GetNotificationDocument(w, req)
	if notificationPtr == nil || bodyContentPtr == nil {
		return //the error message has already been output
	}
	routingKey := fmt.Sprintf("vidispine.job.%s.%s", notificationPtr.GetType("unknown"), notificationPtr.GetAction())

	log.Printf("DEBUG VidispineMessageHandler.ServeHTTP received message for %s, body: %s, exchangeName: %s", routingKey, bodyContentPtr, h.ExchangeName)
	sendErr := h.ConnectionPool.Send(h.ExchangeName, routingKey, bodyContentPtr)

	if sendErr == nil {
		w.WriteHeader(200)
	} else {
		log.Print("ERROR VidispineMessageHandler.ServeHTTP could not send message at all!")
		w.WriteHeader(500)
	}
}
