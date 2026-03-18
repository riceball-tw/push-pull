package main

import (
	"log"

	"push-pull/internal/game"
)

func main() {
	if err := game.Run(); err != nil {
		log.Fatal(err)
	}
}
