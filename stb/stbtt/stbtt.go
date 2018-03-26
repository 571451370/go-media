package stb

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
