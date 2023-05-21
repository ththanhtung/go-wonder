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

type BoardCastGameStateEvent struct {
	EventType string `json:"type"`
	Desc      string `json:"desc"`
	Revealed  string `json:"revealed"`
}

type WinningEvent struct {
	EventType string `json:"type"`
	Winner    string `json:"winner"`
	Score     int    `json:"score"`
}
