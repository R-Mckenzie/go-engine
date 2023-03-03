package engine

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	window *glfw.Window
	Width  int
	Height int

	// These are to fix incorrect rendering on macOS
	windowMoved bool
	moveDir     int
}

func CreateWindow(width, height int) *Window {
	if err := glfw.Init(); err != nil {
		glfw.Terminate()
		panic(err)
	}
	log.Println("Initialised glfw")

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Go Game Engine", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	log.Println("Created window")
	return &Window{
		window:      window,
		Width:       width,
		Height:      height,
		windowMoved: false,
		moveDir:     1,
	}
}

func (w *Window) SetTitle(title string) {
	w.window.SetTitle(title)
}

func (w *Window) Redraw() {
	// These are to fix incorrect rendering on macOS
	if !w.windowMoved {
		x, y := w.window.GetPos()
		w.window.SetPos(x+w.moveDir, y)
		w.moveDir *= -1
		w.windowMoved = true
	}
	w.window.SwapBuffers()
}

func (w *Window) processInput() {
	glfw.PollEvents()
	Input()
}

func (w *Window) Closed() bool {
	return w.window.ShouldClose()
}
