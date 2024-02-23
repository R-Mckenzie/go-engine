package engine

import (
	"bufio"
	"image"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/gltext"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	ttf           truetype.Font
	glyphs        map[int][]glyph
	atlas         map[int]Image
	renderDatas   map[int]map[string]stringRenderItemSize
	renderedSizes []int
}

type stringRenderItemSize struct {
	ri   renderItem
	size mgl32.Vec2
}

func LoadFont(path string) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Println("Error loading font: ", err)
		return nil, err
	}

	ttf, err := freetype.ParseFont(data)
	if err != nil {
		log.Println("Error loading font: ", err)
		return nil, err
	}

	return &Font{
		ttf:           *ttf,
		glyphs:        make(map[int][]glyph),
		atlas:         make(map[int]Image),
		renderDatas:   make(map[int]map[string]stringRenderItemSize),
		renderedSizes: []int{},
	}, nil
}

func (f *Font) genNewFontSize(size int) {
	low, high := rune(32), rune(127)
	glyphs := make([]glyph, high-low+1)

	gc := int32(len(glyphs))
	glyphsPerRow := int32(16)
	glyphsPerCol := (gc / glyphsPerRow) + 1

	gb := f.ttf.Bounds(fixed.Int26_6(size))
	gw := int32(gb.Max.X - gb.Min.X)
	gh := int32((gb.Max.Y - gb.Min.Y) + 5)
	iw := gltext.Pow2(uint32(gw * glyphsPerRow))
	ih := gltext.Pow2(uint32(gh * glyphsPerCol))

	rect := image.Rect(0, 0, int(iw), int(ih))
	img := image.NewRGBA(rect)

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(&f.ttf)
	c.SetFontSize(float64(size))
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)

	var gi int
	var gx, gy int32

	for ch := low; ch <= high; ch++ {
		index := f.ttf.Index(ch)
		metric := f.ttf.HMetric(fixed.Int26_6(size), index)

		glyphs[gi].advance = int(metric.AdvanceWidth)
		glyphs[gi].x = int(gx)
		glyphs[gi].y = int(gy) - int(gh)/2
		glyphs[gi].width = int(gw)
		glyphs[gi].height = int(gh)

		pt := freetype.Pt(int(gx), int(gy)+int(c.PointToFixed(float64(size))>>8))
		c.DrawString(string(ch), pt)

		if gi%16 == 0 {
			gx = 0
			gy += gh
		} else {
			gx += gw
		}
		gi++
	}

	//== Build up OpenGL texture data
	img = gltext.Pow2Image(img).(*image.RGBA)
	ib := img.Bounds()
	texWidth, texHeight := float32(ib.Dx()), float32(ib.Dy())

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(texWidth), int32(texHeight), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	atlas := Image{
		id:     tex,
		width:  texWidth,
		height: texHeight,
	}

	for i, glyph := range glyphs {
		glyphs[i].texture = NewTextureFromAtlas(atlas, float32(glyph.x), float32(glyph.y), float32(glyph.width), float32(glyph.height), false)
	}

	f.glyphs[size] = glyphs
	f.atlas[size] = atlas
	f.renderDatas[size] = make(map[string]stringRenderItemSize)
	f.renderedSizes = append(f.renderedSizes, size)
}

func toPNG(img image.Image) {
	// Save that RGBA image to disk.
	outFile, err := os.Create("out2.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, img)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func (f *Font) isPrintable(char rune) bool {
	if len(f.renderedSizes) == 0 {
		return false
	}
	if int(char)-32 < len(f.glyphs[f.renderedSizes[0]]) {
		return true
	}
	return false
}

// Returns the renderItem to be drawn, and the width in pixels from left edge to right edge
func (f *Font) renderItem(x, y float32, size int, str string) stringRenderItemSize {
	// Check if we have already rendered font atlas for the desired size
	_, ok := f.atlas[size]
	if !ok {
		f.genNewFontSize(size)
		println("Gen size: ", size)
	}

	// Use existing renderItem
	rd, ok := f.renderDatas[size][str]
	if ok {
		rd.ri.transform = NewTransform(x, y, 9)
		return rd
	}

	// Create new renderItem
	vertices := make([]float32, 0, 5*4*len(str))
	indices := make([]uint32, 0, 6*len(str))
	offset := uint32(0)
	currentX := float32(0)
	currentY := float32(0)
	maxHeight := 0.0

	for _, v := range str {
		g := f.glyphs[size][v-32]
		maxHeight = math.Max(maxHeight, float64(g.height))
		vertices = append(vertices,
			currentX, currentY, 0, g.texture.texCoords[0], g.texture.texCoords[2],
			currentX+float32(g.width-1), currentY, 0, g.texture.texCoords[1], g.texture.texCoords[2],
			currentX+float32(g.width-1), currentY+float32(g.height-1), 0, g.texture.texCoords[1], g.texture.texCoords[3],
			currentX, currentY+float32(g.height-1), 0, g.texture.texCoords[0], g.texture.texCoords[3],
		)
		currentX += float32(g.advance)

		indices = append(indices,
			0+offset, 1+offset, 3+offset,
			1+offset, 2+offset, 3+offset,
		)
		offset += 4
	}

	var vbo, vao, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, 3*4)
	gl.EnableVertexAttribArray(1)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	renderItem := renderItem{
		vao:       vao,
		indices:   int32(len(indices)),
		image:     f.atlas[size],
		transform: NewTransform(x, y, 9),
	}

	stringPrintData := stringRenderItemSize{ri: renderItem, size: mgl32.Vec2{currentX, float32(maxHeight)}}
	f.renderDatas[size][str] = stringPrintData
	return stringPrintData
}

type glyph struct {
	x       int
	y       int
	width   int
	height  int
	advance int
	texture Texture
}
