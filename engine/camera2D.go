package engine

import "github.com/go-gl/mathgl/mgl32"

type Camera interface {
	ViewMatrix() mgl32.Mat4
}

type Camera2D struct {
	X   int
	Y   int
	dir mgl32.Vec4
}

func NewCamera2D(x, y int) Camera2D {
	return Camera2D{
		X:   x,
		Y:   y,
		dir: mgl32.Vec4{0, 0, 1, 0},
	}
}

func (c *Camera2D) SetPos(x, y int) {
	c.X = x
	c.Y = y
}

func (c Camera2D) ViewMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(-float32(c.X), -float32(c.Y), -10)
}
