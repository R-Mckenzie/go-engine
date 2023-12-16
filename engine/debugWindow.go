package engine

import (
	"fmt"
	"sort"

	imgui "github.com/AllenDang/cimgui-go"
)

var debug map[string]string = make(map[string]string)

func debugLists() ([]string, []string) {
	labels := make([]string, 0)
	for k := range debug {
		labels = append(labels, k)
	}
	sort.Strings(labels)

	values := make([]string, 0)
	for _, l := range labels {
		values = append(values, debug[l])
	}
	return labels, values
}

func AddDebugInfo(label string, value interface{}) {
	debug[label] = fmt.Sprint(value)
}

func RemoveDebugInfo(label string) {
	delete(debug, label)
}

func displayDebug() {
	imgui.SetNextWindowSize(imgui.ImVec2{X: 200, Y: 200}, imgui.ImGuiCond(imgui.ImGuiCond_FirstUseEver))
	imgui.Begin("Stats", nil, 0)
	imgui.PushItemWidth(-100)
	k, v := debugLists()
	for i := range k {
		imgui.LabelText(k[i], v[i])
		imgui.Separator()
	}
	imgui.PopItemWidth()
	imgui.End()
}
