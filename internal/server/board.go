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
	if y > 5 {
		return errors.New("Row out of range")
	}

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

func (b *Board) StageTiles(tiles []Tile, y uint8) error {
	colors := map[Tile]bool{}

	for _, t := range tiles {
		if t == FirstPlayerTile {
			continue
		}
		colors[t] = true
	}
	if len(colors) > 1 {
		return errors.New("Cannot stage tiles of different types")
	}

	for _, t := range tiles {
		if err := b.StageTile(t, y); err != nil {
			return err
		}
	}

	return nil
}

// // Returns the score of placing that tile or an error if it is an invalid operation
// func (b *Board) PlaceTile(t Tile, x uint8, y uint8) (uint8, error) {
// 	if b.GetTile(x, y) != EmptyTile {
// 		return 0, errors.New("Position is not empty")
// 	}
// 	if b.GetBluePrint()[y][x] != t {
// 		return 0, errors.New("Invalid tile position")
// 	}
// 	b.SetTile(t, x, y)

// 	return b.ScoreTile(x, y), nil
// }

// func (b *Board) GetTile(x uint8, y uint8) Tile {
// 	return (*b)[y][x]
// }

// func (b *Board) SetTile(t Tile, x uint8, y uint8) {
// 	(*b)[y][x] = t
// }

// func (b *Board) ScoreTile(x uint8, y uint8) uint8 {
// 	var scoreX uint8
// 	var scoreY uint8

// 	// score row
// 	for i, t := range (*b)[y] {
// 		// handle empty tile
// 		if t == EmptyTile {
// 			if uint8(i) < x {
// 				scoreX = 0
// 				continue
// 			} else {
// 				break
// 			}
// 		}
// 		scoreX += 1
// 	}
// 	// score column

// 	for i, r := range *b {
// 		t := r[x]

// 		// handle empty tile
// 		if t == EmptyTile {
// 			if uint8(i) < y {
// 				scoreY = 0
// 				continue
// 			} else {
// 				break
// 			}
// 		}
// 		scoreY += 1

// 	}
// 	return scoreX + scoreY
// }

func (b *Board) GetBluePrint() [][]Tile {
	return [][]Tile{
		{BlueTile, YellowTile, RedTile, BlackTile, WhiteTile},
		{WhiteTile, BlueTile, YellowTile, RedTile, BlackTile},
		{BlackTile, WhiteTile, BlueTile, YellowTile, RedTile},
		{RedTile, BlackTile, WhiteTile, BlueTile, YellowTile},
		{YellowTile, RedTile, BlackTile, WhiteTile, BlueTile},
	}
}
