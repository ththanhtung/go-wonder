package main

import (
	"encoding/json"
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

		i := 0

		for {
			if i > c.state.GetPlayersCount()-1 {
				i = 0
			}
			log.Println("It's your turn, player:", players[i].userID)
			c.networkManager.BoardcastGameState(c.state, players[i])
			c.networkManager.BoardcastCurrentPlayerState(players[i])
			if players[i].MakeTurn() {
				c.ProcessInput(players[i], players)
				log.Println("It's your turn, player 2:", players[i].userID)
				c.networkManager.BoardcastGameState(c.state, players[i])
				c.networkManager.BoardcastCurrentPlayerState(players[i])
				if c.state.wonderWordGame.CheckIfWinning() {
					winner := WinningEvent{
						EventType: "win",
						Winner:    players[i].username,
						Score:     players[i].score,
					}
					WinningEventJson, _ := json.Marshal(winner)
					c.networkManager.broadcast <- []byte(WinningEventJson)
					c.state.wonderWordGame.Start()
				}
				i++
			} else {
				c.networkManager.BoardcastGameState(c.state, players[i])
				c.networkManager.BoardcastCurrentPlayerState(players[i])
				log.Println("It's your turn, player 2:", players[i].userID)
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
			p.guess = event.Payload
			p.Update(players, c.state)
			log.Println(event.Payload)
			log.Println("word", c.state.wonderWordGame.Challenge.Desc)
			log.Println("word", c.state.wonderWordGame.RevealedWord)
			log.Println("score", p.score)
			return
		case "msg":
			log.Println(event.Payload)
			return
		}
	}
}
