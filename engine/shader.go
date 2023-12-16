package engine

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	vertexShaderSource   = "shaders/vertexShader.glsl"
	fragmentShaderSource = "shaders/fragmentShader.glsl"
)

type Shader struct {
	id       uint32
	uniforms map[string]int32
}

func NewShader(vertexPath, fragmentPath string) (Shader, error) {
	// Compile shader src
	vShader := loadShader(vertexPath, gl.VERTEX_SHADER)
	fShader := loadShader(fragmentPath, gl.FRAGMENT_SHADER)

	// Create program and bind shaders
	id := gl.CreateProgram()
	gl.AttachShader(id, vShader)
	gl.AttachShader(id, fShader)
	gl.LinkProgram(id)

	// Log errors
	var status int32
	gl.GetProgramiv(id, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(id, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(id, logLength, nil, gl.Str(log))
		return Shader{}, fmt.Errorf("shader linking error: %v", log)
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)
	return Shader{id, make(map[string]int32)}, nil
}

func (s Shader) Use() {
	gl.UseProgram(s.id)
}

func (s Shader) SetBool(name string, value bool) {
	var v int32 = gl.TRUE
	if !value {
		v = gl.FALSE
	}
	gl.Uniform1i(s.UniformLoc(name), v)
}

func (s Shader) SetInt(name string, v int32) {
	gl.Uniform1i(s.UniformLoc(name), v)
}

func (s Shader) SetFloat(name string, v float32) {
	gl.Uniform1f(s.UniformLoc(name), v)
}

func (s Shader) SetVec2(name string, v mgl32.Vec2) {
	gl.Uniform2f(s.UniformLoc(name), v[0], v[1])
}
func (s Shader) SetVec4(name string, v mgl32.Vec4) {
	gl.Uniform4f(s.UniformLoc(name), v[0], v[1], v[2], v[3])
}

func (s Shader) SetMatrix(name string, value mgl32.Mat4) {
	mat := [16]float32(value)
	gl.UniformMatrix4fv(s.UniformLoc(name), 1, false, &mat[0])
}

func (s *Shader) UniformLoc(name string) int32 {
	n := gl.Str(name + "\x00") // OpgenGL requires null termination character
	loc, ok := s.uniforms[name]
	if !ok {
		s.uniforms[name] = gl.GetUniformLocation(s.id, n)
	}

	return loc
}

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

func loadShader(filepath string, sType uint32) uint32 {
	src := readFile(filepath)
	shader, err := compileShader(src, sType)
	if err != nil {
		panic(err)
	}
	return shader
}

func readFile(filepath string) string {
	text, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	text = append(text, '\x00')
	return string(text)
}
