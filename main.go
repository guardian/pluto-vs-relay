package main

import (
	"github.com/streadway/amqp"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/sender"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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

func setUpNotifications(vidispine_url *url.URL, requestor *vidispine.VSRequestor, callbackUrl *url.URL) {
	expectedEntityTypes := []string{"job", "metadata", "item"}
	expectedNotificationTypes := [][]string{
		{"stop", "update", "create"},
		{"modify"},
		{"create", "delete"},
	}
	requiredSubpaths := []string{"/job", "/item/metadata", "/item"}
	log.Print("Checking for our notifications in ", vidispine_url.String())

	for entityIndex, et := range expectedEntityTypes {
		for _, nt := range expectedNotificationTypes[entityIndex] {
			notificationPresent, check_err := SearchForMyNotification(requestor, callbackUrl.String()+requiredSubpaths[entityIndex], et, nt)
			if check_err != nil {
				log.Fatal("Could not check for notification: ", check_err)
			}

			if !notificationPresent {
				log.Printf("INFO setUpNotifications missing %s %s notification", et, nt)
				createErr := CreateNotification(requestor, callbackUrl.String()+requiredSubpaths[entityIndex], et, nt)
				if createErr != nil {
					log.Fatal("Could not create notification: ", createErr)
				}
			}
		}
	}
}

/**
sets up a signal handler to terminate cleanly if we receive SIGINT or SIGTERM
*/
func handleSignals(rmq mocks.AmqpConnectionInterface) {
	sigChan := make(chan os.Signal, 1) //buffer of 1 signal in case we're not in receiving state when it comes through
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		receivedSig := <-sigChan
		log.Printf("INFO handleSignals received %s, shutting down in 5s", receivedSig.String())
		time.Sleep(5 * time.Second)
		closeErr := rmq.Close()
		if closeErr != nil {
			log.Print("ERROR handleSignals could not close connection: ", closeErr)
		}
		os.Exit(0)
	}()
}

func main() {
	vidispine_url_str := os.Getenv("VIDISPINE_URL")
	vidispine_user := os.Getenv("VIDISPINE_USER")
	vidispine_passwd := os.Getenv("VIDISPINE_PASSWORD")
	callback_uri_str := os.Getenv("CALLBACK_BASE")
	rabbitmq_uri_str := os.Getenv("RABBITMQ_URI")
	exchangeName := os.Getenv("RABBITMQ_EXCHANGE")

	if vidispine_url_str == "" || vidispine_user == "" || vidispine_passwd == "" {
		log.Fatal("Please set VIDISPINE_URL, VIDISPINE_USER and VIDISPINE_PASSWORD in the environment")
	}

	if callback_uri_str == "" {
		log.Fatal("Please set CALLBACK_BASE in the environment")
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
		log.Fatal("CALLBACK_BASE is not valid: ", url_parse_err)
	}

	rmq, rmqErr := amqp.Dial(rabbitmq_uri_str)
	if rmqErr != nil {
		log.Fatal("Could not connect to rabbitmq: ", rmqErr)
	}

	conn := &mocks.AmqpConnectionShim{
		Connection: rmq,
	}

	/*
		ensure that rabbitmq connection is terminated cleanly even if program exits uncleanly
	*/
	defer func() {
		if r := recover(); r != nil {
			log.Print("WARNING main Program is existing due to panic")
			if rmq != nil && !rmq.IsClosed() {
				log.Print("INFO main Shutting down broker connection")
				closeErr := rmq.Close()
				if closeErr != nil {
					log.Print("ERROR main Could not shut down broker connection but terminating anyway ", closeErr)
				}
			}
			os.Exit(0xFF)
		}
	}()

	/*
		ensure that rabbiqmq connection is terminated cleanly if we receive termination signal
	*/
	handleSignals(conn)

	requestor := vidispine.NewVSRequestor(*vidispine_url, vidispine_user, vidispine_passwd)

	setUpNotifications(vidispine_url, requestor, callback_url)

	setUpExchange(rmq, exchangeName)
	amqpPool := sender.NewAmqpConnectionPool(conn)

	jobMessageHandler := VidispineMessageHandler{
		ConnectionPool: amqpPool,
		ExchangeName:   exchangeName,
		ChannelTimeout: 45 * time.Second,
	}
	itemMessageHandler := VidispineItemMessageHandler{
		ConnectionPool: amqpPool,
		ExchangeName:   exchangeName,
		ChannelTimeout: 45 * time.Second,
	}
	metaMessageHandler := VidispineMetadataMessageHandler{
		ConnectionPool: amqpPool,
		ExchangeName:   exchangeName,
		ChannelTimeout: 45 * time.Second,
	}
	healthcheckHandler := HealthcheckHandler{}

	log.Printf("Callback URL path is %s", callback_url.Path)
	http.Handle(callback_url.Path+"/job", jobMessageHandler)
	http.Handle(callback_url.Path+"/item/metadata", metaMessageHandler)
	http.Handle(callback_url.Path+"/item", itemMessageHandler)

	http.Handle("/healthcheck", healthcheckHandler)

	log.Printf("Starting up on port 8080...")
	startServeErr := http.ListenAndServe(":8080", nil)
	if startServeErr != nil {
		log.Fatal(startServeErr)
	}
}
