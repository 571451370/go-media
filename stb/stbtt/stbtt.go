package stbtt

/*
#include "stb_truetype.h"

#cgo LDFLAGS: -lm
*/
import "C"
import "fmt"

type (
	PackContext C.stbtt_pack_context
	PackRange   C.stbtt_pack_range
	PackedChar  C.stbtt_packedchar
	FontInfo    C.stbtt_fontinfo
	Rect        C.stbrp_rect
	AlignedQuad C.stbtt_aligned_quad
)

func (p *PackContext) Begin(pixels []byte, width, height, stride_in_bytes, padding int) error {
	rc := C.stbtt_PackBegin((*C.stbtt_pack_context)(p), (*C.uchar)(&pixels[0]), C.int(width), C.int(height), C.int(stride_in_bytes), C.int(padding), nil)
	if rc == 0 {
		return fmt.Errorf("out of memory")
	}
	return nil
}

func (p *PackContext) End() {
	C.stbtt_PackEnd((*C.stbtt_pack_context)(p))
}

func (p *PackContext) PackFontRangesRenderIntoRects(info *FontInfo, ranges []PackRange, rects []Rect) {
	C.stbtt_PackFontRangesRenderIntoRects((*C.stbtt_pack_context)(p), (*C.stbtt_fontinfo)(info), (*C.stbtt_pack_range)(&ranges[0]), C.int(len(ranges)), (*C.struct_stbrp_rect)(&rects[0]))
}

func (p *PackContext) PackFontRangesPackRects(rects []Rect) {
	C.stbtt_PackFontRangesPackRects((*C.stbtt_pack_context)(p), (*C.stbrp_rect)(&rects[0]), C.int(len(rects)))
}

func (p *PackContext) SetOversampling(h_oversample, v_oversample uint) {
	C.stbtt_PackSetOversampling((*C.stbtt_pack_context)(p), C.uint(h_oversample), C.uint(v_oversample))
}

func (p *PackedChar) GetPackedQuad(pw, ph, char_index int, xpos, ypos *float64, q *AlignedQuad, align_to_integer int) {
	var cxpos, cypos C.float
	C.stbtt_GetPackedQuad((*C.stbtt_packedchar)(p), C.int(pw), C.int(ph), C.int(char_index), &cxpos, &cypos, (*C.stbtt_aligned_quad)(q), C.int(align_to_integer))
	*xpos = float64(cxpos)
	*ypos = float64(cypos)
}

func (f *FontInfo) Init(data []byte, offset int) error {
	rc := C.stbtt_InitFont((*C.stbtt_fontinfo)(f), (*C.uchar)(&data[0]), C.int(offset))
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

func (f *FontInfo) GetFontVMetrics() (ascent, descent, lineGap int) {
	var cascent, cdescent, clineGap C.int
	C.stbtt_GetFontVMetrics((*C.stbtt_fontinfo)(f), &cascent, &cdescent, &clineGap)
	return int(cascent), int(cdescent), int(clineGap)
}

func GetFontOffsetForIndex(data []byte, index int) int {
	return int(C.stbtt_GetFontOffsetForIndex((*C.uchar)(&data[0]), C.int(index)))
}

func GetNumberOfFonts(data []byte) int {
	return int(C.stbtt_GetNumberOfFonts((*C.uchar)(&data[0])))
}