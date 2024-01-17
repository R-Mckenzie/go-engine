package engine

import (
	"github.com/go-gl/mathgl/mgl32"
)

type ui struct {
	input *input
	skin  Texture
	font  *Font

	hotItem    int // item hovered over by mouse
	activeItem int // the item currently selected by clicking mouse
}

func (ui *ui) Begin() {
	ui.hotItem = 0
}

func (ui *ui) End() {
	if !ui.input.KeyDown(MouseLeft) {
		ui.activeItem = 0
	} else if ui.activeItem == 0 {
		ui.activeItem = -1
	}
}

func (ui *ui) Button(x, y, w, h float32, id int, label string, colour mgl32.Vec4) bool {
	x = ScreenW / x
	y = ScreenH / y

	text := ui.font.renderItem(x, y, 64, label)
	text.colour = mgl32.Vec4{0, 0, 0, 1}

	if ui.regionhit(x, y, w, h) {
		if (colour[0]+colour[1]+colour[2])/3 < 0.5 {
			colour = colour.Add(mgl32.Vec4{0.3, 0.3, 0.3, 0})
			y -= 2
			text.transform.Pos[1] -= 2
		} else {
			colour = colour.Sub(mgl32.Vec4{0.3, 0.3, 0.3, 0})
			y -= 2
			text.transform.Pos[1] -= 2
		}

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

	Renderer.PushUI(text)

	if !ui.input.KeyDown(MouseLeft) && ui.hotItem == id && ui.activeItem == id {
		return true
	}
	return false
}

func (ui *ui) regionhit(x, y, w, h float32) bool {
	mouse := ui.input.MousePosition()

	if mouse.X() < x || mouse.Y() < y || mouse.X() >= x+w || mouse.Y() >= y+h {
		return false
	}
	return true
}

func (ui *ui) Label(text string) {

}

func (ui *ui) TextInput(label, hint string, buf *string) {

}

func (ui *ui) Checkbox(label string, val *bool) {

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
