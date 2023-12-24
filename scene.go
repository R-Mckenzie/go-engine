package main

import (
	"github.com/R-Mckenzie/game-engine/engine"
	"github.com/go-gl/mathgl/mgl32"

	imgui "github.com/AllenDang/cimgui-go"
)

type Player struct {
	engine.Sprite
	c     engine.Collider
	speed float32
}

func NewPlayer() Player {
	texture := engine.NewTexture("res/man.png")
	return Player{
		Sprite: engine.NewSprite(64, 64, 500, 200, 5, texture, nil),
		c:      engine.NewCollider(64, 64, 500, 200),
		speed:  5,
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
	p       Player
	camera  engine.Camera2D
	tileMap engine.Tilemap
	font    *engine.Font
}

var animator engine.Animator

func newScene() *testScene {
	p := NewPlayer()
	tilemap := engine.LoadTilemap("res/test.tmx", "res/atlas.png", "res/atlas_n.png", 2)
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

	runRight := engine.NewAnimation(run, 7, 6, 1, 32, 32, false)
	runLeft := engine.NewAnimation(run, 7, 6, 1, 32, 32, true)
	idleAnim := engine.NewAnimation(idle, 7, 4, 1, 32, 32, false)

	animator = engine.NewAnimator()
	animator.Add(runRight, "run_right")
	animator.Add(runLeft, "run_left")
	animator.Add(idleAnim, "idle")

	engine.LoadShader("shaders/postprocessVertex.glsl", "shaders/funkyEdgesFragment.glsl", "funky lines")

	return &testScene{
		p:       p,
		tileMap: tilemap,
		camera:  engine.NewCamera2D(0, 0),
		font:    font,
	}
}

var invOpen bool = false
var isfunky = false

var r, g, b float32 = 0.3, 0.3, 0.3
var intensity int32 = 30
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

	if engine.Input().KeyOnce(engine.KeyP) {
		engine.LoopSound("bg", -1)
	}
	if engine.Input().KeyOnce(engine.KeyO) {
		engine.PauseLoop("bg")
	}
	if engine.Input().KeyOnce(engine.KeyV) {
		engine.PlaySound("shot", 1)
	}

	if engine.Input().KeyOnce(engine.KeyI) {
		invOpen = !invOpen
	}

	if engine.Input().KeyOnce(engine.KeyC) {
		if isfunky {
			engine.Renderer2D().SetPostShader("null")
			isfunky = !isfunky
		} else {
			engine.Renderer2D().SetPostShader("funky lines")
			isfunky = !isfunky
		}
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
	engine.Renderer2D().SetExposure(exposure)
	engine.Renderer2D().BeginScene(s.camera, mgl32.Vec3{r, g, b})
	engine.Renderer2D().PushItem(s.tileMap.StaticRenderItem())
	engine.Renderer2D().PushItem(s.tileMap.AnimatedRenderItem())
	engine.Renderer2D().PushItem(s.p.RenderItem())
	engine.Renderer2D().PushLight(engine.NewLight(s.p.Pos[0], s.p.Pos[1], 50, 1, 1, 1, 1))
	engine.Renderer2D().PushLight(engine.NewLight(600, 400, 50, 1, 0, 0, 1))
	engine.Renderer2D().PushLight(engine.NewLight(600, 600, 50, 0, 0, 1, 1))
}

func (s *testScene) DebugGUI() {
	imgui.SetNextWindowSize(imgui.ImVec2{X: 400, Y: 200}, imgui.ImGuiCond(imgui.ImGuiCond_FirstUseEver))
	imgui.Begin("SCENE DATA", nil, 0)
	if imgui.Button("shoot", imgui.NewImVec2(0, 0)) {
		engine.PlaySound("shot", 1)
	}
	imgui.ColorEdit3("lighting", [3]*float32{&r, &g, &b}, 0)
	imgui.SliderFloat("speed", &s.p.speed, 1, 10, "%.3f", 0)
	imgui.SliderFloat("exposure", &exposure, 0, 10, "%.3f", 0)
	imgui.SliderInt("intensity", &intensity, 1, 200, "%d", 0)
	imgui.End()
}
