package engine

import "github.com/go-gl/mathgl/mgl32"

type Camera interface {
	ViewMatrix() mgl32.Mat4
}

type camera2D struct {
	X   float32
	Y   float32
	dir mgl32.Vec4
}

func NewCamera2D(x, y int) Camera {
	return &camera2D{
		X:   float32(x),
		Y:   float32(y),
		dir: mgl32.Vec4{0, 0, 1, 0},
	}
}

func (c camera2D) ViewMatrix() mgl32.Mat4 {
	return mgl32.Translate3D(-c.X, -c.Y, -3)
}
