package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Controller struct {
	networkManager *NetworkManager
	state          *GameState
}

func NewController(networkManager *NetworkManager, state *GameState) *Controller {
	return &Controller{
		networkManager: networkManager,
		state:          state,
	}
}

func (c *Controller) Init() {
	c.networkManager.Start()
}

func (c *Controller) RegisterPlayer(conn *websocket.Conn) {
	playerID := uuid.NewString()
	player := NewPlayer(playerID, playerID, conn)
	playersNumber := c.state.AddPlayer(player)
	log.Println("current players number:", playersNumber)
	c.networkManager.Register(player)
	c.CheckStartCondition()
}

func (c *Controller) CheckStartCondition() {
	if c.state.ReadyToStart() {
		c.state.Start()
		c.GameLoop()
	}
	log.Println("game cannot start now...")
}

func (c *Controller) GameLoop() {
	log.Println("starting game loop")
	ticker := time.NewTicker(33 * time.Millisecond)

	for range ticker.C {
		players := c.state.GetPlayers()
		log.Println("players:", players)

		i := 0

		for {
			if i > c.state.GetPlayersCount()-1 {
				i = 0
			}
			log.Println("It's your turn, player:", players[i].userID)
			if players[i].MakeTurn() {
				c.ProcessInput(players[i], players)
				i++
			} else {
				i++
			}
		}
	}

}

func (c *Controller) ProcessInput(p *Player, players []*Player) {
	for {
		event := DecodeEvent(<-p.client.msgIn)
		switch event.EventType {
		case "guess":
			p.Update(players, c.state)
			log.Println(event.Payload)
			return
		case "msg":
			log.Println(event.Payload)
			return
		}
	}

	// for {
	// 	select {
	// 	case msg := <-p.client.msgIn:
	// 		log.Println(string(msg))
	// 		return
	// 	}
	// }
}
