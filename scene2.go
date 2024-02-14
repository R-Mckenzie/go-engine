package main

import (
	"fmt"

	"github.com/R-Mckenzie/go-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type scene2 struct {
	game *engine.Game
}

func newScene2(game *engine.Game) *scene2 {
	return &scene2{
		game: game,
	}
}

func (s *scene2) Update() {
	engine.UI.Begin()
	if engine.UI.Button(400, 100, 300, 100, "Button", mgl32.Vec4{1, 0.3, 0.2, 1}) {
		fmt.Printf("clicked scene 2\n")
	}
	engine.UI.End()
}
