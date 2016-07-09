package imageutil

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "github.com/qeedquan/go-media/image/psd"
	_ "github.com/qeedquan/go-media/image/tga"
	_ "golang.org/x/image/bmp"
)

func LoadFile(name string) (*image.RGBA, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	if err != nil {
		return nil, &os.PathError{Op: "decode", Path: name, Err: err}
	}

	if p, _ := m.(*image.RGBA); p != nil {
		return p, nil
	}

	r := m.Bounds()
	p := image.NewRGBA(r)
	draw.Draw(p, p.Bounds(), m, r.Min, draw.Src)
	return p, nil
}

func ColorKey(m image.Image, c color.Color) image.Image {
	p := image.NewRGBA(m.Bounds())
	b := p.Bounds()

	cr, cg, cb, _ := c.RGBA()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			col := m.At(x, y)
			r, g, b, _ := col.RGBA()
			if cr == r && cg == g && cb == b {
				p.Set(x, y, color.RGBA{})
			} else {
				p.Set(x, y, col)
			}
		}
	}
	return p
}
