package main

import (
	"github.com/R-Mckenzie/game-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	engine.Transform
	engine.Sprite
}

func NewPlayer() Player {
	texture := engine.NewTexture("res/man.png")
	return Player{
		Transform: engine.NewTransform(100, 100, 0),
		Sprite:    engine.NewSprite(64, 64, 100, 100, texture),
	}
}

func (p *Player) Update() {
	move := mgl32.Vec3{0, 0, 0}
	if engine.Input().KeyDown(engine.KeyW) {
		move[1] -= 1
	}
	if engine.Input().KeyDown(engine.KeyA) {
		move[0] -= 1
	}
	if engine.Input().KeyDown(engine.KeyS) {
		move[1] += 1
	}
	if engine.Input().KeyDown(engine.KeyD) {
		move[0] += 1
	}

	if move.Len() > 0 {
		move = move.Normalize()
	}

	p.Pos = p.Transform.Pos.Add(move.Mul(10))
	p.Sprite.Transform = p.Transform
}

type testScene struct {
	p      Player
	camera engine.Camera
	tiles  []engine.Sprite
}

func newScene() *testScene {
	p := NewPlayer()
	atlas, _ := engine.LoadImage("res/atlas.png")
	water := engine.TextureFromAtlas(atlas, 0, 0, 32, 32)
	grass := engine.TextureFromAtlas(atlas, 32, 0, 32, 32)
	sand := engine.TextureFromAtlas(atlas, 64, 0, 32, 32)
	dirt := engine.TextureFromAtlas(atlas, 96, 0, 32, 32)
	placeholder := engine.TextureFromAtlas(atlas, 128, 0, 32, 32)
	stone := engine.TextureFromAtlas(atlas, 0, 32, 32, 32)

	tiles := []engine.Sprite{
		engine.NewSprite(32, 32, 100, 100, water),
		engine.NewSprite(32, 32, 300, 400, grass),
		engine.NewSprite(32, 32, 500, 10, placeholder),
		engine.NewSprite(32, 32, 234, 200, stone),
		engine.NewSprite(32, 32, 634, 500, dirt),
		engine.NewSprite(32, 32, 634, 300, sand),
		engine.NewSprite(32, 32, 400, 400, water),
	}
	return &testScene{
		p:      p,
		tiles:  tiles,
		camera: engine.NewCamera2D(0, 0),
	}
}

func (s *testScene) Update() {
	s.p.Update()
	engine.Renderer2D().BeginScene(s.camera)
	for _, t := range s.tiles {
		engine.Renderer2D().PushItem(t.RenderItem())
	}
	engine.Renderer2D().PushItem(s.p.RenderItem())
}
