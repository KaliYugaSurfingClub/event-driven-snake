package snake

import (
	"time"
)

type Direction int8

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Game struct {
	board             Board
	direction         Direction
	directionConsumer chan Direction
	cellsProducer     chan [][]Cell
	isDone            chan struct{}
}

func NewGame() *Game {
	return &Game{
		board:             NewBoard(10, 10),
		directionConsumer: make(chan Direction),
		cellsProducer:     make(chan [][]Cell),
		isDone:            make(chan struct{}),
	}
}

func (g *Game) Start() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case g.direction = <-g.directionConsumer:
			var c Cell //todo
			if g.direction == DirectionUp {
				c = SnakeCell
			} else {
				c = EmptyCell
			}
			g.board.Assign(1, 1, c)
		default:
		}

		g.cellsProducer <- g.board.Cells()
	}

	g.isDone <- struct{}{}
}

func (g *Game) Done() <-chan struct{} {
	return g.isDone
}

func (g *Game) MoveSnake(d Direction) {
	g.directionConsumer <- d
}

func (g *Game) ConsumeCells() <-chan [][]Cell {
	return g.cellsProducer
}
