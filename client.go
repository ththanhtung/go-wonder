package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	msgIn chan []byte
	msgOut chan []byte
}


func NewClient(conn *websocket.Conn)*Client {
	return &Client{
		conn: conn,
		msgIn: make(chan []byte),
		msgOut: make(chan []byte),
	}
}

