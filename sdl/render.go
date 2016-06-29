package sdl

/*
#include "sdl.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type (
	Renderer        C.SDL_Renderer
	RendererFlags   C.SDL_RendererFlags
	RendererFlip    C.SDL_RendererFlip
	Texture         C.SDL_Texture
	TextureAccess   C.SDL_TextureAccess
	TextureModulate C.SDL_TextureModulate
)

type RendererInfo struct {
	Name             string
	Flags            uint32
	TextureFormats   []uint32
	MaxTextureWidth  int
	MaxTextureHeight int
}

const (
	RENDERER_SOFTWARE      RendererFlags = C.SDL_RENDERER_SOFTWARE
	RENDERER_ACCELERATED   RendererFlags = C.SDL_RENDERER_ACCELERATED
	RENDERER_PRESENTVSYNC  RendererFlags = C.SDL_RENDERER_PRESENTVSYNC
	RENDERER_TARGETTEXTURE RendererFlags = C.SDL_RENDERER_TARGETTEXTURE
)

const (
	TEXTUREACCESS_STATIC    TextureAccess = C.SDL_TEXTUREACCESS_STATIC
	TEXTUREACCESS_STREAMING TextureAccess = C.SDL_TEXTUREACCESS_STREAMING
	TEXTUREACCESS_TARGET    TextureAccess = C.SDL_TEXTUREACCESS_TARGET
)

const (
	TEXTUREMODULATE_NONE  TextureModulate = C.SDL_TEXTUREMODULATE_NONE
	TEXTUREMODULATE_COLOR TextureModulate = C.SDL_TEXTUREMODULATE_COLOR
	TEXTUREMODULATE_ALPHA TextureModulate = C.SDL_TEXTUREMODULATE_ALPHA
)

const (
	FLIP_NONE       RendererFlip = C.SDL_FLIP_NONE
	FLIP_HORIZONTAL RendererFlip = C.SDL_FLIP_HORIZONTAL
	FLIP_VERTICAL   RendererFlip = C.SDL_FLIP_VERTICAL
)

func makeRendererInfo(info *C.SDL_RendererInfo) RendererInfo {
	r := RendererInfo{
		C.GoString(info.name),
		uint32(info.flags),
		nil,
		int(info.max_texture_width),
		int(info.max_texture_height),
	}
	r.TextureFormats = make([]uint32, info.num_texture_formats)
	for i := range r.TextureFormats {
		r.TextureFormats[i] = uint32(info.texture_formats[i])
	}
	return r
}

func GetNumRenderDrivers() int {
	return int(C.SDL_GetNumRenderDrivers())
}

func GetRenderDriverInfo(index int) (RendererInfo, error) {
	var info C.SDL_RendererInfo
	rc := C.SDL_GetRenderDriverInfo(C.int(index), &info)
	if rc < 0 {
		return RendererInfo{}, GetError()
	}
	return makeRendererInfo(&info), nil
}

func CreateWindowAndRenderer(width, height int, windowFlags WindowFlags) (*Window, *Renderer, error) {
	var window *C.SDL_Window
	var renderer *C.SDL_Renderer
	err := ek(C.SDL_CreateWindowAndRenderer(C.int(width), C.int(height), C.Uint32(windowFlags), &window, &renderer))
	return (*Window)(window), (*Renderer)(renderer), err
}

func CreateRenderer(window *Window, index int, rendererFlags RendererFlags) (*Renderer, error) {
	renderer := C.SDL_CreateRenderer(window, C.int(index), C.Uint32(rendererFlags))
	if renderer == nil {
		return nil, GetError()
	}
	return (*Renderer)(renderer), nil
}

func CreateSoftwareRenderer(surface *Surface) (*Renderer, error) {
	re := C.SDL_CreateSoftwareRenderer(surface)
	if re == nil {
		return nil, GetError()
	}
	return (*Renderer)(re), nil
}

func (w *Window) Renderer() *Renderer {
	return (*Renderer)(C.SDL_GetRenderer(w))
}

func (re *Renderer) Info() (RendererInfo, error) {
	var info C.SDL_RendererInfo
	rc := C.SDL_GetRendererInfo(re, &info)
	return makeRendererInfo(&info), ek(rc)
}

func (re *Renderer) OutputSize() (width, height int, err error) {
	var cw, ch C.int
	rc := C.SDL_GetRendererOutputSize(re, &cw, &ch)
	return int(cw), int(ch), ek(rc)
}

func (re *Renderer) CreateTexture(format uint32, access TextureAccess, width, height int) (*Texture, error) {
	t := C.SDL_CreateTexture(re, C.Uint32(format), C.int(access), C.int(width), C.int(height))
	if t == nil {
		return nil, GetError()
	}
	return (*Texture)(t), nil
}

func (re *Renderer) CreateTextureFromSurface(surface *Surface) (*Texture, error) {
	t := C.SDL_CreateTextureFromSurface(re, surface)
	if t == nil {
		return nil, GetError()
	}
	return (*Texture)(t), nil
}

func (t *Texture) Query() (format uint32, access TextureAccess, width, height int, err error) {
	var cformat C.Uint32
	var caccess, cw, ch C.int
	rc := C.SDL_QueryTexture(t, &cformat, &caccess, &cw, &ch)
	return uint32(cformat), TextureAccess(caccess), int(cw), int(ch), ek(rc)
}

func (t *Texture) SetColorMod(c Color) error {
	return ek(C.SDL_SetTextureColorMod(t, C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B)))
}

func (t *Texture) ColorMod() (Color, error) {
	var cr, cg, cb C.Uint8
	rc := C.SDL_GetTextureColorMod(t, &cr, &cg, &cb)
	return Color{uint8(cr), uint8(cg), uint8(cb), 255}, ek(rc)
}

func (t *Texture) SetAlphaMod(alpha uint8) error {
	return ek(C.SDL_SetTextureAlphaMod(t, C.Uint8(alpha)))
}

func (t *Texture) AlphaMod() (uint8, error) {
	var calpha C.Uint8
	rc := C.SDL_GetTextureAlphaMod(t, &calpha)
	return uint8(calpha), ek(rc)
}

func (t *Texture) Lock(rect *Rect) ([]byte, error) {
	_, _, _, height, err := t.Query()
	if err != nil {
		return nil, err
	}

	var pixels unsafe.Pointer
	var pitch C.int
	rc := C.SDL_LockTexture(t, (*C.SDL_Rect)(unsafe.Pointer(rect)), &pixels, &pitch)
	if rc < 0 {
		return nil, GetError()
	}

	var buf []byte
	sl := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sl.Data = uintptr(pixels)
	if rect == nil {
		sl.Len = int(pitch) * int(height)
	} else {
		sl.Len = int(pitch) * int(rect.H)
	}
	sl.Cap = sl.Len
	return buf, nil
}

func (t *Texture) Unlock() {
	C.SDL_UnlockTexture(t)
}

func (re *Renderer) SetDrawColor(c Color) error {
	return ek(C.SDL_SetRenderDrawColor(re, C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)))
}

func (re *Renderer) DrawColor() (Color, error) {
	var cr, cg, cb, ca C.Uint8
	rc := C.SDL_GetRenderDrawColor(re, &cr, &cg, &cb, &ca)
	return Color{uint8(cr), uint8(cg), uint8(cb), uint8(ca)}, ek(rc)
}

func (re *Renderer) DrawPoint(x, y int) error {
	return ek(C.SDL_RenderDrawPoint(re, C.int(x), C.int(y)))
}

func (re *Renderer) DrawPoints(pts []Point) error {
	return ek(C.SDL_RenderDrawPoints(re, (*C.SDL_Point)(unsafe.Pointer(&pts[0])), C.int(len(pts))))
}

func (re *Renderer) DrawLine(x1, y1, x2, y2 int) error {
	return ek(C.SDL_RenderDrawLine(re, C.int(x1), C.int(y1), C.int(x2), C.int(y2)))
}

func (re *Renderer) DrawLines(points []Point) error {
	return ek(C.SDL_RenderDrawLines(re, (*C.SDL_Point)(unsafe.Pointer(&points[0])), C.int(len(points))))
}

func (re *Renderer) DrawRect(rect *Rect) error {
	return ek(C.SDL_RenderDrawRect(re, (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func (re *Renderer) DrawRects(rects []Rect) error {
	return ek(C.SDL_RenderDrawRects(re, (*C.SDL_Rect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) FillRect(rect *Rect) error {
	return ek(C.SDL_RenderFillRect(re, (*C.SDL_Rect)(unsafe.Pointer(rect))))
}

func (re *Renderer) FillRects(rects []Rect) error {
	return ek(C.SDL_RenderFillRects(re, (*C.SDL_Rect)(unsafe.Pointer(&rects[0])), C.int(len(rects))))
}

func (re *Renderer) CopyEx(texture *Texture, src, dst *Rect, angle float64, center *Point, flip RendererFlip) error {
	return ek(C.SDL_RenderCopyEx(re, texture, (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_Rect)(unsafe.Pointer(dst)),
		C.double(angle), (*C.SDL_Point)(unsafe.Pointer(center)), C.SDL_RendererFlip(flip)))
}

func (re *Renderer) Copy(texture *Texture, src, dst *Rect) error {
	return ek(C.SDL_RenderCopy(re, texture, (*C.SDL_Rect)(unsafe.Pointer(src)), (*C.SDL_Rect)(unsafe.Pointer(dst))))
}

func (re *Renderer) ReadPixels(rect *Rect, format uint32, pixels []byte, pitch int) error {
	return ek(C.SDL_RenderReadPixels(re, (*C.SDL_Rect)(unsafe.Pointer(rect)), C.Uint32(format), unsafe.Pointer(&pixels[0]), C.int(pitch)))
}

func (re *Renderer) SetTarget(texture *Texture) error {
	return ek(C.SDL_SetRenderTarget(re, (*C.SDL_Texture)(texture)))
}

func (re *Renderer) Present() {
	C.SDL_RenderPresent(re)
}

func (t *Texture) Destroy() {
	C.SDL_DestroyTexture(t)
}

func (re *Renderer) Destroy() {
	C.SDL_DestroyRenderer(re)
}

func (re *Renderer) SetLogicalSize(width, height int) {
	C.SDL_RenderSetLogicalSize(re, C.int(width), C.int(height))
}

func (re *Renderer) LogicalSize() (width, height int) {
	var cw, ch C.int
	C.SDL_RenderGetLogicalSize(re, &cw, &ch)
	return int(cw), int(ch)
}

func (re *Renderer) Clear() error {
	return ek(C.SDL_RenderClear(re))
}

func (re *Renderer) SetBlendMode(blendMode BlendMode) error {
	return ek(C.SDL_SetRenderDrawBlendMode(re, C.SDL_BlendMode(blendMode)))
}

func (re *Renderer) BlendMode() (BlendMode, error) {
	var mode C.SDL_BlendMode
	rc := C.SDL_GetRenderDrawBlendMode(re, &mode)
	return BlendMode(mode), ek(rc)
}
