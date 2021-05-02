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
func CreateNotificationDoc(callbackUri string, objectClass ObjectClass, jobType string) (string, error) {
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
    <%s>
      <%s/>
    </%s>
  </trigger>
</NotificationDocument>`
	finalDoc := fmt.Sprintf(basestring, callbackUri, objectClass, jobType, objectClass)

	_, testErr := xmlquery.Parse(strings.NewReader(finalDoc))
	if testErr != nil {
		log.Print("Failing document was: ", finalDoc)
		return "", testErr
	}
	return finalDoc, nil
}

func CreateNotification(r *vidispine.VSRequestor, callback_uri string, objectClass ObjectClass, notificationType string) error {
	newdoc, build_err := CreateNotificationDoc(callback_uri, objectClass, notificationType)
	if build_err != nil {
		log.Print("ERROR CreateNotification could not build a valid xml document: ", build_err)
		return errors.New("could not build valid xml")
	}

	urlPath := fmt.Sprintf("/API/%s/notification", objectClass)
	response, serverErr := r.Post(urlPath,
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
