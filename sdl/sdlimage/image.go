package sdlimage

/*
#include <SDL.h>
*/
import "C"

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"reflect"
	"unsafe"

	"github.com/qeedquan/go-sdl/sdl"
	_ "github.com/qeedquan/go-sdl/sdl/sdlimage/tga"
)

func LoadTextureFile(re *sdl.Renderer, name string) (*sdl.Texture, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return LoadTextureReader(re, f)
}

func LoadTextureReader(re *sdl.Renderer, r io.Reader) (*sdl.Texture, error) {
	m, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return LoadTextureImage(re, m)
}

func LoadTextureImage(re *sdl.Renderer, m image.Image) (*sdl.Texture, error) {
	r := m.Bounds()
	w, h := r.Dx(), r.Dy()

	var p []byte
	b := C.malloc(C.size_t(w * h * 4))
	l := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	l.Data = uintptr(b)
	l.Len = w * h * 4
	l.Cap = l.Len

	n := &image.NRGBA{p, w * 4, r}
	draw.Draw(n, n.Bounds(), m, image.ZP, draw.Src)

	s, err := sdl.CreateRGBSurfaceFrom(n.Pix[:], w, h, 32, w*4, 0xff, 0xff00, 0xff0000, 0xff000000)
	if err != nil {
		C.free(b)
		return nil, err
	}
	defer s.Free()

	return re.CreateTextureFromSurface(s)
}
