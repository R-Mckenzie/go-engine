package engine

type Atlas struct {
	texture    Image // The atlas itself
	tileSize   int   // Pixel size of each individual sprite
	atlasWidth int   // Number of tiles along X and Y axis
}

func NewAtlas(filepath string, tileSize, atlasWidth int) Atlas {
	tex, err := LoadImage(filepath)
	if err != nil {
		panic(err)
	}
	return Atlas{
		texture:    tex,
		tileSize:   tileSize,
		atlasWidth: atlasWidth,
	}
}

func (a Atlas) offset(index int) (float32, float32) {
	xOff := float32((index % a.atlasWidth)) / float32(a.atlasWidth)
	row := index / a.atlasWidth
	yOff := float32(row) / float32(a.atlasWidth)
	return float32(xOff), float32(yOff)
}
