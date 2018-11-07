// ported from https://github.com/ddunbar/PDSample
package pdsample

import (
	"math"
	"math/rand"

	"github.com/qeedquan/go-media/math/f64"
)

type Sampler struct {
	neighbors    []int
	grid         [][]int
	gridSize     int
	gridCellSize float64
	isTiled      bool
	Points       []f64.Vec2
	Radius       float64
}

func NewSampler(radius float64, isTiled, usesGrid bool) *Sampler {
	const maxPointsPerCell = 9

	var (
		gridSize     int
		gridCellSize float64
		grid         [][]int
	)
	// grid size is chosen so that 4*radius search only
	// requires search adjacent cells, also determine
	// max points per cells
	if usesGrid {
		gridSize = int(math.Ceil(2 / (4 * radius)))
		if gridSize < 2 {
			gridSize = 2
		}
		gridCellSize = 2 / float64(gridSize)
		grid = make([][]int, gridSize*gridSize)
		for i := range grid {
			grid[i] = make([]int, maxPointsPerCell)
		}
	}

	return &Sampler{
		grid:         grid,
		gridSize:     gridSize,
		gridCellSize: gridCellSize,
		isTiled:      isTiled,
		Radius:       radius,
	}
}

func (s *Sampler) pointInDomain(p f64.Vec2) bool {
	return -1 <= p.X && p.X <= 1 &&
		-1 <= p.Y && p.Y <= 1
}

func (s *Sampler) randomPoint() f64.Vec2 {
	return f64.Vec2{
		2*rand.Float64() - 1,
		2*rand.Float64() - 1,
	}
}

func (s *Sampler) getDistanceSquared(a, b f64.Vec2) float64 {
	v := s.getTiled(b.Sub(a))
	return v.X*v.X + v.Y*v.Y
}

func (s *Sampler) getTiled(v f64.Vec2) f64.Vec2 {
	x, y := v.X, v.Y
	if s.isTiled {
		if x < -1 {
			x += 2
		} else if x > 1 {
			x -= 2
		}

		if y < -1 {
			y += 2
		} else if y > 1 {
			y -= 2
		}
	}
	return f64.Vec2{x, y}
}

func (s *Sampler) getGridXY(v f64.Vec2) (gx, gy int) {
	gx = int(math.Floor(0.5 * (v.X + 1) * float64(s.gridSize)))
	gy = int(math.Floor(0.5 * (v.Y + 1) * float64(s.gridSize)))
	return
}

func (s *Sampler) addPoint(pt f64.Vec2) {
	s.Points = append(s.Points, pt)
	if s.grid != nil {
		gx, gy := s.getGridXY(pt)
		cell := s.grid[gy*s.gridSize+gx]
		for i := range cell {
			if cell[i] == -1 {
				cell[i] = len(s.Points) - 1
				return
			}
		}
		panic("internal error: overflowed max points per cell")
	}
}

func (s *Sampler) findNeighbors(pt f64.Vec2, distance float64) int {
	if s.grid == nil {
		panic("internal error: sampler cannot search without grid")
	}

	distanceSquared := distance * distance
	N := int(math.Ceil(distance / s.gridCellSize))
	if N > s.gridSize>>1 {
		N = s.gridSize >> 1
	}

	s.neighbors = s.neighbors[:0]
	gx, gy := s.getGridXY(pt)
	for j := -N; j <= N; j++ {
		for i := -N; i <= N; i++ {
			cx := (gx + i + s.gridSize) % s.gridSize
			cy := (gy + j + s.gridSize) % s.gridSize
			cell := s.grid[cy*s.gridSize+cx]

			for k := range cell {
				if cell[k] == -1 {
					break
				}
				if s.getDistanceSquared(pt, s.Points[cell[k]]) < distanceSquared {
					s.neighbors = append(s.neighbors, cell[k])
				}
			}
		}
	}
	return len(s.neighbors)
}

func (s *Sampler) findClosestNeighbor(pt f64.Vec2, distance float64) float64 {
	if s.grid == nil {
		panic("internal error: sampler cannot search without grid")
	}

	closestSquared := distance * distance
	N := int(math.Ceil(distance / s.gridCellSize))
	if N > s.gridSize>>1 {
		N = s.gridSize >> 1
	}

	gx, gy := s.getGridXY(pt)
	for j := -N; j <= N; j++ {
		for i := -N; i <= N; i++ {
			cx := (gx + i + s.gridSize) % s.gridSize
			cy := (gy + j + s.gridSize) % s.gridSize
			cell := s.grid[cy*s.gridSize+cx]

			for k := range cell {
				if cell[k] == -1 {
					break
				}
				d := s.getDistanceSquared(pt, s.Points[cell[k]])
				if d < closestSquared {
					closestSquared = d
				}
			}
		}
	}
	return math.Sqrt(closestSquared)
}

func (s *Sampler) findNeighborRanges(index int, rl *RangeList) {
	if s.grid == nil {
		panic("internal error, sampler cannot search without grid")
	}

}

func (s *Sampler) maximize() {
	rl := NewRangeList(0, 0)
	N := len(s.Points)

	for i := 0; i < N; i++ {
		candidate := &s.Points[i]
		rl.Reset(0, 2*math.Pi)
		s.findNeighborRanges(i, rl)
		for rl.NumRanges != 0 {
			re := rl.Ranges[rand.Intn(rl.NumRanges)]
			angle := re.Min + (re.Max-re.Min)*rand.Float64()
			pt := s.getTiled(f64.Vec2{
				candidate.X + math.Cos(angle)*2*s.Radius,
				candidate.Y + math.Sin(angle)*2*s.Radius,
			})
			s.addPoint(pt)
			rl.Subtract(angle-math.Pi/3, angle+math.Pi/3)
		}
	}
}

type DartThrowing struct {
	*Sampler
	minMaxThrows  int
	maxThrowsMult int
}

func NewDartThrowing(radius float64, isTiled bool, minMaxThrows, maxThrowsMult int) *DartThrowing {
	return &DartThrowing{}
}

func (d *DartThrowing) Complete() {
	for {
	}
}
