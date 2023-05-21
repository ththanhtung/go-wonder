package main

func main() {
	transport := NewTransport()
	wonderWordGame := NewWonderWordGame()
	state := NewGameState(wonderWordGame)
	networkManager := NewNetworkManager(transport)
	controller := NewController(networkManager, state)
	transport.RegisterNewConnHandler(controller.RegisterPlayer)
	controller.Init()
}