package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"snake-game/snake"
	"strings"
	"unicode"
)

type ConsoleGame struct {
	isDone chan struct{}
	errs   chan error
	game   *snake.Game //todo interface
	//todo listen terminal size changes
}

func NewConsoleGame(game *snake.Game) (*ConsoleGame, error) {
	//todo game box must fit in terminal box

	cg := ConsoleGame{
		isDone: make(chan struct{}),
		errs:   make(chan error),
		game:   game,
	}

	return &cg, nil
}

func (cg *ConsoleGame) RedirectKeyboardToGame(r io.Reader) {
	buf := bufio.NewReaderSize(r, unicode.MaxRune)

	for {
		key, _, err := buf.ReadRune()
		if err != nil {
			cg.errs <- fmt.Errorf("ConsoleGame.RedirectKeyboardToGame: cannot read rune: %w", err)
			continue
		}

		switch key {
		case 'w':
			cg.game.MoveSnake(snake.DirectionUp)
		case 'd':
			cg.game.MoveSnake(snake.DirectionRight)
		case 's':
			cg.game.MoveSnake(snake.DirectionDown)
		case 'a':
			cg.game.MoveSnake(snake.DirectionLeft)
		case 27: //esc
			cg.isDone <- struct{}{}
		}

		buf.Reset(r)
	}
}

func (cg *ConsoleGame) DisplayCells(w io.Writer, cells [][]snake.Cell) {
	var sb strings.Builder

	for _, row := range cells {
		for _, cell := range row {
			switch cell {
			case snake.BorderCell:
				sb.WriteRune('#')
			case snake.EmptyCell:
				sb.WriteRune(' ')
			case snake.SnakeCell:
				sb.WriteRune('S')
			case snake.AppleCell:
				sb.WriteRune('A')
			}
		}
		sb.WriteRune('\n')
	}

	if _, err := fmt.Fprint(w, sb.String()); err != nil {
		cg.errs <- fmt.Errorf("ConsoleGame.DispalyBoard: cannot write cells: %w", err)
	}
}

func (cg *ConsoleGame) Start(ctx context.Context) {
	go cg.game.Start()

	go cg.RedirectKeyboardToGame(os.Stdin)

	go func() {
		for cells := range cg.game.ConsumeCells() {
			cg.DisplayCells(os.Stdout, cells)
		}
	}()

	select {
	case <-ctx.Done():
	case <-cg.game.Done():
	case <-cg.isDone:
	}
}
