package main


type GameState struct {
	players     map[string]*Player
	inProgress   bool
	playerCount int
	wonderWordGame *WonderWordGame
}

func NewGameState(wonderWordGame *WonderWordGame) *GameState {
	return &GameState{
		players:     make(map[string]*Player),
		inProgress:   false,
		playerCount: 0,
		wonderWordGame: wonderWordGame,
	}
}

func (g *GameState) ReadyToStart() bool{
	return g.playerCount > 1 && !g.inProgress
}

func (g *GameState) Start(){
	g.wonderWordGame.Start()
	g.inProgress = true
}

func (g *GameState) End(){
	g.inProgress = false
}

func (g *GameState) AddPlayer(p *Player) int{

	g.playerCount++

	g.players[p.userID] = p
	return int(g.playerCount)
}

func (g *GameState) InProgress() bool {
	return g.inProgress
}

func (g *GameState) GetPlayers() []*Player {
	players := make([]*Player, 0)

	for _, p := range g.players {
		players = append(players, p)
	}
	return players
}

func (g *GameState) GetPlayersCount() int {
	return g.playerCount
}