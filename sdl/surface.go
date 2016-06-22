package sdl

/*
#include "sdl.h"
*/
import "C"
import "unsafe"

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
