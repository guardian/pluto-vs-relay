package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

type TemplateContent struct {
	Url              string
	EntityType       string
	NotificationType string
}

/**
create and test a simple notification document.
if the result does not parse as xml, returns an error
*/
func CreateNotificationDoc(content *TemplateContent) (string, error) {
	templateContent := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<NotificationDocument xmlns="http://xml.vidispine.com/schema/vidispine">
  <action>
    <http synchronous="false">
      <retry>5</retry>
      <contentType>application/json</contentType>
      <url>{{.Url}}</url>
      <method>POST</method>
      <timeout>10</timeout>
    </http>
  </action>
  <trigger>
    <{{.EntityType}}>
      <{{.NotificationType}}/>
    </{{.EntityType}}>
  </trigger>
</NotificationDocument>`
	tmpl, err := template.New("notificationDoc").Parse(templateContent)
	if err != nil {
		return "", err
	}
	buffer := bytes.Buffer{}
	err = tmpl.Execute(&buffer, content)
	if err != nil {
		return "", err
	}

	finalDoc := buffer.String()

	log.Printf("DEBUG generated document is %s", finalDoc)
	_, testErr := xmlquery.Parse(strings.NewReader(finalDoc))
	if testErr != nil {
		log.Print("Failing document was: ", finalDoc)
		return "", testErr
	}
	return finalDoc, nil
}

func CreateNotification(r *vidispine.VSRequestor, callback_uri string, entityType string, notificationType string) error {
	newdoc, build_err := CreateNotificationDoc(&TemplateContent{
		Url:              callback_uri,
		EntityType:       entityType,
		NotificationType: notificationType,
	})
	if build_err != nil {
		log.Print("ERROR CreateNotification could not build a valid xml document: ", build_err)
		return errors.New("could not build valid xml")
	}

	urlEntityType := entityType
	if entityType == "metadata" {
		urlEntityType = "item" //metadata updates get sent to to /item endpoint
	}
	response, serverErr := r.Post(
		fmt.Sprintf("/API/%s/notification", urlEntityType),
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
