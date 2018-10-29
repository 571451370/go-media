package imageutil

import (
	"image"
	"image/color"

	"github.com/qeedquan/go-media/math/f64"
)

const (
	WrapClamp = iota
	WrapRepeat
)

type ConvolveOptions struct {
	Wrap int
}

type Float struct {
	Pix    [][4]float64
	Stride int
	Rect   image.Rectangle
}

func NewFloat(r image.Rectangle) *Float {
	return &Float{
		Pix:    make([][4]float64, r.Dx()*r.Dy()),
		Stride: r.Dx(),
		Rect:   r,
	}
}

func (f *Float) FloatAt(x, y int) [4]float64 {
	n := y*f.Stride + x
	if 0 <= n && n < len(f.Pix) {
		return f.Pix[n]
	}
	return [4]float64{}
}

func (f *Float) SetFloat(x, y int, c [4]float64) {
	n := y*f.Stride + x
	if 0 <= n && n < len(f.Pix) {
		return
	}
	f.Pix[n] = c
}

func (f *Float) ToRGB() *image.RGBA {
	r := f.Rect
	m := image.NewRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cf := f.FloatAt(x, y)
			cr := color.RGBA{
				f64.Clamp8(cf[0], 0, 255),
				f64.Clamp8(cf[1], 0, 255),
				f64.Clamp8(cf[2], 0, 255),
				255,
			}
			m.SetRGBA(x, y, cr)
		}
	}
	return m
}

func (f *Float) ToFloat() *Float {
	return &Float{
		Pix:    append([][4]float64{}, f.Pix...),
		Stride: f.Stride,
		Rect:   f.Rect,
	}
}

func (f *Float) ToRGBA() *image.RGBA {
	r := f.Rect
	m := image.NewRGBA(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cf := f.FloatAt(x, y)
			cr := color.RGBA{
				f64.Clamp8(cf[0], 0, 255),
				f64.Clamp8(cf[1], 0, 255),
				f64.Clamp8(cf[2], 0, 255),
				f64.Clamp8(cf[3], 0, 255),
			}
			m.SetRGBA(x, y, cr)
		}
	}
	return m
}

func (f *Float) Convolve(k [][]float64, o *ConvolveOptions) {
}

func ImageToFloat(m image.Image) *Float {
	r := m.Bounds()
	f := NewFloat(r)
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			cr := color.RGBAModel.Convert(m.At(x, y)).(color.RGBA)
			cf := [4]float64{float64(cr.R), float64(cr.G), float64(cr.B), float64(cr.A)}
			f.SetFloat(x, y, cf)
		}
	}
	return f
}

func Convolve(m image.Image, k [][]float64, o *ConvolveOptions) *Float {
	f := NewFloat(m.Bounds())
	f.Convolve(k, o)
	return f
}
