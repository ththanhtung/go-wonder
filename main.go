package main

func main() {
	transport := NewTransport()
	state := NewGameState()
	networkManager := NewNetworkManager(transport)
	controller := NewController(networkManager, state)
	transport.RegisterNewConnHandler(controller.RegisterPlayer)
	controller.Init()
}