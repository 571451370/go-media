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
	ColorPopupBg
	ColorBorder
	ColorBorderShadow
	ColorFrameBg
	ColorFrameBgHovered
	ColorFrameBgActive
	ColorTitleBg
	ColorTitleBgActive
	ColorTitleBgCollapsed
	ColorMenuBarBg
	ColorScrollbarBg
	ColorScrollbarGrab
	ColorScrollbarGrabHovered
	ColorScrollbarGrabActive
	ColorCheckMark
	ColorSliderGrab
	ColorSliderGrabActive
	ColorButton
	ColorButtonHovered
	ColorButtonActive
	ColorHeader
	ColorHeaderHovered
	ColorHeaderActive
	ColorSeparator
	ColorSeparatorHovered
	ColorSeparatorActive
	ColorResizeGrip
	ColorResizeGripHovered
	ColorResizeGripActive
	ColorCloseButton
	ColorCloseButtonHovered
	ColorCloseButtonActive
	ColorPlotLines
	ColorPlotLinesHovered
	ColorPlotHistogram
	ColorPlotHistogramHovered
	ColorTextSelectedBg
	ColorModalWindowDarkening
	ColorDragDropTarget
	ColorNavHighlight
	ColorNavWindowingHighlight
	ColorMax
)

func StyleColorsDark(s *Style) {
	c := s.Colors[:]
	c[ColorText] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	c[ColorTextDisabled] = f64.Vec4{0.50, 0.50, 0.50, 1.00}
	c[ColorWindowBg] = f64.Vec4{0.06, 0.06, 0.06, 0.94}
	c[ColorChildBg] = f64.Vec4{1.00, 1.00, 1.00, 0.00}
	c[ColorPopupBg] = f64.Vec4{0.08, 0.08, 0.08, 0.94}
	c[ColorBorder] = f64.Vec4{0.43, 0.43, 0.50, 0.50}
	c[ColorBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	c[ColorFrameBg] = f64.Vec4{0.16, 0.29, 0.48, 0.54}
	c[ColorFrameBgHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	c[ColorFrameBgActive] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	c[ColorTitleBg] = f64.Vec4{0.04, 0.04, 0.04, 1.00}
	c[ColorTitleBgActive] = f64.Vec4{0.16, 0.29, 0.48, 1.00}
	c[ColorTitleBgCollapsed] = f64.Vec4{0.00, 0.00, 0.00, 0.51}
	c[ColorMenuBarBg] = f64.Vec4{0.14, 0.14, 0.14, 1.00}
	c[ColorScrollbarBg] = f64.Vec4{0.02, 0.02, 0.02, 0.53}
	c[ColorScrollbarGrab] = f64.Vec4{0.31, 0.31, 0.31, 1.00}
	c[ColorScrollbarGrabHovered] = f64.Vec4{0.41, 0.41, 0.41, 1.00}
	c[ColorScrollbarGrabActive] = f64.Vec4{0.51, 0.51, 0.51, 1.00}
	c[ColorCheckMark] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	c[ColorSliderGrab] = f64.Vec4{0.24, 0.52, 0.88, 1.00}
	c[ColorSliderGrabActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	c[ColorButton] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	c[ColorButtonHovered] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	c[ColorButtonActive] = f64.Vec4{0.06, 0.53, 0.98, 1.00}
	c[ColorHeader] = f64.Vec4{0.26, 0.59, 0.98, 0.31}
	c[ColorHeaderHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.80}
	c[ColorHeaderActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	c[ColorSeparator] = c[ColorBorder]
	c[ColorSeparatorHovered] = f64.Vec4{0.10, 0.40, 0.75, 0.78}
	c[ColorSeparatorActive] = f64.Vec4{0.10, 0.40, 0.75, 1.00}
	c[ColorResizeGrip] = f64.Vec4{0.26, 0.59, 0.98, 0.25}
	c[ColorResizeGripHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	c[ColorResizeGripActive] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	c[ColorCloseButton] = f64.Vec4{0.41, 0.41, 0.41, 0.50}
	c[ColorCloseButtonHovered] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	c[ColorCloseButtonActive] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	c[ColorPlotLines] = f64.Vec4{0.61, 0.61, 0.61, 1.00}
	c[ColorPlotLinesHovered] = f64.Vec4{1.00, 0.43, 0.35, 1.00}
	c[ColorPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	c[ColorPlotHistogramHovered] = f64.Vec4{1.00, 0.60, 0.00, 1.00}
	c[ColorTextSelectedBg] = f64.Vec4{0.26, 0.59, 0.98, 0.35}
	c[ColorModalWindowDarkening] = f64.Vec4{0.80, 0.80, 0.80, 0.35}
	c[ColorDragDropTarget] = f64.Vec4{1.00, 1.00, 0.00, 0.90}
	c[ColorNavHighlight] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	c[ColorNavWindowingHighlight] = f64.Vec4{1.00, 1.00, 1.00, 0.70}
}
