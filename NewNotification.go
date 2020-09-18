package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"io/ioutil"
	"log"
	"strings"
)

/**
create and test a simple notification document.
if the result does not parse as xml, returns an error
*/
func CreateNotificationDoc(callback_uri string) (string, error) {
	basestring := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<NotificationDocument xmlns="http://xml.vidispine.com/schema/vidispine">
  <action>
    <http synchronous="false">
      <retry>5</retry>
      <contentType>application/json</contentType>
      <url>%s</url>
      <method>POST</method>
      <timeout>10</timeout>
    </http>
  </action>
  <trigger>
    <job>
      <stop/>
	  <create/>
	  <update/>
    </job>
  </trigger>
</NotificationDocument>`
	_, test_err := xmlquery.Parse(strings.NewReader(basestring))
	if test_err != nil {
		return "", test_err
	}
	return fmt.Sprintf(basestring, callback_uri), nil
}

func CreateNotification(r *vidispine.VSRequestor, callback_uri string) error {
	newdoc, build_err := CreateNotificationDoc(callback_uri)
	if build_err != nil {
		log.Print("ERROR CreateNotification could not build a valid xml document: ", build_err)
		return errors.New("could not build valid xml")
	}

	response, serverErr := r.Post("/API/job/notification",
		"application/xml",
		"application/xml",
		strings.NewReader(newdoc),
	)

	if serverErr != nil {
		return serverErr
	}

	serverResponseBytes, _ := ioutil.ReadAll(response)
	log.Printf("Notification created succesfully: %s", string(serverResponseBytes))
	return nil
}
