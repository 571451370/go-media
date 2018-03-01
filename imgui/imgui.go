package imgui

import "github.com/qeedquan/go-media/math/f64"

type Context struct {
}

func New() *Context {
	return &Context{}
}

type Style struct {
	Alpha                  float64
	WindowPadding          f64.Vec2
	WindowRounding         float64
	WindowBorderSize       float64
	WindowMinSize          f64.Vec2
	WindowTitleAlign       f64.Vec2
	ChildRounding          float64
	ChildBorderSize        float64
	PopupRounding          float64
	PopupBorderSize        float64
	FramePadding           f64.Vec2
	FrameRounding          float64
	FrameBorderSize        float64
	ItemSpacing            f64.Vec2
	ItemInnerSpacing       f64.Vec2
	TouchExtraPadding      f64.Vec2
	IndentSpacing          float64
	ColumnsMinSpacing      float64
	ScrollbarSize          float64
	ScrollbarRounding      float64
	GrabMinSize            float64
	GrabRounding           float64
	ButtonTextAlign        f64.Vec2
	DisplayWindowPadding   f64.Vec2
	DisplaySafeAreaPadding f64.Vec2
	MouseCursorScale       float64
	AntiAliasedLines       bool
	AntiAliasedFill        bool
	CurveTesselationTol    float64
	Colors                 [ColorMax]f64.Vec4
}

const (
	ColorText = iota
	ColorTextDisabled
	ColorWindowBg
	ColorChildBg
	ColorPoupBg
	ColorBorder
	ColorBorderShadow
	ColorFrameBg
	ColorMax
)
