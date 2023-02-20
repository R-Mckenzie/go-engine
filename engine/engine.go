package engine

import (
	"log"
	"runtime"
	"time"

	"github.com/R-Mckenzie/game-engine/graphics"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	idealFPS   = 60
	idealDelta = time.Second / idealFPS
)

type Game struct {
	window   graphics.Window
	renderer *graphics.Renderer
}

// setup the game
func NewGame() *Game {
	runtime.LockOSThread()
	win := graphics.CreateWindow(800, 600)
	return &Game{
		window:   win,
		renderer: graphics.NewRenderer(win),
	}
}

func (g *Game) Run() {
	// Track FPS
	fps := 0 // increments on every render call
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()

	// sTicker pings every second, when it does print current fps and reset
	go func() {
		for range tick.C {
			log.Printf("FPS: %d\n", fps)
			fps = 0
			tick.Reset(time.Second * 1)
		}
	}()

	prev := time.Now()
	var acc float64

	for {
		if g.window.Closed() {
			break
		}

		delta := time.Since(prev).Seconds()
		prev = time.Now()
		acc += delta

		if acc >= idealDelta.Seconds() {
			for acc >= idealDelta.Seconds() {
				g.update()
				g.renderer.Draw()
				acc -= idealDelta.Seconds()
				fps++
			}
		} else {
			time.Sleep(time.Millisecond * 1)
		}
	}
	g.Quit()
}

func (g *Game) update() {

}

func (g *Game) Quit() {
	glfw.Terminate()
	runtime.UnlockOSThread()
}
