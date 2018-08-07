package f64

import (
	"fmt"
	"image"
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

func (p Vec2) AddScale(q Vec2, t float64) Vec2 {
	return p.Add(q.Scale(t))
}

func (p Vec2) SubScale(q Vec2, t float64) Vec2 {
	return p.Sub(q.Scale(t))
}

func (p Vec2) Neg() Vec2 {
	return Vec2{-p.X, -p.Y}
}

func (p Vec2) Abs() Vec2 {
	return Vec2{math.Abs(p.X), math.Abs(p.Y)}
}

func (p Vec2) Dot(q Vec2) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Vec2) MinComp() float64 {
	return math.Min(p.X, p.Y)
}

func (p Vec2) MaxComp() float64 {
	return math.Max(p.X, p.Y)
}

func (p Vec2) Max(q Vec2) Vec2 {
	return Vec2{
		math.Max(p.X, q.X),
		math.Max(p.Y, q.Y),
	}
}

func (p Vec2) Min(q Vec2) Vec2 {
	return Vec2{
		math.Min(p.X, q.X),
		math.Min(p.Y, q.Y),
	}
}

func (p Vec2) Floor() Vec2 {
	return Vec2{
		math.Floor(p.X),
		math.Floor(p.Y),
	}
}

func (p Vec2) Ceil() Vec2 {
	return Vec2{
		math.Ceil(p.X),
		math.Ceil(p.Y),
	}
}

func (p Vec2) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Vec2) LenSquared() float64 {
	return p.Dot(p)
}

func (p Vec2) Normalize() Vec2 {
	l := p.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{p.X / l, p.Y / l}
}

func (p Vec2) Scale2(q Vec2) Vec2 {
	return Vec2{p.X * q.X, p.Y * q.Y}
}

func (p Vec2) Scale(k float64) Vec2 {
	return Vec2{p.X * k, p.Y * k}
}

func (p Vec2) Shear(k float64) Vec2 {
	return Vec2{p.X + k*p.Y, p.Y + k*p.X}
}

func (p Vec2) Shearv(q Vec2) Vec2 {
	return Vec2{p.X + q.X*p.Y, p.Y + q.Y*p.X}
}

func (p Vec2) Lerp(t float64, q Vec2) Vec2 {
	return Vec2{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
	}
}

func (p Vec2) Lerp2(t, q Vec2) Vec2 {
	return Vec2{
		Lerp(t.X, p.X, q.X),
		Lerp(t.Y, p.Y, q.Y),
	}
}

func (p Vec2) Wrap(s, e float64) Vec2 {
	p.X = Wrap(p.X, s, e)
	p.Y = Wrap(p.Y, s, e)
	return p
}

func (p Vec2) Distance(q Vec2) float64 {
	return p.Sub(q).Len()
}

func (p Vec2) DistanceSquared(q Vec2) float64 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec2) Polar() Polar {
	return Polar{p.Len(), math.Atan2(p.Y, p.X)}
}

func (p Vec2) MinScalar(k float64) Vec2 {
	return Vec2{
		math.Min(p.X, k),
		math.Min(p.Y, k),
	}
}

func (p Vec2) MaxScalar(k float64) Vec2 {
	return Vec2{
		math.Max(p.X, k),
		math.Max(p.Y, k),
	}
}

func (p Vec2) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (p Vec2) Rotate(r float64) Vec2 {
	si, co := math.Sincos(r)
	return Vec2{
		p.X*co - p.Y*si,
		p.X*si + p.Y*co,
	}
}

func (p Vec2) Shrink(k float64) Vec2 {
	return Vec2{p.X / k, p.Y / k}
}

func (p Vec2) Shrink2(q Vec2) Vec2 {
	return Vec2{p.X / q.X, p.Y / q.Y}
}

func (p Vec2) YX() Vec2 {
	return Vec2{p.Y, p.X}
}

func (p Vec2) Clamp(s, e float64) Vec2 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	return p
}

func (p Vec2) Clamp2(s, e Vec2) Vec2 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	return p
}

func (p Vec2) Equals(q Vec2) bool {
	const eps = 1e-6
	return math.Abs(p.X-q.X) <= eps && math.Abs(p.Y-q.Y) <= eps
}

func (p Vec2) OnLine(a, b Vec2) bool {
	sx := math.Min(a.X, b.X)
	sy := math.Min(a.Y, b.Y)
	ex := math.Max(a.X, b.X)
	ey := math.Max(a.Y, b.Y)
	return sx <= p.X && p.X <= ex &&
		sy <= p.Y && p.Y <= ey
}

func (p Vec2) Angle(q Vec2) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
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

func (p Vec3) AddScale(q Vec3, t float64) Vec3 {
	return p.Add(q.Scale(t))
}

func (p Vec3) SubScale(q Vec3, t float64) Vec3 {
	return p.Sub(q.Scale(t))
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

func (p Vec3) CrossNormalize(q Vec3) Vec3 {
	return p.Cross(q).Normalize()
}

func (p Vec3) Neg() Vec3 {
	return Vec3{-p.X, -p.Y, -p.Z}
}

func (p Vec3) Reflect(q Vec3) Vec3 {
	q = q.Scale(2 * p.Dot(q))
	return p.Sub(q)
}

func (p Vec3) Refract(q Vec3, eta float64) Vec3 {
	x := p.Dot(q)
	k := 1 - eta*eta*(1-x*x)
	if k < 0 {
		return Vec3{}
	}
	a := q.Scale(eta)
	b := p.Scale(eta*x + math.Sqrt(k))
	return a.Sub(b)
}

func (p Vec3) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Vec3) LenSquared() float64 {
	return p.Dot(p)
}

func (p Vec3) Scale3(q Vec3) Vec3 {
	return Vec3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

func (p Vec3) Scale(k float64) Vec3 {
	return Vec3{p.X * k, p.Y * k, p.Z * k}
}

func (p Vec3) Shrink(k float64) Vec3 {
	return Vec3{
		p.X / k,
		p.Y / k,
		p.Z / k,
	}
}

func (p Vec3) Shrink3(q Vec3) Vec3 {
	return Vec3{
		p.X / q.X,
		p.Y / q.Y,
		p.Z / q.Z,
	}
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

func (p Vec3) DistanceSquared(q Vec3) float64 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec3) Lerp(t float64, q Vec3) Vec3 {
	return Vec3{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
		Lerp(t, p.Z, q.Z),
	}
}

func (p Vec3) Abs() Vec3 {
	return Vec3{
		math.Abs(p.X),
		math.Abs(p.Y),
		math.Abs(p.Z),
	}
}

func (p Vec3) MaxScalar(k float64) Vec3 {
	return Vec3{
		math.Max(p.X, k),
		math.Max(p.Y, k),
		math.Max(p.Z, k),
	}
}

func (p Vec3) MinScalar(k float64) Vec3 {
	return Vec3{
		math.Min(p.X, k),
		math.Min(p.Y, k),
		math.Min(p.Z, k),
	}
}

func (p Vec3) Max(q Vec3) Vec3 {
	return Vec3{
		math.Max(p.X, q.X),
		math.Max(p.Y, q.Y),
		math.Max(p.Z, q.Z),
	}
}

func (p Vec3) Min(q Vec3) Vec3 {
	return Vec3{
		math.Min(p.X, q.X),
		math.Min(p.Y, q.Y),
		math.Min(p.Z, q.Z),
	}
}

func (p Vec3) MinComp() float64 {
	return math.Min(p.X, math.Min(p.Y, p.Z))
}

func (p Vec3) MaxComp() float64 {
	return math.Max(p.X, math.Max(p.Y, p.Z))
}

func (p Vec3) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	return color.RGBA{
		uint8(Clamp(p.X+0.5, 0, 255)),
		uint8(Clamp(p.Y+0.5, 0, 255)),
		uint8(Clamp(p.Z+0.5, 0, 255)),
		255,
	}

}

func (p Vec3) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec3) Spherical() Spherical {
	l := p.Len()
	return Spherical{
		R: l,
		T: math.Acos(p.Z / l),
		P: math.Atan2(p.Y, p.X),
	}
}

func (p Vec3) Equals(q Vec3) bool {
	const eps = 1e-6
	return math.Abs(p.X-q.X) <= eps && math.Abs(p.Y-q.Y) <= eps &&
		math.Abs(p.Z-q.Z) <= eps
}

func (p Vec3) XY() Vec2 { return Vec2{p.X, p.Y} }
func (p Vec3) XZ() Vec2 { return Vec2{p.X, p.Z} }
func (p Vec3) YX() Vec2 { return Vec2{p.Y, p.X} }
func (p Vec3) YZ() Vec2 { return Vec2{p.Y, p.Z} }
func (p Vec3) ZX() Vec2 { return Vec2{p.Z, p.X} }
func (p Vec3) ZY() Vec2 { return Vec2{p.Z, p.Y} }

func (p Vec3) Clamp(s, e float64) Vec3 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	p.Z = Clamp(p.Z, s, e)
	return p
}

func (p Vec3) Clamp3(s, e Vec3) Vec3 {
	p.X = Clamp(p.X, s.X, e.X)
	p.Y = Clamp(p.Y, s.Y, e.Y)
	p.Z = Clamp(p.Z, s.Z, e.Z)
	return p
}

func (p Vec3) Angle(q Vec3) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
}

type Vec4 struct {
	X, Y, Z, W float64
}

func (p Vec4) Add(q Vec4) Vec4 {
	return Vec4{p.X + q.X, p.Y + q.Y, p.Z + q.Z, p.W}
}

func (p Vec4) Sub(q Vec4) Vec4 {
	return Vec4{p.X - q.X, p.Y - q.Y, p.Z - q.Z, p.W}
}

func (p Vec4) AddScale(q Vec4, k float64) Vec4 {
	return p.Add(q.Scale(k))
}

func (p Vec4) SubScale(q Vec4, k float64) Vec4 {
	return p.Sub(q.Scale(k))
}

func (p Vec4) Scale(k float64) Vec4 {
	return Vec4{p.X * k, p.Y * k, p.Z * k, p.W}
}

func (p Vec4) Scale3(q Vec3) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W}
}

func (p Vec4) Scale4(q Vec4) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W * q.W}
}

func (p Vec4) Shrink(k float64) Vec4 {
	return Vec4{p.X / k, p.Y / k, p.Z / k, p.W}
}

func (p Vec4) Shrink4(q Vec4) Vec4 {
	return Vec4{p.X / q.X, p.Y / q.Y, p.Z / q.Z, p.W}
}

func (p Vec4) Dot(q Vec4) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

func (p Vec4) Len() float64 {
	return math.Sqrt(p.Dot(p))
}

func (p Vec4) Dot3(q Vec4) float64 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p Vec4) Len3() float64 {
	return math.Sqrt(p.Dot3(p))
}

func (p Vec4) XYZ() Vec3 { return Vec3{p.X, p.Y, p.Z} }
func (p Vec4) XZY() Vec3 { return Vec3{p.X, p.Z, p.Y} }
func (p Vec4) YXZ() Vec3 { return Vec3{p.Y, p.X, p.Z} }
func (p Vec4) YZX() Vec3 { return Vec3{p.Y, p.Z, p.X} }

func (p Vec4) Normalize() Vec4 {
	l := p.Len()
	if l == 0 {
		return Vec4{0, 0, 0, p.W}
	}
	return Vec4{
		p.X / l,
		p.Y / l,
		p.Z / l,
		p.W / l,
	}
}

func (p Vec4) Normalize3() Vec4 {
	l := p.Len3()
	if l == 0 {
		return Vec4{0, 0, 0, p.W}
	}
	return Vec4{
		p.X / l,
		p.Y / l,
		p.Z / l,
		p.W,
	}
}

func (p Vec4) Floor() Vec4 {
	return Vec4{
		math.Floor(p.X),
		math.Floor(p.Y),
		math.Floor(p.Z),
		math.Floor(p.W),
	}
}

func (p Vec4) ToRGBA() color.RGBA {
	if 0 <= p.X && p.X <= 1 {
		p.X *= 255
	}
	if 0 <= p.Y && p.Y <= 1 {
		p.Y *= 255
	}
	if 0 <= p.Z && p.Z <= 1 {
		p.Z *= 255
	}
	if 0 <= p.W && p.W <= 1 {
		p.W *= 255
	}
	c := color.RGBA{
		uint8(Clamp(p.X+0.5, 0, 255)),
		uint8(Clamp(p.Y+0.5, 0, 255)),
		uint8(Clamp(p.Z+0.5, 0, 255)),
		uint8(Clamp(p.W+0.5, 0, 255)),
	}
	return c
}

func (p Vec4) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

func (p Vec4) Angle(q Vec4) float64 {
	a := p.Len()
	b := q.Len()
	d := p.Dot(q)
	return math.Acos(d / (a * b))
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

func (m *Mat3) Transform(p Vec3) Vec3 {
	return Vec3{
		m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z,
		m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z,
		m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z,
	}
}

func (m *Mat3) Transform2(p Vec2) Vec2 {
	v := m.Transform(Vec3{p.X, p.Y, 1})
	if v.Z != 0 {
		v.X /= v.Z
		v.Y /= v.Z
	}
	return Vec2{v.X, v.Y}
}

func (m *Mat3) Transpose() *Mat3 {
	var p Mat3
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
	return m
}

func (m *Mat3) Det() float64 {
	A1 := Vec3{m[0][0], m[0][1], m[0][2]}
	A2 := Vec3{m[1][1], m[1][2], m[1][0]}
	A3 := Vec3{m[2][2], m[2][0], m[2][1]}
	A4 := Vec3{m[1][2], m[1][0], m[1][1]}
	A5 := Vec3{m[2][1], m[2][2], m[2][0]}

	X := A2.Scale3(A3)
	Y := A4.Scale3(A5)
	A6 := X.Sub(Y)
	return A1.Dot(A6)
}

func (m *Mat3) Trace() float64 {
	return m[0][0] + m[1][1] + m[2][2]
}

func (m *Mat3) Mat4() Mat4 {
	return Mat4{
		{m[0][0], m[1][0], m[2][0], 0},
		{m[0][1], m[1][1], m[2][1], 0},
		{m[0][2], m[1][2], m[2][2], 0},
		{0, 0, 0, 1},
	}
}

func (m Mat3) String() string {
	return fmt.Sprintf(`
Mat3[% 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
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

func (m *Mat4) Trace() float64 {
	return m[0][0] + m[1][1] + m[2][2] + m[3][3]
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

func (m *Mat4) LookAt(eye, center, up Vec3) *Mat4 {
	f := center
	f = f.Sub(eye)
	f = f.Normalize()

	s := f.Cross(up)
	s = s.Normalize()
	u := s.Cross(f)

	*m = Mat4{
		{s.X, s.Y, s.Z, 0},
		{u.X, u.Y, u.Z, 0},
		{-f.X, -f.Y, -f.Z, 0},
		{0, 0, 0, 1},
	}

	var t Mat4
	t.Translate(-eye.X, -eye.Y, -eye.Z)
	m.Mul(&t, m)
	return m
}

func (m *Mat4) Viewport(x, y, w, h float64) {
	*m = Mat4{
		{w / 2, 0, 0, w/2 + x},
		{0, h / 2, 0, h/2 + y},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m *Mat4) RotateX(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{1, 0, 0, 0},
		{0, co, -si, 0},
		{0, si, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotateY(r float64) *Mat4 {
	si, co := math.Sincos(r)
	*m = Mat4{
		{co, 0, si, 0},
		{0, 1, 0, 0},
		{-si, 0, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotateZ(r float64) *Mat4 {
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

	s = s.Scale(invDet)
	t = t.Scale(invDet)
	v := c.Scale(invDet)

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

func (m *Mat4) Transpose() *Mat4 {
	var p Mat4
	for i := range m {
		for j := range m[i] {
			p[j][i] = m[i][j]
		}
	}
	*m = p
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
	s := m[3][0]*v.X + m[3][1]*v.Y + m[3][2]*v.Z + m[3][3]
	switch s {
	case 0:
		return Vec3{}
	default:
		p := m.Transform(Vec4{v.X, v.Y, v.Z, 1})
		return Vec3{p.X, p.Y, p.Z}
	}

}

func (m Mat4) String() string {
	return fmt.Sprintf(`
Mat4[% 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
     % 0.3f, % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3])
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

func (q Quat) Neg() Quat {
	return Quat{-q.X, -q.Y, -q.Z, -q.W}
}

func (q Quat) Dot(p Quat) float64 {
	return q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
}

func (q Quat) Scale(k float64) Quat {
	return Quat{
		q.X * k,
		q.Y * k,
		q.Z * k,
		q.W * k,
	}
}

func (q Quat) Len() float64 {
	return math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
}

func (q Quat) Normalize() Quat {
	l := q.Len()
	if l == 0 {
		return Quat{}
	}
	return Quat{
		q.X / l,
		q.Y / l,
		q.Z / l,
		q.W / l,
	}
}

func (q Quat) Conj() Quat {
	return Quat{-q.X, -q.Y, -q.Z, q.W}
}

func (q Quat) FromAxis(v Vec3, r float64) Quat {
	r *= 0.5
	vn := v.Normalize()
	si, co := math.Sincos(r)
	return Quat{
		vn.X * si,
		vn.Y * si,
		vn.Z * si,
		co,
	}
}

func (q Quat) FromEuler(pitch, yaw, roll float64) Quat {
	p := pitch / 2
	y := yaw / 2
	r := roll / 2

	sinp := math.Sin(p)
	siny := math.Sin(y)
	sinr := math.Sin(r)
	cosp := math.Cos(p)
	cosy := math.Cos(y)
	cosr := math.Cos(r)

	return Quat{
		sinr*cosp*cosy - cosr*sinp*siny,
		cosr*sinp*cosy + sinr*cosp*siny,
		cosr*cosp*siny - sinr*sinp*cosy,
		cosr*cosp*cosy + sinr*sinp*siny,
	}.Normalize()
}

func (q Quat) Mat3() Mat3 {
	x, y, z, w := q.X, q.Y, q.Z, q.W
	x2 := x * x
	y2 := y * y
	z2 := z * z
	xy := x * y
	xz := x * z
	yz := y * z
	wx := w * x
	wy := w * y
	wz := w * z

	return Mat3{
		{1.0 - 2.0*(y2+z2), 2.0 * (xy - wz), 2.0 * (xz + wy)},
		{2.0 * (xy + wz), 1.0 - 2.0*(x2+z2), 2.0 * (yz - wx)},
		{2.0 * (xz - wy), 2.0 * (yz + wx), 1.0 - 2.0*(x2+y2)},
	}
}

func (q Quat) Mat4() Mat4 {
	m := q.Mat3()
	return m.Mat4()
}

func (q Quat) Transform3(v Vec3) Vec3 {
	m := q.Mat3()
	return m.Transform(v)
}

func (q Quat) Transform4(v Vec4) Vec4 {
	m := q.Mat4()
	return m.Transform(v)
}

func (q Quat) Axis() (v Vec3, r float64) {
	s := math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	v = Vec3{q.X / s, q.Y / s, q.Z / s}
	r = math.Acos(q.W) * 2
	return
}

func (q Quat) Lerp(t float64, p Quat) Quat {
	return q.Add(p.Sub(q).Scale(t))
}

type Spherical struct {
	R, T, P float64
}

func (s Spherical) Euclidean() Vec3 {
	sint := math.Sin(s.T)
	cost := math.Sin(s.T)
	sinp := math.Sin(s.P)
	cosp := math.Sin(s.P)
	r := s.R

	return Vec3{
		r * sint * cosp,
		r * sint * sinp,
		r * cost,
	}
}

func Slerp(v0, v1 Quat, t float64) Quat {
	v0 = v0.Normalize()
	v1 = v1.Normalize()

	const threshold = 0.9995
	dot := v0.Dot(v1)
	if dot > threshold {
		return v0.Lerp(t, v1).Normalize()
	}

	if dot < 0 {
		v1 = v1.Neg()
		dot = -dot
	}

	dot = Clamp(dot, -1, 1)
	theta0 := math.Acos(dot)
	theta := theta0 * t

	v2 := v1.Sub(v0.Scale(dot))
	v2 = v2.Normalize()

	v3 := v0.Scale(math.Cos(theta))
	v4 := v2.Scale(math.Sin(theta))
	return v3.Add(v4)
}

func Lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func Unlerp(t, a, b float64) float64 {
	return (t - a) / (b - a)
}

func LinearRemap(x, a, b, c, d float64) float64 {
	return Lerp(Unlerp(x, a, b), c, d)
}

func Smoothstep(a, b, x float64) float64 {
	t := Clamp((x-a)/(b-a), 0, 1)
	return t * t * (3 - 2*t)
}

func CubicBezier1D(t, p0, p1, p2, p3 float64) float64 {
	it := 1 - t
	return it*it*it*p0 + 3*it*it*t*p1 + 3*it*t*t*p2 + t*t*t*p3
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

func Saturate(x float64) float64 {
	return math.Max(0, math.Min(1, x))
}

func SignStrict(x float64) float64 {
	if x >= 0 {
		return 1
	}
	return -1
}

func Sign(x float64) float64 {
	if x < 0 {
		return -1
	}
	if x == 0 {
		return 0
	}
	return 1
}

const (
	Radians = math.Pi / 180
	Degrees = 180 / math.Pi
)

type Circle struct {
	X, Y, R float64
}

func (c Circle) InPoint(x, y float64) bool {
	return (x-c.X)*(x-c.X)+(y-c.Y)*(y-c.Y) <= c.R
}

func (c Circle) InRect(r Rectangle) bool {
	dx := c.X - math.Max(r.Min.X, math.Min(c.X, r.Max.X))
	dy := c.Y - math.Max(r.Min.Y, math.Min(c.Y, r.Max.Y))
	return dx*dx+dy*dy <= c.R
}

type Rectangle struct {
	Min, Max Vec2
}

func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Vec2{x0, y0}, Vec2{x1, y1}}
}

func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

func (r Rectangle) Overlaps(s Rectangle) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}

	if r.Empty() {
		return Rectangle{}
	}
	return r
}

func (r Rectangle) Union(s Rectangle) Rectangle {
	if r.Empty() {
		return s
	}
	if s.Empty() {
		return r
	}
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

func (r Rectangle) Int() image.Rectangle {
	return image.Rect(int(r.Min.X), int(r.Min.Y), int(r.Max.X), int(r.Max.Y))
}

func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

func (r Rectangle) Scale(s Vec2) Rectangle {
	r.Min.X *= s.X
	r.Max.X *= s.X
	r.Min.Y *= s.Y
	r.Max.Y *= s.Y
	return r
}

func (r Rectangle) Add(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X + p.X, r.Min.Y + p.Y},
		Vec2{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

func (r Rectangle) Sub(p Vec2) Rectangle {
	return Rectangle{
		Vec2{r.Min.X - p.X, r.Min.Y - p.Y},
		Vec2{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

func (r Rectangle) Size() Vec2 {
	return Vec2{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

func (r Rectangle) Floor() Rectangle {
	return Rectangle{
		r.Min.Floor(),
		r.Max.Floor(),
	}
}

func (r Rectangle) Expand(x, y float64) Rectangle {
	r.Min.X -= x
	r.Min.Y -= y
	r.Max.X += x
	r.Max.Y += y
	return r
}

func (r Rectangle) Expand2(v Vec2) Rectangle {
	return r.Expand(v.X, v.Y)
}

func (r Rectangle) Inset(n float64) Rectangle {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Center() Vec2 {
	return Vec2{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rectangle) Diagonal() float64 {
	x := r.Max.X - r.Min.X
	y := r.Max.Y - r.Min.Y
	return math.Sqrt(x*x + y*y)
}

func (r Rectangle) PosSize() (x, y, w, h float64) {
	return r.Min.X, r.Min.Y, r.Dx(), r.Dy()
}

func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}

	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

func (r Rectangle) TL() Vec2 { return r.Min }
func (r Rectangle) TR() Vec2 { return Vec2{r.Max.X, r.Min.Y} }
func (r Rectangle) BL() Vec2 { return Vec2{r.Min.X, r.Max.Y} }
func (r Rectangle) BR() Vec2 { return r.Max }

func (r Rectangle) Inverted() bool {
	return r.Min.X > r.Max.X || r.Min.Y > r.Max.Y
}

func RoundPrec(v float64, prec int) float64 {
	if prec < 0 {
		return v
	}

	tab := [...]float64{
		1, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8, 1e-9, 1e-10,
	}
	step := 0.0
	if prec < len(tab) {
		step = tab[prec]
	} else {
		step = math.Pow(10, float64(-prec))
	}

	neg := v < 0
	v = math.Abs(v)
	rem := math.Mod(v, step)
	if rem <= step*0.5 {
		v -= rem
	} else {
		v += step - rem
	}

	if neg {
		v = -v
	}

	return v
}

func Sinc(x float64) float64 {
	x *= math.Pi
	if x < 0.01 && x > -0.01 {
		return 1 + x*x*((-1.0/6)+x*x*1.0/120)
	}
	return math.Sin(x) / x
}

func LinearController(curpos *float64, targetpos, acc, deacc, dt float64) {
	sign := 1.0
	p := 0.0
	cp := *curpos
	if cp == targetpos {
		return
	}
	if targetpos < cp {
		targetpos = -targetpos
		cp = -cp
		sign = -1
	}

	// first decelerate
	if cp < 0 {
		p = cp + deacc*dt
		if p > 0 {
			p = 0
			dt = dt - p/deacc
			if dt < 0 {
				dt = 0
			}
		} else {
			dt = 0
		}
		cp = p
	}

	// now accelerate
	p = cp + acc*dt
	if p > targetpos {
		p = targetpos
	}
	*curpos = p * sign
}

func Multiple(a, m float64) float64 {
	return math.Ceil(a/m) * m
}

func Wrap(x, s, e float64) float64 {
	if x < s {
		x += e
	}
	if x >= e {
		x -= e
	}
	return x
}
