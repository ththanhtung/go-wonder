package main

import (
	"ww.com/game"
	"ww.com/models"
	"ww.com/transport"
)

func main() {
	transport := transport.NewTransport()
	wonderWordGame := models.NewWonderWordGame()
	state := models.NewGameState(wonderWordGame)
	networkManager := game.NewNetworkManager(transport)
	controller := game.NewController(networkManager, state)
	transport.RegisterNewConnHandler(controller.RegisterPlayer)
	controller.Init()
}