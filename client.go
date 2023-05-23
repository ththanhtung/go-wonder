package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	msgIn  chan []byte
	msgOut chan []byte
	msg    chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn:   conn,
		msgIn:  make(chan []byte,100),
		msgOut: make(chan []byte, 100),
		msg:    make(chan []byte,100),
	}
}
