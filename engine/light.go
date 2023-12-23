package engine

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Light struct {
	Colour    mgl32.Vec4
	Falloffs  mgl32.Vec3
	transform Transform
}

func NewLight(x, y, z, r, g, b float32, intensity float32) Light {
	transform := NewTransform(x, y, z)
	colour := mgl32.Vec4{r, g, b, intensity}

	return Light{
		transform: transform,
		Colour:    colour,
		Falloffs:  mgl32.Vec3{0.3, 4, 20},
	}
}

func (l Light) position(view, projection mgl32.Mat4) mgl32.Vec3 {
	pos := projection.Mul4x1(l.transform.Pos.Vec4(1))
	pos = view.Mul4x1(pos)
	pos = GetMatrix(l.transform).Mul4x1(pos)
	pos[1] = DispH - pos[1]
	return pos.Vec3()
}
