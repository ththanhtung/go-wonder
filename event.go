package main

import "encoding/json"

type Event struct {
	EventType string `json:"type"`
	Payload   string `json:"payload"`
}

func DecodeEvent(rawEvent []byte) *Event {
	var decodedEvent Event
	json.Unmarshal(rawEvent, &decodedEvent)
	return &decodedEvent
}