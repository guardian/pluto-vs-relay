package main

import (
	"errors"
	"github.com/antchfx/xmlquery"
	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/vidispine"
	"log"
)

func CheckUrlPath(expectedUrlPath *string, doc *xmlquery.Node) bool {
	urlNodes := xmlquery.Find(doc, "//action//http//url")
	//log.Println("DEBUG CheckUrlPath got ", len(urlNodes), " url nodes")

	for _, node := range urlNodes {
		if node.InnerText() == *expectedUrlPath {
			return true
		}
	}
	return false
}

func TestDocument(r *vidispine.VSRequestor, docurl string, expectedUriPtr *string) (bool, error) {
	notificationDoc, serverErr := r.Get(docurl, "application/xml")
	if serverErr != nil {
		return false, serverErr
	}

	parsedNotification, parseErr := xmlquery.Parse(notificationDoc)
	if parseErr != nil {
		return false, parseErr
	}

	return CheckUrlPath(expectedUriPtr, parsedNotification), nil
}

/**
Searches all available notifications to find our one
*/
func SearchForMyNotification(r *vidispine.VSRequestor, expectedUri string) (bool, error) {
	listResponse, serverErr := r.Get("/API/job/notification", "application/xml")
	if serverErr != nil {
		return false, serverErr
	}

	parsedResponse, parseErr := xmlquery.Parse(listResponse)
	if parseErr != nil {
		log.Printf("ERROR SearchForMyNotification could not parse server response: %s", parseErr)
		return false, errors.New("Invalid server response")
	}

	urinodes := xmlquery.Find(parsedResponse, "//uri")
	for _, node := range urinodes {
		log.Print("INFO SearchForMyNotification checking ", node.InnerText())
		result, err := TestDocument(r, node.InnerText(), &expectedUri)
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
