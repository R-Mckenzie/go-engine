package engine

import (
	"log"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type renderer struct {
	renderBuffer map[Image][]renderItem
	uiBuffer     map[Image][]renderItem
	activeCam    Camera
	projection   mgl32.Mat4
}

type renderItem struct {
	vao       uint32
	indices   int32
	shader    Shader
	image     Image
	transform Transform
}

type Renderer interface {
	BeginScene(Camera)
	PushItem(renderItem)
	PushUI(renderItem)
	render()
}

var rendererOnce sync.Once
var rendererSingleton Renderer
var defaultShader Shader

func Renderer2DInit(width, height float32) {
	rendererOnce.Do(func() {
		orthoProjection := mgl32.Ortho(0, width, height, 0, -0.1, 10.1)
		r := renderer{
			renderBuffer: make(map[Image][]renderItem),
			uiBuffer:     make(map[Image][]renderItem),
			projection:   orthoProjection,
			activeCam:    Camera2D{},
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
	r.renderBuffer = make(map[Image][]renderItem)
	r.uiBuffer = make(map[Image][]renderItem)
	r.activeCam = c
}
func (r *renderer) PushItem(ri renderItem) {
	r.renderBuffer[ri.image] = append(r.renderBuffer[ri.image], ri)
}

func (r *renderer) PushUI(ri renderItem) {
	r.uiBuffer[ri.image] = append(r.uiBuffer[ri.image], ri)

}

func (r *renderer) PushBatch() {}

// TODO (Ross): Filter by shaders, types etc
func (r *renderer) render() {
	drawCalls := 0
	textureSwaps := 0
	shaderSwaps := 0

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, v := range r.renderBuffer {
		textureSwaps++
		gl.ActiveTexture(gl.TEXTURE0)
		v[0].image.Use()
		for _, ri := range v {
			shaderSwaps++
			ri.shader.Use()
			ri.shader.SetMatrix("u_projection", r.projection)
			ri.shader.SetMatrix("u_view", r.activeCam.ViewMatrix())
			ri.shader.SetMatrix("u_model", GetMatrix(ri.transform))
			gl.BindVertexArray(ri.vao)
			gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
			drawCalls++
		}
	}

	for _, v := range r.uiBuffer {
		textureSwaps++
		gl.ActiveTexture(gl.TEXTURE0)
		v[0].image.Use()
		for _, ri := range v {
			shaderSwaps++
			ri.shader.Use()
			ri.shader.SetMatrix("u_projection", r.projection)
			ri.shader.SetMatrix("u_view", mgl32.Ident4())
			ri.shader.SetMatrix("u_model", mgl32.Ident4())
			gl.BindVertexArray(ri.vao)
			gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
			drawCalls++
		}
	}

	AddDebugInfo("Draw Calls", drawCalls)
	AddDebugInfo("Texture Swaps", textureSwaps)
	AddDebugInfo("Shader Swaps", shaderSwaps)
}
