package game

import (
	"math/rand"
	"time"
)

type Cell struct {
	Alive     bool
	AliveNext bool

	X int
	Y int
}

func MakeCells(rows, columns int, threshold float64) [][]*Cell {
	rand.Seed(time.Now().UnixNano())

	cells := make([][]*Cell, rows)
	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)
			c.Alive = rand.Float64() < threshold
			c.AliveNext = c.Alive
			cells[x] = append(cells[x], c)
		}
	}
	return cells
}

func newCell(x, y int) *Cell {
	return &Cell{
		X: x,
		Y: y,
	}
}

func (c *Cell) CheckState(cells [][]*Cell) {
	c.Alive = c.AliveNext
	c.AliveNext = c.Alive

	liveCount := c.liveNeighbors(cells)
	if c.Alive {
		if liveCount < 2 {
			c.AliveNext = false
		}
		if liveCount == 2 || liveCount == 3 {
			c.AliveNext = true
		}
		if liveCount > 3 {
			c.AliveNext = false
		}
	} else {
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

	add(c.X-1, c.Y)
	add(c.X+1, c.Y)
	add(c.X, c.Y+1)
	add(c.X, c.Y-1)
	add(c.X-1, c.Y+1)
	add(c.X+1, c.Y+1)
	add(c.X-1, c.Y-1)
	add(c.X+1, c.Y-1)

	return liveCount
}
