package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	UserID   string
	Username string
	Client   *Client
	Score    int
	Guess    string
}

func NewPlayer(id, name string, conn *websocket.Conn) *Player {
	return &Player{
		UserID:   id,
		Username: name,
		Client:   NewClient(conn),
		Score:    0,
		Guess:    "",
	}
}

func (p *Player) Update(players []*Player, state *GameState) {
	scoreForOneGuess, _ := state.WonderWordGame.ScoreCalculator(p.Guess)
	p.Score += scoreForOneGuess
}

func (p *Player) MakeTurn() bool {
	timer := time.NewTimer(10 * time.Second)

	select {
	case <-timer.C:
		return false
	case <-p.Client.MsgIn:
		return true
	}
}
