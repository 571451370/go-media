package sdl

/*
#include "sdl.h"

Uint32 getPixel(SDL_Surface *surface, size_t x, size_t y) {
	SDL_LockSurface(surface);
    int bpp = surface->format->BytesPerPixel;
    Uint8 *p = (Uint8 *)surface->pixels + y * surface->pitch + x * bpp;
	Uint32 c;

    switch(bpp) {
    case 1:
        c = *p;
        break;

    case 2:
        c = *(Uint16 *)p;
        break;

    case 3:
        if(SDL_BYTEORDER == SDL_BIG_ENDIAN)
            c = p[0] << 16 | p[1] << 8 | p[2];
        else
            c = p[0] | p[1] << 8 | p[2] << 16;
        break;

    case 4:
        c = *(Uint32 *)p;
        break;

    default:
    	c = 0;
	}
	SDL_UnlockSurface(surface);
	return c;
}

void setPixel(SDL_Surface *surface, size_t x, size_t y, Uint32 pixel) {
	SDL_LockSurface(surface);
    int bpp = surface->format->BytesPerPixel;
    Uint8 *p = (Uint8 *)surface->pixels + y * surface->pitch + x * bpp;

    switch(bpp) {
    case 1:
        *p = pixel;
        break;

    case 2:
        *(Uint16 *)p = pixel;
        break;

    case 3:
        if(SDL_BYTEORDER == SDL_BIG_ENDIAN) {
            p[0] = (pixel >> 16) & 0xff;
            p[1] = (pixel >> 8) & 0xff;
            p[2] = pixel & 0xff;
        } else {
            p[0] = pixel & 0xff;
            p[1] = (pixel >> 8) & 0xff;
            p[2] = (pixel >> 16) & 0xff;
        }
        break;

    case 4:
        *(Uint32 *)p = pixel;
        break;
    }
	SDL_UnlockSurface(surface);
}
*/
import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

type Surface C.SDL_Surface

func (s *Surface) Flags() uint32             { return uint32(s.flags) }
func (s *Surface) Size() (width, height int) { return int(s.w), int(s.h) }

func CreateRGBSurface(flags uint32, width, height, depth int, rmask, gmask, bmask, amask uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurface(C.Uint32(flags), C.int(width), C.int(height), C.int(depth),
		C.Uint32(rmask), C.Uint32(gmask), C.Uint32(bmask), C.Uint32(amask)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func CreateRGBSurfaceFrom(pixels []byte, width, height, depth, pitch int, rmask, gmask, bmask, amask uint32) (*Surface, error) {
	s := (*Surface)(C.SDL_CreateRGBSurfaceFrom(unsafe.Pointer(&pixels[0]), C.int(width), C.int(height), C.int(depth), C.int(pitch),
		C.Uint32(rmask), C.Uint32(gmask), C.Uint32(bmask), C.Uint32(amask)))
	if s == nil {
		return nil, GetError()
	}
	return s, nil
}

func (s *Surface) SetClipRect(rect Rect) {
	C.SDL_SetClipRect(s, (*C.SDL_Rect)(unsafe.Pointer(&rect)))
}

func (s *Surface) ClipRect() Rect {
	var rect Rect
	C.SDL_GetClipRect(s, (*C.SDL_Rect)(unsafe.Pointer(&rect)))
	return rect
}

func (s *Surface) Lock() error {
	return ek(C.SDL_LockSurface(s))
}

func (s *Surface) Unlock() {
	C.SDL_UnlockSurface(s)
}

func (s *Surface) Free() {
	C.SDL_FreeSurface(s)
}

func (s *Surface) ColorModel() color.Model {
	return color.NRGBAModel
}

func (s *Surface) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(s.w), int(s.h))
}

func (s *Surface) At(x, y int) color.Color {
	var cr, cg, cb, ca C.Uint8
	pixel := C.getPixel(s, C.size_t(x), C.size_t(y))
	C.SDL_GetRGBA(pixel, s.format, &cr, &cg, &cb, &ca)
	return color.NRGBA{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}
}

func (s *Surface) Set(x, y int, c color.Color) {
	p := color.NRGBAModel.Convert(c).(color.NRGBA)
	C.setPixel(s, C.size_t(x), C.size_t(y), C.SDL_MapRGBA(s.format, C.Uint8(p.R), C.Uint8(p.G), C.Uint8(p.B), C.Uint8(p.A)))
}
