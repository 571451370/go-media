package f64

import (
	"image/color"
	"math"
)

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

func (p Vec2) Normalize() Vec2 {
	l := p.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{p.X / l, p.Y / l}
}

func (p Vec2) Scale(kx, ky float64) Vec2 {
	return Vec2{p.X * kx, p.Y * ky}
}

func (p Vec2) Lerp(q Vec2, t float64) Vec2 {
	return Vec2{
		Lerp(p.X, q.X, t),
		Lerp(p.Y, q.Y, t),
	}
}

func (p Vec2) Distance(q Vec2) float64 {
	return p.Sub(q).Len()
}

func (p Vec2) Polar() Polar {
	return Polar{p.Len(), math.Atan2(p.Y, p.X)}
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

func (p Vec3) Cross(q Vec3) Vec3 {
	return Vec3{
		p.Y*q.Z - p.Z*q.Y,
		p.Z*q.X - p.X*q.Z,
		p.X*q.Y - p.Y*q.X,
	}
}

func (p Vec3) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Vec3) Scale(kx, ky, kz float64) Vec3 {
	return Vec3{p.X * kx, p.Y * ky, p.Z * kz}
}

func (p Vec3) Normalize() Vec3 {
	l := p.Len()
	if l == 0 {
		return Vec3{}
	}
	return Vec3{p.X / l, p.Y / l, p.Z / l}
}

func (p Vec3) Distance(q Vec3) float64 {
	return p.Sub(q).Len()
}

func (p Vec3) Lerp(q Vec3, t float64) Vec3 {
	return Vec3{
		Lerp(p.X, q.X, t),
		Lerp(p.Y, q.Y, t),
		Lerp(p.Z, q.Z, t),
	}
}

func (p Vec3) RGBA() (r, g, b, a uint32) {
	c := color.RGBA{
		uint8(Clamp(p.X*255, 0, 255)),
		uint8(Clamp(p.Y*255, 0, 255)),
		uint8(Clamp(p.Z*255, 0, 255)),
		255,
	}
	return c.RGBA()
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (p Vec4) Dot(q Vec4) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

func (p Vec4) RGBA() (r, g, b, a uint32) {
	c := color.RGBA{
		uint8(Clamp(p.X*255, 0, 255)),
		uint8(Clamp(p.Y*255, 0, 255)),
		uint8(Clamp(p.Z*255, 0, 255)),
		uint8(Clamp(p.W*255, 0, 255)),
	}
	return c.RGBA()
}

type Mat3 [3][3]float64

func (m *Mat3) Identity() *Mat3 {
	*m = Mat3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
	return m
}

func (m *Mat3) Mul(a, b *Mat3) *Mat3 {
	var p Mat3
	for i := range a {
		for j := range a[i] {
			for k := range a[j] {
				p[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	*m = p
	return m
}

type Aff2 Mat3

func (m *Aff2) Identity() *Aff2 {
	return (*Aff2)((*Mat3)(m).Identity())
}

func (m *Aff2) Mul(a, b *Aff2) *Aff2 {
	return (*Aff2)((*Mat3)(m).Mul((*Mat3)(a), (*Mat3)(b)))
}

func (m *Aff2) Translate(v Vec2) *Aff2 {
	t := &Aff2{
		{1, 0, v.X},
		{0, 1, v.Y},
		{0, 0, 1},
	}
	return m.Mul(m, t)
}

func (m *Aff2) Scale(v Vec2) *Aff2 {
	s := &Aff2{
		{v.X, 0, 0},
		{0, v.Y, 0},
		{0, 0, 1},
	}
	return m.Mul(m, s)
}

func (m *Aff2) Shear(v Vec2) *Aff2 {
	s := &Aff2{
		{0, v.X, 0},
		{v.Y, 0, 0},
		{0, 0, 1},
	}
	return m.Mul(m, s)
}

func (m *Aff2) Rotate(rad float64) *Aff2 {
	s, c := math.Sincos(rad)
	r := &Aff2{
		{c, -s, 0},
		{s, c, 0},
		{0, 0, 1},
	}
	return m.Mul(m, r)
}

func (m *Aff2) Transform(v Vec2) Vec2 {
	return Vec2{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2],
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2],
	}
}

type Mat4 [4][4]float64

type Polar struct {
	R, P float64
}

func (p *Polar) Cartesian() Vec2 {
	s, c := math.Sincos(p.P)
	return Vec2{p.R * c, p.R * s}
}

type Quat struct {
	R, I, J, K float64
}

func Lerp(a, b, t float64) float64 {
	return a*t + (1-t)*b
}

func Smoothstep(a, b, x float64) float64 {
	t := Clamp((x-a)/(b-a), 0, 1)
	return t * t * (3 - 2*t)
}

func Clamp(x, s, e float64) float64 {
	if x < s {
		x = s
	}
	if x > e {
		x = e
	}
	return x
}
