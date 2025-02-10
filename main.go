package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"snake-game/snake"
)

func main() {
	game := snake.NewGame()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	_, rollback, err := SetTerminalRowMod()
	if err != nil {
		log.Fatal(err)
	}
	defer rollback()

	consoleGame, err := NewConsoleGame(game)
	if err != nil {
		log.Fatal(err)
	}

	consoleGame.Start(ctx)
}
