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

	idCount int
}

func (ui *ui) Begin() {
	Renderer.beginUI()
	ui.hotItem = 0
	ui.idCount = 1
}

func (ui *ui) End() {
	if !ui.input.KeyDown(MouseLeft) {
		ui.activeItem = 0
	} else if ui.activeItem == 0 {
		ui.activeItem = -1
	}
}

func (ui *ui) Button(x, y, w, h float32, label string, colour mgl32.Vec4) bool {
	printData := ui.font.renderItem(x, y, 64, label)
	printData.ri.transform.Pos = printData.ri.transform.Pos.Add(mgl32.Vec3{(w / 2) - (printData.size[0] / 2), (h / 2) - (printData.size[1] / 2)})
	printData.ri.colour = mgl32.Vec4{0, 0, 0, 1}

	id := ui.idCount
	ui.idCount++

	if ui.regionhit(x, y, w, h) {
		if (colour[0]+colour[1]+colour[2])/3 < 0.5 {
			colour = colour.Add(mgl32.Vec4{0.3, 0.3, 0.3, 0})
			y -= 2
			printData.ri.transform.Pos[1] -= 2
		} else {
			colour = colour.Sub(mgl32.Vec4{0.3, 0.3, 0.3, 0})
			y -= 2
			printData.ri.transform.Pos[1] -= 2
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
	Renderer.PushUI(printData.ri)

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

func (ui *ui) Label(label string, x, y float32, size int, colour mgl32.Vec4) {
	printData := ui.font.renderItem(x, y, size, label)
	printData.ri.colour = colour
	Renderer.PushUI(printData.ri)
}

func (ui *ui) TextInput(label, hint string, buf *string) {

}

func (ui *ui) Checkbox(label string, val *bool) {

}
