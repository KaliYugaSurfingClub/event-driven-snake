package snake

type Cell int8

const (
	CellBorder Cell = iota + 1
	CellEmpty
	CellSnake
	CellApple
)
