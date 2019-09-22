package opengl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func GlInit() {
	err := gl.Init()

	if err != nil {
		panic(err)
	}
}

func Render(vaos []*VAO) {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	var vao *VAO

	for i := 0; i < len(vaos); i++ {
		vao = vaos[i]

		vertNum := vao.PrepRender()
		gl.DrawArrays(gl.TRIANGLES, 0, vertNum)
		vao.FinishRender()
	}
}
