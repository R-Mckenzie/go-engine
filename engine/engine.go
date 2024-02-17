package engine

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	updatesPS   = 60
	targetDelta = time.Second / updatesPS
)

type Scene interface {
	Update()
}

type Game struct {
	window *window
	scene  Scene

	quit bool
}

// Engine Systems
var Renderer Renderer2D
var Input *input
var UI *ui

var ScreenW, ScreenH float32
var dispW, dispH float32

// setup the game
func CreateGame(width, height float32) *Game {
	runtime.LockOSThread()

	// Create GLFW window and input
	win := createWindow(int(width), int(height))
	Input = initInput()
	Input.setWindow(win)

	// Init OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 1, 1)

	Renderer = Renderer2DInit(width, height)
	uifont, _ := LoadFont("res/ProggyClean.ttf")

	UI = &ui{
		font:  uifont,
		input: Input,
		skin:  NewTexture("res/ui9slice.png"),
	}

	dispW, dispH = win.getFramebuffer()
	ScreenW, ScreenH = width, height

	return &Game{
		window: win,
		quit:   false,
	}
}

func (g *Game) Run() {
	// Track FPS
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	fps, ups := 0, 0 // increments on every render call

	// sTicker pings every second, when it does print current fps and reset
	go func() {
		for range tick.C {
			fmt.Printf("FPS: %d, UPS: %d\n", fps, ups)
			fps = 0
			ups = 0
			tick.Reset(time.Second * 1)
		}
	}()

	prev := time.Now()
	var acc float64 //LAG

	for {
		dispW, dispH = g.window.getFramebuffer()
		delta := time.Since(prev).Seconds()
		prev = time.Now()
		acc += delta

		g.window.pollEvents()
		for acc >= targetDelta.Seconds() {
			g.update()
			Input.update()
			ups++
			acc -= targetDelta.Seconds()
		}

		// Rendering
		Renderer.render()
		g.window.redraw()

		fps++

		// Slow down if running too fast
		if fps > 120 {
			time.Sleep(time.Millisecond * 20)
		}

		if g.quit {
			g.window.close()
			break
		}
	}

	g.window.Terminate()
	runtime.UnlockOSThread()
}

func (g *Game) SetScene(s Scene) {
	Input.pauseInput(time.Millisecond * 300)
	g.scene = s
}

func (g *Game) update() {
	if Input.KeyDown(KeyEscape) {
		g.Quit()
		return
	}
	g.scene.Update()
}

func (g *Game) Quit() {
	g.quit = true
}
