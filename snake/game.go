package snake

import (
	"math/rand"
	"slices"
	"time"
)

type Game struct {
	snake         []Point
	board         [][]Cell
	direction     Direction
	ticker        Ticker
	directionChan chan Direction
	stateChan     chan State
	isDone        chan struct{}
}

func NewGame(width, height int, startDirection Direction, ticker Ticker) *Game {
	board := make([][]Cell, height)

	for y := range board {
		board[y] = make([]Cell, width)
		for x := range board[y] {
			if x == 0 || y == 0 || x == width-1 || y == height-1 {
				board[y][x] = CellBorder
			} else {
				board[y][x] = CellEmpty
			}
		}
	}

	return &Game{
		direction:     startDirection,
		board:         board,
		ticker:        ticker,
		directionChan: make(chan Direction),
		stateChan:     make(chan State),
		isDone:        make(chan struct{}),
	}
}

func (g *Game) Start() {
	defer close(g.stateChan)
	defer close(g.isDone)

	g.generateApple()
	g.generateSnakePosition()

	g.produceState()
	g.waitFirstMove()

	for {
		g.readDirection()

		head := g.snake[len(g.snake)-1]
		tail := g.snake[0]

		switch g.direction {
		case DirectionUp:
			head.Y -= 1
		case DirectionRight:
			head.X += 1
		case DirectionDown:
			head.Y += 1
		case DirectionLeft:
			head.X -= 1
		}

		switch g.board[head.Y][head.X] {
		case CellEmpty:
			g.board[head.Y][head.X] = CellSnake
			g.board[tail.Y][tail.X] = CellEmpty
			g.snake = g.snake[1:]
		case CellApple:
			g.board[head.Y][head.X] = CellSnake
			g.generateApple()
			g.ticker.ReduceInterval()
		case CellBorder, CellSnake:
			return
		}

		g.snake = append(g.snake, head)
		g.produceState()

		<-time.After(g.ticker.Interval())
	}
}

func (g *Game) readDirection() {
	select {
	case d := <-g.directionChan:
		if !isOppositeDirections(d, g.direction) || len(g.snake) == 1 {
			g.direction = d
		}
	default:
	}
}

func (g *Game) generateApple() {
	empty := make([]Point, 0, len(g.board)*len(g.board[0]))

	for y := range g.board {
		for x, c := range g.board[y] {
			if c == CellEmpty {
				empty = append(empty, Point{x, y})
			}
		}
	}

	n := rand.Intn(len(empty) - 1)
	p := empty[n]

	g.board[p.Y][p.X] = CellApple
}

func (g *Game) generateSnakePosition() {
	p := randomPoint(len(g.board)-2, len(g.board[0])-2)
	g.board[p.Y][p.X] = CellSnake
	g.snake = append(g.snake, p)
}

func (g *Game) produceState() {
	cells := make([][]Cell, len(g.board))
	for i, row := range g.board {
		cells[i] = slices.Clone(row)
	}

	speed := 1 / float32(g.ticker.Interval().Seconds())

	state := State{
		Cells:    cells,
		Speed:    speed,
		SnakeLen: len(g.snake),
	}

	g.stateChan <- state
}

func (g *Game) waitFirstMove() {
	g.direction = <-g.directionChan
}

func (g *Game) Done() <-chan struct{} {
	return g.isDone
}

func (g *Game) MoveSnake(d Direction) {
	g.directionChan <- d
}

func (g *Game) State() <-chan State {
	return g.stateChan
}

// todo use certain seed
func randomPoint(width, height int) Point {
	return Point{
		X: rand.Intn(width) + 1,
		Y: rand.Intn(height) + 1,
	}
}
