package engine

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Light struct {
	Transform Transform
	Colour    mgl32.Vec3
	Radius    int32
}

// func NewLight(x, y, r, g, b float32, radius int32) Light {
// make quad
// vao, vbo, i := newQuadVAO(float32(radius), float32(radius))

// return Light{
// 	Pos:     mgl32.Vec2{x, y},
// 	Colour:  mgl32.Vec3{r, g, b},
// 	Radius:  radius,
// 	vao:     vao,
// 	vbo:     vbo,
// 	indices: i,
// }
// }
