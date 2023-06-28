package game

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"ww.com/events"
	"ww.com/models"
	"ww.com/transport"
)

type NetworkManager struct {
	transport *transport.Transport
	broadcast chan []byte
	register  chan *models.Client
	Clients   map[*models.Client]bool
}

func NewNetworkManager(transport *transport.Transport) *NetworkManager {
	return &NetworkManager{
		transport: transport,
		broadcast: make(chan []byte, 100),
		register:  make(chan *models.Client),
		Clients:   map[*models.Client]bool{},
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
				case client.MsgOut <- msg:
				}
			}
		}
	}
}

func (n *NetworkManager) Register(p *models.Player) {
	n.register <- p.Client
}

func (n *NetworkManager) BoardcastGameState(state *models.GameState, p *models.Player) {
	revealed := strings.Join(state.WonderWordGame.RevealedWord, "")

	currentState := events.BoardCastGameStateEvent{
		EventType:         "boardcast_game_state",
		Desc:              state.WonderWordGame.Challenge.Desc,
		Revealed:          revealed,
		CurrentPlayerID:   p.UserID,
		CurrentPlayerName: p.Username,
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
	msgEvent := events.BoardcastMessageEvent{
		EventType:  "boardcast_msg",
		SenderName: senderName,
		SenderId:   senderId,
		Payload:    string(msg),
	}

	msgEventJson, _ := json.Marshal(msgEvent)
	n.broadcast <- []byte(msgEventJson)
}

func (n *NetworkManager) SendToClient(c *models.Client, rawMsg []byte) {
	c.MsgOut <- rawMsg
}

func (n *NetworkManager) BoardcastCurrentPlayerState(p *models.Player) {
	updatePlayerState := events.UpdatePlayerStateEvent{
		EventType: "update_player_state",
		UserId:    p.UserID,
		Username:  p.Username,
		Score:     p.Score,
	}
	updatePlayerStateJson, _ := json.Marshal(&updatePlayerState)
	n.SendToClient(p.Client, []byte(updatePlayerStateJson))
}

func (n *NetworkManager) ReadMsg(c *models.Client) {
	log.Println("client reading msg")
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error reading message:", err.Error())
			}
			break
		}

		log.Println("msg:", string(payload))

		event := events.DecodeEvent(payload)

		if event.EventType == "guess" {
			log.Println("msg wonder:", event.Payload)
			c.MsgIn <- payload
		}
		if event.EventType == "msg" {
			log.Println("msg:", event.Payload)
			n.BoardcastMsg(event.SenderName, event.SenderId, event.Payload)
		}
	}
}

func (n *NetworkManager) WriteMsg(c *models.Client) {
	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.MsgOut:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Fatal("error sending message", err.Error())
				}
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Fatal("error sending message", err.Error())
			}
			log.Println("message sent")
		}
	}
}
