package main

import (
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/sender"
	"log"
	"net/http"
	"time"
)

type VidispineItemMessageHandler struct {
	ConnectionPool sender.AmqpConnectionPool
	ChannelTimeout time.Duration
	ExchangeName   string
}

/**
receive the message from VS and pass it on
*/
func (h VidispineItemMessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	notificationPtr, bodyContentPtr := GetNotificationDocument(w, req)
	if notificationPtr == nil || bodyContentPtr == nil {
		return //the error message has already been output
	}
	routingKey := fmt.Sprintf("vidispine.item.%s", notificationPtr.GetAction())

	log.Printf("DEBUG VidispineItemMessageHandler.ServeHTTP received message for %s", routingKey)
	sendErr := h.ConnectionPool.Send(h.ExchangeName, routingKey, bodyContentPtr)

	if sendErr == nil {
		w.WriteHeader(200)
	} else {
		log.Print("ERROR VidispineItemMessageHandler.ServeHTTP could not send message at all!")
		w.WriteHeader(500)
	}
}
