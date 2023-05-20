package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	userID   string
	username string
	client   *Client
	score    int16
}

func NewPlayer(id, name string, conn *websocket.Conn)*Player {
	return &Player{
		userID:   id,
		username: name,
		client:   NewClient(conn),
		score:    0,
	}
}

func (p *Player) Update(players []*Player, state *GameState){
	
}

func (p *Player) MakeTurn() bool{
	timer := time.NewTimer(10*time.Second)

	select {
	case <- timer.C:
		return false
	case <-p.client.msgIn:
		return true
	}
}