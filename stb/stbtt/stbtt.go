package stbtt

/*
#define STBRP_LARGE_RECTS 1
#define STB_RECT_PACK_IMPLEMENTATION
#include "stb_rect_pack.h"
#define STB_TRUETYPE_IMPLEMENTATION
#include "stb_truetype.h"

#cgo LDFLAGS: -lm
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

type (
	PackContext C.stbtt_pack_context
	PackRange   C.stbtt_pack_range
	PackedChar  C.stbtt_packedchar
	FontInfo    C.stbtt_fontinfo
	Rect        C.stbrp_rect
	AlignedQuad C.stbtt_aligned_quad
)

func NewPackContext() *PackContext {
	return (*PackContext)(C.calloc(1, C.sizeof_stbtt_pack_context))
}

func NewFontInfo() *FontInfo {
	return (*FontInfo)(C.calloc(1, C.sizeof_stbtt_fontinfo))
}

func MakePackedChars(n int) []PackedChar {
	p := C.calloc(C.size_t(n), C.sizeof_stbtt_packedchar)
	s := ((*[1 << 30]PackedChar)(unsafe.Pointer(p)))[:n:n]
	return s
}

func MakeRects(n int) []Rect {
	p := C.calloc(C.size_t(n), C.sizeof_stbrp_rect)
	s := ((*[1 << 30]Rect)(unsafe.Pointer(p)))[:n:n]
	return s
}

func MakePackRanges(n int) []PackRange {
	p := C.calloc(C.size_t(n), C.sizeof_stbtt_pack_range)
	s := ((*[1 << 30]PackRange)(unsafe.Pointer(p)))[:n:n]
	return s
}

func FreePackContext(p *PackContext) {
	C.free(unsafe.Pointer(p))
}

func FreeFontInfo(p *FontInfo) {
	C.free(unsafe.Pointer(p))
}

func FreePackedChars(p []PackedChar) {
	C.free(unsafe.Pointer(&p[0]))
}

func FreePackRanges(p []PackRange) {
	C.free(unsafe.Pointer(&p[0]))
}

func FreeRects(p []Rect) {
	C.free(unsafe.Pointer(&p[0]))
}

func (p *PackContext) Pixels() []byte {
	len := p.width * p.height * p.stride_in_bytes
	buf := ((*[1 << 30]byte)(unsafe.Pointer(p.pixels)))[:len:len]
	return buf
}

func (p *PackContext) StrideInBytes() int {
	return int(p.stride_in_bytes)
}

func (p *PackContext) Begin(pixels []byte, width, height, stride_in_bytes, padding int) error {
	var ptr *byte
	if pixels != nil {
		ptr = &pixels[0]
	}
	rc := C.stbtt_PackBegin((*C.stbtt_pack_context)(p), (*C.uchar)(ptr), C.int(width), C.int(height), C.int(stride_in_bytes), C.int(padding), nil)
	if rc == 0 {
		return fmt.Errorf("out of memory")
	}
	return nil
}

func (p *PackContext) End() {
	C.stbtt_PackEnd((*C.stbtt_pack_context)(p))
}

func (p *PackContext) SetPixels(pixels []byte) {
	p.pixels = (*C.uchar)(&pixels[0])
}

func (p *PackContext) SetHeight(height int) {
	p.height = C.int(height)
}

func (p *PackContext) FontRangesRenderIntoRects(info *FontInfo, ranges []PackRange, rects []Rect) {
	C.stbtt_PackFontRangesRenderIntoRects((*C.stbtt_pack_context)(p), (*C.stbtt_fontinfo)(info), (*C.stbtt_pack_range)(&ranges[0]), C.int(len(ranges)), (*C.struct_stbrp_rect)(&rects[0]))
}

func (p *PackContext) FontRangesPackRects(rects []Rect) {
	C.stbtt_PackFontRangesPackRects((*C.stbtt_pack_context)(p), (*C.stbrp_rect)(&rects[0]), C.int(len(rects)))
}

func (p *PackContext) SetOversampling(h_oversample, v_oversample uint) {
	C.stbtt_PackSetOversampling((*C.stbtt_pack_context)(p), C.uint(h_oversample), C.uint(v_oversample))
}

func (p *PackContext) FontRangesGatherRects(info *FontInfo, ranges []PackRange, rects []Rect) int {
	return int(C.stbtt_PackFontRangesGatherRects((*C.stbtt_pack_context)(p), (*C.stbtt_fontinfo)(info), (*C.stbtt_pack_range)(&ranges[0]), C.int(len(ranges)), (*C.stbrp_rect)(&rects[0])))
}

func GetPackedQuad(p []PackedChar, pw, ph, char_index int, xpos, ypos *float64, q *AlignedQuad, align_to_integer int) {
	var cxpos, cypos C.float
	C.stbtt_GetPackedQuad((*C.stbtt_packedchar)(&p[0]), C.int(pw), C.int(ph), C.int(char_index), &cxpos, &cypos, (*C.stbtt_aligned_quad)(q), C.int(align_to_integer))
	*xpos = float64(cxpos)
	*ypos = float64(cypos)
}

func (p *PackedChar) XAdvance() float64 {
	return float64(p.xadvance)
}

func (p *PackedChar) X0() int {
	return int(p.x0)
}

func (p *PackedChar) Y0() int {
	return int(p.y0)
}

func (p *PackedChar) X1() int {
	return int(p.x1)
}

func (p *PackedChar) Y1() int {
	return int(p.y1)
}

func (p *PackRange) SetFontSize(font_size float64) {
	p.font_size = C.float(font_size)
}

func (p *PackRange) FontSize() float64 {
	return float64(p.font_size)
}

func (p *PackRange) SetFirstUnicodeCodepointInRange(range_ int) {
	p.first_unicode_codepoint_in_range = C.int(range_)
}

func (p *PackRange) FirstUnicodeCodepointInRange() int {
	return int(p.first_unicode_codepoint_in_range)
}

func (p *PackRange) SetNumChars(num_chars int) {
	p.num_chars = C.int(num_chars)
}

func (p *PackRange) NumChars() int {
	return int(p.num_chars)
}

func (p *PackRange) SetChardataForRange(chardata_for_range []PackedChar) {
	p.chardata_for_range = (*C.stbtt_packedchar)(&chardata_for_range[0])
}

func (p *PackRange) CharDataForRange() []PackedChar {
	c := ((*[1 << 30]PackedChar)(unsafe.Pointer(p.chardata_for_range)))[:p.num_chars:p.num_chars]
	return c
}

func (p *PackRange) FirstUnicodepointInRange() int {
	return int(p.first_unicode_codepoint_in_range)
}

func (f *FontInfo) Init(data []byte, offset int) error {
	rc := C.stbtt_InitFont((*C.stbtt_fontinfo)(f), (*C.uchar)(unsafe.Pointer(&data[0])), C.int(offset))
	if rc == 0 {
		return fmt.Errorf("failed to load font")
	}
	return nil
}

func (f *FontInfo) ScaleForPixelHeight(height float64) float64 {
	return float64(C.stbtt_ScaleForPixelHeight((*C.stbtt_fontinfo)(f), C.float(height)))
}

func (f *FontInfo) ScaleForMappingEmToPixels(pixels float64) float64 {
	return float64(C.stbtt_ScaleForMappingEmToPixels((*C.stbtt_fontinfo)(f), C.float(pixels)))
}

func (f *FontInfo) FontVMetrics() (ascent, descent, lineGap int) {
	var cascent, cdescent, clineGap C.int
	C.stbtt_GetFontVMetrics((*C.stbtt_fontinfo)(f), &cascent, &cdescent, &clineGap)
	return int(cascent), int(cdescent), int(clineGap)
}

func (f *FontInfo) CodepointHMetrics(codepoint rune) (advanceWidth, leftSideBearing int) {
	var cadvanceWidth, cleftSideBearing C.int
	C.stbtt_GetCodepointHMetrics((*C.stbtt_fontinfo)(f), C.int(codepoint), &cadvanceWidth, &cleftSideBearing)
	return int(cadvanceWidth), int(cleftSideBearing)
}

func (f *FontInfo) CodepointKernAdvance(ch1, ch2 rune) int {
	return int(C.stbtt_GetCodepointKernAdvance((*C.stbtt_fontinfo)(f), C.int(ch1), C.int(ch2)))
}

func (f *FontInfo) CodepointBox(codepoint rune) image.Rectangle {
	var x0, y0, x1, y1 C.int
	C.stbtt_GetCodepointBox((*C.stbtt_fontinfo)(f), C.int(codepoint), &x0, &y0, &x1, &y1)
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func (f *FontInfo) GlyphHMetrics(glyphIndex int) (advanceWidth, leftSideBearing int) {
	var cadvanceWidth, cleftSideBearing C.int
	C.stbtt_GetGlyphHMetrics((*C.stbtt_fontinfo)(f), C.int(glyphIndex), &cadvanceWidth, &cleftSideBearing)
	return int(cadvanceWidth), int(cleftSideBearing)
}

func (f *FontInfo) GlyphKernAdvance(glyph1, glyph2 rune) int {
	return int(C.stbtt_GetGlyphKernAdvance((*C.stbtt_fontinfo)(f), C.int(glyph1), C.int(glyph2)))
}

func (f *FontInfo) BoundingBox() image.Rectangle {
	var x0, y0, x1, y1 C.int
	C.stbtt_GetFontBoundingBox((*C.stbtt_fontinfo)(f), &x0, &y0, &x1, &y1)
	return image.Rect(int(x0), int(y0), int(x1), int(y1))
}

func GetFontOffsetForIndex(data []byte, index int) int {
	return int(C.stbtt_GetFontOffsetForIndex((*C.uchar)(&data[0]), C.int(index)))
}

func GetNumberOfFonts(data []byte) int {
	return int(C.stbtt_GetNumberOfFonts((*C.uchar)(&data[0])))
}

func (r *Rect) ID() int {
	return int(r.id)
}

func (r *Rect) X() int {
	return int(r.x)
}

func (r *Rect) Y() int {
	return int(r.y)
}

func (r *Rect) W() int {
	return int(r.w)
}

func (r *Rect) H() int {
	return int(r.h)
}

func (r *Rect) SetX(x int) {
	r.x = C.stbrp_coord(x)
}

func (r *Rect) SetY(y int) {
	r.y = C.stbrp_coord(y)
}

func (r *Rect) SetW(w int) {
	r.w = C.stbrp_coord(w)
}

func (r *Rect) SetH(h int) {
	r.h = C.stbrp_coord(h)
}

func (r *Rect) WasPacked() int {
	return int(r.was_packed)
}

func (q *AlignedQuad) X0() float64 {
	return float64(q.x0)
}

func (q *AlignedQuad) Y0() float64 {
	return float64(q.y0)
}

func (q *AlignedQuad) X1() float64 {
	return float64(q.x1)
}

func (q *AlignedQuad) Y1() float64 {
	return float64(q.y1)
}

func (q *AlignedQuad) S0() float64 {
	return float64(q.s0)
}

func (q *AlignedQuad) T0() float64 {
	return float64(q.t0)
}

func (q *AlignedQuad) S1() float64 {
	return float64(q.s1)
}

func (q *AlignedQuad) T1() float64 {
	return float64(q.t1)
}

type (
	BakedChar C.stbtt_bakedchar
)

type Bitmap struct {
	*image.RGBA
	*FontInfo
	Chardata    []BakedChar
	Firstchar   rune
	PixelHeight float64
	FG, BG      color.RGBA
}

func BakeFontBitmap(data []byte, offset int, pixel_height float64, pw, ph int, first_char rune, num_chars int) (bmp *Bitmap, numfits int, err error) {
	pixels := make([]byte, pw*ph)
	chardata := make([]BakedChar, num_chars)
	numfits = int(C.stbtt_BakeFontBitmap((*C.uchar)(unsafe.Pointer(&data[0])), C.int(offset), C.float(pixel_height), (*C.uchar)(unsafe.Pointer(&pixels[0])), C.int(pw), C.int(ph), C.int(first_char), C.int(num_chars), (*C.stbtt_bakedchar)(unsafe.Pointer(&chardata[0]))))
	img := image.NewRGBA(image.Rect(0, 0, pw, ph))
	for y := 0; y < ph; y++ {
		for x := 0; x < pw; x++ {
			if pixels[y*pw+x] != 0 {
				img.Set(x, y, color.White)
			}
		}
	}
	bmp = &Bitmap{
		RGBA:        img,
		FontInfo:    &FontInfo{},
		Chardata:    chardata,
		Firstchar:   first_char,
		PixelHeight: pixel_height,
		FG:          color.RGBA{255, 255, 255, 255},
	}
	err = bmp.Init(data, offset)
	return
}

func BakedQuad(cdata []BakedChar, pw, ph, char_index int, xpos, ypos *float64, q *AlignedQuad, opengl_fillrule int) {
	cxpos := C.float(*xpos)
	cypos := C.float(*ypos)
	C.stbtt_GetBakedQuad((*C.stbtt_bakedchar)(unsafe.Pointer(&cdata[0])), C.int(pw), C.int(ph), C.int(char_index), &cxpos, &cypos, (*C.stbtt_aligned_quad)(q), C.int(opengl_fillrule))
	*xpos = float64(cxpos)
	*ypos = float64(cypos)
}

func (b *Bitmap) Print(m *image.RGBA, x, y int, args ...interface{}) {
	b.print(m, x, y, fmt.Sprint(args...))
}

func (b *Bitmap) Printf(m *image.RGBA, x, y int, format string, args ...interface{}) {
	b.print(m, x, y, fmt.Sprintf(format, args...))
}

func (b *Bitmap) print(m *image.RGBA, x, y int, s string) {
	r := b.Bounds()
	px := float64(x)
	py := float64(y)
	sx := px
	for _, c := range s {
		if c == '\n' {
			px = sx
			py += b.PixelHeight
			continue
		}
		c -= rune(b.Firstchar)

		var q AlignedQuad
		BakedQuad(b.Chardata, r.Dx(), r.Dy(), int(c), &px, &py, &q, 1)

		dr := image.Rect(int(q.X0()), int(q.Y0()), int(q.X1()), int(q.Y1()))
		dr = dr.Add(image.Pt(0, int(b.PixelHeight)))

		s0 := q.S0() * float64(r.Dx())
		t0 := q.T0() * float64(r.Dy())
		sp := image.Pt(int(s0), int(t0))

		for y, ty := dr.Min.Y, 0; y < dr.Max.Y; y, ty = y+1, ty+1 {
			for x, tx := dr.Min.X, 0; x < dr.Max.X; x, tx = x+1, tx+1 {
				col := b.RGBAAt(sp.X+tx, sp.Y+ty)
				if col == (color.RGBA{}) {
					if b.BG != (color.RGBA{}) {
						m.Set(x, y, b.BG)
					}
				} else {
					m.Set(x, y, b.FG)
				}
			}
		}
	}
}

func (b *Bitmap) StringSize(text string) (width, height float64) {
	w, h := 0.0, b.PixelHeight
	mw := w
	for _, c := range text {
		if c == '\n' {
			w = 0
			h += b.PixelHeight
			continue
		}

		a, _ := b.CodepointHMetrics(c)
		w += float64(a)
		if mw < w {
			mw = w
		}
	}
	w *= b.ScaleForPixelHeight(b.PixelHeight)
	return mw, h
}
