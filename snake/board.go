package snake

import (
	"slices"
)

type Cell int8

const (
	BorderCell Cell = iota
	EmptyCell
	SnakeCell
	AppleCell
)

type Board struct {
	cells [][]Cell
}

func NewBoard(width, height int) Board {
	//for border
	width += 2
	height += 2

	cells := make([][]Cell, height)

	for i := range cells {
		if i == 0 || i == height-1 {
			cells[i] = fill(BorderCell, width)
		} else {
			cells[i] = slices.Concat([]Cell{BorderCell}, fill(EmptyCell, width-2), []Cell{BorderCell})
		}
	}

	return Board{
		cells: cells,
	}
}

func (b Board) Width() int {
	return len(b.cells[0])
}

func (b Board) Height() int {
	return len(b.cells)
}

func (b Board) Assign(x, y int, c Cell) bool {
	//user cannot change boraders
	if x > b.Width()-2 {
		return false
	}
	if y > b.Height()-2 {
		return false
	}

	y = len(b.cells) - y - 1
	b.cells[y][x] = c

	return true
}

func (b Board) At(x, y int) (Cell, bool) {
	if x > b.Width() {
		return 0, false
	}
	if y > b.Height() {
		return 0, false
	}

	// Flip the Y-axis to match the coordinate system where (0,0) is at the bottom-left.
	y = len(b.cells) - y - 1

	return b.cells[y][x], true
}

func (b Board) Cells() [][]Cell {
	res := make([][]Cell, len(b.cells))

	for i := range b.cells {
		res[i] = slices.Clone(b.cells[i])
	}

	return res
}

func fill(c Cell, width int) []Cell {
	row := make([]Cell, width)
	for i := range row {
		row[i] = c
	}
	return row
}
