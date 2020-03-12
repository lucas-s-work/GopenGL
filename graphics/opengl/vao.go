package opengl

/*
Textured VAO implementation
TODO :
 - Add efficient update methods
 - Add support for multiple texture units.
*/

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const DEFAULT_VECTOR_SIZE = 2
const DEFAULT_TEXS_SIZE = 2

type VAO struct {
	ID                        uint32
	vertID                    uint32
	texID                     uint32
	rotGroupID                uint32
	windowWidth, windowHeight float32
	verts                     []float32
	texs                      []float32
	rotGroups                 []mgl32.Vec4 // Grouped rotations
	rot                       mgl32.Vec4   // Global VAO rotation
	trans                     mgl32.Vec2   // Global VAO translation, individual translation should be performed on each vertex
	vertNum                   int32
	shader                    *Program
	created, defaultShader    bool
	Texture                   *Texture
	uniforms                  map[string]interface{}
}

/*
VBO creation and modification functions
*/

//CreateVAO ... size of vao in vertices.
func CreateVAO(size uint32, textureSource string, defaultShader bool, width float32, height float32) *VAO {
	var vaoID, vertID, rotGroupID, texID uint32

	gl.GenVertexArrays(1, &vaoID)
	gl.GenBuffers(1, &vertID)
	gl.GenBuffers(1, &texID)
	gl.GenBuffers(1, &rotGroupID)

	var program *Program

	texture := LoadTexture(textureSource)
	vao := &VAO{
		vaoID,
		vertID,
		texID,
		rotGroupID,
		width,
		height,
		make([]float32, size*DEFAULT_VECTOR_SIZE),
		make([]float32, size*DEFAULT_TEXS_SIZE),
		make([]mgl32.Vec4, size),
		mgl32.Vec4{},
		mgl32.Vec2{},
		int32(size),
		program,
		false,
		defaultShader,
		texture,
		make(map[string]interface{}),
	}

	vao.DefaultShader()

	return vao
}

/*
VBO creation and setup
*/

func (vao *VAO) CreateBuffers() {
	if vao.created {
		panic(fmt.Errorf("Cannot recreate created VAO"))
	}

	vao.created = true

	gl.BindVertexArray(vao.ID)

	//vertex buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.vertID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vao.verts), gl.Ptr(vao.verts), gl.DYNAMIC_DRAW)
	vertAttrib := vao.shader.EnableAttribute("vert")
	gl.VertexAttribPointer(vertAttrib, DEFAULT_VECTOR_SIZE, gl.FLOAT, false, 0, nil)

	//grouped rotation buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.rotGroupID)
	vao.ResetGroupedRotation()
	rotGroups := destructureVecArray(vao.rotGroups)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(rotGroups), gl.Ptr(rotGroups), gl.DYNAMIC_DRAW)
	rotGroupAttrib := vao.shader.EnableAttribute("rotgroup")
	gl.VertexAttribPointer(rotGroupAttrib, 4, gl.FLOAT, false, 0, nil)

	//texture buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.texID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vao.texs), gl.Ptr(vao.texs), gl.DYNAMIC_DRAW)
	texAttrib := vao.shader.EnableAttribute("verttexcoord")
	gl.VertexAttribPointer(texAttrib, DEFAULT_TEXS_SIZE, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (vao *VAO) UpdateBuffers() {
	if !vao.created {
		vao.CreateBuffers()

		return
	}

	gl.BindVertexArray(vao.ID)
	// Verts
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.vertID)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(vao.verts), gl.Ptr(vao.verts))

	//Grouped rotations
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.rotGroupID)
	rotGroups := destructureVecArray(vao.rotGroups)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(rotGroups), gl.Ptr(rotGroups))

	// Texs
	gl.BindBuffer(gl.ARRAY_BUFFER, vao.texID)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(vao.texs), gl.Ptr(vao.texs))

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (vao *VAO) UpdateBufferIndex(index int, vert_data []float32, tex_data []float32) {
	if !vao.created {
		vao.CreateBuffers()
	}

	vao.UpdateVertBufferIndex(index, vert_data)
	vao.UpdateTexBufferIndex(index, tex_data)
}

// UpdateBufferData ... set the vert data of the vao and update the buffer
func (vao *VAO) UpdateBufferData(vertData []float32, texData []float32, rotGroupData []mgl32.Vec4) {
	vao.SetData(vertData, texData, rotGroupData)
	vao.UpdateBuffers()
}

func (vao *VAO) UpdateVertBufferData(vertData []float32) {
	vao.verts = vertData
	vao.UpdateBuffers()
}

func (vao *VAO) UpdateRotGroupBufferData(rotGroupData []mgl32.Vec4) {
	vao.rotGroups = rotGroupData
	vao.UpdateBuffers()
}

func (vao *VAO) UpdateTexBufferData(texData []float32) {
	vao.texs = texData
	vao.UpdateBuffers()
}

func (vao *VAO) UpdateVertBufferIndex(index int, vertData []float32) {
	index *= DEFAULT_VECTOR_SIZE

	for i, val := range vertData {
		vao.verts[index+i] = val
	}

	vao.UpdateBuffers()
}

// TODO add proper index update handling, currently non functional

func (vao *VAO) UpdateTexBufferIndex(index int, texData []float32) {
	index *= DEFAULT_TEXS_SIZE

	for i, val := range texData {
		vao.texs[index+i] = val
	}

	vao.UpdateBuffers()
}

// SetData ... set the vert/tex data of the vao, does not update the buffer
func (vao *VAO) SetData(vertData []float32, texData []float32, rotGroupData []mgl32.Vec4) {
	vao.verts = vertData
	vao.texs = texData
	vao.rotGroups = rotGroupData
}

/*
Global rotation
*/

func (vao *VAO) SetRotation(x, y, rad float32) {
	vao.rot = mgl32.Vec4{x, y, float32(math.Cos(float64(rad))), float32(math.Sin(float64(rad)))}
	vao.shader.SetUniform("rot", vao.rot)
}

/*
Per group rotation, start and end are the start and end vertices
*/

func (vao *VAO) SetGroupedRotation(x, y, rad float32, start, end int) {
	c := float32(math.Cos(float64(rad)))
	s := float32(math.Sin(float64(rad)))

	for i := start; i <= end; i++ {
		vao.rotGroups[i] = mgl32.Vec4{x, y, c, s}
	}
}

func (vao *VAO) SetAllGroupedRotation(x, y, rad float32) {
	vao.SetGroupedRotation(x, y, rad, 0, len(vao.rotGroups))
}

func (vao *VAO) ResetGroupedRotation() {
	for i := 0; i < len(vao.rotGroups); i++ {
		vao.rotGroups[i] = mgl32.Vec4{0, 0, 1, 0}
	}
}

func (vao *VAO) SetTranslation(x, y float32) {
	vao.trans = mgl32.Vec2{x, y}
	vao.shader.SetUniform("trans", vao.trans)
}

func (vao *VAO) Delete() {
	gl.DeleteBuffers(1, &vao.vertID)
	gl.DeleteBuffers(1, &vao.texID)
	gl.DeleteVertexArrays(1, &vao.ID)
}

/*
Render handling
*/

func (vao *VAO) PrepRender() int32 {
	vao.shader.Use()
	vao.PrepUniforms()
	gl.BindVertexArray(vao.ID)
	vao.Texture.Use()

	return vao.vertNum
}

func (vao *VAO) FinishRender() {
	// TODO: figure out if this needs removing at some point
}

/*
Utility
*/

func (vao *VAO) DefaultShader() Program {
	program := CreateProgram(0)
	vao.AttachProgram(program)

	program.LoadVertShader("./shaders/vertex.vert")
	program.LoadFragShader("./shaders/fragment.frag")
	program.Link()

	program.AddAttribute("vert")
	// Currently unusued, optimized out by the shader compiler so will fail
	program.AddAttribute("rotgroup")
	program.AddAttribute("verttexcoord")

	// Add and set rotation uniform
	vao.AddUniform("rot", mgl32.Vec4{})
	vao.SetRotation(0, 0, 0)

	// Other uniforms can use default values.
	vao.AddUniform("trans", mgl32.Vec2{})
	vao.AddUniform("dim", mgl32.Vec2{vao.windowWidth, vao.windowHeight})

	return *program
}

func (vao *VAO) AttachProgram(program *Program) {
	vao.shader = program
}

/*
Shader uniform implementation
*/

func (vao *VAO) AddUniform(name string, value interface{}) {
	vao.shader.AddUniform(name, value)

	vao.uniforms[name] = value
}

func (vao *VAO) PrepUniforms() {
	for id, uni := range vao.shader.uniforms {
		vao.shader.SetUniform(id, uni.Value())
	}
}

/*
converts an array of Vec3's into a float32 array for use by vbo's
*/
func destructureVecArray(vecs []mgl32.Vec4) []float32 {
	vecFloats := make([]float32, len(vecs)*4)

	for i, vec := range vecs {
		vecFloats[i*4] = vec.X()
		vecFloats[i*4+1] = vec.Y()
		vecFloats[i*4+2] = vec.Z()
		vecFloats[i*4+3] = vec.W()
	}

	return vecFloats
}
