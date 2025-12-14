package main

import (
	"fmt"
	"runtime"
	"time"

	"go-conway/game"
	"go-conway/gfx"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width     = 500
	height    = 500
	rows      = 300
	columns   = 300
	threshold = 0.15
	fps       = 10
)

func main() {
	runtime.LockOSThread()

	window := gfx.InitGlfw(width, height, "Conway's Game of Life")
	defer glfw.Terminate()

	program := gfx.InitOpenGL()

	cells := game.MakeCells(rows, columns, threshold)

	// --- FPS Variables ---
	previousTime := glfw.GetTime()
	frameCount := 0

	for !window.ShouldClose() {
		t := time.Now()

		// Update logic
		for x := range cells {
			for _, c := range cells[x] {
				c.CheckState(cells)
			}
		}

		// Draw logic
		draw(cells, window, program)

		// --- FPS Calculation ---
		currentTime := glfw.GetTime()
		frameCount++

		if currentTime-previousTime >= 1.0 {
			window.SetTitle(fmt.Sprintf("Conway's Game of Life | FPS: %d", frameCount))
			frameCount = 0
			previousTime = currentTime
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func draw(cells [][]*game.Cell, window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	for x := range cells {
		for _, c := range cells[x] {
			c.Draw()
		}
	}

	glfw.PollEvents()
	window.SwapBuffers()
}
