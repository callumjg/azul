package server

type GameStatus string

const (
	PENDING  GameStatus = "PENDING"
	STARTED  GameStatus = "STARTED"
	COMPLETE GameStatus = "COMPLETE"
)

const (
	MAX_PLAYERS = 4
)

type Player struct {
	*Client
	Score uint8
	Board
}

type Game struct {
	Status              GameStatus `json:"status"`
	Players             []Player   `json:"players"`
	CurrentPlayerIndex  uint8      `json:"currentPlayerIndex"`
	StartingPlayerIndex uint8      `json:"-"`
}

func (g *Game) ScoreRound() {

}

func (g Game) isGameOver() bool {
	return false
}
