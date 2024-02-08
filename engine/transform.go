package engine

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Pos   mgl32.Vec3
	Scale mgl32.Vec3
	Rot   mgl32.Vec3
}

func NewTransform(xP, yP, zP float32) Transform {
	return Transform{
		Pos:   mgl32.Vec3{xP, yP, zP},
		Rot:   mgl32.Vec3{0, 0, 0},
		Scale: mgl32.Vec3{1, 1, 1},
	}
}

func GetMatrix(t Transform) mgl32.Mat4 {
	scaleMat := mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2])

	rotationMat := mgl32.HomogRotate3DX(t.Rot[0]).Mul4(mgl32.HomogRotate3DY(t.Rot[1])).Mul4(mgl32.HomogRotate3DZ(t.Rot[2]))
	translationMat := mgl32.Translate3D(t.Pos[0], t.Pos[1], t.Pos[2])
	m := scaleMat.Mul4(translationMat).Mul4(rotationMat)
	return m
}
