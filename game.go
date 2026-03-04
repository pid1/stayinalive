package main

import "math/rand"

// Grid holds the 2D cell state for Conway's Game of Life.
type Grid struct {
	cells  [][]bool
	width  int
	height int
}

// NewGrid allocates a w x h grid with all cells dead.
// If w or h is less than 1, a 1x1 grid is returned.
func NewGrid(w, h int) *Grid {
	if w <= 0 {
		w = 1
	}
	if h <= 0 {
		h = 1
	}
	cells := make([][]bool, h)
	for y := range cells {
		cells[y] = make([]bool, w)
	}
	return &Grid{cells: cells, width: w, height: h}
}

// Randomize sets each cell to alive with the given probability density.
// Density is clamped to [0.0, 1.0].
func (g *Grid) Randomize(density float64) {
	if density < 0.0 {
		density = 0.0
	}
	if density > 1.0 {
		density = 1.0
	}
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.cells[y][x] = rand.Float64() < density
		}
	}
}

// Get returns the cell state at (x, y) with toroidal wrapping.
func (g *Grid) Get(x, y int) bool {
	x = ((x % g.width) + g.width) % g.width
	y = ((y % g.height) + g.height) % g.height
	return g.cells[y][x]
}

// Set sets the cell state at (x, y) with toroidal wrapping.
func (g *Grid) Set(x, y int, alive bool) {
	x = ((x % g.width) + g.width) % g.width
	y = ((y % g.height) + g.height) % g.height
	g.cells[y][x] = alive
}

// CountNeighbors counts the alive Moore neighbors of cell (x, y).
func (g *Grid) CountNeighbors(x, y int) int {
	count := 0
	offsets := [8][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}
	for _, o := range offsets {
		if g.Get(x+o[0], y+o[1]) {
			count++
		}
	}
	return count
}

// Tick advances the simulation by one generation using the double-buffer pattern.
// Conway's rules: alive cell survives with 2 or 3 neighbors; dead cell is born with exactly 3.
func (g *Grid) Tick() *Grid {
	next := NewGrid(g.width, g.height)
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			neighbors := g.CountNeighbors(x, y)
			alive := g.cells[y][x]
			if alive && (neighbors == 2 || neighbors == 3) {
				next.cells[y][x] = true
			} else if !alive && neighbors == 3 {
				next.cells[y][x] = true
			}
		}
	}
	return next
}

// CountAlive returns the number of alive cells in the grid.
func (g *Grid) CountAlive() int {
	count := 0
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.cells[y][x] {
				count++
			}
		}
	}
	return count
}
