package game

import (
	"math/rand"
	"time"

	"go-conway/gfx" // Import the local graphics package

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
)

type Cell struct {
	drawable uint32

	Alive     bool
	AliveNext bool

	x int
	y int
}

// MakeCells creates the grid
func MakeCells(rows, columns int, threshold float64) [][]*Cell {
	rand.Seed(time.Now().UnixNano())

	cells := make([][]*Cell, rows)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y, rows, columns)

			c.Alive = rand.Float64() < threshold
			c.AliveNext = c.Alive

			cells[x] = append(cells[x], c)
		}
	}

	return cells
}

func newCell(x, y, rows, columns int) *Cell {
	points := make([]float32, len(square))
	copy(points, square)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(columns)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(rows)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &Cell{
		drawable: gfx.MakeVao(points),
		x:        x,
		y:        y,
	}
}

// Draw performs the OpenGL draw call
func (c *Cell) Draw() {
	if !c.Alive {
		return
	}
	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

// CheckState determines the state of the cell for the next tick of the game.
func (c *Cell) CheckState(cells [][]*Cell) {
	c.Alive = c.AliveNext
	c.AliveNext = c.Alive

	liveCount := c.liveNeighbors(cells)
	if c.Alive {
		// 1. Underpopulation
		if liveCount < 2 {
			c.AliveNext = false
		}
		// 2. Survive
		if liveCount == 2 || liveCount == 3 {
			c.AliveNext = true
		}
		// 3. Overpopulation
		if liveCount > 3 {
			c.AliveNext = false
		}
	} else {
		// 4. Reproduction
		if liveCount == 3 {
			c.AliveNext = true
		}
	}
}

func (c *Cell) liveNeighbors(cells [][]*Cell) int {
	var liveCount int
	rows := len(cells)
	cols := len(cells[0])

	add := func(x, y int) {
		// Wrap around logic
		if x == rows {
			x = 0
		} else if x == -1 {
			x = rows - 1
		}
		if y == cols {
			y = 0
		} else if y == -1 {
			y = cols - 1
		}

		if cells[x][y].Alive {
			liveCount++
		}
	}

	add(c.x-1, c.y)   // Left
	add(c.x+1, c.y)   // Right
	add(c.x, c.y+1)   // Up
	add(c.x, c.y-1)   // Down
	add(c.x-1, c.y+1) // Top-Left
	add(c.x+1, c.y+1) // Top-Right
	add(c.x-1, c.y-1) // Bottom-Left
	add(c.x+1, c.y-1) // Bottom-Right

	return liveCount
}
