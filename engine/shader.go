package engine

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

//go:embed shaders/vertexShader.glsl
var vertexShaderSource string

//go:embed shaders/fragmentShader.glsl
var fragmentShaderSource string

//go:embed shaders/postprocessFragment.glsl
var ppFragmentShaderSource string

//go:embed shaders/postprocessVertex.glsl
var ppVertexShaderSource string

//go:embed shaders/uiFragment.glsl
var uiFragmentSource string

type Shader struct {
	id       uint32
	uniforms map[string]int32
}

var shaderMap map[string]Shader

func loadShader(vertex, fragment string) Shader {
	shader := NewShaderFromFile(vertex, fragment)
	return shader
}

func LoadShader(vertex, fragment, name string) {
	shader := NewShaderFromFile(vertex, fragment)
	shaderMap[name] = shader
}

func createGLShader(vShader, fShader uint32) (Shader, error) {
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

func NewShaderFromString(vertexSrc, fragmentSrc string) Shader {
	vShader, err := compileShader(vertexSrc, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fShader, err := compileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	shader, err := createGLShader(vShader, fShader)
	if err != nil {
		panic(err)
	}
	return shader
}

func NewShaderFromFile(vertexPath, fragmentPath string) Shader {
	// Compile shader src
	vShader := loadShaderFile(vertexPath, gl.VERTEX_SHADER)
	fShader := loadShaderFile(fragmentPath, gl.FRAGMENT_SHADER)
	shader, err := createGLShader(vShader, fShader)
	if err != nil {
		panic(err)
	}
	return shader
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

func (s Shader) SetFloatArray(name string, i int32, v []float32) {
	gl.Uniform1fv(s.UniformLoc(name), i, &v[0])
}

func (s Shader) SetVec2(name string, v mgl32.Vec2) {
	gl.Uniform2f(s.UniformLoc(name), v[0], v[1])
}

func (s Shader) SetVec2Array(name string, i int32, v []float32) {
	gl.Uniform2fv(s.UniformLoc(name), i, &v[0])
}

func (s Shader) SetVec3(name string, v mgl32.Vec3) {
	gl.Uniform3f(s.UniformLoc(name), v[0], v[1], v[2])
}

func (s Shader) SetVec3Array(name string, i int32, v []float32) {
	gl.Uniform3fv(s.UniformLoc(name), i, &v[0])
}

func (s Shader) SetVec4(name string, v mgl32.Vec4) {
	gl.Uniform4f(s.UniformLoc(name), v[0], v[1], v[2], v[3])
}

func (s Shader) SetVec4Array(name string, i int32, v []float32) {
	gl.Uniform4fv(s.UniformLoc(name), i, &v[0])
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
	if src[len(src)-1] != '\x00' {
		src += string('\x00')
	}
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

func loadShaderFile(filepath string, sType uint32) uint32 {
	// Read source file
	text, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	text = append(text, '\x00')
	src := string(text)

	// Compile shader
	shader, err := compileShader(src, sType)
	if err != nil {
		panic(err)
	}
	return shader
}
