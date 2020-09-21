package main

import (
	"github.com/antchfx/xmlquery"
	"strings"
	"testing"
)

func TestCheckUrlPath(t *testing.T) {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<NotificationDocument xmlns="http://xml.vidispine.com/schema/vidispine">
  <action>
    <http synchronous="false">
      <retry>5</retry>
      <contentType>application/json</contentType>
      <url>http://pathto/api/notify/</url>
      <method>POST</method>
      <timeout>10</timeout>
    </http>
  </action>
  <trigger>
    <job>
      <stop/>
    </job>
  </trigger>
</NotificationDocument>`
	stringReader := strings.NewReader(content)
	doc, err := xmlquery.Parse(stringReader)

	if err != nil {
		t.Error("Could not parse test data: ", err)
		t.FailNow()
	}
	expectedUrlPath := "http://pathto/api/notify/"

	if CheckUrlPath(&expectedUrlPath, doc) != true {
		t.Error("Did not find expected URL path")
	}

	unexpectedUrlPath := "http://fdsjkhsfhjksfhjk"
	if CheckUrlPath(&unexpectedUrlPath, doc) != false {
		t.Error("Found unexpected URL path")
	}

}

func TestCheckNotificationType(t *testing.T) {
	content := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<NotificationDocument xmlns="http://xml.vidispine.com/schema/vidispine">
  <action>
    <http synchronous="false">
      <retry>5</retry>
      <contentType>application/json</contentType>
      <url>http://pathto/api/notify/</url>
      <method>POST</method>
      <timeout>10</timeout>
    </http>
  </action>
  <trigger>
    <job>
      <stop/>
    </job>
  </trigger>
</NotificationDocument>`

	stringReader := strings.NewReader(content)
	doc, err := xmlquery.Parse(stringReader)

	if err != nil {
		t.Error("Could not parse test data: ", err)
		t.FailNow()
	}

	expectedNotificationType := "stop"
	if !CheckNotificationType(&expectedNotificationType, doc) {
		t.Error("Could not find existing notification type of 'stop'")
	}

	unExpectedNotificationType := "rhubarb"
	if CheckNotificationType(&unExpectedNotificationType, doc) {
		t.Error("Got unexpected notification type")
	}
}
