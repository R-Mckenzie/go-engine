package engine

import (
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const MAX_LIGHTS = 15

type renderer struct {
	renderBuffer map[Image][]renderItem
	uiBuffer     map[Image][]renderItem
	ambientLight mgl32.Vec3
	activeCam    Camera
	projection   mgl32.Mat4
	postShader   Shader
	lights       []Light
	lightingFB   frameBuffer
	postFB       frameBuffer
	exposure     float32

	screenTransform Transform
}

type renderItem struct {
	vao        uint32
	indices    int32
	shader     Shader
	image      Image
	useNormals bool
	normals    Image
	transform  Transform
}

type Renderer interface {
	BeginScene(Camera, mgl32.Vec3)
	PushItem(renderItem)
	PushLight(Light)
	PushUI(renderItem)
	SetPostShader(string)
	SetExposure(float32)
	render()
}

var screenVAO uint32
var screenInd int32

var rendererOnce sync.Once
var rendererSingleton Renderer

var objectShader defaultShader
var postShader postprocessShader

type defaultShader struct {
	Shader
}

func (ds defaultShader) loadUniforms(model, view, projection mgl32.Mat4) {
	ds.SetMatrix("u_projection", projection)
	ds.SetMatrix("u_view", view)
	ds.SetMatrix("u_model", model)
}

type postprocessShader struct {
	Shader
}

func (ps postprocessShader) loadUniforms() {

}

// Initialises a 2D renderer. Takes in the width and height of the window render target
func Renderer2DInit(width, height float32) {
	rendererOnce.Do(func() {
		shaderMap = make(map[string]Shader)
		objectShader = defaultShader{loadShader(vertexShaderSource, fragmentShaderSource)}
		postShader = postprocessShader{loadShader("shaders/postprocessVertex.glsl", "shaders/postprocessFragment.glsl")}

		screenVAO, _, screenInd = screenQuadVAO()

		fb := newFrameBuffer(int32(width), int32(height))
		lb := newFrameBuffer(int32(width), int32(height))

		transform := NewTransform(0, 0, 0)
		transform.Scale = mgl32.Vec3{width, -height, 1}

		orthoProjection := mgl32.Ortho(0, width, height, 0, -0.1, 10.1)
		r := renderer{
			renderBuffer: make(map[Image][]renderItem),
			uiBuffer:     make(map[Image][]renderItem),
			projection:   orthoProjection,
			activeCam:    Camera2D{},
			postShader:   postShader.Shader,
			postFB:       fb,

			lightingFB:      lb,
			ambientLight:    mgl32.Vec3{1, 1, 1},
			screenTransform: transform,
		}

		rendererSingleton = &r
	})
}

func (r *renderer) SetExposure(e float32) {
	r.exposure = e

	r.postShader.Use()
	r.postShader.SetFloat("exposure", e)
}

func Renderer2D() Renderer {
	return rendererSingleton
}

func pushLightUniforms(lights []Light, view, projection mgl32.Mat4) {
	var positions []float32 // vec3
	var falloffs []float32  // vec3
	var colours []float32   // vec4

	for i := 0; i < MAX_LIGHTS; i++ {
		if i < len(lights) {
			light := lights[i]
			position := light.position(view, projection)

			positions = append(positions, position[0], position[1], position[2])
			falloffs = append(falloffs, light.Falloffs[0], light.Falloffs[1], light.Falloffs[2])
			colours = append(colours, light.Colour[0], light.Colour[1], light.Colour[2], light.Colour[3])
		} else {
			positions = append(positions, 0, 0, 0)
			falloffs = append(falloffs, 0, 0, 0)
			colours = append(colours, 0, 0, 0, 0)
		}
	}

	objectShader.SetVec3Array("LightPos", MAX_LIGHTS, positions)
	objectShader.SetVec4Array("LightColor", MAX_LIGHTS, colours)
	objectShader.SetVec3Array("Falloff", MAX_LIGHTS, falloffs)
}

func (r *renderer) BeginScene(c Camera, ambientLight mgl32.Vec3) {
	r.renderBuffer = make(map[Image][]renderItem)
	r.lights = []Light{}
	r.uiBuffer = make(map[Image][]renderItem)
	r.activeCam = c
	r.ambientLight = ambientLight

	objectShader.Use()
	objectShader.SetInt("u_texture", 0) //GL_TEXTURE0
	objectShader.SetInt("u_normals", 1) //GL_TEXTURE1
	objectShader.SetVec4("AmbientColor", r.ambientLight.Vec4(1))
	objectShader.SetVec2("Resolution", mgl32.Vec2{DispW, DispH})
}

func (r *renderer) PushItem(ri renderItem) {
	r.renderBuffer[ri.image] = append(r.renderBuffer[ri.image], ri)
}

func (r *renderer) PushLight(light Light) {
	r.lights = append(r.lights, light)
}

func (r *renderer) PushUI(ri renderItem) {
	r.uiBuffer[ri.image] = append(r.uiBuffer[ri.image], ri)
}

func (r *renderer) SetPostShader(name string) {
	shader, ok := shaderMap[name]
	if !ok {
		r.postShader = postShader.Shader
		return
	}
	r.postShader = shader
}

func (r *renderer) PushBatch() {}

// TODO (Ross): Filter by shaders, types etc
func (r *renderer) render() {
	// Bind scene framebuffer, render to texture
	r.postFB.use()
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	w, h := WindowSize()

	objectShader.Use()

	pushLightUniforms(r.lights, r.activeCam.ViewMatrix(), r.projection)

	for _, v := range r.renderBuffer {
		gl.ActiveTexture(gl.TEXTURE0)
		v[0].image.Use()
		for _, ri := range v {
			if ri.useNormals {
				// load normal uniforms
				objectShader.SetBool("UseNormals", true)
				gl.ActiveTexture(gl.TEXTURE1)
				ri.normals.Use()
			} else {
				objectShader.SetBool("UseNormals", false)
			}

			objectShader.loadUniforms(GetMatrix(ri.transform), r.activeCam.ViewMatrix(), r.projection)
			gl.BindVertexArray(ri.vao)
			gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
			gl.ActiveTexture(gl.TEXTURE0)
		}
	}

	// Render UI on top
	// objectShader.Use()
	// for _, v := range r.uiBuffer {
	// 	gl.ActiveTexture(gl.TEXTURE0)
	// 	v[0].image.Use()
	// 	for _, ri := range v {
	// 		objectShader.loadUniforms(mgl32.Ident4(), mgl32.Ident4(), r.projection)
	// 		gl.BindVertexArray(ri.vao)
	// 		gl.DrawElements(gl.TRIANGLES, ri.indices, gl.UNSIGNED_INT, nil)
	// 	}
	// }

	// now bind back to default framebuffer and draw a quad plane with the attached framebuffer color texture
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Disable(gl.DEPTH_TEST) // disable depth test so screen-space quad isn't discarded due to depth test.
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.Viewport(0, 0, int32(w*2), int32(h*2))

	r.postShader.Use()
	r.postShader.SetInt("u_texture", 0) //GL_TEXTURE0
	r.postFB.tex.image.Use()
	gl.BindVertexArray(screenVAO)
	gl.BindTexture(gl.TEXTURE_2D, r.postFB.tex.image.id)
	gl.DrawElements(gl.TRIANGLES, screenInd, gl.UNSIGNED_INT, nil)
}

type frameBuffer struct {
	id     uint32
	rbo    uint32
	quad   uint32
	tex    Texture
	width  int32
	height int32
}

func newFrameBuffer(w, h int32) frameBuffer {
	var fbo uint32
	gl.GenFramebuffers(1, &fbo)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

	tex := NewBlankTexture(float32(w), float32(h))

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tex.image.id, 0)

	var rbo uint32
	gl.GenRenderbuffers(1, &rbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, w, h)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER, rbo)

	// put the default buffer back
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return frameBuffer{
		id:     fbo,
		rbo:    rbo,
		tex:    tex,
		width:  w,
		height: h,
	}
}

func (f frameBuffer) use() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.id)
	gl.Viewport(0, 0, f.width, f.height)
}

func screenQuadVAO() (uint32, uint32, int32) {
	p, i := []float32{ // vertices
		-1, -1, 0.0, 0, 0,
		1, -1, 0.0, 1, 0,
		1, 1, 0.0, 1, 1,
		-1, 1, 0.0, 0, 1,
	}, []uint32{ // indices
		0, 1, 3,
		1, 2, 3,
	}
	return genVAO(p, i)
}

func quad(width, height float32, uv mgl32.Vec4) ([]float32, []uint32) {
	w2, h2 := width/2, height/2
	return []float32{ // vertices
			-w2, -h2, 0.0, uv[0], uv[2],
			w2, -h2, 0.0, uv[1], uv[2],
			w2, h2, 0.0, uv[1], uv[3],
			-w2, h2, 0.0, uv[0], uv[3],
		}, []uint32{ // indices
			0, 1, 3,
			1, 2, 3,
		}
}

func newQuadVAO(width, height float32, uv mgl32.Vec4) (uint32, uint32, int32) {
	p, i := quad(width, height, uv)
	return genVAO(p, i)
}

func genVAO(p []float32, i []uint32) (uint32, uint32, int32) {
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
	return vao, vbo, int32(len(i))
}
