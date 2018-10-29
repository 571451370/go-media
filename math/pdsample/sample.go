package pdsample

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type Sampler struct {
	neighbors []int
	grid      [][]int
	gridSize  int
	cellSize  float64
	Points    []f64.Vec2
	Radius    float64
	Tiled     bool
}

func NewSampler(radius float64, tiled, usesGrid bool) *Sampler {
	var (
		gridSize int
		cellSize int
		grid     [][]int
	)
	// grid size is chosen so that 4*radius search only
	// requires search adjacent cells, also determine
	// max points per cells
	if usesGrid {
		gridSize = int(math.Ceil(2 / (4 * radius)))
		if gridSize < 2 {
			gridSize = 2
		}
		cellSize = int(2 / float64(gridSize))
	}

	return &Sampler{
		grid:     grid,
		gridSize: gridSize,
		cellSize: cellSize,
		Radius:   radius,
		Tiled:    tiled,
	}
}
