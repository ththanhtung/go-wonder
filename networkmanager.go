package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type NetworkManager struct {
	transport *Transport
	broadcast chan []byte
	register  chan *Client
	Clients   map[*Client]bool
}

func NewNetworkManager(transport *Transport) *NetworkManager {
	return &NetworkManager{
		transport: transport,
		broadcast: make(chan []byte, 100),
		register:  make(chan *Client),
		Clients:   map[*Client]bool{},
	}
}

func (n *NetworkManager) Start() {
	go n.Run()
	n.transport.Run()
}

func (n *NetworkManager) Run() {
	log.Printf("NetworkManager: Listening for incoming Network traffic ...")
	for {
		select {
		case client := <-n.register:
			n.Clients[client] = true
			go n.ReadMsg(client)
			go n.WriteMsg(client)
			log.Println("new user is connected")
		case msg := <-n.broadcast:
			for client := range n.Clients {
				select {
				case client.msgOut <- msg:
				}
			}
		}
	}
}

func (n *NetworkManager) Register(p *Player) {
	n.register <- p.client
}

func (n *NetworkManager) BoardcastGameState(state *GameState, p *Player) {
	revealed := strings.Join(state.wonderWordGame.RevealedWord, "")

	currentState := BoardCastGameStateEvent{
		EventType:         "boardcast_game_state",
		Desc:              state.wonderWordGame.Challenge.Desc,
		Revealed:          revealed,
		CurrentPlayerID:   p.userID,
		CurrentPlayerName: p.username,
	}

	currentStateJson, err := json.Marshal(currentState)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(currentStateJson)
	n.broadcast <- []byte(currentStateJson)
}

func (n *NetworkManager) BoardcastMsg(senderName, senderId string, msg string) {
	msgEvent := BoardcastMessageEvent{
		EventType:  "boardcast_msg",
		SenderName: senderName,
		SenderId:   senderId,
		Payload:    string(msg),
	}

	msgEventJson, _ := json.Marshal(msgEvent)
	n.broadcast <- []byte(msgEventJson)
}

func (n *NetworkManager) SendToClient(c *Client, rawMsg []byte) {
	c.msgOut <- rawMsg
}

func (n *NetworkManager) BoardcastCurrentPlayerState(p *Player) {
	updatePlayerState := UpdatePlayerStateEvent{
		EventType: "update_player_state",
		UserId:    p.userID,
		Username:  p.username,
		Score:     p.score,
	}
	updatePlayerStateJson, _ := json.Marshal(&updatePlayerState)
	n.SendToClient(p.client, []byte(updatePlayerStateJson))
}

func (n *NetworkManager) ReadMsg(c *Client) {
	log.Println("client reading msg")
	defer func() {
		c.conn.Close()
	}()

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message:", err.Error())
			}
			break
		}

		log.Println("msg:", string(payload))

		event := DecodeEvent(payload)

		if event.EventType == "guess" {
			log.Println("msg wonder:", event.Payload)
			c.msgIn <- payload
		}
		if event.EventType == "msg" {
			log.Println("msg:", event.Payload)
			n.BoardcastMsg(event.SenderName, event.SenderId, event.Payload)
		}
	}
}

func (n *NetworkManager) WriteMsg(c *Client) {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.msgOut:
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Fatal("error sending message", err.Error())
				}
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Fatal("error sending message", err.Error())
			}
			log.Println("message sent")
		}
	}
}
