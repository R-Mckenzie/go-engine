package main

import (
	"github.com/R-Mckenzie/game-engine/engine"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	engine.Sprite
	c engine.Collider
}

func NewPlayer() Player {
	texture := engine.NewTexture("res/man.png")
	return Player{
		Sprite: engine.NewSprite(64, 64, 500, 200, 5, texture),
		c:      engine.NewCollider(64, 64, 500, 200),
	}
}

func (p *Player) Update(t engine.Tilemap) {
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

	speed := float32(5.0)

	oldPos := p.Pos
	p.Pos[0] = p.Pos[0] + move[0]*speed
	p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	if engine.CollidesMapCollider(t, p.c) {
		p.Pos = oldPos
		p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	}

	oldPos = p.Pos
	p.Pos[1] = p.Pos[1] + move[1]*speed
	p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	if engine.CollidesMapCollider(t, p.c) {
		p.Pos = oldPos
		p.c = p.c.SetPos(int(p.Pos[0]), int(p.Pos[1]))
	}
}

// =====
type testScene struct {
	p       Player
	camera  engine.Camera2D
	tiles   []engine.Sprite
	tileMap engine.Tilemap
	font    *engine.Font
}

var animator engine.Animator

func newScene() *testScene {
	p := NewPlayer()
	tilemap := engine.LoadTilemap("res/test.tmx", "res/atlas.png", 2)
	font, _ := engine.LoadFont("res/ProggyClean.ttf", 64)

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

	runRight := engine.NewAnimation(run, 5, 6, 1, 32, 32, false)
	runLeft := engine.NewAnimation(run, 5, 6, 1, 32, 32, true)
	idleAnim := engine.NewAnimation(idle, 4, 4, 1, 32, 32, false)

	animator = engine.NewAnimator()
	animator.Add(runRight, "run_right")
	animator.Add(runLeft, "run_left")
	animator.Add(idleAnim, "idle")

	return &testScene{
		p:       p,
		tileMap: tilemap,
		camera:  engine.NewCamera2D(0, 0),
		font:    font,
	}
}

func (s *testScene) Update() {
	s.p.Update(s.tileMap)
	camX, camY := int(s.p.Pos[0]-400), int(s.p.Pos[1]-300)

	mw, mh := s.tileMap.PixelSize()
	ww, wh := engine.WindowSize()

	if camX < 0 {
		camX = 0
	}
	if camY < 0 {
		camY = 0
	}
	if camX+ww > mw {
		camX = mw - ww
	}
	if camY+wh > mh {
		camY = mh - wh
	}

	if engine.Input().KeyUp(engine.KeyP) {
		engine.LoopSound("bg", -1)
	}
	if engine.Input().KeyUp(engine.KeyO) {
		engine.PauseLoop("bg")
	}
	if engine.Input().KeyUp(engine.KeyI) {
		engine.StopLoop("bg")
	}

	if engine.Input().KeyUp(engine.KeyV) {
		engine.PlaySound("shot", 1)
	}

	if engine.Input().KeyDown(engine.KeyA) {
		animator.Trigger("run_left")
	} else if engine.Input().KeyDown(engine.KeyD) {
		animator.Trigger("run_right")
	} else {
		animator.Trigger("idle")
	}

	if animator.Current.Changed {
		s.p.Sprite.SetTexture(animator.Frame())
	}

	s.camera.SetPos(camX, camY)
	engine.Renderer2D().BeginScene(s.camera)
	engine.Renderer2D().PushItem(s.tileMap.RenderItem())
	engine.Renderer2D().PushItem(s.p.RenderItem())
	s.font.Print(100, 0, "Hello World")
}
