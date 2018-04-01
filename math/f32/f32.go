package f32

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

type Vec2 struct {
	X, Y float32
}

func (p Vec2) Add(q Vec2) Vec2 {
	return Vec2{p.X + q.X, p.Y + q.Y}
}

func (p Vec2) Sub(q Vec2) Vec2 {
	return Vec2{p.X - q.X, p.Y - q.Y}
}

func (p Vec2) AddScale(q Vec2, t float32) Vec2 {
	return p.Add(q.Scale(t))
}

func (p Vec2) SubScale(q Vec2, t float32) Vec2 {
	return p.Sub(q.Scale(t))
}

func (p Vec2) Neg() Vec2 {
	return Vec2{-p.X, -p.Y}
}

func (p Vec2) Abs() Vec2 {
	return Vec2{Abs(p.X), Abs(p.Y)}
}

func (p Vec2) Dot(q Vec2) float32 {
	return p.X*q.X + p.Y*q.Y
}

func (p Vec2) MinComp() float32 {
	return Min(p.X, p.Y)
}

func (p Vec2) MaxComp() float32 {
	return Max(p.X, p.Y)
}

func (p Vec2) Max(q Vec2) Vec2 {
	return Vec2{
		Max(p.X, q.X),
		Max(p.Y, q.Y),
	}
}

func (p Vec2) Min(q Vec2) Vec2 {
	return Vec2{
		Min(p.X, q.X),
		Min(p.Y, q.Y),
	}
}

func (p Vec2) Floor() Vec2 {
	return Vec2{
		Floor(p.X),
		Floor(p.Y),
	}
}

func (p Vec2) Ceil() Vec2 {
	return Vec2{
		Ceil(p.X),
		Ceil(p.Y),
	}
}

func (p Vec2) Len() float32 {
	return Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Vec2) LenSquared() float32 {
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

func (p Vec2) Scale(k float32) Vec2 {
	return Vec2{p.X * k, p.Y * k}
}

func (p Vec2) Shear(k float32) Vec2 {
	return Vec2{p.X + k*p.Y, p.Y + k*p.X}
}

func (p Vec2) Shearv(q Vec2) Vec2 {
	return Vec2{p.X + q.X*p.Y, p.Y + q.Y*p.X}
}

func (p Vec2) Lerp(t float32, q Vec2) Vec2 {
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

func (p Vec2) Distance(q Vec2) float32 {
	return p.Sub(q).Len()
}

func (p Vec2) DistanceSquared(q Vec2) float32 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec2) Polar() Polar {
	return Polar{p.Len(), Atan2(p.Y, p.X)}
}

func (p Vec2) MinScalar(k float32) Vec2 {
	return Vec2{
		Min(p.X, k),
		Min(p.Y, k),
	}
}

func (p Vec2) MaxScalar(k float32) Vec2 {
	return Vec2{
		Max(p.X, k),
		Max(p.Y, k),
	}
}

func (p Vec2) In(r Rectangle) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

func (p Vec2) Rotate(r float32) Vec2 {
	si, co := Sincos(r)
	return Vec2{
		p.X*co - p.Y*si,
		p.X*si + p.Y*co,
	}
}

func (p Vec2) Shrink(k float32) Vec2 {
	return Vec2{p.X / k, p.Y / k}
}

func (p Vec2) Shrink2(q Vec2) Vec2 {
	return Vec2{p.X / q.X, p.Y / q.Y}
}

func (p Vec2) YX() Vec2 {
	return Vec2{p.Y, p.X}
}

func (p Vec2) Clamp(s, e float32) Vec2 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	return p
}

func (p Vec2) Equals(q Vec2) bool {
	const eps = 1e-6
	return Abs(p.X-q.X) <= eps && Abs(p.Y-q.Y) <= eps
}

func (p Vec2) OnLine(a, b Vec2) bool {
	sx := Min(a.X, b.X)
	sy := Min(a.Y, b.Y)
	ex := Max(a.X, b.X)
	ey := Max(a.Y, b.Y)
	return sx <= p.X && p.X <= ex &&
		sy <= p.Y && p.Y <= ey
}

type Vec3 struct {
	X, Y, Z float32
}

func (p Vec3) Add(q Vec3) Vec3 {
	return Vec3{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

func (p Vec3) Sub(q Vec3) Vec3 {
	return Vec3{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

func (p Vec3) AddScale(q Vec3, t float32) Vec3 {
	return p.Add(q.Scale(t))
}

func (p Vec3) SubScale(q Vec3, t float32) Vec3 {
	return p.Sub(q.Scale(t))
}

func (p Vec3) Dot(q Vec3) float32 {
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

func (p Vec3) Refract(q Vec3, eta float32) Vec3 {
	x := p.Dot(q)
	k := 1 - eta*eta*(1-x*x)
	if k < 0 {
		return Vec3{}
	}
	a := q.Scale(eta)
	b := p.Scale(eta*x + Sqrt(k))
	return a.Sub(b)
}

func (p Vec3) Len() float32 {
	return Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p Vec3) LenSquared() float32 {
	return p.Dot(p)
}

func (p Vec3) Scale3(q Vec3) Vec3 {
	return Vec3{p.X * q.X, p.Y * q.Y, p.Z * q.Z}
}

func (p Vec3) Scale(k float32) Vec3 {
	return Vec3{p.X * k, p.Y * k, p.Z * k}
}

func (p Vec3) Shrink(k float32) Vec3 {
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

func (p Vec3) Distance(q Vec3) float32 {
	return p.Sub(q).Len()
}

func (p Vec3) DistanceSquared(q Vec3) float32 {
	r := p.Sub(q)
	return r.Dot(r)
}

func (p Vec3) Lerp(t float32, q Vec3) Vec3 {
	return Vec3{
		Lerp(t, p.X, q.X),
		Lerp(t, p.Y, q.Y),
		Lerp(t, p.Z, q.Z),
	}
}

func (p Vec3) Abs() Vec3 {
	return Vec3{
		Abs(p.X),
		Abs(p.Y),
		Abs(p.Z),
	}
}

func (p Vec3) MaxScalar(k float32) Vec3 {
	return Vec3{
		Max(p.X, k),
		Max(p.Y, k),
		Max(p.Z, k),
	}
}

func (p Vec3) MinScalar(k float32) Vec3 {
	return Vec3{
		Min(p.X, k),
		Min(p.Y, k),
		Min(p.Z, k),
	}
}

func (p Vec3) Max(q Vec3) Vec3 {
	return Vec3{
		Max(p.X, q.X),
		Max(p.Y, q.Y),
		Max(p.Z, q.Z),
	}
}

func (p Vec3) Min(q Vec3) Vec3 {
	return Vec3{
		Min(p.X, q.X),
		Min(p.Y, q.Y),
		Min(p.Z, q.Z),
	}
}

func (p Vec3) MinComp() float32 {
	return Min(p.X, Min(p.Y, p.Z))
}

func (p Vec3) MaxComp() float32 {
	return Max(p.X, Max(p.Y, p.Z))
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
		uint8(Clamp(p.X, 0, 255)),
		uint8(Clamp(p.Y, 0, 255)),
		uint8(Clamp(p.Z, 0, 255)),
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
		T: Acos(p.Z / l),
		P: Atan2(p.Y, p.X),
	}
}

func (p Vec3) Equals(q Vec3) bool {
	const eps = 1e-6
	return Abs(p.X-q.X) <= eps && Abs(p.Y-q.Y) <= eps &&
		Abs(p.Z-q.Z) <= eps
}

func (p Vec3) XY() Vec2 { return Vec2{p.X, p.Y} }
func (p Vec3) XZ() Vec2 { return Vec2{p.X, p.Z} }
func (p Vec3) YX() Vec2 { return Vec2{p.Y, p.X} }
func (p Vec3) YZ() Vec2 { return Vec2{p.Y, p.Z} }
func (p Vec3) ZX() Vec2 { return Vec2{p.Z, p.X} }
func (p Vec3) ZY() Vec2 { return Vec2{p.Z, p.Y} }

func (p Vec3) Clamp(s, e float32) Vec3 {
	p.X = Clamp(p.X, s, e)
	p.Y = Clamp(p.Y, s, e)
	p.Z = Clamp(p.Z, s, e)
	return p
}

type Vec4 struct {
	X, Y, Z, W float32
}

func (p Vec4) Add(q Vec4) Vec4 {
	return Vec4{p.X + q.X, p.Y + q.Y, p.Z + q.Z, p.W}
}

func (p Vec4) Sub(q Vec4) Vec4 {
	return Vec4{p.X - q.X, p.Y - q.Y, p.Z - q.Z, p.W}
}

func (p Vec4) AddScale(q Vec4, k float32) Vec4 {
	return p.Add(q.Scale(k))
}

func (p Vec4) SubScale(q Vec4, k float32) Vec4 {
	return p.Sub(q.Scale(k))
}

func (p Vec4) Scale(k float32) Vec4 {
	return Vec4{p.X * k, p.Y * k, p.Z * k, p.W}
}

func (p Vec4) Scale3(q Vec3) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W}
}

func (p Vec4) Scale4(q Vec4) Vec4 {
	return Vec4{p.X * q.X, p.Y * q.Y, p.Z * q.Z, p.W * q.W}
}

func (p Vec4) Shrink(k float32) Vec4 {
	return Vec4{p.X / k, p.Y / k, p.Z / k, p.W}
}

func (p Vec4) Shrink4(q Vec4) Vec4 {
	return Vec4{p.X / q.X, p.Y / q.Y, p.Z / q.Z, p.W}
}

func (p Vec4) Dot(q Vec4) float32 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z + p.W*q.W
}

func (p Vec4) Len() float32 {
	return Sqrt(p.Dot(p))
}

func (p Vec4) Dot3(q Vec4) float32 {
	return p.X*q.X + p.Y*q.Y + p.Z*q.Z
}

func (p Vec4) Len3() float32 {
	return Sqrt(p.Dot3(p))
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
	c := color.RGBA{
		uint8(Clamp(p.X, 0, 255)),
		uint8(Clamp(p.Y, 0, 255)),
		uint8(Clamp(p.Z, 0, 255)),
		uint8(Clamp(p.W, 0, 255)),
	}
	return c
}

func (p Vec4) RGBA() (r, g, b, a uint32) {
	c := p.ToRGBA()
	return c.RGBA()
}

type Mat3 [3][3]float32

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

func (m *Mat3) Det() float32 {
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

func (m *Mat3) Trace() float32 {
	return m[0][0] + m[1][1] + m[2][2]
}

func (m Mat3) String() string {
	return fmt.Sprintf(`Mat3[% 0.3f, % 0.3f, % 0.3f,
		     % 0.3f, % 0.3f, % 0.3f,
			      % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2],
		m[1][0], m[1][1], m[1][2],
		m[2][0], m[2][1], m[2][2])
}

type Mat4 [4][4]float32

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

func (m *Mat4) Trace() float32 {
	return m[0][0] + m[1][1] + m[2][2] + m[3][3]
}

func (m *Mat4) Translate(tx, ty, tz float32) *Mat4 {
	*m = Mat4{
		{1, 0, 0, tx},
		{0, 1, 0, ty},
		{0, 0, 1, tz},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Scale(sx, sy, sz float32) *Mat4 {
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

func (m *Mat4) Viewport(x, y, w, h float32) {
	*m = Mat4{
		{w / 2, 0, 0, w/2 + x},
		{0, h / 2, 0, h/2 + y},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func (m *Mat4) RotX(r float32) *Mat4 {
	si, co := Sincos(r)
	*m = Mat4{
		{1, 0, 0, 0},
		{0, co, -si, 0},
		{0, si, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotY(r float32) *Mat4 {
	si, co := Sincos(r)
	*m = Mat4{
		{co, 0, si, 0},
		{0, 1, 0, 0},
		{-si, 0, co, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) RotZ(r float32) *Mat4 {
	si, co := Sincos(r)
	*m = Mat4{
		{co, -si, 0, 0},
		{si, co, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
	return m
}

func (m *Mat4) Perspective(fovy, aspect, near, far float32) *Mat4 {
	f := Tan(fovy / 2)
	z := near - far
	*m = Mat4{
		{1 / (f * aspect), 0, 0, 0},
		{0, 1 / f, 0, 0},
		{0, 0, (-near - far) / z, 2 * far * near / z},
		{0, 0, 1, 0},
	}
	return m
}

func (m *Mat4) Ortho(l, r, b, t, n, f float32) *Mat4 {
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
	return fmt.Sprintf(`Mat4[% 0.3f, % 0.3f, % 0.3f, % 0.3f,
		     % 0.3f, % 0.3f, % 0.3f, % 0.3f,
			      % 0.3f, % 0.3f, % 0.3f, % 0.3f,
				       % 0.3f, % 0.3f, % 0.3f, % 0.3f]`,
		m[0][0], m[0][1], m[0][2], m[0][3],
		m[1][0], m[1][1], m[1][2], m[1][3],
		m[2][0], m[2][1], m[2][2], m[2][3],
		m[3][0], m[3][1], m[3][2], m[3][3])
}

type Polar struct {
	R, P float32
}

func (p *Polar) Cartesian() Vec2 {
	s, c := Sincos(p.P)
	return Vec2{p.R * c, p.R * s}
}

type Quat struct {
	X, Y, Z, W float32
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

func (q Quat) Dot(p Quat) float32 {
	return q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
}

func (q Quat) Scale(k float32) Quat {
	return Quat{
		q.X * k,
		q.Y * k,
		q.Z * k,
		q.W * k,
	}
}

func (q Quat) Len() float32 {
	return Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
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

func (q Quat) FromAxis(v Vec3, r float32) Quat {
	r *= 0.5
	vn := v.Normalize()
	si, co := Sincos(r)
	return Quat{
		vn.X * si,
		vn.Y * si,
		vn.Z * si,
		co,
	}
}

func (q Quat) FromEuler(pitch, yaw, roll float32) Quat {
	p := pitch / 2
	y := yaw / 2
	r := roll / 2

	sinp := Sin(p)
	siny := Sin(y)
	sinr := Sin(r)
	cosp := Cos(p)
	cosy := Cos(y)
	cosr := Cos(r)

	return Quat{
		sinr*cosp*cosy - cosr*sinp*siny,
		cosr*sinp*cosy + sinr*cosp*siny,
		cosr*cosp*siny - sinr*sinp*cosy,
		cosr*cosp*cosy + sinr*sinp*siny,
	}.Normalize()
}

func (q Quat) Matrix() Mat4 {
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

	return Mat4{
		{1.0 - 2.0*(y2+z2), 2.0 * (xy - wz), 2.0 * (xz + wy), 0.0},
		{2.0 * (xy + wz), 1.0 - 2.0*(x2+z2), 2.0 * (yz - wx), 0.0},
		{2.0 * (xz - wy), 2.0 * (yz + wx), 1.0 - 2.0*(x2+y2), 0.0},
		{0.0, 0.0, 0.0, 1.0},
	}
}

func (q Quat) Axis() (v Vec3, r float32) {
	s := Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	v = Vec3{q.X / s, q.Y / s, q.Z / s}
	r = Acos(q.W) * 2
	return
}

func (q Quat) Lerp(t float32, p Quat) Quat {
	return q.Add(p.Sub(q).Scale(t))
}

type Spherical struct {
	R, T, P float32
}

func (s Spherical) Euclidean() Vec3 {
	sint := Sin(s.T)
	cost := Sin(s.T)
	sinp := Sin(s.P)
	cosp := Sin(s.P)
	r := s.R

	return Vec3{
		r * sint * cosp,
		r * sint * sinp,
		r * cost,
	}
}

func Slerp(v0, v1 Quat, t float32) Quat {
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
	theta0 := Acos(dot)
	theta := theta0 * t

	v2 := v1.Sub(v0.Scale(dot))
	v2 = v2.Normalize()

	v3 := v0.Scale(Cos(theta))
	v4 := v2.Scale(Sin(theta))
	return v3.Add(v4)
}

func Lerp(t, a, b float32) float32 {
	return a + t*(b-a)
}

func Unlerp(t, a, b float32) float32 {
	return (t - a) / (b - a)
}

func LinearRemap(x, a, b, c, d float32) float32 {
	return Lerp(c, d, Unlerp(x, a, b))
}

func Smoothstep(a, b, x float32) float32 {
	t := Clamp((x-a)/(b-a), 0, 1)
	return t * t * (3 - 2*t)
}

func CubicBezier1D(t, p0, p1, p2, p3 float32) float32 {
	it := 1 - t
	return it*it*it*p0 + 3*it*it*t*p1 + 3*it*t*t*p2 + t*t*t*p3
}

func Clamp(x, s, e float32) float32 {
	if x < s {
		x = s
	}
	if x > e {
		x = e
	}
	return x
}

func Saturate(x float32) float32 {
	return Max(0, Min(1, x))
}

func Sign(x float32) float32 {
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
	X, Y, R float32
}

func (c Circle) InPoint(x, y float32) bool {
	return (x-c.X)*(x-c.X)+(y-c.Y)*(y-c.Y) <= c.R
}

func (c Circle) InRect(r Rectangle) bool {
	dx := c.X - Max(r.Min.X, Min(c.X, r.Max.X))
	dy := c.Y - Max(r.Min.Y, Min(c.Y, r.Max.Y))
	return dx*dx+dy*dy <= c.R
}

type Rectangle struct {
	Min, Max Vec2
}

func Rect(x0, y0, x1, y1 float32) Rectangle {
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

func (r Rectangle) Expand(x, y float32) Rectangle {
	r.Min.X -= x
	r.Min.Y -= y
	r.Max.X += x
	r.Max.Y += y
	return r
}

func (r Rectangle) Expand2(v Vec2) Rectangle {
	return r.Expand(v.X, v.Y)
}

func (r Rectangle) Inset(n float32) Rectangle {
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

func (r Rectangle) Dx() float32 {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Dy() float32 {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) Center() Vec2 {
	return Vec2{
		(r.Min.X + r.Max.X) / 2,
		(r.Min.Y + r.Max.Y) / 2,
	}
}

func (r Rectangle) Diagonal() float32 {
	x := r.Max.X - r.Min.X
	y := r.Max.Y - r.Min.Y
	return Sqrt(x*x + y*y)
}

func (r Rectangle) XYWH() (x, y, w, h float32) {
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

func RoundPrec(v float32, prec int) float32 {
	if prec < 0 {
		return v
	}

	tab := [...]float32{
		1, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8, 1e-9, 1e-10,
	}
	step := float32(0.0)
	if prec < len(tab) {
		step = tab[prec]
	} else {
		step = Pow(10, float32(-prec))
	}

	neg := v < 0
	v = Abs(v)
	rem := Mod(v, step)
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

func Sinc(x float32) float32 {
	x *= math.Pi
	if x < 0.01 && x > -0.01 {
		return 1 + x*x*((-1.0/6)+x*x*1.0/120)
	}
	return Sin(x) / x
}

func LinearController(curpos *float32, targetpos, acc, deacc, dt float32) {
	sign := float32(1.0)
	p := float32(0.0)
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

func Multiple(a, m float32) float32 {
	return Ceil(a/m) * m
}

func Abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func Min(a, b float32) float32 {
	return float32(math.Min(float64(a), float64(b)))
}

func Max(a, b float32) float32 {
	return float32(math.Max(float64(a), float64(b)))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Tan(x float32) float32 {
	return float32(math.Tan(float64(x)))
}

func Floor(x float32) float32 {
	return float32(math.Floor(float64(x)))
}

func Ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

func Atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Sincos(x float32) (si, co float32) {
	s, c := math.Sincos(float64(x))
	return float32(s), float32(c)
}

func Acos(x float32) float32 {
	return float32(math.Acos(float64(x)))
}

func Pow(x, y float32) float32 {
	return float32(math.Pow(float64(x), float64(y)))
}

func Mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}
