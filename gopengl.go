package main

import (
	"fmt"
	"gopengl/graphics"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	//Create opengl context

	window := graphics.CreateWindow(800, 600, "test window")
	graphics.Init()
	graphics.SetWindow(window)
	graphics.SetWindowSize(800, 600)
	go startTick()
	graphics.Listen()
}

func startTick() {
	obj := graphics.CreateRenderObjectJob(12, "./sprites/test.png", true)
	vert := obj.AddSquareJob(400, 400, 0, 0, 100, 16)
	obj.AddSquareJob(600, 400, 0, 0, 100, 16)
	fmt.Println(*vert)
	// fmt.Println(*vert2)
}
