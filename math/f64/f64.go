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

func (p Vec2) Scalar(k float64) Vec2 {
	return Vec2{p.X * k, p.Y * k}
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

func (p Vec3) Scalar(k float64) Vec3 {
	return Vec3{p.X * k, p.Y * k, p.Z * k}
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

func (p Vec3) Abs() Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
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

type Mat4 [4][4]float64

func (m *Mat4) Identity() *Mat4 {
	*m = Mat4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Add(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] + b[i][j]
		}
	}
	return m
}

func (m *Mat4) Sub(a, b *Mat4) *Mat4 {
	for i := range a {
		for j := range a[i] {
			m[i][j] = a[i][j] - b[i][j]
		}
	}
	return m
}

func (m *Mat4) Mul(a, b *Mat4) *Mat4 {
	var p Mat4
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

func (m *Mat4) Translate(tx, ty, tz float64) *Mat4 {
	*m = Mat4{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Scale(sx, sy, sz float64) *Mat4 {
	*m = Mat4{
		{sx, 0, 0, 0},
		{0, sy, 0, 0},
		{0, 0, sz, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) LookAt(eye, at, up Vec3) *Mat4 {
	z := at.Sub(eye).Normalize()
	x := up.Cross(z).Normalize()
	y := z.Cross(x)

	var r, t Mat4
	r = Mat4{
		{x.X, y.X, z.X, 0},
		{x.Y, y.Y, z.Y, 0},
		{x.Z, y.Z, z.Z, 0},
		{0, 0, 0, 1},
	}
	t.Translate(-eye.X, -eye.Y, -eye.Z)
	m.Mul(&r, &t)
	return m
}

func (m *Mat4) RotX(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{1, 0, 0, 0},
		{0, co, -si, 0},
		{0, si, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotY(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{co, 0, si, 0},
		{0, 1, 0, 0},
		{-si, 0, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotZ(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{co, -si, 0, 0},
		{si, co, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Perspective(fovy, aspect, near, far float64) *Mat4 {
	f := math.Tan(fovy / 2)
	z := near - far
	*m = Mat4{
		{1 / (f * aspect), 0, 0, 0},
		{0, 1 / f, 0, 0},
		{0, 0, (-near - far) / z, 2 * far * near / z},
		{0, 0, 1, 0},
	}
	return m
}

func (m *Mat4) Ortho(l, r, b, t, n, f float64) *Mat4 {
	sx := 2 / (r - l)
	sy := 2 / (t - b)
	sz := -2 / (f - n)

	tx := -(r + l) / (r - l)
	ty := -(t + b) / (t - b)
	tz := -(f + n) / (f - n)

	*m = Mat4{
		{sx, 0, 0, tx},
		{0, sy, 0, ty},
		{0, 0, sz, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Inverse() *Mat4 {
	a := Vec3{m[0][0], m[1][0], m[2][0]}
	b := Vec3{m[0][1], m[1][1], m[2][1]}
	c := Vec3{m[0][2], m[1][2], m[2][2]}
	d := Vec3{m[0][3], m[1][3], m[2][3]}

	s := a.Cross(b)
	t := c.Cross(d)

	invDet := 1 / s.Dot(c)

	s = s.Scalar(invDet)
	t = t.Scalar(invDet)
	v := c.Scalar(invDet)

	r0 := b.Cross(v)
	r1 := v.Cross(a)

	*m = Mat4{
		{r0.X, r0.Y, r0.Z, -b.Dot(t)},
		{r1.X, r1.Y, r1.Z, a.Dot(t)},
		{s.X, s.Y, s.Z, -d.Dot(s)},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Transform(v Vec4) Vec4 {
	v = Vec4{
		m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z + m[0][3]*v.W,
		m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z + m[1][3]*v.W,
		m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z + m[2][3]*v.W,
		m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]*v.W,
	}
	if v.W != 0 {
		v.X /= v.W
		v.Y /= v.W
		v.Z /= v.W
		v.W = 1
	}
	return v
}

func (m *Mat4) Transform3(v Vec3) Vec3 {
	p := Vec4{v.X, v.Y, v.Z, 1}
	p = m.Transform(p)
	return Vec3{p.X, p.Y, p.Z}
}

type Polar struct {
	R, P float64
}

func (p *Polar) Cartesian() Vec2 {
	s, c := math.Sincos(p.P)
	return Vec2{p.R * c, p.R * s}
}

type Quat struct {
	X, Y, Z, W float64
}

func (q Quat) Add(r Quat) Quat {
	return Quat{
		q.X + r.X,
		q.Y + r.Y,
		q.Z + r.Z,
		q.W + r.W,
	}
}

func (q Quat) Sub(r Quat) Quat {
	return Quat{
		q.X - r.X,
		q.Y - r.Y,
		q.Z - r.Z,
		q.W - r.W,
	}
}

func (q Quat) Mul(r Quat) Quat {
	return Quat{
		q.X*r.W + q.Y*r.Z - q.Z*r.Y + q.W*r.X,
		q.Y*r.W + q.Z*r.X + q.W*r.Y - q.X*r.Z,
		q.Z*r.W + q.W*r.Z + q.X*r.Y - q.Y*r.X,
		q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z,
	}
}

func Lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func Unlerp(a, b, t float64) float64 {
	return (t - a) / (b - a)
}

func LinearRemap(x, a, b, c, d float64) float64 {
	return Lerp(c, d, Unlerp(a, b, x))
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
