package main

import (
	"github.com/R-Mckenzie/game-engine/engine"
)

const width, height = 32.0 * 30, 32.0 * 20

func main() {
	game := engine.CreateGame(width, height)
	s := newScene()
	game.SetScene(s)

	game.Run()
}
