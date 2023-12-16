package engine

import (
	"log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	Transform
	Width   int
	Height  int
	vao     uint32
	vbo     uint32
	texture Texture
}

func NewSprite(width, height, x, y, z int, texture Texture) Sprite {
	var p []float32
	var i []uint32

	if width == 1 && height == 1 {
		p, i = screenQuad(float32(width), float32(height), mgl32.Vec4{0, 1, 0, 1})
	} else {
		p, i = quad(float32(width), float32(height), texture.texCoords)
	}

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

	return Sprite{
		Width:     width,
		Height:    height,
		Transform: NewTransform(x, y, z),
		texture:   texture,
		vao:       vao,
		vbo:       vbo,
	}
}

func (s *Sprite) SetTexture(texture Texture) {
	s.texture = texture
	v, _ := quad(float32(s.Width), float32(s.Height), texture.texCoords)

	gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(v), gl.Ptr(v), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (s Sprite) RenderItem() renderItem {
	return renderItem{
		vao:       s.vao,
		shader:    shaderMap[DEFAULT_SHADER],
		indices:   6,
		transform: s.Transform,
		image:     s.texture.image,
	}
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

func screenQuad(w, h float32, uv mgl32.Vec4) ([]float32, []uint32) {
	return []float32{ // vertices
			-w, -h, 0.0, uv[0], uv[2],
			w, -h, 0.0, uv[1], uv[2],
			w, h, 0.0, uv[1], uv[3],
			-w, h, 0.0, uv[0], uv[3],
		}, []uint32{ // indices
			0, 1, 3,
			1, 2, 3,
		}
}

type Animator struct {
	Current    *Animation
	animations map[string]*Animation
}

func NewAnimator() Animator {
	animator := Animator{
		animations: make(map[string]*Animation),
	}
	return animator
}

func (a *Animator) Add(animation *Animation, name string) {
	if a.animations == nil {
		a.animations = make(map[string]*Animation)
	}
	a.animations[name] = animation
}

func (a *Animator) Trigger(name string) {
	anim, ok := a.animations[name]
	if !ok {
		log.Println("animation does not exits: ", name)
	}
	// If we're already playing this animation, return
	if a.Current == anim {
		return
	}
	a.Current = anim
	a.Current.Play()
}

func (a *Animator) Frame() Texture {
	return a.Current.Frame()
}

type Animation struct {
	frameTime time.Duration
	IsPlaying bool
	frameTick *time.Ticker
	// Spritesheet info
	textures  []Texture
	frames    int
	currIndex int
	Changed   bool
}

func NewAnimation(sheet Image, fps, sheetWidth, sheetHeight, spriteWidth, spriteHeight int, flipX bool) *Animation {
	frames := sheetWidth * sheetHeight
	textures := make([]Texture, frames)
	for i := range textures {
		col := float32((i % sheetWidth))
		row := float32(i / sheetWidth)
		textures[i] = NewTextureFromAtlas(sheet, col*float32(spriteWidth), row*float32(spriteHeight), float32(spriteWidth), float32(spriteHeight), flipX)
	}

	animation := &Animation{
		IsPlaying: false,
		frameTime: time.Second / time.Duration(fps),
		frameTick: time.NewTicker(time.Second),
		textures:  textures,
		frames:    sheetWidth * sheetHeight,
		currIndex: 0,
	}

	go func() {
		for range animation.frameTick.C {
			if !animation.IsPlaying {
				continue
			}

			animation.currIndex++
			if animation.currIndex > animation.frames-1 {
				animation.currIndex = 0
			}
			animation.Changed = true
			animation.frameTick.Reset(animation.frameTime)
		}
	}()

	return animation
}

func (a *Animation) Play() {
	a.currIndex = 0
	a.frameTick.Reset(a.frameTime)
	a.IsPlaying = true
	a.Changed = true
}

func (a *Animation) Stop() {
	a.IsPlaying = false
}

func (a *Animation) Frame() Texture {
	a.Changed = false
	return a.textures[a.currIndex]
}
