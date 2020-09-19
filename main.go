package main

import (
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	vidispine_url_str := os.Getenv("VIDISPINE_URL")
	vidispine_user := os.Getenv("VIDISPINE_USER")
	vidispine_passwd := os.Getenv("VIDISPINE_PASSWORD")
	callback_uri_str := os.Getenv("CALLBACK_URI")

	if vidispine_url_str == "" || vidispine_user == "" || vidispine_passwd == "" {
		log.Fatal("Please set VIDISPINE_URL, VIDISPINE_USER and VIDISPINE_PASSWORD in the environment")
	}

	if callback_uri_str == "" {
		log.Fatal("Please set CALLBACK_URI int the environment")
	}

	vidispine_url, url_parse_err := url.Parse(vidispine_url_str)
	if url_parse_err != nil {
		log.Fatal("VIDISPINE_URL is not valid: ", url_parse_err)
	}

	callback_url, url_parse_err := url.Parse(callback_uri_str)
	if url_parse_err != nil {
		log.Fatal("CALLBACK_URI is not valid: ", url_parse_err)
	}

	requestor := vidispine.NewVSRequestor(*vidispine_url, vidispine_user, vidispine_passwd)

	log.Print("Checking for our notification in ", vidispine_url.String())
	have_notification, check_err := SearchForMyNotification(requestor, callback_url.String())
	if check_err != nil {
		log.Fatal("Could not check for notification: ", check_err)
	}

	if !have_notification {
		log.Print("Notification not found, adding one...")
		createErr := CreateNotification(requestor, callback_url.String())
		if createErr != nil {
			log.Fatal("Could not create notification: ", createErr)
		}
	}

	messageHandler := VidispineMessageHandler{}

	log.Printf("Callback URL path is %s", callback_url.Path)
	http.Handle(callback_url.Path, messageHandler)

	log.Printf("Starting up on port 8080...")
	startServeErr := http.ListenAndServe(":8080", nil)
	if startServeErr != nil {
		log.Fatal(startServeErr)
	}
}
