package engine

import "github.com/go-gl/mathgl/mgl32"

type Box struct {
	colour mgl32.Vec4
	sprite Sprite
}

// x and y point to the top left point of the box
func NewGUIBox(w, h, x, y float32, colour mgl32.Vec4) Box {
	// load textures
	tex := NewTexture("res/ui9slice.png")
	sprite := NewSprite(w, h, x+w/2, y+h/2, 9, tex, nil)

	return Box{colour: colour, sprite: sprite}
}

func (b Box) renderItem() []renderItem {
	ri := b.sprite.renderItem()
	ri[0].shader = uiShader.Shader
	ri[0].colour = b.colour
	return ri
}

type TextField struct {
	x, y   float32
	font   *Font
	Text   string
	Colour mgl32.Vec4
}

func NewTextField(text string, fontSize int, x, y float32, font *Font, colour mgl32.Vec4) *TextField {
	return &TextField{font: font, Text: text, x: x, y: y, Colour: colour}
}

func (t TextField) renderItem() []renderItem {
	ri := t.font.renderItem(t.x, t.y, t.Text)
	ri.colour = t.Colour
	return []renderItem{ri}
}
