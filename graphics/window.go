package graphics

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type Window struct {
	glfwWindow *glfw.Window
	width      int
	height     int
}

func CreateWindow(width, height int) Window {
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

	window, err := glfw.CreateWindow(width, height, "Go Game Engine", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	log.Println("Created window")
	return Window{
		glfwWindow: window,
		width:      width,
		height:     height,
	}
}

func (w Window) swapBuffers() {
	w.glfwWindow.SwapBuffers()
	glfw.PollEvents()
}

func (w Window) Closed() bool {
	return w.glfwWindow.ShouldClose()
}
