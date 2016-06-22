package sdlgfx

/*
#include "SDL2_gfxPrimitives.h"
#include "SDL2_framerate.h"
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/qeedquan/go-sdl/sdl"
)

func Pixel(re *sdl.Renderer, x, y int, r, g, b, a uint8) error {
	return ek(C.pixelRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Hline(re *sdl.Renderer, x1, x2, y int, r, g, b, a uint8) error {
	return ek(C.hlineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(x2), C.Sint16(y), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Vline(re *sdl.Renderer, x, y1, y2 int, r, g, b, a uint8) error {
	return ek(C.vlineRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y1), C.Sint16(y2), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Rectangle(re *sdl.Renderer, x1, y1, x2, y2 int, r, g, b, a uint8) error {
	return ek(C.rectangleRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func RoundedRectangle(re *sdl.Renderer, x1, y1, x2, y2, rad int, r, g, b, a uint8) error {
	return ek(C.roundedRectangleRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(rad), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Box(re *sdl.Renderer, x1, y1, x2, y2 int, r, g, b, a uint8) error {
	return ek(C.boxRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func RoundedBox(re *sdl.Renderer, x1, y1, x2, y2, rad int, r, g, b, a uint8) error {
	return ek(C.roundedBoxRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(rad), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Line(re *sdl.Renderer, x1, y1, x2, y2 int, r, g, b, a uint8) error {
	return ek(C.lineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func AALine(re *sdl.Renderer, x1, y1, x2, y2 int, r, g, b, a uint8) error {
	return ek(C.aalineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func ThickLine(re *sdl.Renderer, x1, y1, x2, y2, width int, r, g, b, a uint8) error {
	return ek(C.thickLineRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Uint8(width), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Circle(re *sdl.Renderer, x, y, rad int, r, g, b, a uint8) error {
	return ek(C.circleRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Arc(re *sdl.Renderer, x, y, rad, start, end, int, r, g, b, a uint8) error {
	return ek(C.arcRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func AACircle(re *sdl.Renderer, x, y, rad, int, r, g, b, a uint8) error {
	return ek(C.aacircleRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func FilledCircle(re *sdl.Renderer, x, y, rad, int, r, g, b, a uint8) error {
	return ek(C.filledCircleRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Ellipse(re *sdl.Renderer, x, y, rx, ry int, r, g, b, a uint8) error {
	return ek(C.ellipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func AAEllipse(re *sdl.Renderer, x, y, rx, ry int, r, g, b, a uint8) error {
	return ek(C.aaellipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func FilledEllipse(re *sdl.Renderer, x, y, rx, ry int, r, g, b, a uint8) error {
	return ek(C.filledEllipseRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rx), C.Sint16(ry), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Pie(re *sdl.Renderer, x, y, rad, start, end int, r, g, b, a uint8) error {
	return ek(C.pieRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func FilledPie(re *sdl.Renderer, x, y, rad, start, end int, r, g, b, a uint8) error {
	return ek(C.filledPieRGBA((*C.SDL_Renderer)(re), C.Sint16(x), C.Sint16(y), C.Sint16(rad), C.Sint16(start), C.Sint16(end), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Trigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, r, g, b, a uint8) error {
	return ek(C.trigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func AATrigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, r, g, b, a uint8) error {
	return ek(C.aatrigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func FilledTrigon(re *sdl.Renderer, x1, y1, x2, y2, x3, y3 int, r, g, b, a uint8) error {
	return ek(C.filledTrigonRGBA((*C.SDL_Renderer)(re), C.Sint16(x1), C.Sint16(y1), C.Sint16(x2), C.Sint16(y2), C.Sint16(x3), C.Sint16(y3), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func Polygon(re *sdl.Renderer, vx, vy []int16, r, g, b, a uint8) error {
	return ek(C.polygonRGBA((*C.SDL_Renderer)(re), (*C.Sint16)(unsafe.Pointer(&vx[0])), (*C.Sint16)(unsafe.Pointer(&vy[0])), C.int(len(vx)), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func AAPolygon(re *sdl.Renderer, vx, vy []int16, r, g, b, a uint8) error {
	return ek(C.aapolygonRGBA((*C.SDL_Renderer)(re), (*C.Sint16)(unsafe.Pointer(&vx[0])), (*C.Sint16)(unsafe.Pointer(&vy[0])), C.int(len(vx)), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func FilledPolygon(re *sdl.Renderer, vx, vy []int16, r, g, b, a uint8) error {
	return ek(C.filledPolygonRGBA((*C.SDL_Renderer)(re), (*C.Sint16)(unsafe.Pointer(&vx[0])), (*C.Sint16)(unsafe.Pointer(&vy[0])), C.int(len(vx)), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func TexturedPolygon(re *sdl.Renderer, vx, vy []int16, texture *sdl.Surface, dx, dy int) error {
	return ek(C.texturedPolygon((*C.SDL_Renderer)(re), (*C.Sint16)(unsafe.Pointer(&vx[0])), (*C.Sint16)(unsafe.Pointer(&vy[0])), C.int(len(vx)), (*C.SDL_Surface)(unsafe.Pointer(&texture)), C.int(dx), C.int(dy)))
}

func Bezier(re *sdl.Renderer, vx, vy []int16, s int, r, g, b, a uint8) error {
	return ek(C.bezierRGBA((*C.SDL_Renderer)(re), (*C.Sint16)(unsafe.Pointer(&vx[0])), (*C.Sint16)(unsafe.Pointer(&vy[0])), C.int(len(vx)), C.int(s), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func SetFont(font []byte, cw, ch uint32) {
	C.gfxPrimitivesSetFont(unsafe.Pointer(&font[0]), C.Uint32(cw), C.Uint32(ch))
}

func SetFontRotation(rotation uint32) {
	C.gfxPrimitivesSetFontRotation(C.Uint32(rotation))
}

func Character(re *sdl.Renderer, x, y int, c rune, r, g, b, a uint8) error {
	return ek(C.characterRGBA(re, C.Sint16(x), C.Sint16(y), C.char(c), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func String(re *sdl.Renderer, x, y int, s string, r, g, b, a uint8) error {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return ek(C.stringRGBA(re, C.Sint16(x), C.Sint16(y), cs, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a)))
}

func ek(rc C.int) error {
	if rc < 0 {
		return errors.New("invalid parameter")
	}
	return nil
}

type FPSManager C.FPSmanager

func (m *FPSManager) Init() {
	C.SDL_initFramerate((*C.FPSmanager)(unsafe.Pointer(m)))
}

func (m *FPSManager) SetRate(rate uint32) error {
	return ek(C.SDL_setFramerate((*C.FPSmanager)(unsafe.Pointer(m)), C.Uint32(rate)))
}

func (m *FPSManager) Rate() (int, error) {
	rc := C.SDL_getFramerate((*C.FPSmanager)(unsafe.Pointer(m)))
	return int(rc), ek(rc)
}

func (m *FPSManager) Count() (int, error) {
	rc := C.SDL_getFramecount((*C.FPSmanager)(unsafe.Pointer(m)))
	return int(rc), ek(rc)
}

func (m *FPSManager) Delay() uint32 {
	return uint32(C.SDL_framerateDelay((*C.FPSmanager)(unsafe.Pointer(m))))
}
