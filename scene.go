package main

import (
	"fmt"

	"github.com/R-Mckenzie/game-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	engine.Sprite
	c     engine.Collider
	speed float32
}

func NewPlayer() Player {
	texture := engine.NewTexture("res/man.png")
	return Player{
		Sprite: engine.NewSprite(64, 64, 500, 200, 10, texture, nil),
		c:      engine.NewCollider(64, 64, 500, 200),
		speed:  5,
	}
}

func (p *Player) Update(t engine.Tilemap) {
	move := mgl32.Vec3{0, 0, 0}
	if engine.Input.KeyDown(engine.KeyW) {
		move[1] -= 1
	}
	if engine.Input.KeyDown(engine.KeyA) {
		move[0] -= 1
	}
	if engine.Input.KeyDown(engine.KeyS) {
		move[1] += 1
	}
	if engine.Input.KeyDown(engine.KeyD) {
		move[0] += 1
	}
	if move.Len() > 0 {
		move = move.Normalize()
	}

	oldPos := p.Pos
	p.Pos[0] = p.Pos[0] + move[0]*p.speed
	p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	if engine.CollidesMapCollider(t, p.c) {
		p.Pos = oldPos
		p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	}

	oldPos = p.Pos
	p.Pos[1] = p.Pos[1] + move[1]*p.speed
	p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	if engine.CollidesMapCollider(t, p.c) {
		p.Pos = oldPos
		p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	}
}

// =====
type testScene struct {
	game    engine.Game
	p       Player
	camera  engine.Camera2D
	tileMap engine.Tilemap
	font    *engine.Font
	sprite  engine.Sprite
}

var animator engine.Animator

func newScene(game engine.Game) *testScene {
	p := NewPlayer()
	tilemap := engine.LoadTilemap("res/test.tmx", "res/atlas.png", "res/atlas_n.png", 2)
	font, _ := engine.LoadFont("res/ProggyClean.ttf")

	engine.LoadSound("res/music.mp3", "bg")
	engine.LoadSound("res/shot.mp3", "shot")

	run, err := engine.NewImage("res/3 Dude_Monster/Dude_Monster_Run_6.png")
	if err != nil {
		panic(err)
	}

	idle, err := engine.NewImage("res/3 Dude_Monster/Dude_Monster_Idle_4.png")
	if err != nil {
		panic(err)
	}

	runRight := engine.NewAnimation(run, 7, 6, 1, 32, 32, false)
	runLeft := engine.NewAnimation(run, 7, 6, 1, 32, 32, true)
	idleAnim := engine.NewAnimation(idle, 7, 4, 1, 32, 32, false)

	animator = engine.NewAnimator()
	animator.Add(runRight, "run_right")
	animator.Add(runLeft, "run_left")
	animator.Add(idleAnim, "idle")

	engine.LoadShader("shaders/postprocessVertex.glsl", "shaders/funkyEdgesFragment.glsl", "funky lines")
	norm := engine.NewTexture("res/atlascobble_n.png")

	return &testScene{
		game:    game,
		p:       p,
		tileMap: tilemap,
		camera:  engine.NewCamera2D(0, 0),
		sprite:  engine.NewSprite(64, 64, 800, 300, 10, engine.NewTexture("res/man.png"), &norm),
		font:    font,
	}
}

var invOpen bool = false
var isfunky = false

var r, g, b float32 = 0.3, 0.3, 0.3
var f1, f2, f3 float32 = 0.3, 0.3, 0.3
var intensity float32 = 2
var exposure float32 = 1

func (s *testScene) Update() {
	s.p.Update(s.tileMap)
	camX, camY := int(s.p.Pos[0]-400), int(s.p.Pos[1]-300)

	mw, mh := s.tileMap.PixelSize()

	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX+int(engine.ScreenW) > mw {
		camX = mw - int(engine.ScreenW)
	}
	if camY+int(engine.ScreenH) > mh {
		camY = mh - int(engine.ScreenH)
	}

	if engine.Input.KeyOnce(engine.KeyP) {
		engine.LoopSound("bg", -1)
	}
	if engine.Input.KeyOnce(engine.KeyO) {
		engine.PauseLoop("bg")
	}
	if engine.Input.KeyOnce(engine.KeyV) {
		engine.PlaySound("shot", 1)
	}

	if engine.Input.KeyOnce(engine.KeyI) {
		invOpen = !invOpen
	}

	if engine.Input.KeyOnce(engine.KeyC) {
		if isfunky {
			engine.Renderer.SetPostShader("null")
			isfunky = !isfunky
		} else {
			engine.Renderer.SetPostShader("funky lines")
			isfunky = !isfunky
		}
	}

	if engine.Input.KeyDown(engine.KeyA) {
		animator.Trigger("run_left")
	} else if engine.Input.KeyDown(engine.KeyD) {
		animator.Trigger("run_right")
	} else {
		animator.Trigger("idle")
	}

	if animator.Current.Changed {
		s.p.Sprite.SetTexture(animator.Frame())
	}

	s.camera.SetPos(camX, camY)
	engine.Renderer.BeginScene(s.camera, mgl32.Vec3{r, g, b}, exposure)
	engine.Renderer.PushItem(&s.tileMap)
	engine.Renderer.PushItem(s.p)
	engine.Renderer.PushItem(s.sprite)
	engine.Renderer.PushLight(engine.NewLight(s.p.Pos[0], s.p.Pos[1], 50, 1, 1, 1, f1, f2, f3, intensity))
	engine.Renderer.PushLight(engine.NewLight(600, 400, 50, 1, 0, 0, f1, f2, f3, 1))
	engine.Renderer.PushLight(engine.NewLight(600, 600, 50, 0, 0, 1, f1, f2, f3, 1))

	engine.UI.Begin()
	if engine.UI.Button(100, 100, 300, 100, 1, "Button", mgl32.Vec4{0.2, 0.3, 1, 1}) {
		fmt.Printf("clicked\n")
	}
	engine.UI.Label("Test Label xya 1 2 3 4 5", 100, 100, 32, mgl32.Vec4{0, 0, 0, 1})
	engine.UI.End()
}
