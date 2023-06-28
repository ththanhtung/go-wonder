package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"ww.com/events"
	"ww.com/models"
)

type Controller struct {
	networkManager *NetworkManager
	state          *models.GameState
}

func NewController(networkManager *NetworkManager, state *models.GameState) *Controller {
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
	player := models.NewPlayer(playerID, playerID, conn)
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
			log.Println("It's your turn, player:", players[i].UserID)
			c.networkManager.BoardcastGameState(c.state, players[i])
			c.networkManager.BoardcastCurrentPlayerState(players[i])
			if players[i].MakeTurn() {
				c.ProcessInput(players[i], players)
				c.networkManager.BoardcastGameState(c.state, players[i])
				c.networkManager.BoardcastCurrentPlayerState(players[i])
				if c.state.WonderWordGame.CheckIfEndGame() {
					highestScoringPlayer := c.FindHighestScoringPlayer()

					winner := events.WinningEvent{
						EventType: "win",
						Winner:    highestScoringPlayer.Username,
						Score:     highestScoringPlayer.Score,
					}
					WinningEventJson, _ := json.Marshal(winner)
					c.networkManager.broadcast <- []byte(WinningEventJson)
					c.state.WonderWordGame.Start()
				}
				i++
			} else {
				c.networkManager.BoardcastGameState(c.state, players[i])
				c.networkManager.BoardcastCurrentPlayerState(players[i])
				log.Println("It's your turn, player 2:", players[i].UserID)
				i++
			}

		}
	}
}

func (c Controller) FindHighestScoringPlayer() *models.Player {
	players := c.state.GetPlayers()

	highestScoringPlayer := &models.Player{
		Score: 0,
	}

	for _, p := range players {
		if p.Score > highestScoringPlayer.Score {
			highestScoringPlayer = p
		}
	}

	return highestScoringPlayer
}

func (c *Controller) ProcessInput(p *models.Player, players []*models.Player) {
	for {
		event := events.DecodeEvent(<-p.Client.MsgIn)
		switch event.EventType {
		case "guess":
			p.Guess = event.Payload
			p.Update(players, c.state)
			log.Println(event.Payload)
			log.Println("word", c.state.WonderWordGame.Challenge.Desc)
			log.Println("word", c.state.WonderWordGame.RevealedWord)
			log.Println("score", p.Score)
			return
		case "msg":
			log.Println(event.Payload)
			return
		}
	}
}
