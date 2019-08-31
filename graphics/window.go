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
	pollInputs(window)
}

/*
Input handling
*/

var (
	KeyMap  = map[string]bool{}
	MouseX  float64
	MouseY  float64
	RButton = false
	LButton = false
)

func pollInputs(window *glfw.Window) {
	glfw.PollEvents()
	pollKeys(window)
	pollMouse(window)
}

func pollKeys(window *glfw.Window) {
	KeyMap["w"] = window.GetKey(glfw.KeyW) == glfw.Press
	KeyMap["a"] = window.GetKey(glfw.KeyA) == glfw.Press
	KeyMap["s"] = window.GetKey(glfw.KeyS) == glfw.Press
	KeyMap["d"] = window.GetKey(glfw.KeyD) == glfw.Press
}

func KeyComboPressed(keys []string) bool {
	for _, val := range keys {
		if !KeyMap[val] {
			return false
		}
	}

	return true
}

func pollMouse(window *glfw.Window) {
	RButton = window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press
	LButton = window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
	MouseX, MouseY = window.GetCursorPos()
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func DestroyWindow(window *glfw.Window) {
	window.Destroy()
}
