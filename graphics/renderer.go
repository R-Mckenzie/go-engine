package graphics

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer struct {
	prog uint32
	win  Window
}

func NewRenderer(w Window) *Renderer {
	p := initOpenGL(int32(w.width), int32(w.height))
	return &Renderer{
		prog: p,
		win:  w,
	}
}

func (r *Renderer) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.prog)
	vao := makeVAO(triangle())
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle())/3))
	r.win.swapBuffers()
}

func initOpenGL(w, h int32) uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	gl.ClearColor(1, 0, 0, 1)

	vShader := LoadVertexShader(vertexShaderSource)
	fShader := LoadFragmentShader(fragmentShaderSource)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vShader)
	gl.AttachShader(prog, fShader)
	gl.LinkProgram(prog)
	log.Println("OpenGL initialised")
	return prog
}

func triangle() []float32 {
	return []float32{
		0.0, 0.5, 0.0, //top
		-0.5, -0.5, 0.0, //left
		0.5, -0.5, 0.0, //right
	}
}

func makeVAO(p []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(p), gl.Ptr(p), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return vao
}
