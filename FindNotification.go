package main

import (
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
)

type ObjectClass string

const (
	Job      = "job"
	Metadata = "metadata"
)

/**
check if the url given in the document is the one we expect
*/
func CheckUrlPath(expectedUrlPath *string, doc *xmlquery.Node) bool {
	urlNodes := xmlquery.Find(doc, "//action//http//url")

	for _, node := range urlNodes {
		if node.InnerText() == *expectedUrlPath {
			return true
		}
	}
	return false
}

/**
check if the notification type in the document is the one we expect
*/
func CheckNotificationType(nt *string, doc *xmlquery.Node) bool {
	triggerNodes := xmlquery.Find(doc, fmt.Sprintf("//trigger//job//%s", *nt))

	return len(triggerNodes) == 1 //multiple nodes don't work
}

/**
test the notification whose spec is at the given url to see if the notification type and url match.
returns "true" if the doc matches or "false" otherwise
*/
func TestDocument(r *vidispine.VSRequestor, docurl string, expectedUriPtr *string, notificationTypePtr *string) (bool, error) {
	notificationDoc, serverErr := r.Get(docurl, "application/xml")
	if serverErr != nil {
		return false, serverErr
	}

	parsedNotification, parseErr := xmlquery.Parse(notificationDoc)
	if parseErr != nil {
		return false, parseErr
	}

	urlMatch := CheckUrlPath(expectedUriPtr, parsedNotification)
	notificationTypeMatch := CheckNotificationType(notificationTypePtr, parsedNotification)
	return urlMatch && notificationTypeMatch, nil
}

/**
Searches all available notifications to find our ones.
Returns a list of the notification types that are _missing_.
*/
func SearchForMyNotification(r *vidispine.VSRequestor, expectedUri string, objectClass ObjectClass, notificationType string) (bool, error) {
	baseUrl := fmt.Sprintf("/API/%s/notification", objectClass)
	listResponse, serverErr := r.Get(baseUrl, "application/xml")
	if serverErr != nil {
		return false, serverErr
	}

	parsedResponse, parseErr := xmlquery.Parse(listResponse)
	if parseErr != nil {
		log.Printf("ERROR SearchForMyNotification could not parse server response: %s", parseErr)
		return false, errors.New("invalid server response")
	}

	urinodes := xmlquery.Find(parsedResponse, "//uri")
	for _, node := range urinodes {
		log.Print("INFO SearchForMyNotification checking ", node.InnerText())
		result, err := TestDocument(r, node.InnerText(), &expectedUri, &notificationType)
		if err != nil {
			log.Print("ERROR SearchForMyNotification could not process: ", err)
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}
