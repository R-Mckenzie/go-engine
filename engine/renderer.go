package engine

import (
	"log"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Light struct {
	X      int
	Y      int
	Radius int
}

type renderer struct {
	renderBuffer map[Image][]renderItem
	uiBuffer     map[Image][]renderItem
	lightBuffer  []Light
	activeCam    Camera
	projection   mgl32.Mat4
	postShader   PostShader
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
	PushLight(Light)
	PushUI(renderItem)
	SetPostShader(PostShader)
	render()
}

var rendererOnce sync.Once
var rendererSingleton Renderer
var defaultShader Shader
var fb frameBuffer

func Renderer2DInit(width, height float32) {
	rendererOnce.Do(func() {
		// Create and set default shader
		ds, err := NewShader(vertexShaderSource, fragmentShaderSource)
		if err != nil {
			log.Println("error loading shader")
			panic(err)
		}
		defaultShader = ds

		orthoProjection := mgl32.Ortho(0, width, height, 0, -0.1, 10.1)
		r := renderer{
			renderBuffer: make(map[Image][]renderItem),
			uiBuffer:     make(map[Image][]renderItem),
			projection:   orthoProjection,
			activeCam:    Camera2D{},
			postShader:   newDefaultPostShader(),
		}

		rendererSingleton = &r

		fb = newFrameBuffer(int32(width), int32(height))
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
	r.lightBuffer = []Light{}
	r.uiBuffer = make(map[Image][]renderItem)
	r.activeCam = c
}

func (r *renderer) PushItem(ri renderItem) {
	r.renderBuffer[ri.image] = append(r.renderBuffer[ri.image], ri)
}

func (r *renderer) PushLight(light Light) {
	r.lightBuffer = append(r.lightBuffer, light)
}

func (r *renderer) PushUI(ri renderItem) {
	r.uiBuffer[ri.image] = append(r.uiBuffer[ri.image], ri)
}

func (r *renderer) SetPostShader(ps PostShader) {
	r.postShader = ps
}

func (r *renderer) PushBatch() {}

// TODO (Ross): Filter by shaders, types etc
func (r *renderer) render() {
	drawCalls := 0
	textureSwaps := 0
	shaderSwaps := 0
	uniformSets := 0

	gl.BindFramebuffer(gl.FRAMEBUFFER, fb.id)
	gl.Enable(gl.DEPTH_TEST)
	// gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	w, h := WindowSize()
	x := -(w - w) / 2
	y := -(h - h) / 2
	gl.Viewport(int32(x), int32(y), int32(w), int32(h))

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
			uniformSets += 3
			gl.BindVertexArray(ri.vao)
			gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
			drawCalls++
		}
	}

	// now bind back to default framebuffer and draw a quad plane with the attached framebuffer color texture
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Disable(gl.DEPTH_TEST) // disable depth test so screen-space quad isn't discarded due to depth test.
	// clear all relevant buffers
	gl.ClearColor(1.0, 1.0, 1.0, 1.0) // set clear color to white (not really necessary actually, since we won't be able to see behind the quad anyways)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Viewport(0, 0, int32(w*2), int32(h*2))

	r.postShader.use()
	gl.BindVertexArray(fb.sprite.vao)
	gl.BindTexture(gl.TEXTURE_2D, fb.sprite.texture.image.id) // use the color attachment texture as the texture of the quad plane
	gl.DrawElements(gl.TRIANGLES, fb.sprite.RenderItem().indices, gl.UNSIGNED_INT, nil)

	// Lighting pass
	if len(r.lightBuffer) > 0 {
		// add lighting

	}

	// Render UI on top
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
			uniformSets += 3
			gl.BindVertexArray(ri.vao)
			gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
			drawCalls++
		}
	}

	AddDebugInfo("Draw Calls", drawCalls)
	AddDebugInfo("Texture Swaps", textureSwaps)
	AddDebugInfo("Shader Swaps", shaderSwaps)
	AddDebugInfo("Uniform Sets", uniformSets)
}

type frameBuffer struct {
	id     uint32
	rbo    uint32
	sprite Sprite
}

func newFrameBuffer(texWidth, texHeight int32) frameBuffer {
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

	screenSprite := NewSprite(1, 1, 0, 0, 0, NewBlankTexture(float32(texWidth), float32(texHeight)))

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, screenSprite.texture.image.id, 0)

	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, texWidth, texHeight)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	return frameBuffer{
		id:     fbo,
		rbo:    rbo,
		sprite: screenSprite,
	}
}
