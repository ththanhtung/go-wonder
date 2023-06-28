package events

import "encoding/json"

type Event struct {
	EventType  string `json:"type"`
	Payload    string `json:"payload"`
	SenderName string `json:"senderName"`
	SenderId   string `json:"senderId"`
}

func DecodeEvent(rawEvent []byte) *Event {
	var decodedEvent Event
	json.Unmarshal(rawEvent, &decodedEvent)
	return &decodedEvent
}

type BoardCastGameStateEvent struct {
	EventType         string `json:"type"`
	Desc              string `json:"desc"`
	Revealed          string `json:"revealed"`
	CurrentPlayerID   string `json:"userId"`
	CurrentPlayerName string `json:"username"`
}

type WinningEvent struct {
	EventType string `json:"type"`
	Winner    string `json:"winner"`
	Score     int    `json:"score"`
}

type UpdatePlayerStateEvent struct {
	EventType string `json:"type"`
	UserId    string `json:"userId"`
	Username  string `json:"username"`
	Score     int    `json:"score"`
}

type BoardcastMessageEvent struct {
	EventType  string `json:"type"`
	Payload    string `json:"payload"`
	SenderName string `json:"senderName"`
	SenderId   string `json:"senderId"`
}
