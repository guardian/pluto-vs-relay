package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

/** GetNotificationDocument
Checks that the incoming request to a notification endpoint is of the type we expect, stream in the content and parse
it into a VSNotificationDocument
*/
func GetNotificationDocument(w http.ResponseWriter, req *http.Request) (*VSNotificationDocument, *[]byte) {
	if req.Method != "POST" {
		w.Header().Add("Content-Type", "application/json")
		resp := GenericResponse{Status: "error", Detail: "Expected a POST"}
		responseBytes, _ := json.Marshal(resp)

		w.WriteHeader(405)
		io.Copy(w, bytes.NewReader(responseBytes))
		return nil, nil
	}

	bodyContent, readErr := ioutil.ReadAll(req.Body)
	if readErr != nil {
		log.Print("ERROR VidispineMessageHandler could not read content sent by server: ", readErr)
		w.WriteHeader(400)
		return nil, nil
	}

	log.Printf("DEBUG received notification content %s", string(bodyContent))

	var notification VSNotificationDocument
	parseErr := json.Unmarshal(bodyContent, &notification)
	if parseErr != nil {
		log.Printf("ERROR Could not parse content from server (expecting JSON). Offending content was: ")
		log.Printf(string(bodyContent))
		w.WriteHeader(400)
		return nil, nil
	}
	return &notification, &bodyContent
}
