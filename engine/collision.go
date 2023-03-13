package engine

import ()

type Collider struct {
	width  int
	height int
	X      int
	Y      int
}

func NewCollider(width, height, x, y int) Collider {
	return Collider{
		width:  width,
		height: height,
		X:      x - width/2,
		Y:      y - height/2,
	}
}

func (c Collider) SetPos(x, y int) Collider {
	return Collider{
		c.width,
		c.height,
		x - c.width/2,
		y - c.height/2,
	}
}

// Returns true if a and b intersect
func Collides(a, b Collider) bool {
	return a.X < b.X+b.width && a.X+a.width > b.X && a.Y < b.Y+b.height && a.Y+a.height > b.Y
}

func getTileIndex(t Tilemap, x, y int) int {
	gx := x / t.tileSize
	gy := y / t.tileSize
	idx := gy*t.width + gx
	return idx
}

func CollidesMapPoint(t Tilemap, x, y int) bool {
	if x < 0 || x > t.width*t.tileSize || y < 0 || y > t.height*t.tileSize {
		return true
	}
	return t.collision[getTileIndex(t, x, y)] != -1
}

// True if the collider collides with tilemap collision layer
func CollidesMapCollider(t Tilemap, c Collider) bool {
	if c.X < 0 || c.X+c.width > t.width*t.tileSize || c.Y < 0 || c.Y+c.height > t.height*t.tileSize {
		return true
	}

	tlx, tly := c.X, c.Y
	trx, try := (c.X + c.width), c.Y
	blx, bly := c.X, (c.Y + c.height)
	brx, bry := (c.X + c.width), (c.Y + c.height)

	if t.collision[getTileIndex(t, tlx, tly)] != -1 || t.collision[getTileIndex(t, trx, try)] != -1 || t.collision[getTileIndex(t, brx, bry)] != -1 || t.collision[getTileIndex(t, blx, bly)] != -1 {
		return true
	}
	return false
}
