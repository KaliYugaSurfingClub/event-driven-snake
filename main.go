package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"snake-game/snake"
	"snake-game/tui"
	"time"
)

func main() {
	logs, _ := os.OpenFile("log/logs.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)

	log.SetOutput(logs)
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %v", r)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	_, rollback, err := SetTerminalRowMod()
	if err != nil {
		log.Fatal(err)
	}
	defer rollback()

	//catch panic and do rollback

	ticker := snake.NewTicker(time.Second, time.Second/6, 0.9)
	game := snake.NewGame(11, 11, snake.DirectionUp, *ticker)
	consoleGame := tui.NewConsoleGame(game)

	consoleGame.Start(ctx)

	fmt.Print("good bay")
}
