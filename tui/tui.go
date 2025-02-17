package tui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"snake-game/snake"
	"strings"
	"unicode"
)

type ConsoleGame struct {
	isDone chan struct{}
	errs   chan error
	game   *snake.Game //todo interface
}

func NewConsoleGame(game *snake.Game) *ConsoleGame {
	cg := ConsoleGame{
		isDone: make(chan struct{}),
		errs:   make(chan error),
		game:   game,
	}

	return &cg
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
		case 'w': //todo keys
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
			case snake.CellBorder:
				sb.WriteRune('#')
			case snake.CellEmpty:
				sb.WriteRune(' ')
			case snake.CellSnake:
				sb.WriteRune('S')
			case snake.CellApple:
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
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("panic: %v", r)
			}
		}()

		cg.game.Start() //todo pass ctx
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("panic: %v", r)
			}
		}()

		cg.RedirectKeyboardToGame(os.Stdin)
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("panic: %v", r)
			}
		}()

		for state := range cg.game.State() {
			cg.DisplayCells(os.Stdout, state.Cells)
		}
	}()

	select {
	case <-ctx.Done():
	case <-cg.game.Done():
	case <-cg.isDone:
	}
}
