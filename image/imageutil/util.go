package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/qeedquan/go-media/image/pnm"
	_ "github.com/qeedquan/go-media/image/psd"
	"github.com/qeedquan/go-media/image/tga"
	"golang.org/x/image/bmp"
)

func LoadRGBAVFS(fs FS, name string) (*image.RGBA, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadRGBAReader(f)
}

func LoadRGBAFile(name string) (*image.RGBA, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m, err := LoadReader(f)
	if err != nil {
		return nil, &os.PathError{Op: "decode", Path: name, Err: err}
	}
	return m, nil
}

func LoadRGBAReader(rd io.Reader) (*image.RGBA, error) {
	m, _, err := image.Decode(rd)
	if err != nil {
		return nil, err
	}

	if p, _ := m.(*image.RGBA); p != nil {
		return p, nil
	}

	r := m.Bounds()
	p := image.NewRGBA(r)
	draw.Draw(p, p.Bounds(), m, r.Min, draw.Src)
	return p, nil
}

func WriteRGBAFile(img image.Image, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	case ".pbm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 1})
	case ".pgm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 2})
	case ".ppm":
		err = pnm.Encode(f, img, &pnm.Options{Format: 3})
	case ".gif":
		err = gif.Encode(f, img, &gif.Options{
			NumColors: 256,
		})
	case ".tga":
		err = tga.Encode(f, img)
	case ".bmp":
		err = bmp.Encode(f, img)
	case ".png":
		fallthrough
	default:
		err = png.Encode(f, img)
	}

	xerr := f.Close()
	if err == nil {
		err = xerr
	}
	return err
}

func ColorKey(m image.Image, c color.Color) *image.RGBA {
	p := image.NewRGBA(m.Bounds())
	b := p.Bounds()

	cr, cg, cb, _ := c.RGBA()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			col := m.At(x, y)
			r, g, b, _ := col.RGBA()
			if cr == r && cg == g && cb == b {
				p.Set(x, y, color.Transparent)
			} else {
				p.Set(x, y, col)
			}
		}
	}
	return p
}

func ParseColor(s string) (color.RGBA, error) {
	var r, g, b, a uint8
	n, _ := fmt.Sscanf(s, "rgb(%v,%v,%v)", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}
	}

	n, _ = fmt.Sscanf(s, "rgba(%v,%v,%v,%v)", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}
	}

	n, _ = fmt.Sscanf("#%02x%02x%02x%02x", &r, &g, &b, &a)
	if n == 4 {
		return color.RGBA{r, g, b, a}
	}

	n, _ = fmt.Sscanf("#%02x%02x%02x", &r, &g, &b)
	if n == 3 {
		return color.RGBA{r, g, b, 255}
	}

	n, _ = fmt.Sscanf("#%02x", &r)
	if n == 1 {
		return color.RGBA{r, r, r, 255}
	}

	return fmt.Errorf("failed to parse color %q, unknown format", s)
}