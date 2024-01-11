package engine

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	updatesPS  = 60
	idealDelta = time.Second / updatesPS
)

type Scene interface {
	Update()
}

type Game struct {
	window *window
	scene  Scene
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

	uiElementIDs := 0
	UI = &ui{
		input:        Input,
		skin:         NewTexture("res/ui9slice.png"),
		currentGenID: &uiElementIDs,
	}

	dispW, dispH = width, height
	ScreenW, ScreenH = width, height

	return &Game{
		window: win,
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
		dispW, dispH = g.window.getSize()
		delta := time.Since(prev).Seconds()
		prev = time.Now()
		acc += delta

		for acc >= idealDelta.Seconds() {
			g.window.processInput()
			g.update()
			Input.update()
			ups++
			acc -= idealDelta.Seconds()
		}

		// Rendering
		Renderer.render()
		g.window.redraw()

		fps++

		// Slow down if running too fast
		if fps > 120 {
			time.Sleep(time.Millisecond * 20)
		}

		if g.window.closed() {
			break
		}
	}
	g.Quit()
}

func (g *Game) SetScene(s Scene) {
	g.scene = s
}

func (g *Game) update() {
	if Input.KeyDown(KeyEscape) {
		g.window.win.SetShouldClose(true)
	}
	g.scene.Update()
}

func (g *Game) Quit() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}
