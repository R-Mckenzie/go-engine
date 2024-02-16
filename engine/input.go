package engine

import (
	"runtime"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type input struct {
	window      *window
	textInput   []rune
	keysDown    [KeyLast]bool
	currentKeys [KeyLast]bool
	keysOnce    [KeyLast]bool
	keysUp      [KeyLast]bool
	mouseDown   [3]bool

	screenScale float32 // On mac, we need to divide the framebuffer size by 2 to get the accurate screen position for the mouse input
}

func initInput() *input {
	keysDown := [KeyLast]bool{}
	currentKeys := [KeyLast]bool{}
	keysOnce := [KeyLast]bool{}
	keysUp := [KeyLast]bool{}
	mouseDown := [3]bool{}

	var screenScale float32
	if runtime.GOOS == "darwin" {
		screenScale = 2
	} else {
		screenScale = 1
	}

	return &input{
		keysDown:    keysDown,
		currentKeys: currentKeys,
		keysOnce:    keysOnce,
		keysUp:      keysUp,
		mouseDown:   mouseDown,
		screenScale: screenScale,
	}
}

func (i *input) setCallBacks(w *window) {
	w.win.SetMouseButtonCallback(func(win *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if int(button) < len(i.mouseDown) {
			switch action {
			case glfw.Press:
				i.keysDown[button] = true
				i.keysOnce[button] = true
			case glfw.Release:
				i.keysDown[button] = false
				i.keysOnce[button] = false
				i.keysUp[button] = true
			}
		}
	})

	w.win.SetKeyCallback(func(win *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		switch action {
		case glfw.Press:
			i.keysDown[key] = true
		case glfw.Release:
			i.keysDown[key] = false
			i.keysUp[key] = true
		}
	})

	w.win.SetScrollCallback(func(window *glfw.Window, x, y float64) {
	})

	w.win.SetCharCallback(func(w *glfw.Window, char rune) {
		i.textInput = append(i.textInput, char)
	})
}

func (i *input) setWindow(w *window) {
	i.setCallBacks(w)
	i.window = w
}

func (i *input) KeyDown(key int) bool {
	return i.keysDown[key]
}

func (i *input) KeyOnce(key int) bool {
	return i.keysOnce[key]
}

func (i *input) KeyUp(key int) bool {
	return i.keysUp[key]
}

func (i *input) MousePosition() mgl32.Vec2 {
	x, y := i.window.win.GetCursorPos()

	// TODO: scaling input based on resolution
	xRatio := ScreenW / (dispW / i.screenScale)
	yRatio := ScreenH / (dispH / i.screenScale)
	return mgl32.Vec2{float32(x) * xRatio, float32(y) * yRatio}
}

func (i *input) update() {
	for x := 0; x < KeyLast; x++ {
		i.keysUp[x] = false
		i.keysOnce[x] = false
		if i.keysDown[x] && !i.currentKeys[x] {
			i.keysOnce[x] = true
		}
		i.currentKeys[x] = i.keysDown[x]
	}
	i.textInput = []rune{}
}

// List of all keyboard buttons.
const (
	KeyUnknown      = int(glfw.KeyUnknown)
	KeySpace        = int(glfw.KeySpace)
	KeyApostrophe   = int(glfw.KeyApostrophe)
	KeyComma        = int(glfw.KeyComma)
	KeyMinus        = int(glfw.KeyMinus)
	KeyPeriod       = int(glfw.KeyPeriod)
	KeySlash        = int(glfw.KeySlash)
	Key0            = int(glfw.Key0)
	Key1            = int(glfw.Key1)
	Key2            = int(glfw.Key2)
	Key3            = int(glfw.Key3)
	Key4            = int(glfw.Key4)
	Key5            = int(glfw.Key5)
	Key6            = int(glfw.Key6)
	Key7            = int(glfw.Key7)
	Key8            = int(glfw.Key8)
	Key9            = int(glfw.Key9)
	KeySemicolon    = int(glfw.KeySemicolon)
	KeyEqual        = int(glfw.KeyEqual)
	KeyA            = int(glfw.KeyA)
	KeyB            = int(glfw.KeyB)
	KeyC            = int(glfw.KeyC)
	KeyD            = int(glfw.KeyD)
	KeyE            = int(glfw.KeyE)
	KeyF            = int(glfw.KeyF)
	KeyG            = int(glfw.KeyG)
	KeyH            = int(glfw.KeyH)
	KeyI            = int(glfw.KeyI)
	KeyJ            = int(glfw.KeyJ)
	KeyK            = int(glfw.KeyK)
	KeyL            = int(glfw.KeyL)
	KeyM            = int(glfw.KeyM)
	KeyN            = int(glfw.KeyN)
	KeyO            = int(glfw.KeyO)
	KeyP            = int(glfw.KeyP)
	KeyQ            = int(glfw.KeyQ)
	KeyR            = int(glfw.KeyR)
	KeyS            = int(glfw.KeyS)
	KeyT            = int(glfw.KeyT)
	KeyU            = int(glfw.KeyU)
	KeyV            = int(glfw.KeyV)
	KeyW            = int(glfw.KeyW)
	KeyX            = int(glfw.KeyX)
	KeyY            = int(glfw.KeyY)
	KeyZ            = int(glfw.KeyZ)
	KeyLeftBracket  = int(glfw.KeyLeftBracket)
	KeyBackslash    = int(glfw.KeyBackslash)
	KeyRightBracket = int(glfw.KeyRightBracket)
	KeyGraveAccent  = int(glfw.KeyGraveAccent)
	KeyWorld1       = int(glfw.KeyWorld1)
	KeyWorld2       = int(glfw.KeyWorld2)
	KeyEscape       = int(glfw.KeyEscape)
	KeyEnter        = int(glfw.KeyEnter)
	KeyTab          = int(glfw.KeyTab)
	KeyBackspace    = int(glfw.KeyBackspace)
	KeyInsert       = int(glfw.KeyInsert)
	KeyDelete       = int(glfw.KeyDelete)
	KeyRight        = int(glfw.KeyRight)
	KeyLeft         = int(glfw.KeyLeft)
	KeyDown         = int(glfw.KeyDown)
	KeyUp           = int(glfw.KeyUp)
	KeyPageUp       = int(glfw.KeyPageUp)
	KeyPageDown     = int(glfw.KeyPageDown)
	KeyHome         = int(glfw.KeyHome)
	KeyEnd          = int(glfw.KeyEnd)
	KeyCapsLock     = int(glfw.KeyCapsLock)
	KeyScrollLock   = int(glfw.KeyScrollLock)
	KeyNumLock      = int(glfw.KeyNumLock)
	KeyPrintScreen  = int(glfw.KeyPrintScreen)
	KeyPause        = int(glfw.KeyPause)
	KeyF1           = int(glfw.KeyF1)
	KeyF2           = int(glfw.KeyF2)
	KeyF3           = int(glfw.KeyF3)
	KeyF4           = int(glfw.KeyF4)
	KeyF5           = int(glfw.KeyF5)
	KeyF6           = int(glfw.KeyF6)
	KeyF7           = int(glfw.KeyF7)
	KeyF8           = int(glfw.KeyF8)
	KeyF9           = int(glfw.KeyF9)
	KeyF10          = int(glfw.KeyF10)
	KeyF11          = int(glfw.KeyF11)
	KeyF12          = int(glfw.KeyF12)
	KeyF13          = int(glfw.KeyF13)
	KeyF14          = int(glfw.KeyF14)
	KeyF15          = int(glfw.KeyF15)
	KeyF16          = int(glfw.KeyF16)
	KeyF17          = int(glfw.KeyF17)
	KeyF18          = int(glfw.KeyF18)
	KeyF19          = int(glfw.KeyF19)
	KeyF20          = int(glfw.KeyF20)
	KeyF21          = int(glfw.KeyF21)
	KeyF22          = int(glfw.KeyF22)
	KeyF23          = int(glfw.KeyF23)
	KeyF24          = int(glfw.KeyF24)
	KeyF25          = int(glfw.KeyF25)
	KeyKP0          = int(glfw.KeyKP0)
	KeyKP1          = int(glfw.KeyKP1)
	KeyKP2          = int(glfw.KeyKP2)
	KeyKP3          = int(glfw.KeyKP3)
	KeyKP4          = int(glfw.KeyKP4)
	KeyKP5          = int(glfw.KeyKP5)
	KeyKP6          = int(glfw.KeyKP6)
	KeyKP7          = int(glfw.KeyKP7)
	KeyKP8          = int(glfw.KeyKP8)
	KeyKP9          = int(glfw.KeyKP9)
	KeyKPDecimal    = int(glfw.KeyKPDecimal)
	KeyKPDivide     = int(glfw.KeyKPDivide)
	KeyKPMultiply   = int(glfw.KeyKPMultiply)
	KeyKPSubtract   = int(glfw.KeyKPSubtract)
	KeyKPAdd        = int(glfw.KeyKPAdd)
	KeyKPEnter      = int(glfw.KeyKPEnter)
	KeyKPEqual      = int(glfw.KeyKPEqual)
	KeyLeftShift    = int(glfw.KeyLeftShift)
	KeyLeftControl  = int(glfw.KeyLeftControl)
	KeyLeftAlt      = int(glfw.KeyLeftAlt)
	KeyLeftSuper    = int(glfw.KeyLeftSuper)
	KeyRightShift   = int(glfw.KeyRightShift)
	KeyRightControl = int(glfw.KeyRightControl)
	KeyRightAlt     = int(glfw.KeyRightAlt)
	KeyRightSuper   = int(glfw.KeyRightSuper)
	KeyMenu         = int(glfw.KeyMenu)
	KeyLast         = int(glfw.KeyLast)

	MouseLeft   = int(glfw.MouseButtonLeft)
	MouseRight  = int(glfw.MouseButtonRight)
	MouseMiddle = int(glfw.MouseButtonMiddle)
)

var mouseNames = map[int]string{
	MouseLeft:   "LMB",
	MouseRight:  "RMB",
	MouseMiddle: "MMB",
}

var buttonNames = map[int]string{
	KeyUnknown:      "Unknown",
	KeySpace:        "Space",
	KeyApostrophe:   "Apostrophe",
	KeyComma:        "Comma",
	KeyMinus:        "Minus",
	KeyPeriod:       "Period",
	KeySlash:        "Slash",
	Key0:            "0",
	Key1:            "1",
	Key2:            "2",
	Key3:            "3",
	Key4:            "4",
	Key5:            "5",
	Key6:            "6",
	Key7:            "7",
	Key8:            "8",
	Key9:            "9",
	KeySemicolon:    "Semicolon",
	KeyEqual:        "Equal",
	KeyA:            "A",
	KeyB:            "B",
	KeyC:            "C",
	KeyD:            "D",
	KeyE:            "E",
	KeyF:            "F",
	KeyG:            "G",
	KeyH:            "H",
	KeyI:            "I",
	KeyJ:            "J",
	KeyK:            "K",
	KeyL:            "L",
	KeyM:            "M",
	KeyN:            "N",
	KeyO:            "O",
	KeyP:            "P",
	KeyQ:            "Q",
	KeyR:            "R",
	KeyS:            "S",
	KeyT:            "T",
	KeyU:            "U",
	KeyV:            "V",
	KeyW:            "W",
	KeyX:            "X",
	KeyY:            "Y",
	KeyZ:            "Z",
	KeyLeftBracket:  "LeftBracket",
	KeyBackslash:    "Backslash",
	KeyRightBracket: "RightBracket",
	KeyGraveAccent:  "GraveAccent",
	KeyWorld1:       "World1",
	KeyWorld2:       "World2",
	KeyEscape:       "Escape",
	KeyEnter:        "Enter",
	KeyTab:          "Tab",
	KeyBackspace:    "Backspace",
	KeyInsert:       "Insert",
	KeyDelete:       "Delete",
	KeyRight:        "Right",
	KeyLeft:         "Left",
	KeyDown:         "Down",
	KeyUp:           "Up",
	KeyPageUp:       "PageUp",
	KeyPageDown:     "PageDown",
	KeyHome:         "Home",
	KeyEnd:          "End",
	KeyCapsLock:     "CapsLock",
	KeyScrollLock:   "ScrollLock",
	KeyNumLock:      "NumLock",
	KeyPrintScreen:  "PrintScreen",
	KeyPause:        "Pause",
	KeyF1:           "F1",
	KeyF2:           "F2",
	KeyF3:           "F3",
	KeyF4:           "F4",
	KeyF5:           "F5",
	KeyF6:           "F6",
	KeyF7:           "F7",
	KeyF8:           "F8",
	KeyF9:           "F9",
	KeyF10:          "F10",
	KeyF11:          "F11",
	KeyF12:          "F12",
	KeyF13:          "F13",
	KeyF14:          "F14",
	KeyF15:          "F15",
	KeyF16:          "F16",
	KeyF17:          "F17",
	KeyF18:          "F18",
	KeyF19:          "F19",
	KeyF20:          "F20",
	KeyF21:          "F21",
	KeyF22:          "F22",
	KeyF23:          "F23",
	KeyF24:          "F24",
	KeyF25:          "F25",
	KeyKP0:          "KP0",
	KeyKP1:          "KP1",
	KeyKP2:          "KP2",
	KeyKP3:          "KP3",
	KeyKP4:          "KP4",
	KeyKP5:          "KP5",
	KeyKP6:          "KP6",
	KeyKP7:          "KP7",
	KeyKP8:          "KP8",
	KeyKP9:          "KP9",
	KeyKPDecimal:    "KPDecimal",
	KeyKPDivide:     "KPDivide",
	KeyKPMultiply:   "KPMultiply",
	KeyKPSubtract:   "KPSubtract",
	KeyKPAdd:        "KPAdd",
	KeyKPEnter:      "KPEnter",
	KeyKPEqual:      "KPEqual",
	KeyLeftShift:    "LeftShift",
	KeyLeftControl:  "LeftControl",
	KeyLeftAlt:      "LeftAlt",
	KeyLeftSuper:    "LeftSuper",
	KeyRightShift:   "RightShift",
	KeyRightControl: "RightControl",
	KeyRightAlt:     "RightAlt",
	KeyRightSuper:   "RightSuper",
	KeyMenu:         "Menu",
}
