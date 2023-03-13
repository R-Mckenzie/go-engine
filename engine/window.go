package engine

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type window struct {
	win    *glfw.Window
	Width  int
	Height int

	// These are to fix incorrect rendering on macOS
	windowMoved bool
	moveDir     int
}

func createWindow(width, height int) *window {
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

	win, err := glfw.CreateWindow(width, height, "Go Game Engine", nil, nil)
	if err != nil {
		panic(err)
	}
	win.MakeContextCurrent()
	glfw.SwapInterval(0)

	log.Println("Created window")
	return &window{
		win:         win,
		Width:       width,
		Height:      height,
		windowMoved: false,
		moveDir:     1,
	}
}

func (w *window) getSize() [2]float32 {
	width, height := w.win.GetSize()
	return [2]float32{float32(width), float32(height)}
}

func (w *window) getFramebuffer() [2]float32 {
	width, height := w.win.GetFramebufferSize()
	return [2]float32{float32(width), float32(height)}
}

func (w *window) setTitle(title string) {
	w.win.SetTitle(title)
}

func (w *window) redraw() {
	// These are to fix incorrect rendering on macOS
	if !w.windowMoved {
		x, y := w.win.GetPos()
		w.win.SetPos(x+w.moveDir, y)
		w.moveDir *= -1
		w.windowMoved = true
	}
	w.win.SwapBuffers()
}

func (w *window) processInput() {
	glfw.PollEvents()
}

func (w *window) closed() bool {
	return w.win.ShouldClose()
}

func WindowSize() (int, int) {
	return glfw.GetCurrentContext().GetSize()
}
