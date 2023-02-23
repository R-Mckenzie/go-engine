package graphics

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	X   float32
	Y   float32
	dir mgl32.Vec4
}

func NewCamera2D(x, y int) Camera {
	return Camera{
		X:   float32(x),
		Y:   float32(y),
		dir: mgl32.Vec4{0, 0, 1, 0},
	}
}
