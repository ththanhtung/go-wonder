package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Transport struct {
	register   func(*websocket.Conn)
	upgrader websocket.Upgrader
}

func NewTransport() *Transport{
	return &Transport{
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (t *Transport) RegisterNewConnHandler(register func(*websocket.Conn)){
	t.register = register
}

func (t *Transport) Run(){
	server := gin.Default()

	server.GET("/start", t.init)
	
	server.Run(":8080")
}

func (t *Transport) init(c *gin.Context){
	conn, err := t.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	t.register(conn)
}