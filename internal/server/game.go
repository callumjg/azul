package server

import (
	"math/rand"
	"time"
)

type GameStatus string

const (
	PENDING  GameStatus = "PENDING"
	STARTED  GameStatus = "STARTED"
	COMPLETE GameStatus = "COMPLETE"
)

const (
	MAX_PLAYERS = 4
)

type Game struct {
	Status              GameStatus `json:"status"`
	Players             []Player   `json:"players"`
	CurrentPlayerIndex  uint8      `json:"currentPlayerIndex"`
	StartingPlayerIndex uint8      `json:"startingPlayerIndex"`
	FactoryDisplays     [][]Tile   `json:"factoryDisplays"`
	TableCentre         []Tile     `json:"tableCentre"`
	TileBag             []Tile     `json:"tileBag"`
	UsedTiles           []Tile     `json:"usedTiles"`
}

func NewGame() Game {

	bag := []Tile{}

	//add 100 tiles to the bag
	for i := 0; i < 20; i++ {
		bag = append(bag, WhiteTile, BlackTile, BlueTile, YellowTile, RedTile)
	}

	shuffle(bag)

	return Game{
		Status:          PENDING,
		Players:         []Player{},
		FactoryDisplays: [][]Tile{},
		TableCentre:     []Tile{FirstPlayerTile},
		TileBag:         bag,
	}
}

func shuffle(deck []Tile) {
	n := 5
	rand.Seed(time.Now().UnixNano())
	for n > 0 {
		rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
		n = n - 1
	}
}

func (g *Game) Start() {
	numOfPlayers := len(g.Players)
	for i := 0; i < numOfPlayers*2+1; i++ {
		g.FactoryDisplays = append(g.FactoryDisplays, []Tile{})
	}

}

type FactoryDisplay struct {
}

func (g *Game) ScoreRound() {

}

func (g Game) isGameOver() bool {
	return false
}
