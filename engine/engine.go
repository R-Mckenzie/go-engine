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
	idealFPS   = 60
	idealDelta = time.Second / idealFPS
)

type Entity interface {
	Update()
}

type Game struct {
	window *Window
	scene  Entity
}

// setup the game
func CreateGame(width, height float32) *Game {
	runtime.LockOSThread()
	win := CreateWindow(800, 600)
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.ClearColor(0.5, 0.5, 1, 1)
	Input().setWindow(win)
	Renderer2DInit(width, height)
	game := &Game{
		window: win,
	}
	return game
}

func (g *Game) Run() {
	// Track FPS
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	fps := 0 // increments on every render call
	fpsChan := make(chan int)

	// sTicker pings every second, when it does print current fps and reset
	go func() {
		for range tick.C {
			fpsChan <- fps
			fps = 0
			tick.Reset(time.Second * 1)
		}
	}()

	prev := time.Now()
	var acc float64 //LAG

	for {
		delta := time.Since(prev).Seconds()
		prev = time.Now()
		acc += delta
		g.window.processInput()
		for acc >= idealDelta.Seconds() {
			g.update()
			g.window.Redraw()
			acc -= idealDelta.Seconds()
		}
		Renderer2D().render()
		fps++
		// Slow down if running too fast
		if fps > 120 {
			time.Sleep(time.Millisecond)
		}

		// glfwSetTitle needs to run on main thread, so we set it here instead of the goroutine
		select {
		case val := <-fpsChan:
			g.window.SetTitle(fmt.Sprintf("Go Game Engine | FPS: %d", val))
		default:
		}

		if g.window.Closed() {
			break
		}
	}
	g.Quit()
}

func (g *Game) SetScene(s Entity) {
	g.scene = s
}

func (g *Game) update() {
	glfw.PollEvents()
	if Input().KeyDown(KeyEscape) {
		g.window.window.SetShouldClose(true)
	}
	g.scene.Update()
}

func (g *Game) Quit() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}
