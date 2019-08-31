package graphics

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

//TODO implement GLFWError
func CreateWindow(width, height int, name string) *glfw.Window {
	err := glfw.Init()

	checkerr(err)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)

	checkerr(err)

	window.MakeContextCurrent()

	return window
}

func Poll(window *glfw.Window) {
	window.SwapBuffers()
	glfw.PollEvents()
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func DestroyWindow(window *glfw.Window) {
	window.Destroy()
}
