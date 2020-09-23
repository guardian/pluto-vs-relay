package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/sender"
	"io"
	"io/ioutil"
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
	routing_key := fmt.Sprintf("vidispine.job.%s.%s", notificationPtr.GetType("unknown"), notificationPtr.GetAction())

	log.Printf("DEBUG VidispineMessageHandler.ServeHTTP I will send the content %s to the routing key %s", string(bodyContent), routing_key)

	sendErr := h.ConnectionPool.Send(h.ExchangeName, routing_key, &bodyContent)

	if sendErr == nil {
		w.WriteHeader(200)
	} else {
		log.Print("ERROR VidispineMessageHandler.ServeHTTP could not send message at all!")
		w.WriteHeader(500)
	}
}
