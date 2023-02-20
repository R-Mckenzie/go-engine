package graphics

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	vertexShaderSource   = "graphics/shaders/vertexShader.glsl"
	fragmentShaderSource = "graphics/shaders/fragmentShader.glsl"
)

func compileShader(src string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", src, log)
	}
	return shader, nil
}

func LoadVertexShader(filepath string) uint32 {
	src := readFile(filepath)
	vShader, err := compileShader(src, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	return vShader
}

func LoadFragmentShader(filepath string) uint32 {
	src := readFile(filepath)
	fShader, err := compileShader(src, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	return fShader
}

func readFile(filepath string) string {
	text, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return string(text)
}
