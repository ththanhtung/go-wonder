package models


type GameState struct {
	Players     map[string]*Player
	InProgress   bool
	PlayerCount int
	WonderWordGame *WonderWordGame
}

func NewGameState(wonderWordGame *WonderWordGame) *GameState {
	return &GameState{
		Players:     make(map[string]*Player, 100),
		InProgress:   false,
		PlayerCount: 0,
		WonderWordGame: wonderWordGame,
	}
}

func (g *GameState) ReadyToStart() bool{
	return g.PlayerCount > 2 && !g.InProgress
}

func (g *GameState) Start(){
	g.WonderWordGame.Start()
	g.InProgress = true
}

func (g *GameState) End(){
	g.InProgress = false
}

func (g *GameState) AddPlayer(p *Player) int{

	g.PlayerCount++

	g.Players[p.UserID] = p
	return int(g.PlayerCount)
}

func (g *GameState) IsInProgress() bool {
	return g.InProgress
}

func (g *GameState) GetPlayers() []*Player {
	players := make([]*Player, 0)

	for _, p := range g.Players {
		players = append(players, p)
	}
	return players
}

func (g *GameState) GetPlayersCount() int {
	return g.PlayerCount
}