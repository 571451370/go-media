package nvg

/*
#define NANOVG_GL3_IMPLEMENTATION

#include <stdlib.h>
#include <GL/glew.h>
#include "nanovg.h"
#include "nanovg_gl.h"
*/
// #cgo LDFLAGS: -lm -lGLEW -lGL -ldl
import "C"

import (
	"fmt"
	"image/color"
	"unsafe"
)

type (
	Context       C.NVGcontext
	TextRow       C.NVGtextRow
	GlyphPosition C.NVGglyphPosition
	Paint         C.NVGpaint

	CreateFlags C.enum_NVGcreateFlags
	LineCap     C.enum_NVGlineCap
	Align       C.enum_NVGalign
	BlendFactor C.enum_NVGblendFactor
	CompositeOp C.enum_NVGcompositeOperation
	Winding     C.enum_NVGwinding
	Solidity    C.enum_NVGsolidity
	ImageFlags  C.enum_NVGimageFlags
)

const (
	ALIGN_LEFT   Align = C.NVG_ALIGN_LEFT
	ALIGN_CENTER Align = C.NVG_ALIGN_CENTER
	ALIGN_RIGHT  Align = C.NVG_ALIGN_RIGHT

	ALIGN_TOP      Align = C.NVG_ALIGN_TOP
	ALIGN_MIDDLE   Align = C.NVG_ALIGN_MIDDLE
	ALIGN_BOTTOM   Align = C.NVG_ALIGN_BOTTOM
	ALIGN_BASELINE Align = C.NVG_ALIGN_BASELINE
)

const (
	ZERO                BlendFactor = C.NVG_ZERO
	ONE                 BlendFactor = C.NVG_ONE
	SRC_COLOR           BlendFactor = C.NVG_SRC_COLOR
	ONE_MINUS_SRC_COLOR BlendFactor = C.NVG_ONE_MINUS_SRC_COLOR
	DST_COLOR           BlendFactor = C.NVG_DST_COLOR
	ONE_MINUS_DST_COLOR BlendFactor = C.NVG_ONE_MINUS_DST_COLOR
	SRC_ALPHA           BlendFactor = C.NVG_SRC_ALPHA
	ONE_MINUS_SRC_ALPHA BlendFactor = C.NVG_ONE_MINUS_SRC_ALPHA
	DST_ALPHA           BlendFactor = C.NVG_DST_ALPHA
	ONE_MINUS_DST_ALPHA BlendFactor = C.NVG_ONE_MINUS_DST_ALPHA
	SRC_ALPHA_SATURATE  BlendFactor = C.NVG_SRC_ALPHA_SATURATE
)

const (
	SOURCE_OVER      CompositeOp = C.NVG_SOURCE_OVER
	SOURCE_IN        CompositeOp = C.NVG_SOURCE_IN
	SOURCE_OUT       CompositeOp = C.NVG_SOURCE_OUT
	ATOP             CompositeOp = C.NVG_ATOP
	DESTINATION_OVER CompositeOp = C.NVG_DESTINATION_OVER
	DESTINATION_IN   CompositeOp = C.NVG_DESTINATION_IN
	DESTINATION_OUT  CompositeOp = C.NVG_DESTINATION_OUT
	DESTINATION_ATOP CompositeOp = C.NVG_DESTINATION_ATOP
	LIGHTER          CompositeOp = C.NVG_LIGHTER
	COPY             CompositeOp = C.NVG_COPY
	XOR              CompositeOp = C.NVG_XOR
)

const (
	BUTT   LineCap = C.NVG_BUTT
	ROUND  LineCap = C.NVG_ROUND
	SQUARE LineCap = C.NVG_SQUARE
	BEVEL  LineCap = C.NVG_BEVEL
	MITER  LineCap = C.NVG_MITER
)

const (
	CCW Winding = C.NVG_CCW
	CW  Winding = C.NVG_CW
)

const (
	SOLID Solidity = C.NVG_SOLID
	HOLE  Solidity = C.NVG_HOLE
)

const (
	IMAGE_GENERATE_MIPMAPS ImageFlags = C.NVG_IMAGE_GENERATE_MIPMAPS
	IMAGE_REPEATX          ImageFlags = C.NVG_IMAGE_REPEATX
	IMAGE_REPEATY          ImageFlags = C.NVG_IMAGE_REPEATY
	IMAGE_FLIPY            ImageFlags = C.NVG_IMAGE_FLIPY
	IMAGE_PREMULTIPLIED    ImageFlags = C.NVG_IMAGE_PREMULTIPLIED
	IMAGE_NEAREST          ImageFlags = C.NVG_IMAGE_NEAREST
)

const (
	ANTIALIAS       CreateFlags = C.NVG_ANTIALIAS
	STENCIL_STROKES CreateFlags = C.NVG_STENCIL_STROKES
	DEBUG           CreateFlags = C.NVG_DEBUG
)

func CreateGL3(flags CreateFlags) (*Context, error) {
	ctx := (*Context)(C.nvgCreateGL3(C.int(flags)))
	if ctx == nil {
		return nil, fmt.Errorf("failed to create nvg context")
	}
	return ctx, nil
}

func (c *Context) GlobalAlpha(alpha float64) {
	C.nvgGlobalAlpha((*C.NVGcontext)(c), C.float(alpha))
}

func (c *Context) BeginFrame(width, height int, aspect float64) {
	C.nvgBeginFrame((*C.NVGcontext)(c), C.int(width), C.int(height), C.float(aspect))
}

func (c *Context) EndFrame() {
	C.nvgEndFrame((*C.NVGcontext)(c))
}

func (c *Context) ResetTransform() {
	C.nvgResetTransform((*C.NVGcontext)(c))
}

func (c *Context) Transform(a, b, c_, d, e, f float64) {
	C.nvgTransform((*C.NVGcontext)(c), C.float(a), C.float(b), C.float(c_), C.float(d), C.float(e), C.float(f))
}

func (c *Context) Translate(x, y float64) {
	C.nvgTranslate((*C.NVGcontext)(c), C.float(x), C.float(y))
}

func (c *Context) Rotate(angle float64) {
	C.nvgRotate((*C.NVGcontext)(c), C.float(angle))
}

func (c *Context) SkewX(angle float64) {
	C.nvgSkewX((*C.NVGcontext)(c), C.float(angle))
}

func (c *Context) SkewY(angle float64) {
	C.nvgSkewY((*C.NVGcontext)(c), C.float(angle))
}

func (c *Context) Scale(x, y float64) {
	C.nvgScale((*C.NVGcontext)(c), C.float(x), C.float(y))
}

func (c *Context) Save() {
	C.nvgSave((*C.NVGcontext)(c))
}

func (c *Context) Restore() {
	C.nvgRestore((*C.NVGcontext)(c))
}

func (c *Context) Scissor(x, y, w, h float64) {
	C.nvgScissor((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h))
}

func (c *Context) IntersectScissor(x, y, w, h float64) {
	C.nvgIntersectScissor((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h))
}

func (c *Context) ResetScissor() {
	C.nvgResetScissor((*C.NVGcontext)(c))
}

func (c *Context) CreateImage(name string, flags ImageFlags) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return int(C.nvgCreateImage((*C.NVGcontext)(c), cname, C.int(flags)))
}

func (c *Context) ImageSize(image_ int) (w, h int) {
	var cw, ch C.int
	C.nvgImageSize((*C.NVGcontext)(c), C.int(image_), &cw, &ch)
	return int(cw), int(ch)
}

func (c *Context) DeleteImage(image_ int) {
	C.nvgDeleteImage((*C.NVGcontext)(c), C.int(image_))
}

func (c *Context) LinearGradient(sx, sy, ex, ey float64, icol, ocol color.RGBA) Paint {
	return Paint(C.nvgLinearGradient((*C.NVGcontext)(c), C.float(sx), C.float(sy), C.float(ex), C.float(ey), rgba(icol), rgba(ocol)))
}

func (c *Context) BoxGradient(x, y, w, h, r, f float64, icol, ocol color.RGBA) Paint {
	return Paint(C.nvgBoxGradient((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h), C.float(r), C.float(f), rgba(icol), rgba(ocol)))
}

func (c *Context) RadialGradient(cx, cy, inr, outr float64, icol, ocol color.RGBA) Paint {
	return Paint(C.nvgRadialGradient((*C.NVGcontext)(c), C.float(cx), C.float(cy), C.float(inr), C.float(outr), rgba(icol), rgba(ocol)))
}

func (c *Context) ImagePattern(ox, oy, ex, ey, angle float64, image_ int, alpha float64) Paint {
	return Paint(C.nvgImagePattern((*C.NVGcontext)(c), C.float(ox), C.float(oy), C.float(ex), C.float(ey), C.float(angle), C.int(image_), C.float(alpha)))
}

func (c *Context) BeginPath() {
	C.nvgBeginPath((*C.NVGcontext)(c))
}

func (c *Context) MoveTo(x, y float64) {
	C.nvgMoveTo((*C.NVGcontext)(c), C.float(x), C.float(y))
}

func (c *Context) LineTo(x, y float64) {
	C.nvgLineTo((*C.NVGcontext)(c), C.float(x), C.float(y))
}

func (c *Context) BezierTo(c1x, c1y, c2x, c2y, x, y float64) {
	C.nvgBezierTo((*C.NVGcontext)(c), C.float(c1x), C.float(c1y), C.float(c2x), C.float(c2y), C.float(x), C.float(y))
}

func (c *Context) QuadTo(cx, cy, x, y float64) {
	C.nvgQuadTo((*C.NVGcontext)(c), C.float(cx), C.float(cy), C.float(x), C.float(y))
}

func (c *Context) ArcTo(x1, y1, x2, y2, radius float64) {
	C.nvgArcTo((*C.NVGcontext)(c), C.float(x1), C.float(y1), C.float(x2), C.float(y2), C.float(radius))
}

func (c *Context) ClosePath() {
	C.nvgClosePath((*C.NVGcontext)(c))
}

func (c *Context) Arc(cx, cy, r, a0, a1 float64, dir int) {
	C.nvgArc((*C.NVGcontext)(c), C.float(cx), C.float(cy), C.float(r), C.float(a0), C.float(a1), C.int(dir))
}

func (c *Context) RoundRect(x, y, w, h, r float64) {
	C.nvgRoundedRect((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h), C.float(r))
}

func (c *Context) RoundRectVarying(x, y, w, h, radTopLeft, radTopRight, radBottomRight, radBottomLeft float64) {
	C.nvgRoundedRectVarying((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h), C.float(radTopLeft),
		C.float(radTopRight), C.float(radBottomRight), C.float(radTopLeft))
}

func (c *Context) Rect(x, y, w, h float64) {
	C.nvgRect((*C.NVGcontext)(c), C.float(x), C.float(y), C.float(w), C.float(h))
}

func (c *Context) Ellipse(cx, cy, rx, ry float64) {
	C.nvgEllipse((*C.NVGcontext)(c), C.float(cx), C.float(cy), C.float(rx), C.float(ry))
}

func (c *Context) Circle(cx, cy, r float64) {
	C.nvgCircle((*C.NVGcontext)(c), C.float(cx), C.float(cy), C.float(r))
}

func (c *Context) FillColor(p color.RGBA) {
	C.nvgFillColor((*C.NVGcontext)(c), rgba(p))
}

func (c *Context) StrokeColor(p color.RGBA) {
	C.nvgStrokeColor((*C.NVGcontext)(c), rgba(p))
}

func (c *Context) StrokePaint(p Paint) {
	C.nvgStrokePaint((*C.NVGcontext)(c), (C.NVGpaint)(p))
}

func (c *Context) FillPaint(p Paint) {
	C.nvgFillPaint((*C.NVGcontext)(c), (C.NVGpaint)(p))
}

func (c *Context) MiterLimit(limit float64) {
	C.nvgMiterLimit((*C.NVGcontext)(c), C.float(limit))
}

func (c *Context) StrokeWidth(size float64) {
	C.nvgStrokeWidth((*C.NVGcontext)(c), C.float(size))
}

func (c *Context) LineCap(cap LineCap) {
	C.nvgLineCap((*C.NVGcontext)(c), C.int(cap))
}

func (c *Context) LineJoin(join LineCap) {
	C.nvgLineJoin((*C.NVGcontext)(c), C.int(join))
}

func (c *Context) Fill() {
	C.nvgFill((*C.NVGcontext)(c))
}

func (c *Context) FontSize(size float64) {
	C.nvgFontSize((*C.NVGcontext)(c), C.float(size))
}

func (c *Context) FontFace(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.nvgFontFace((*C.NVGcontext)(c), cname)
}

func (c *Context) TextAlign(align Align) {
	C.nvgTextAlign((*C.NVGcontext)(c), C.int(align))
}

func rgba(c color.RGBA) C.NVGcolor {
	return C.nvgRGBA(C.uchar(c.R), C.uchar(c.G), C.uchar(c.B), C.uchar(c.A))
}
