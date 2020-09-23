package main

//
//import (
//	"errors"
//	"gitlab.com/codmill/customer-projects/guardian/pluto-vs-relay/mocks"
//	"testing"
//	"time"
//)
//
//func TestEstablishChannel_ok(t *testing.T) {
//	mockChannel := mocks.NewAmqpChannelMock()
//	mockConnection := mocks.NewAmqpConnectionMock(mockChannel)
//
//	toTest := VidispineMessageHandler{
//		Connection:     mockConnection,
//		ChannelTimeout: 2 * time.Second,
//		ExchangeName:   "test-exchange-name",
//	}
//
//	result, err := toTest.EstablishChannel()
//	if err != nil {
//		t.Error("EstablishChannel returned unexpected error - ", err)
//		t.FailNow()
//	}
//
//	if result != mockChannel {
//		t.Error("Expexcted mockChannel to be returned but got ", result)
//	}
//}
//
//func TestEstablishChannel_failure(t *testing.T) {
//	mockChannel := mocks.NewAmqpChannelMock()
//	mockConnection := mocks.NewAmqpConnectionMock(mockChannel)
//	mockConnection.ChannelError = errors.New("kaboom")
//
//	toTest := VidispineMessageHandler{
//		Connection:     mockConnection,
//		ChannelTimeout: 10 * time.Second,
//		ExchangeName:   "test-exchange-name",
//	}
//
//	result, err := toTest.EstablishChannel()
//	if err == nil {
//		t.Error("EstablishChannel returned no error")
//	}
//
//	if result != nil {
//		t.Error("Expexcted nil channel to be returned but got ", result)
//	}
//
//	if mockConnection.ChannelCallCount != 2 {
//		t.Error("Expected Channel() to have been called 2 times, got ", mockConnection.ChannelCallCount)
//	}
//}
//
///**
//backgroundRetryUntilSend should publish the given message to the rabbitmq channel and return a true value
//via the provided go channel
//*/
//func TestBackgroundRetryUntilSend_ok(t *testing.T) {
//	mockChannel := mocks.NewAmqpChannelMock()
//	mockConnection := mocks.NewAmqpConnectionMock(mockChannel)
//
//	mockChannel.ConfirmPublishImmediate = true
//
//	toTest := VidispineMessageHandler{
//		Connection:     mockConnection,
//		ChannelTimeout: 10 * time.Second,
//		ExchangeName:   "test-exchange-name",
//	}
//
//	sendCompletedChan := make(chan bool)
//	abortSignalChan := make(chan bool)
//	contentBuffer := []byte(`{"key":"somekey","value":"somevalue"}`)
//
//	go toTest.backgroundRetryUntilSent(
//		sendCompletedChan,
//		abortSignalChan,
//		"test_routing_key",
//		&contentBuffer,
//		2*time.Second,
//	)
//
//	//blindly wait for completion, the test harness will kill us rather than go on for ever
//	result := <-sendCompletedChan
//	if !result {
//		t.Error("Expected successful send but got failure")
//	}
//
//	if mockChannel.PublishCallCount != 1 {
//		t.Error("Expected 1 call to Publish, got ", mockChannel.PublishCallCount)
//	}
//	if len(mockChannel.PublishedMessages) != 1 {
//		t.Error("Expected 1 message to have been sent but got ", len(mockChannel.PublishedMessages))
//	} else {
//		msg := mockChannel.PublishedMessages[0]
//		if string(msg.Body) != string(contentBuffer) {
//			t.Error("Sent data did not match incoming")
//		}
//		if msg.ContentType != "application/json" {
//			t.Error("Expected application/json content-type got ", msg.ContentType)
//		}
//
//	}
//}
