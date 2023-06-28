package models

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	MsgIn  chan []byte
	MsgOut chan []byte
	Msg    chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		Conn:   conn,
		MsgIn:  make(chan []byte,100),
		MsgOut: make(chan []byte, 100),
		Msg:    make(chan []byte,100),
	}
}
