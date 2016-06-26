package f64

import "math"

type Vec2 struct {
	X, Y float64
}

func (p Vec2) Add(q Vec2) Vec2 {
	return Vec2{p.X + q.X, p.Y + q.Y}
}

func (p Vec2) Sub(q Vec2) Vec2 {
	return Vec2{p.X - q.X, p.Y - q.Y}
}

func (p Vec2) Dot(q Vec2) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Vec2) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

type Vec3 struct {
	X, Y, Z float64
}

func (p Vec3) Add(q Vec3) Vec3 {
	return Vec3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Vec3) Sub(q Vec3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Vec3) Dot(q Vec3) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p Vec3) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (p Vec4) Dot(q Vec4) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

type Mat4 [4][4]float64
