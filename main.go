package main

import (
	"github.com/streadway/amqp"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func setUpExchange(conn *amqp.Connection, exchangeName string) {
	rmqChan, chanErr := conn.Channel()
	if chanErr != nil {
		log.Fatal("Could not establish initial connection to rabbitmq: ", chanErr)
	}
	defer rmqChan.Close()

	declErr := rmqChan.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
	if declErr != nil {
		log.Fatal("Could not declare channel: ", declErr)
	}
}

func setUpNotifications(vidispine_url *url.URL, requestor *vidispine.VSRequestor, callback_url *url.URL) {
	expectedNotificationTypes := []string{"stop", "update", "create"}
	expectedNotificationClasses := []ObjectClass{Job, Metadata}

	log.Print("Checking for our notifications in ", vidispine_url.String())

	for _, cls := range expectedNotificationClasses {
		for _, nt := range expectedNotificationTypes {
			notificationPresent, check_err := SearchForMyNotification(requestor, callback_url.String(), cls, nt)
			if check_err != nil {
				log.Fatal("Could not check for notification: ", check_err)
			}

			if !notificationPresent {
				log.Printf("INFO setUpNotifications missing %s notification", nt)
				createErr := CreateNotification(requestor, callback_url.String(), cls, nt)
				if createErr != nil {
					log.Fatal("Could not create notification: ", createErr)
				}
			}
		}
	}
}

func main() {
	vidispine_url_str := os.Getenv("VIDISPINE_URL")
	vidispine_user := os.Getenv("VIDISPINE_USER")
	vidispine_passwd := os.Getenv("VIDISPINE_PASSWORD")
	callback_uri_str := os.Getenv("CALLBACK_URI")
	rabbitmq_uri_str := os.Getenv("RABBITMQ_URI")
	exchangeName := os.Getenv("RABBITMQ_EXCHANGE")

	if vidispine_url_str == "" || vidispine_user == "" || vidispine_passwd == "" {
		log.Fatal("Please set VIDISPINE_URL, VIDISPINE_USER and VIDISPINE_PASSWORD in the environment")
	}

	if callback_uri_str == "" {
		log.Fatal("Please set CALLBACK_URI int the environment")
	}

	if rabbitmq_uri_str == "" {
		log.Fatal("Please set RABBITMQ_URI in the environment")
	}

	if exchangeName == "" {
		log.Fatal("Please set RABBITMQ_EXCHANGE in the environment")
	}

	vidispine_url, url_parse_err := url.Parse(vidispine_url_str)
	if url_parse_err != nil {
		log.Fatal("VIDISPINE_URL is not valid: ", url_parse_err)
	}

	callback_url, url_parse_err := url.Parse(callback_uri_str)
	if url_parse_err != nil {
		log.Fatal("CALLBACK_URI is not valid: ", url_parse_err)
	}

	rmq, rmqErr := amqp.Dial(rabbitmq_uri_str)
	if rmqErr != nil {
		log.Fatal("Could not connect to rabbitmq: ", rmqErr)
	}

	requestor := vidispine.NewVSRequestor(*vidispine_url, vidispine_user, vidispine_passwd)

	setUpNotifications(vidispine_url, requestor, callback_url)

	setUpExchange(rmq, exchangeName)

	messageHandler := VidispineMessageHandler{
		Connection:     &mocks.AmqpConnectionShim{Connection: rmq},
		ExchangeName:   exchangeName,
		ChannelTimeout: 45 * time.Second,
	}
	healthcheckHandler := HealthcheckHandler{}

	log.Printf("Callback URL path is %s", callback_url.Path)
	http.Handle(callback_url.Path, messageHandler)
	http.Handle("/healthcheck", healthcheckHandler)

	log.Printf("Starting up on port 8080...")
	startServeErr := http.ListenAndServe(":8080", nil)
	if startServeErr != nil {
		log.Fatal(startServeErr)
	}
}
