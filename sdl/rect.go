package sdl

type Point struct {
	X, Y int32
}

type Rect struct {
	X, Y, W, H int32
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}
