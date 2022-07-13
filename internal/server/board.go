package server

import "errors"

type Tile uint8

const (
	EmptyTile Tile = iota
	WhiteTile
	BlackTile
	RedTile
	BlueTile
	YellowTile
	FirstPlayerTile
)

type Board struct {
	Tiles  [][]Tile `json:"tiles"`
	Staged [][]Tile `json:"staged"`
}

func NewBoard() Board {
	tiles := [][]Tile{}
	for i := 0; i < 5; i++ {
		r := []Tile{}
		for j := 0; j < 5; j++ {
			r = append(r, EmptyTile)
		}
		tiles = append(tiles, r)
	}

	staged := [][]Tile{}

	for i := 0; i < 6; i++ {
		staged = append(staged, []Tile{})
	}

	return Board{
		Tiles:  tiles,
		Staged: staged,
	}
}

func (b *Board) StageNegativeTile(t Tile) {
	if len(b.Staged[5]) < 8 {

		b.Staged[5] = append(b.Staged[5], t)
	}
}

func (b *Board) StageTile(t Tile, y uint8) error {

	// stage negatives row
	if y == 5 || t == FirstPlayerTile {
		b.StageNegativeTile(t)
		return nil
	}

	r := &b.Staged[y]
	size := len(*r)

	if size > 0 && (*r)[0] != t {
		return errors.New("Row type is mismatched")
	}

	free := y + 1 - uint8(size)

	if free > 0 {
		*r = append(*r, t)
	} else {
		b.StageNegativeTile(t)
	}
	return nil
}

// StageTiles places tiles on the y row in the staging area
func (b *Board) StageTiles(tiles []Tile, y uint8) error {
	if y < 0 || y > 5 {
		return errors.New("Row out of range")
	}

	// early exit if no tiles
	if len(tiles) == 0 {
		return nil
	}

	var color Tile
	colors := map[Tile]bool{}

	// check all tiles are the same type
	for _, t := range tiles {
		if t == FirstPlayerTile {
			continue
		}
		color = t
		colors[t] = true
	}

	if len(colors) > 1 {
		return errors.New("Cannot stage tiles of multiple types on same row")
	}

	// validate non-negative rows
	if y < 5 {
		// check staged row is not another color
		for _, j := range b.Staged[y] {
			if j != color {
				return errors.New("Cannot stage multiple colours in same row")
			}
		}

		// check board has available place
		for _, p := range b.Tiles[y] {
			if p == color {
				return errors.New("Tile alredy placed on board")
			}
		}
	}

	for _, t := range tiles {
		if err := b.StageTile(t, y); err != nil {
			return err
		}
	}

	return nil
}

// ScoreStaged updates the tile board and returns points awarded
func (b *Board) ScoreStaged() (int8, []Tile) {
	var score int8
	removedTiles := []Tile{}

	for i := 0; i < len(b.Staged); i++ {
		s, rip := b.scoreRow(uint8(i))
		score += s
		removedTiles = append(removedTiles, rip...)
	}

	return score, removedTiles
}

// scoreRow updates the tile board and returns the points awarded for the row
// as well as the discarded tiles
func (b *Board) scoreRow(y uint8) (int8, []Tile) {
	r := b.Staged[y]

	// score negative
	if y == 5 {
		var score int8
		switch len(r) {
		case 0:
			score = 0
		case 1:
			score = -1
		case 2:
			score = -2
		case 3:
			score = -4
		case 4:
			score = -6
		case 5:
			score = -8
		case 6:
			score = -11
		default:
			score = -14
		}
		d := []Tile{}
		for _, t := range b.Staged[5] {
			if t != FirstPlayerTile {
				d = append(d, t)
			}
		}
		b.Staged[5] = []Tile{}
		return score, d
	}

	max := int(y + 1)

	// ignore row if it is incomplete
	if len(r) < max {
		return 0, []Tile{}
	}

	var x uint8
	col := b.Staged[y][0]

	m := b.GetBluePrint()

	for i, v := range m[y] {
		if v == col {
			x = uint8(i)
		}
	}
	// update tile board
	b.Tiles[y][x] = col
	// reset staged
	d := b.Staged[y]
	b.Staged[y] = []Tile{}
	return b.ScoreTile(x, y), d

}

// ScoreTile returns the score of placing a tile on the board
func (b *Board) ScoreTile(x uint8, y uint8) int8 {
	var scoreX int8
	var scoreY int8

	// score row
	for i, t := range b.Tiles[y] {
		// handle empty tile
		if t == EmptyTile {
			if uint8(i) < x {
				scoreX = 0
				continue
			} else {
				break
			}
		}
		scoreX += 1
	}

	// score column
	for i, r := range b.Tiles {
		t := r[x]
		// handle empty tile
		if t == EmptyTile {
			if uint8(i) < y {
				scoreY = 0
				continue
			} else {
				break
			}
		}
		scoreY += 1

	}
	return scoreX + scoreY
}

func (b *Board) GetBluePrint() [][]Tile {
	return [][]Tile{
		{BlueTile, YellowTile, RedTile, BlackTile, WhiteTile},
		{WhiteTile, BlueTile, YellowTile, RedTile, BlackTile},
		{BlackTile, WhiteTile, BlueTile, YellowTile, RedTile},
		{RedTile, BlackTile, WhiteTile, BlueTile, YellowTile},
		{YellowTile, RedTile, BlackTile, WhiteTile, BlueTile},
	}
}
