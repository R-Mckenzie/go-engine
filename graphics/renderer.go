package graphics

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
	vao     uint32
	win     glfw.Window
	shader  Shader
	texture Texture
	camera  Camera
}

func NewRenderer(w glfw.Window) *Renderer {
	width, height := w.GetSize()
	initOpenGL(int32(width), int32(height))

	shader, err := NewShader(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		panic(err)
	}

	texture, err := LoadImage("res/man.png")
	if err != nil {
		panic(err)
	}

	renderer := &Renderer{
		vao:     makeVAO(quad(32, 32)),
		win:     w,
		shader:  shader,
		texture: texture,
		camera:  NewCamera2D(-300, 0),
	}
	return renderer
}

var rotation float32 = 0

func (r *Renderer) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	r.shader.Use()

	w, h := r.win.GetSize()
	pMat := mgl32.Ortho(0, float32(w), float32(h), 0, 0.1, 10)
	vMat := mgl32.Ident4()
	vMat = mgl32.Translate3D(-r.camera.X, -r.camera.Y, -3)

	mMat := mgl32.Ident4()
	translation := mgl32.Translate3D(100, 100, 0)
	rotMat := mgl32.HomogRotate3DZ(mgl32.DegToRad(rotation))
	mMat = translation.Mul4(rotMat)

	rotation += 10

	r.shader.SetMatrix("u_model", mMat)
	r.shader.SetMatrix("u_view", vMat)
	r.shader.SetMatrix("u_projection", pMat)

	gl.ActiveTexture(gl.TEXTURE0)
	r.texture.Use()
	gl.BindVertexArray(r.vao)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	r.win.SwapBuffers()
}

func initOpenGL(w, h int32) {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	log.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	gl.Viewport(0, 0, w, h)
	gl.ClearColor(0.5, 0.5, 1, 1)
}

func quad(width, height float32) ([]float32, []uint32) {
	w2, h2 := width/2, height/2
	return []float32{ // vertices
			-w2, -h2, 0.0, 0.0, 1.0,
			w2, -h2, 0.0, 1.0, 1.0,
			w2, h2, 0.0, 1.0, 0.0,
			-h2, h2, 0.0, 0.0, 0.0,
		}, []uint32{ // indices
			0, 1, 3,
			1, 2, 3,
		}
}

func makeVAO(p []float32, i []uint32) uint32 {
	var vbo, vao, ebo uint32

	// Create GL objects
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(p), gl.Ptr(p), gl.STATIC_DRAW) // The 4 represents the 4 bytes per 32 element in array

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(i), gl.Ptr(i), gl.STATIC_DRAW) // The 4 represents the 4 bytes per 32 element in array

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*5, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return vao
}
