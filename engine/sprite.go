package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	Transform
	Width   int
	Height  int
	vao     uint32
	texture Texture
}

func NewSprite(width, height, x, y int, texture Texture) Sprite {
	return Sprite{
		Width:     width,
		Height:    height,
		Transform: NewTransform(x, y, 0),
		texture:   texture,
		vao:       makeVAO(quad(float32(width), float32(height), texture.texCoords)),
	}
}

func (s Sprite) RenderItem() renderItem {
	return renderItem{
		vao:       s.vao,
		shader:    DefaultShader(),
		transform: s.Transform,
		texture:   s.texture,
	}
}

func quad(width, height float32, uv mgl32.Vec4) ([]float32, []uint32) {
	w2, h2 := width/2, height/2
	return []float32{ // vertices
			-w2, -h2, 0.0, uv[0], uv[2],
			w2, -h2, 0.0, uv[1], uv[2],
			w2, h2, 0.0, uv[1], uv[3],
			-h2, h2, 0.0, uv[0], uv[3],
		}, []uint32{ // indices
			0, 1, 3,
			1, 2, 3,
		}
}

// w2, h2 := width/2, height/2
// return []float32{ // vertices
// 		-w2, -h2, 0.0, 0.0, 0.0,
// 		w2, -h2, 0.0, 1.0, 0.0,
// 		w2, h2, 0.0, 1.0, 1.0,
// 		-h2, h2, 0.0, 0.0, 1.0,
// 	}, []uint32{ // indices
// 		0, 1, 3,
// 		1, 2, 3,
// 	}

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
