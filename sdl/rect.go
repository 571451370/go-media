package sdl

/*
#include <SDL.h>
*/
import "C"

import "unsafe"

type Point struct {
	X, Y int32
}

type Rect struct {
	X, Y, W, H int32
}

func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

func (p Point) In(r Rect) bool {
	return C.SDL_PointInRect((*C.SDL_Point)(unsafe.Pointer(&p)), (*C.SDL_Rect)(unsafe.Pointer(&r))) != 0
}

func (r Rect) Empty() bool {
	return C.SDL_RectEmpty((*C.SDL_Rect)(unsafe.Pointer(&r))) != 0
}

func (r Rect) Equal(p Rect) bool {
	return C.SDL_RectEquals((*C.SDL_Rect)(unsafe.Pointer(&r)), (*C.SDL_Rect)(unsafe.Pointer(&p))) != 0
}

func (r Rect) Enclose(p []Point) Rect {
	var res C.SDL_Rect
	C.SDL_EnclosePoints((*C.SDL_Point)(unsafe.Pointer(&p[0])), C.int(len(p)), (*C.SDL_Rect)(unsafe.Pointer(&r)), &res)
	return Rect{int32(res.x), int32(res.y), int32(res.w), int32(res.h)}
}
