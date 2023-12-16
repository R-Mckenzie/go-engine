package engine

import (
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lafriks/go-tiled"
)

type Tilemap struct {
	width          int
	height         int
	tileSize       int
	staticLayers   [][]int       // first index for layer, second for which tile
	animatedLayers [][]int       // first index for layer, second for which tile
	animatedTiles  map[int][]int // Key is tile index, value is slice of frames
	collision      []int         // Anything not -1 represents a collider
	textures       []Texture
	staticVAO      uint32
	animatedVAO    uint32
	animatedVBO    uint32
	changed        *bool // if the animated tiles have changed
	animIndex      *int
}

var tilemapShader *Shader

// Assumes a square texture atlas
func LoadTilemap(tmxPath, atlasPath string, scale float32) Tilemap {
	m, err := tiled.LoadFile(tmxPath)
	if err != nil {
		panic(err)
	}

	// Load the textures from the given atlas
	textures := atlasToTextures(atlasPath, m.TileHeight, m.Tilesets[0].Columns, m.Tilesets[0].Columns, m.Tilesets[0].TileCount)

	// Build up a list of all animated tile IDs and their frames.
	animatedTiles := make(map[int][]int)
	for _, t := range m.Tilesets[0].Tiles {
		if t.Animation != nil && len(t.Animation) > 0 {
			frames := []int{}
			for _, f := range t.Animation {
				frames = append(frames, int(f.TileID))
			}
			animatedTiles[int(t.ID)] = frames
		}
	}

	// Load the tilemap layers
	static := [][]int{}
	animated := [][]int{}
	var collisionLayer []int = nil
	for _, l := range m.Layers {
		layerTiles := make([]int, m.Width*m.Height)
		animTiles := make([]int, m.Width*m.Height)

		for j, t := range l.Tiles {
			if t.IsNil() {
				layerTiles[j] = -1
				animTiles[j] = -1
			} else if anim, ok := animatedTiles[int(t.ID)]; ok {
				layerTiles[j] = -1
				animTiles[j] = anim[0]
			} else {
				layerTiles[j] = int(t.ID)
				animTiles[j] = -1
			}
		}

		if l.Name == "Collision" {
			collisionLayer = layerTiles
		} else {
			static = append(static, layerTiles)
			animated = append(animated, animTiles)
		}
	}

	changed := true
	i := 0

	tileMap := Tilemap{
		width:          m.Width,
		height:         m.Height,
		tileSize:       int(float32(m.TileHeight) * scale),
		animatedTiles:  animatedTiles,
		animatedLayers: animated,
		staticLayers:   static,
		textures:       textures,
		collision:      collisionLayer,
		changed:        &changed,
		animIndex:      &i,
	}

	if len(animatedTiles) > 0 {
		duration := 0
		duration = int(m.Tilesets[0].Tiles[0].Animation[0].Duration)

		tick := time.NewTicker(time.Duration(duration) * time.Millisecond)
		go func() {
			for range tick.C {
				*tileMap.animIndex++
				if *tileMap.animIndex > 50 {
					*tileMap.animIndex = 0
				}

				*tileMap.changed = true
				tick.Reset(time.Duration(duration) * time.Millisecond)
			}
		}()
	}

	tileMap.staticVAO, tileMap.animatedVAO, tileMap.animatedVBO = tileMap.init()
	return tileMap
}

func atlasToTextures(filepath string, tileSize, atlasWidth, atlasHeight, tileCount int) []Texture {
	image, err := NewImage(filepath)
	if err != nil {
		panic(err)
	}

	textures := make([]Texture, 0, tileCount)
	for i := 0; i < tileCount; i++ {
		col := float32((i % atlasWidth))
		row := float32(i / atlasWidth)

		texture := NewTextureFromAtlas(image, col*float32(tileSize), row*float32(tileSize), float32(tileSize), float32(tileSize), false)
		textures = append(textures, texture)
	}

	return textures
}

func (t *Tilemap) renderItem(vaoID uint32) renderItem {
	return renderItem{
		vao:       vaoID,
		indices:   int32(t.width) * int32(t.height) * 6 * int32(len(t.staticLayers)),
		image:     t.textures[0].image,
		shader:    shaderMap[DEFAULT_SHADER],
		transform: NewTransform(0, 0, 0),
	}
}

func (t *Tilemap) AnimatedRenderItem() renderItem {
	if *t.changed {
		v, _ := t.vertices(*t.animIndex, true)
		gl.BindBuffer(gl.ARRAY_BUFFER, t.animatedVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(t.staticLayers)*t.width*t.height*5*4*4, gl.Ptr(v), gl.STATIC_DRAW)
		*t.changed = false
	}
	return t.renderItem(t.animatedVAO)
}

func (t *Tilemap) StaticRenderItem() renderItem {
	return t.renderItem(t.staticVAO)
}

func (t *Tilemap) vertices(tick int, animated bool) ([]float32, []uint32) {
	layers := t.staticLayers
	if animated {
		layers = t.animatedLayers
	}

	tiles := t.width * t.height
	vertices := make([]float32, 0, tiles*5*4*len(layers)) // tiles * 4 vertices * 5 data points per vertex (x,y,z,u,v) * layers
	indices := make([]uint32, 0, tiles*6*len(layers))     // For every tile we have 2 tris, 3 vertices each for total 6 * layers

	offset := uint32(0)
	for l := 0; l < len(layers); l++ {
		for i := 0; i < (t.width * t.height); i++ {
			tile := layers[l][i]
			if tile == -1 {
				continue
			}

			// get tile position in pixels
			halfTile := (t.tileSize / 2)
			col := i % t.width
			row := i / t.width
			x := col*t.tileSize + halfTile
			y := row*t.tileSize + halfTile

			xf := float32(x)
			yf := float32(y)
			ht := float32(halfTile)

			tex := t.textures[tile]
			if animated {
				anim, ok := t.animatedTiles[tile]
				if !ok {
					continue
				}
				textureIdx := tick % len(anim)
				tex = t.textures[anim[textureIdx]]
			}

			vertices = append(vertices,
				xf-ht, yf-ht, float32(l), tex.texCoords[0], tex.texCoords[2],
				xf+ht, yf-ht, float32(l), tex.texCoords[1], tex.texCoords[2],
				xf+ht, yf+ht, float32(l), tex.texCoords[1], tex.texCoords[3],
				xf-ht, yf+ht, float32(l), tex.texCoords[0], tex.texCoords[3],
			)

			indices = append(indices,
				0+offset, 1+offset, 3+offset,
				1+offset, 2+offset, 3+offset,
			)
			offset += 4
		}
	}
	return vertices, indices
}

func (t *Tilemap) init() (uint32, uint32, uint32) {
	v, i := t.vertices(0, false)
	var staticVBO, staticVAO, staticEBO uint32
	gl.GenVertexArrays(1, &staticVAO)
	gl.GenBuffers(1, &staticVBO)
	gl.GenBuffers(1, &staticEBO)
	gl.BindVertexArray(staticVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, staticVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(t.staticLayers)*t.width*t.height*5*4*4, gl.Ptr(v), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, staticEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(t.staticLayers)*t.width*t.height*6*4, gl.Ptr(i), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	v2, i2 := t.vertices(0, true)
	var animVBO, animVAO, animEBO uint32
	gl.GenVertexArrays(1, &animVAO)
	gl.GenBuffers(1, &animVBO)
	gl.GenBuffers(1, &animEBO)
	gl.BindVertexArray(animVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, animVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(t.animatedLayers)*t.width*t.height*5*4*4, gl.Ptr(v2), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, animEBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(t.animatedLayers)*t.width*t.height*6*4, gl.Ptr(i2), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	return staticVAO, animVAO, animVBO
}

func (t Tilemap) PixelSize() (int, int) {
	return t.width * t.tileSize, t.height * t.tileSize
}
