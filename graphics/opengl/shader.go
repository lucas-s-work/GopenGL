package opengl

import (
	"Gopengl/util"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Shader types

const (
	INVALID    = 0
	VERTSHADER = gl.VERTEX_SHADER
	FRAGSHADER = gl.FRAGMENT_SHADER
	GEOMSHADER = gl.GEOMETRY_SHADER
)

var (
	storedShaders []*shader
)

type shader struct {
	Id   uint32
	file string
}

type Program struct {
	Id         uint32
	attributes map[string]uint32
	uniforms   map[string]uniform
}

func CreateProgram(Id uint32) *Program {
	if Id == 0 {
		Id = gl.CreateProgram()
	}

	return &Program{
		Id,
		make(map[string]uint32),
		make(map[string]uniform),
	}
}

func (program *Program) AttachShader(s *shader) {
	gl.AttachShader(program.Id, s.Id)
}

/*
Load all desired shaders then call program.Link()
*/

func ReadFile(source string) (string, error) {
	data, err := ioutil.ReadFile(util.RelativePath(source))

	fmt.Println(util.RelativePath(source))

	if err != nil {
		return "", err
	}

	return string(data[:]) + "\x00", nil
}

/*
Load and attach shaders, if the shader has already been loaded it is not re-created.
*/
func (program *Program) LoadVertShader(source string) {
	existingShader := findShader(source)

	if existingShader != nil {
		program.AttachShader(existingShader)

		return
	}

	rawData, err := ReadFile(source)

	if err != nil {
		panic(fmt.Errorf("Unable to find vertex shader file: %s, err: %s", source, err.Error()))
	}

	program.loadShader(rawData, VERTSHADER)
}

func (program *Program) LoadFragShader(source string) {
	existingShader := findShader(source)

	if existingShader != nil {
		program.AttachShader(existingShader)

		return
	}

	rawData, err := ReadFile(source)

	if err != nil {
		panic(fmt.Errorf("Unable to find vertex shader file: %s", source))
	}

	program.loadShader(rawData, FRAGSHADER)
}

func (program *Program) loadShader(rawData string, shaderType uint32) {
	shader := gl.CreateShader(shaderType)
	source, free := gl.Strs(rawData)

	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		panic(fmt.Errorf("failed to compile %v: %v", source, log))
	}

	gl.AttachShader(program.Id, shader)
}

func findShader(file string) *shader {
	for _, s := range storedShaders {
		if s.file == file {
			return s
		}
	}

	return nil
}

/*
Shader methods
*/

func (p *Program) Use() {
	gl.UseProgram(p.Id)
}

func (p *Program) UnUse() {
	gl.UseProgram(0)
}

func (p *Program) Link() {
	gl.LinkProgram(p.Id)
}

/*
Attribute implementation
*/

func (p *Program) AddAttribute(attribute string) {
	attrib := gl.GetAttribLocation(p.Id, gl.Str(attribute+"\x00"))

	if attrib == -1 {
		panic("Invalid Attribute given")
	}

	p.attributes[attribute] = uint32(attrib)
}

func (p *Program) EnableAttribute(attribute string) uint32 {
	attributeValue := p.attributes[attribute]

	gl.EnableVertexAttribArray(attributeValue)

	return attributeValue
}

/*
Uniform implementation
*/

type uniform struct {
	id    uint32
	value interface{}
}

func (uni *uniform) ID() uint32 {
	return uni.id
}

func (uni *uniform) Value() interface{} {
	return uni.value
}

func (uni *uniform) Attach() {
	switch uni.value.(type) {
	case float32:
		value := (uni.value).(float32)
		gl.Uniform1f(int32(uni.id), value)
	case mgl32.Vec2:
		value := (uni.value).(mgl32.Vec2)
		gl.Uniform2f(int32(uni.id), value.X(), value.Y())
	case mgl32.Vec3:
		value := (uni.value).(mgl32.Vec3)
		gl.Uniform3f(int32(uni.id), value.X(), value.Y(), value.Z())
	case mgl32.Vec4:
		value := (uni.value).(mgl32.Vec4)
		gl.Uniform4f(int32(uni.id), value.X(), value.Y(), value.Z(), value.W())
	default:
		panic("Unsupported uniform type")
	}
}

func (p *Program) AddUniform(name string, value interface{}) {
	uni := uniform{
		uint32(gl.GetUniformLocation(p.Id, gl.Str(name+"\x00"))),
		value,
	}

	uni.Attach()
	p.uniforms[name] = uni

}

func (p *Program) SetUniform(name string, value interface{}) {
	if _, exists := p.uniforms[name]; !exists {
		panic("Attempting to set non existent uniform")
	}

	uni := uniform{
		p.uniforms[name].id,
		value,
	}

	uni.Attach()
	p.uniforms[name] = uni
}
