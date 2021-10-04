package main

import (
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/sender"
	"log"
	"net/http"
	"time"
)

type VidispineMessageHandler struct {
	ConnectionPool sender.AmqpConnectionPool
	ChannelTimeout time.Duration
	ExchangeName   string
}

/**
receive the message from VS and pass it on
*/
func (h VidispineMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	notificationPtr, bodyContentPtr := GetNotificationDocument(w, req)
	if notificationPtr == nil || bodyContentPtr == nil {
		return //the error message has already been output
	}
	routingKey := fmt.Sprintf("vidispine.job.%s.%s", notificationPtr.GetType("unknown"), notificationPtr.GetAction())

	//log.Printf("DEBUG VidispineMessageHandler.ServeHTTP I will send the content %s to the routing key %s", string(bodyContent), routing_key)
	log.Printf("DEBUG VidispineMessageHandler.ServeHTTP received message for %s", routingKey)
	sendErr := h.ConnectionPool.Send(h.ExchangeName, routingKey, bodyContentPtr)

	if sendErr == nil {
		w.WriteHeader(200)
	} else {
		log.Print("ERROR VidispineMessageHandler.ServeHTTP could not send message at all!")
		w.WriteHeader(500)
	}
}
