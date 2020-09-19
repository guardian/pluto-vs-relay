package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type VidispineMessageHandler struct {
	//will have the rabbitmq client in here

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
		log.Printf("ERROR VidispineMessageHandler could not read content sent by server: ", readErr)
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

	log.Printf("I would send the content %s to the routing key %s", string(bodyContent), routing_key)
	w.WriteHeader(200)
}
