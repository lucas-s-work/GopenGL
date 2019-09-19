package main

import (
	"Gopengl/graphics"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	graphics.Init()
	window := graphics.CreateWindow(800, 600, "test")
	graphics.SetWindow(window)

	go tick()
	graphics.Listen()
}

func tick() {
	for i := 0; i < 14; i++ {
		ro := graphics.CreateRenderObjectJob(168, "./sprites/test.png", true)
		ro.AddSquareJob(0, 0, 0, 0, 10, 10)
	}
}
