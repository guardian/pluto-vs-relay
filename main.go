package main

import (
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
	"net/url"
	"os"
)

func main() {
	vidispine_url_str := os.Getenv("VIDISPINE_URL")
	vidispine_user := os.Getenv("VIDISPINE_USER")
	vidispine_passwd := os.Getenv("VIDISPINE_PASSWORD")
	callback_uri := os.Getenv("CALLBACK_URI")

	if vidispine_url_str == "" || vidispine_user == "" || vidispine_passwd == "" {
		log.Fatal("Please set VIDISPINE_URL, VIDISPINE_USER and VIDISPINE_PASSWORD in the environment")
	}

	if callback_uri == "" {
		log.Fatal("Please set CALLBACK_URI int the environment")
	}

	vidispine_url, url_parse_err := url.Parse(vidispine_url_str)
	if url_parse_err != nil {
		log.Fatal("VIDISPINE_URL is not valid: ", url_parse_err)
	}

	requestor := vidispine.NewVSRequestor(*vidispine_url, vidispine_user, vidispine_passwd)

	log.Print("Checking for our notification in ", vidispine_url.String())
	have_notification, check_err := SearchForMyNotification(requestor, callback_uri)
	if check_err != nil {
		log.Fatal("Could not check for notification: ", check_err)
	}

	if !have_notification {
		log.Print("Notification not found, adding one...")
	}
	log.Fatal("Not implemented uet!")
}
