package engine

import (
	"log"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type renderer struct {
	renderBuffer []renderItem
	activeCam    Camera
	projection   mgl32.Mat4
}

type renderItem struct {
	vao       uint32
	shader    Shader
	texture   Texture
	index     int
	transform Transform
}

type Renderer interface {
	BeginScene(Camera)
	PushItem(renderItem)
	render()
}

var rendererOnce sync.Once
var rendererSingleton Renderer
var defaultShader Shader

func Renderer2DInit(width, height float32) {
	rendererOnce.Do(func() {
		orthoProjection := mgl32.Ortho(0, width, height, 0, 0.1, 10)
		r := renderer{
			renderBuffer: []renderItem{},
			projection:   orthoProjection,
			activeCam:    camera2D{},
		}
		rendererSingleton = &r

		s, err := NewShader(vertexShaderSource, fragmentShaderSource)
		if err != nil {
			log.Println("error loading shader")
			panic(err)
		}
		defaultShader = s
	})
}

func Renderer2D() Renderer {
	return rendererSingleton
}

func DefaultShader() Shader {
	return defaultShader
}

func (r *renderer) BeginScene(c Camera) {
	r.renderBuffer = []renderItem{}
	r.activeCam = c
}
func (r *renderer) PushItem(ri renderItem) {
	r.renderBuffer = append(r.renderBuffer, ri)
}

// TODO (Ross): Filter by shaders, types etc
func (r *renderer) render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, i := range r.renderBuffer {
		i.shader.Use()
		i.shader.SetMatrix("u_projection", r.projection)
		i.shader.SetMatrix("u_view", r.activeCam.ViewMatrix())
		i.shader.SetMatrix("u_model", GetMatrix(i.transform))

		gl.ActiveTexture(gl.TEXTURE0)
		i.texture.use()
		gl.BindVertexArray(i.vao)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	}
}
