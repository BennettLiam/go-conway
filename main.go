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
	rows      = 1000
	columns   = 1000
	threshold = 0.5
	fps       = 60
)

func main() {
	runtime.LockOSThread()

	window := gfx.InitGlfw(width, height, "Conway's Game of Life")
	defer glfw.Terminate()

	program := gfx.InitOpenGL()
	vao, vbo := gfx.CreateBatchObjects()

	cells := game.MakeCells(rows, columns, threshold)

	// Pre-allocate the vertex slice.
	// Max vertices = rows * cols * 18 floats.
	// We make it 0 length, but full capacity to prevent re-allocations.
	batchData := make([]float32, 0, rows*columns*18)

	// Pre-calculate the width/height of a single cell in OpenGL coordinates (-1 to 1)
	// The total width in OpenGL is 2.0 (-1.0 to 1.0).
	cellW := 2.0 / float32(columns)
	cellH := 2.0 / float32(rows)

	previousTime := glfw.GetTime()
	frameCount := 0

	for !window.ShouldClose() {
		t := time.Now()

		batchData = batchData[:0]

		for x := range cells {
			for _, c := range cells[x] {
				c.CheckState(cells)

				if c.Alive {
					// Calculate position on the fly
					// OpenGL coordinates start at -1.0 (left/bottom)
					px := -1.0 + (float32(c.X) * cellW)
					py := -1.0 + (float32(c.Y) * cellH)

					// Append the 2 triangles (6 vertices, x/y/z)
					// We use hardcoded offsets to build the square
					batchData = append(batchData,
						px, py+cellH, 0, // Top-Left
						px, py, 0, // Bottom-Left
						px+cellW, py, 0, // Bottom-Right

						px, py+cellH, 0, // Top-Left
						px+cellW, py+cellH, 0, // Top-Right
						px+cellW, py, 0, // Bottom-Right
					)
				}
			}
		}

		draw(batchData, window, program, vao, vbo)

		currentTime := glfw.GetTime()
		frameCount++

		if currentTime-previousTime >= 1.0 {
			window.SetTitle(fmt.Sprintf("Conway's Game of Life | FPS: %d | RAM: High but optimized", frameCount))
			frameCount = 0
			previousTime = currentTime
		}

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func draw(data []float32, window *glfw.Window, program uint32, vao, vbo uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gfx.UpdateBatch(vbo, data)
	gfx.DrawBatch(vao, int32(len(data)/3))

	glfw.PollEvents()
	window.SwapBuffers()
}
