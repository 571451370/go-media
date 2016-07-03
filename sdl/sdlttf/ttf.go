package sdlttf

/*
#include "SDL_ttf.h"
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/qeedquan/go-media/sdl"
)

type Style int
type Hinting int

const (
	STYLE_NORMAL        Style = C.TTF_STYLE_NORMAL
	STYLE_BOLD          Style = C.TTF_STYLE_BOLD
	STYLE_ITALIC        Style = C.TTF_STYLE_ITALIC
	STYLE_UNDERLINE     Style = C.TTF_STYLE_UNDERLINE
	STYLE_STRIKETHROUGH Style = C.TTF_STYLE_STRIKETHROUGH
)

const (
	HINTING_NORMAL Hinting = C.TTF_HINTING_NORMAL
	HINTING_LIGHT  Hinting = C.TTF_HINTING_LIGHT
	HINTING_MONO   Hinting = C.TTF_HINTING_MONO
	HINTING_NONE   Hinting = C.TTF_HINTING_NONE
)

type Font C.TTF_Font

func GetError() error {
	return errors.New(C.GoString(C.TTF_GetError()))
}

func Init() error {
	err := C.TTF_Init()
	if err < 0 {
		return GetError()
	}
	return nil
}

func Quit() {
	C.TTF_Quit()
}

func WasInit() bool {
	isInit := C.TTF_WasInit()
	if isInit == 0 {
		return false
	}
	return true
}

func OpenFontRW(rwops *sdl.RWOps, freeSrc bool, ptSize int) (*Font, error) {
	var cfreeSrc C.int
	if freeSrc {
		cfreeSrc = 1
	}
	f := C.TTF_OpenFontRW((*C.struct_SDL_RWops)(rwops.Ops), cfreeSrc, C.int(ptSize))
	if f == nil {
		return nil, GetError()
	}
	return (*Font)(f), nil
}

func OpenFont(name string, ptSize int) (*Font, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	f := C.TTF_OpenFont(cs, C.int(ptSize))
	if f == nil {
		return nil, GetError()
	}
	return (*Font)(f), nil
}

func OpenFontIndex(name string, ptSize, index int) (*Font, error) {
	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	f := C.TTF_OpenFontIndex(cs, C.int(ptSize), C.long(index))
	if f == nil {
		return nil, GetError()
	}
	return (*Font)(f), nil
}

func (f *Font) Close() {
	C.TTF_CloseFont((*C.TTF_Font)(f))
}

func (f *Font) SizeUTF8(text string) (w, h int, err error) {
	var cw, ch C.int

	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	xerr := C.TTF_SizeUTF8((*C.TTF_Font)(f), ctext, &cw, &ch)
	w, h = int(cw), int(ch)

	if xerr != 0 {
		err = GetError()
		return
	}

	return
}

func (f *Font) RenderUTF8Solid(text string, fg sdl.Color) (*sdl.Surface, error) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	cc := color(fg)
	surface := C.TTF_RenderUTF8_Solid((*C.TTF_Font)(f), ctext, cc)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) RenderGlyphSolid(r rune, fg sdl.Color) (*sdl.Surface, error) {
	var ch C.Uint16

	ch = C.Uint16(r)
	cc := color(fg)
	surface := C.TTF_RenderGlyph_Solid((*C.TTF_Font)(f), ch, cc)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) RenderUTF8Shaded(text string, fg, bg sdl.Color) (*sdl.Surface, error) {
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))

	cc := color(fg)
	cd := color(bg)
	surface := C.TTF_RenderUTF8_Shaded((*C.TTF_Font)(f), cs, cc, cd)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) RenderGlyphShaded(r rune, fg, bg sdl.Color) (*sdl.Surface, error) {
	var ch C.Uint16

	ch = C.Uint16(r)
	cc := color(fg)
	cd := color(bg)
	surface := C.TTF_RenderGlyph_Shaded((*C.TTF_Font)(f), ch, cc, cd)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) RenderUTF8Blended(text string, fg sdl.Color) (*sdl.Surface, error) {
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))

	cc := color(fg)
	surface := C.TTF_RenderUTF8_Blended((*C.TTF_Font)(f), cs, cc)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) RenderGlyphBlended(r rune, fg sdl.Color) (*sdl.Surface, error) {
	var ch C.Uint16

	ch = C.Uint16(r)
	cc := color(fg)
	surface := C.TTF_RenderGlyph_Blended((*C.TTF_Font)(f), ch, cc)
	if surface == nil {
		return nil, GetError()
	}
	return (*sdl.Surface)(unsafe.Pointer(surface)), nil
}

func (f *Font) GlyphMetrics(r rune) (minx, maxx, miny, maxy, advance int, err error) {
	var cminx, cmaxx, cminy, cmaxy, cadvance C.int

	ch := C.Uint16(r)
	xerr := C.TTF_GlyphMetrics((*C.TTF_Font)(f), ch, &cminx, &cmaxx, &cminy, &cmaxy, &cadvance)
	if xerr < 0 {
		err = GetError()
		return
	}

	minx, maxx, miny, maxy = int(cminx), int(cmaxx), int(cminy), int(cmaxy)
	advance = int(cadvance)
	return
}

func color(c sdl.Color) C.SDL_Color {
	return C.SDL_Color{C.Uint8(c.R), C.Uint8(c.G), C.Uint8(c.B), C.Uint8(c.A)}
}

func (f *Font) Style() Style {
	return Style(C.TTF_GetFontStyle((*C.TTF_Font)(f)))
}

func (f *Font) SetStyle(style Style) {
	C.TTF_SetFontStyle((*C.TTF_Font)(f), C.int(style))
}

func (f *Font) FontOutline() int {
	return int(C.TTF_GetFontOutline((*C.TTF_Font)(f)))
}

func (f *Font) SetFontOutline(outline int) {
	C.TTF_SetFontOutline((*C.TTF_Font)(f), C.int(outline))
}

func (f *Font) FontHinting() Hinting {
	return Hinting(C.TTF_GetFontHinting((*C.TTF_Font)(f)))
}

func (f *Font) SetFontHinting(hinting Hinting) {
	C.TTF_SetFontHinting((*C.TTF_Font)(f), C.int(hinting))
}

func (f *Font) Kerning() bool {
	return C.TTF_GetFontKerning((*C.TTF_Font)(f)) != 0
}

func (f *Font) SetKerning(allowed bool) {
	var v C.int
	if allowed {
		v = 1
	}
	C.TTF_SetFontKerning((*C.TTF_Font)(f), v)
}

func (f *Font) Height() int {
	return int(C.TTF_FontHeight((*C.TTF_Font)(f)))
}

func (f *Font) Ascent() int {
	return int(C.TTF_FontAscent((*C.TTF_Font)(f)))
}

func (f *Font) Descent() int {
	return int(C.TTF_FontDescent((*C.TTF_Font)(f)))
}

func (f *Font) LineSkip() int {
	return int(C.TTF_FontLineSkip((*C.TTF_Font)(f)))
}

func (f *Font) Faces() int {
	return int(C.TTF_FontFaces((*C.TTF_Font)(f)))
}

func (f *Font) FixedWidth() bool {
	return C.TTF_FontFaceIsFixedWidth((*C.TTF_Font)(f)) != 0
}

func (f *Font) FamilyName() string {
	return C.GoString(C.TTF_FontFaceFamilyName((*C.TTF_Font)(f)))
}

func (f *Font) StyleName() string {
	return C.GoString(C.TTF_FontFaceStyleName((*C.TTF_Font)(f)))
}

func (f *Font) GlyphIsProvided(ch rune) int {
	return int(C.TTF_GlyphIsProvided((*C.TTF_Font)(f), C.Uint16(ch)))
}
