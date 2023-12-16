package engine

import "github.com/go-gl/mathgl/mgl32"

type box struct {
	width  int
	height int
	color  mgl32.Vec4
}

type TextField struct {
	x, y float32
	font *Font
	Text string
}

func newTextField(fontFile, text string, fontSize int, x, y float32) *TextField {
	font, err := LoadFont(fontFile, fontSize)
	if err != nil {
		panic(err)
	}
	return &TextField{font: font, Text: text, x: x, y: y}
}

func (t *TextField) draw() {
	t.font.Print(t.x, t.y, t.Text)
}
