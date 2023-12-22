package engine

import (
	"log"
	"runtime"
	"time"

	imgui "github.com/AllenDang/cimgui-go"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	updatesPS  = 60
	idealDelta = time.Second / updatesPS
)

type Scene interface {
	Update()
	DebugGUI()
}

type Game struct {
	window *window
	imgui  *imguiRenderer
	scene  Scene
}

var DispW, DispH float32

var ctx imgui.ImGuiContext

// setup the game
func CreateGame(width, height float32) *Game {
	ctx = imgui.CreateContext(0)
	runtime.LockOSThread()

	// Create GLFW window and input
	win := createWindow(int(width), int(height))
	Input().setWindow(win)

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

	// Init 2d renderer
	Renderer2DInit(width, height)

	// Init imgui renderer
	imguiRenderer := NewImguiRenderer()

	DispW, DispH = width, height

	game := &Game{
		window: win,
		imgui:  imguiRenderer,
	}

	return game
}

func (g *Game) Run() {
	// Track FPS
	tick := time.NewTicker(time.Second * 1)
	defer tick.Stop()
	fps, ups := 0, 0 // increments on every render call

	// sTicker pings every second, when it does print current fps and reset
	go func() {
		for range tick.C {
			AddDebugInfo("FPS", fps)
			AddDebugInfo("UPS", ups)
			fps = 0
			ups = 0
			tick.Reset(time.Second * 1)
		}
	}()

	prev := time.Now()
	var acc float64 //LAG

	for {
		delta := time.Since(prev).Seconds()
		prev = time.Now()
		acc += delta

		// Imgui prep
		Input().imguiPrepFrame(float32(delta))
		imgui.NewFrame()

		for acc >= idealDelta.Seconds() {
			g.window.processInput()
			g.update()
			Input().update()
			ups++
			acc -= idealDelta.Seconds()
		}

		// Rendering
		Renderer2D().render()

		// Imgui
		displayDebug()
		g.scene.DebugGUI()
		imgui.Render()
		g.imgui.Render(g.window.getSize(), g.window.getFramebuffer(), imgui.GetDrawData())
		g.window.redraw()

		fps++

		// Slow down if running too fast
		if fps > 120 {
			time.Sleep(time.Millisecond * 5)
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
	if Input().KeyDown(KeyEscape) {
		g.window.win.SetShouldClose(true)
	}
	g.scene.Update()
}

func (g *Game) Quit() {
	g.imgui.dispose()
	glfw.Terminate()
	ctx.Destroy()
	runtime.UnlockOSThread()
}
