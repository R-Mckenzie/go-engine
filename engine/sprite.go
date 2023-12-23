package engine

import (
	"log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sprite struct {
	Transform
	Width   float32
	Height  float32
	vao     uint32
	vbo     uint32
	texture Texture
	normal  *Texture
}

// normal can be nil
func NewSprite(width, height, x, y, z float32, texture Texture, normal *Texture) Sprite {
	vao, vbo, _ := newQuadVAO(width, height, texture.texCoords)

	return Sprite{
		Width:     width,
		Height:    height,
		Transform: NewTransform(x, y, z),
		texture:   texture,
		normal:    normal,
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

func (s *Sprite) SetNormal(texture *Texture) {
	if texture != nil {
		s.texture = *texture
		v, _ := quad(float32(s.Width), float32(s.Height), texture.texCoords)

		gl.BindBuffer(gl.ARRAY_BUFFER, s.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(v), gl.Ptr(v), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
}

func (s Sprite) RenderItem() renderItem {
	if s.normal != nil {
		return renderItem{
			vao:       s.vao,
			indices:   6,
			transform: s.Transform,
			image:     s.texture.image,
			normals:   s.normal.image,
		}
	} else {
		return renderItem{
			vao:        s.vao,
			indices:    6,
			transform:  s.Transform,
			image:      s.texture.image,
			useNormals: false,
		}
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
