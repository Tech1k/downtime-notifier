package main

import "testing"

func TestSendMsg(t *testing.T) {
	text := "random message"
	SendMsg(text)
	messageReceived = true
	notifiableChatID = 111
	SendMsg(text)
}
