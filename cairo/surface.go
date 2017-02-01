package cairo

/*
#include <cairo.h>
#include <stdlib.h>
*/
import "C"
import (
	"image"
	"image/color"
	"unsafe"
)

type Surface C.cairo_surface_t
type Format C.cairo_format_t
type Content C.cairo_content_t

const (
	FORMAT_INVALID      Format = C.CAIRO_FORMAT_INVALID
	FORMAT_ARGB32       Format = C.CAIRO_FORMAT_ARGB32
	FORMAT_RGB24        Format = C.CAIRO_FORMAT_RGB24
	FORMAT_A8           Format = C.CAIRO_FORMAT_A8
	FORMAT_A1           Format = C.CAIRO_FORMAT_A1
	FORMAT_RGB16_565    Format = C.CAIRO_FORMAT_RGB16_565
	FORMAT_FORMAT_RGB30 Format = C.CAIRO_FORMAT_RGB30
)

const (
	CONTENT_COLOR       Content = C.CAIRO_CONTENT_COLOR
	CONTENT_ALPHA       Content = C.CAIRO_CONTENT_ALPHA
	CONTENT_COLOR_ALPHA Content = C.CAIRO_CONTENT_COLOR_ALPHA
)

func CreateImageSurfaceFromPNG(filename string) *Surface {
	cs := C.CString(filename)
	s := (*Surface)(C.cairo_image_surface_create_from_png(cs))
	C.free(unsafe.Pointer(cs))
	return s
}

func CreateImageSurface(format Format, w, h int) *Surface {
	return (*Surface)(C.cairo_image_surface_create((C.cairo_format_t)(format), C.int(w), C.int(h)))
}

func CreateImageSurfaceForData(data []byte, format Format, width, height, stride int) *Surface {
	return (*Surface)(C.cairo_image_surface_create_for_data((*C.uchar)(&data[0]), C.cairo_format_t(format), C.int(width), C.int(height), C.int(stride)))
}

func (s *Surface) Status() Status {
	return Status(C.cairo_surface_status((*C.cairo_surface_t)(s)))
}

func (s *Surface) Ref() *Surface {
	return (*Surface)(C.cairo_surface_reference((*C.cairo_surface_t)(s)))
}

func (s *Surface) Destroy() {
	C.cairo_surface_destroy((*C.cairo_surface_t)(s))
}

func (s *Surface) MarkDirty() {
	C.cairo_surface_mark_dirty((*C.cairo_surface_t)(s))
}

func (s *Surface) Finish() {
	C.cairo_surface_finish((*C.cairo_surface_t)(s))
}

func (s *Surface) Flush() {
	C.cairo_surface_flush((*C.cairo_surface_t)(s))
}

func (s *Surface) CreateSimilarImage(format Format, w, h int) *Surface {
	return (*Surface)(C.cairo_surface_create_similar_image((*C.cairo_surface_t)(s), C.cairo_format_t(format), C.int(w), C.int(h)))
}

func (s *Surface) WriteToPNG(filename string) {
	cs := C.CString(filename)
	C.cairo_surface_write_to_png((*C.cairo_surface_t)(s), cs)
	C.free(unsafe.Pointer(cs))
}

func (s *Surface) Width() int {
	return int(C.cairo_image_surface_get_width((*C.cairo_surface_t)(s)))
}

func (s *Surface) Height() int {
	return int(C.cairo_image_surface_get_height((*C.cairo_surface_t)(s)))
}

func (s *Surface) Stride() int {
	return int(C.cairo_image_surface_get_stride((*C.cairo_surface_t)(s)))
}

func (s *Surface) Format() Format {
	return Format(C.cairo_image_surface_get_format((*C.cairo_surface_t)(s)))
}

func (s *Surface) Data() []byte {
	n := s.Stride() * s.Height()
	p := ((*[1 << 30]byte)(unsafe.Pointer(C.cairo_image_surface_get_data((*C.cairo_surface_t)(s)))))[:n:n]
	return p
}

func (s *Surface) At(x, y int) color.Color {
	n := s.Stride()
	p := s.Data()
	r := s.Bounds()
	if !image.Pt(x, y).In(r) {
		return color.RGBA{}
	}

	switch s.Format() {
	case FORMAT_ARGB32:
		return color.RGBA{
			p[y*n+x*4],
			p[y*n+x*4+1],
			p[y*n+x*4+2],
			p[y*n+x*4+3],
		}
	}

	return color.RGBA{}
}

func (s *Surface) Set(x, y int, c color.Color) {
	n := s.Stride()
	p := s.Data()
	r := s.Bounds()
	if !image.Pt(x, y).In(r) {
		return
	}

	cr, cg, cb, ca := c.RGBA()
	switch s.Format() {
	case FORMAT_ARGB32:
		p[y*n+x*4] = byte(cr >> 8)
		p[y*n+x*4+1] = byte(cg >> 8)
		p[y*n+x*4+2] = byte(cb >> 8)
		p[y*n+x*4+3] = byte(ca >> 8)
	}
}

func (s *Surface) Bounds() image.Rectangle {
	return image.Rect(0, 0, s.Width(), s.Height())
}

func FormatStrideForWidth(f Format, w int) int {
	return int(C.cairo_format_stride_for_width(C.cairo_format_t(f), C.int(w)))
}
