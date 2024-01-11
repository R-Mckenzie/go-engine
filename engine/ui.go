package engine

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type ui struct {
	input *input
	skin  Texture

	currentGenID *int
	hotItem      int // item hovered over by mouse
	activeItem   int // the item currently selected by clicking mouse
}

func (ui ui) update() {

}

func (ui ui) Button(x, y, w, h float32, label string, colour mgl32.Vec4) bool {
	id := *ui.currentGenID
	*ui.currentGenID += 1

	x = ScreenW / x
	y = ScreenH / y

	if ui.regionhit(x, y, w, h) {
		fmt.Println("MOUSE IN REGION")
		ui.hotItem = id
		if ui.activeItem == 0 && ui.input.KeyDown(MouseLeft) {
			ui.activeItem = id
		}
	}

	vao, _, ind := newQuadVAO(w, h, mgl32.Vec4{0, 1, 0, 1})
	ri := renderItem{
		vao:        vao,
		indices:    ind,
		image:      ui.skin.image,
		shader:     uiShader.Shader,
		useNormals: false,
		transform:  NewTransform(x+(w/2), y+(h/2), 8),
		colour:     colour,
	}
	Renderer.PushUI(ri)

	return false
}

func (ui ui) regionhit(x, y, w, h float32) bool {
	mouse := ui.input.MousePosition()

	if mouse.X() < x || mouse.Y() < y || mouse.X() >= x+w || mouse.Y() >= y+h {
		return false
	}
	return true
}

func (ui ui) Label(text string) {

}

func (ui ui) TextInput(label, hint string, buf *string) {

}

func (ui ui) Checkbox(label string, val *bool) {

}

func (b uiBox) renderItem() []renderItem {
	ri := b.sprite.renderItem()
	ri[0].shader = uiShader.Shader
	ri[0].colour = b.colour
	return ri
}

type uiBox struct {
	colour mgl32.Vec4
	sprite Sprite
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
