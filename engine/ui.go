package engine

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

/*
INPUT

*/

type ui struct {
	input *input
	skin  Texture
	font  *Font

	hotItem    int // item hovered over by mouse
	activeItem int // the id of the currently selected item. 0 means nothing selected

	idCount int
}

func (ui *ui) Begin() {
	Renderer.beginUI()
	ui.hotItem = 0
	ui.idCount = 1
}

func (ui *ui) End() {
	if Input.KeyOnce(MouseLeft) {
		ui.activeItem = 0
	}
}

func (ui *ui) Button(x, y, w, h float32, label string, colour mgl32.Vec4) bool {
	printData := ui.font.renderItem(x, y, 32, label)
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
		if ui.activeItem != id && ui.input.KeyDown(MouseLeft) {
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

	if ui.input.KeyDown(MouseLeft) && ui.hotItem == id && ui.activeItem == id {
		return true
	}
	return false
}

func (ui *ui) Label(label string, x, y float32, fontSize int, colour mgl32.Vec4) mgl32.Vec2 {
	printData := ui.font.renderItem(x, y, fontSize, label)
	printData.ri.colour = colour
	Renderer.PushUI(printData.ri)
	return printData.size
}

var backspaceRepeat = time.Millisecond * 100
var backspaceTimer = time.Now()

func (ui *ui) TextInput(hint string, x, y float32, widthInChars, fontSize int, buf *string) {
	id := ui.idCount
	ui.idCount++

	var w, h float32
	var stringRender stringRenderItemSize
	var strlen int

	if len(*buf) > 0 {
		// Show buf
		stringRender = ui.font.renderItem(x, y, fontSize, *buf)
		stringRender.ri.colour = mgl32.Vec4{0, 0, 0, 1}
		strlen = len(*buf)
	} else {
		// Show hint
		stringRender = ui.font.renderItem(x, y, fontSize, hint)
		stringRender.ri.colour = mgl32.Vec4{0.5, 0.5, 0.5, 1}
		strlen = len(hint)
	}

	w = (stringRender.size[0] / float32(strlen)) * float32(widthInChars)
	h = stringRender.size[1]

	colour := mgl32.Vec4{1, 1, 1, 1}
	if ui.regionhit(x, y, w, h) {
		colour = mgl32.Vec4{0.9, 0.9, 0.9, 1}
		ui.hotItem = id
		if ui.activeItem != id && ui.input.KeyDown(MouseLeft) {
			ui.activeItem = id
		}
	}

	if ui.activeItem == id {
		colour = mgl32.Vec4{0.9, 0.9, 0.9, 1}
		for _, c := range Input.textInput {
			if ui.font.isPrintable(c) && len(*buf) <= widthInChars {
				*buf += string(c)
			}
		}
		if Input.KeyDown(KeyBackspace) && time.Now().Sub(backspaceTimer) > backspaceRepeat {
			backspaceTimer = time.Now()
			if len(*buf) > 0 {
				s := *buf
				s = s[:len(s)-1]
				*buf = s
			}
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
	Renderer.PushUI(stringRender.ri)
}

func (ui *ui) Checkbox(label string, val *bool) {

}

func (ui *ui) regionhit(x, y, w, h float32) bool {
	mouse := ui.input.MousePosition()
	if mouse.X() < x || mouse.Y() < y || mouse.X() >= x+w || mouse.Y() >= y+h {
		return false
	}
	return true
}
