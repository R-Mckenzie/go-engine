package engine

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/lafriks/go-tiled"
)

type Tilemap struct {
	width     int
	height    int
	tileSize  int
	tiles     [][]int // first index for layer, second for which tile
	collision []int   // Anything not -1 represents a collider
	textures  []Texture
	vao       uint32
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

	// Load the tilemap layers
	layers := [][]int{}
	var collisionLayer []int = nil
	for _, l := range m.Layers {
		layerTiles := make([]int, m.Width*m.Height)

		for j, t := range l.Tiles {
			if t.IsNil() {
				layerTiles[j] = -1
			} else {
				layerTiles[j] = int(t.ID)
			}
		}

		if l.Name == "Collision" {
			collisionLayer = layerTiles
		} else {
			layers = append(layers, layerTiles)
		}
	}

	tileMap := Tilemap{
		width:     m.Width,
		height:    m.Height,
		tileSize:  int(float32(m.TileHeight) * scale),
		tiles:     layers,
		textures:  textures,
		collision: collisionLayer,
	}
	tileMap.vao = tileMap.init()
	return tileMap
}

func atlasToTextures(filepath string, tileSize, atlasWidth, atlasHeight, tileCount int) []Texture {
	image, err := LoadImage(filepath)
	if err != nil {
		panic(err)
	}

	textures := make([]Texture, 0, tileCount)
	for i := 0; i < tileCount; i++ {
		col := float32((i % atlasWidth))
		row := float32(i / atlasWidth)

		texture := NewTextureFromAtlas(image, col*float32(tileSize), row*float32(tileSize), float32(tileSize), float32(tileSize))
		textures = append(textures, texture)
	}

	return textures
}

func (t Tilemap) RenderItem() renderItem {
	return renderItem{
		vao:       t.vao,
		indices:   int32(t.width) * int32(t.height) * 6 * int32(len(t.tiles)),
		image:     t.textures[0].image,
		shader:    DefaultShader(),
		transform: NewTransform(0, 0, 0),
	}
}

func (t *Tilemap) generateVertices() ([]float32, []uint32) {
	tiles := t.width * t.height

	vertices := make([]float32, 0, tiles*5*4*len(t.tiles)) // tiles * 4 vertices * 5 data points per vertex (x,y,z,u,v) * layers
	indices := make([]uint32, 0, tiles*6*len(t.tiles))     // For every tile we have 2 tris, 3 vertices each for total 6 * layers

	offset := uint32(0)
	for l := 0; l < len(t.tiles); l++ {
		for i := 0; i < (t.width * t.height); i++ {
			tile := t.tiles[l][i]
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

func (t *Tilemap) init() uint32 {
	v, i := t.generateVertices()
	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(t.tiles)*t.width*t.height*5*4*4, gl.Ptr(v), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(t.tiles)*t.width*t.height*6*4, gl.Ptr(i), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	return vao
}

func (t Tilemap) PixelSize() (int, int) {
	return t.width * t.tileSize, t.height * t.tileSize
}
