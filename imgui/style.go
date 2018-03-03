package imgui

import (
	"github.com/qeedquan/go-media/math/f64"
)

type Style struct {
	Alpha                  float64  // Global alpha applies to everything in ImGui.
	WindowPadding          f64.Vec2 // Padding within a window.
	WindowRounding         float64  // Radius of window corners rounding. Set to 0.0f to have rectangular windows.
	WindowBorderSize       float64  // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize          f64.Vec2 // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign       f64.Vec2 // Alignment for title bar text. Defaults to (0.0f,0.5f) for left-aligned,vertically centered.
	ChildRounding          float64  // Radius of child window corners rounding. Set to 0.0f to have rectangular windows.
	ChildBorderSize        float64  // Thickness of border around child windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	PopupRounding          float64  // Radius of popup window corners rounding.
	PopupBorderSize        float64  // Thickness of border around popup windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	FramePadding           f64.Vec2 // Padding within a framed rectangle (used by most widgets).
	FrameRounding          float64  // Radius of frame corners rounding. Set to 0.0f to have rectangular frame (used by most widgets).
	FrameBorderSize        float64  // Thickness of border around frames. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	ItemSpacing            f64.Vec2 // Horizontal and vertical spacing between widgets/lines.
	ItemInnerSpacing       f64.Vec2 // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label).
	TouchExtraPadding      f64.Vec2 // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	IndentSpacing          float64  // Horizontal indentation when e.g. entering a tree node. Generally == (FontSize + FramePadding.x*2).
	ColumnsMinSpacing      float64  // Minimum horizontal spacing between two columns.
	ScrollbarSize          float64  // Width of the vertical scrollbar, Height of the horizontal scrollbar.
	ScrollbarRounding      float64  // Radius of grab corners for scrollbar.
	GrabMinSize            float64  // Minimum width/height of a grab box for slider/scrollbar.
	GrabRounding           float64  // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	ButtonTextAlign        f64.Vec2 // Alignment of button text when button is larger than text. Defaults to (0.5f,0.5f) for horizontally+vertically centered.
	DisplayWindowPadding   f64.Vec2 // Window positions are clamped to be visible within the display area by at least this amount. Only covers regular windows.
	DisplaySafeAreaPadding f64.Vec2 // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows.
	MouseCursorScale       float64  // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines       bool     // Enable anti-aliasing on lines/borders. Disable if you are really tight on CPU/GPU.
	AntiAliasedFill        bool     // Enable anti-aliasing on filled shapes (rounded rectangles, circles, etc.)
	CurveTessellationTol   float64  // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	Colors                 [ColCount]f64.Vec4
}

type Col uint

const (
	ColText Col = iota
	ColTextDisabled
	ColWindowBg // Background of normal windows
	ColChildBg  // Background of child windows
	ColPopupBg  // Background of popups menus tooltips windows
	ColBorder
	ColBorderShadow
	ColFrameBg // Background of checkbox radio button plot slider text input
	ColFrameBgHovered
	ColFrameBgActive
	ColTitleBg
	ColTitleBgActive
	ColTitleBgCollapsed
	ColMenuBarBg
	ColScrollbarBg
	ColScrollbarGrab
	ColScrollbarGrabHovered
	ColScrollbarGrabActive
	ColCheckMark
	ColSliderGrab
	ColSliderGrabActive
	ColButton
	ColButtonHovered
	ColButtonActive
	ColHeader
	ColHeaderHovered
	ColHeaderActive
	ColSeparator
	ColSeparatorHovered
	ColSeparatorActive
	ColResizeGrip
	ColResizeGripHovered
	ColResizeGripActive
	ColCloseButton
	ColCloseButtonHovered
	ColCloseButtonActive
	ColPlotLines
	ColPlotLinesHovered
	ColPlotHistogram
	ColPlotHistogramHovered
	ColTextSelectedBg
	ColModalWindowDarkening // darken entire screen when a modal window is active
	ColDragDropTarget
	ColNavHighlight          // gamepad/keyboard: current highlighted item
	ColNavWindowingHighlight // gamepad/keyboard: when holding NavMenu to focus/move/resize windows
	ColCount
)

func (c *Context) StyleColorsDark(s *Style) {
	col := s.Colors[:]
	col[ColText] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	col[ColTextDisabled] = f64.Vec4{0.50, 0.50, 0.50, 1.00}
	col[ColWindowBg] = f64.Vec4{0.06, 0.06, 0.06, 0.94}
	col[ColChildBg] = f64.Vec4{1.00, 1.00, 1.00, 0.00}
	col[ColPopupBg] = f64.Vec4{0.08, 0.08, 0.08, 0.94}
	col[ColBorder] = f64.Vec4{0.43, 0.43, 0.50, 0.50}
	col[ColBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	col[ColFrameBg] = f64.Vec4{0.16, 0.29, 0.48, 0.54}
	col[ColFrameBgHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	col[ColFrameBgActive] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	col[ColTitleBg] = f64.Vec4{0.04, 0.04, 0.04, 1.00}
	col[ColTitleBgActive] = f64.Vec4{0.16, 0.29, 0.48, 1.00}
	col[ColTitleBgCollapsed] = f64.Vec4{0.00, 0.00, 0.00, 0.51}
	col[ColMenuBarBg] = f64.Vec4{0.14, 0.14, 0.14, 1.00}
	col[ColScrollbarBg] = f64.Vec4{0.02, 0.02, 0.02, 0.53}
	col[ColScrollbarGrab] = f64.Vec4{0.31, 0.31, 0.31, 1.00}
	col[ColScrollbarGrabHovered] = f64.Vec4{0.41, 0.41, 0.41, 1.00}
	col[ColScrollbarGrabActive] = f64.Vec4{0.51, 0.51, 0.51, 1.00}
	col[ColCheckMark] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	col[ColSliderGrab] = f64.Vec4{0.24, 0.52, 0.88, 1.00}
	col[ColSliderGrabActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	col[ColButton] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	col[ColButtonHovered] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	col[ColButtonActive] = f64.Vec4{0.06, 0.53, 0.98, 1.00}
	col[ColHeader] = f64.Vec4{0.26, 0.59, 0.98, 0.31}
	col[ColHeaderHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.80}
	col[ColHeaderActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	col[ColSeparator] = col[ColBorder]
	col[ColSeparatorHovered] = f64.Vec4{0.10, 0.40, 0.75, 0.78}
	col[ColSeparatorActive] = f64.Vec4{0.10, 0.40, 0.75, 1.00}
	col[ColResizeGrip] = f64.Vec4{0.26, 0.59, 0.98, 0.25}
	col[ColResizeGripHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	col[ColResizeGripActive] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	col[ColCloseButton] = f64.Vec4{0.41, 0.41, 0.41, 0.50}
	col[ColCloseButtonHovered] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	col[ColCloseButtonActive] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	col[ColPlotLines] = f64.Vec4{0.61, 0.61, 0.61, 1.00}
	col[ColPlotLinesHovered] = f64.Vec4{1.00, 0.43, 0.35, 1.00}
	col[ColPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	col[ColPlotHistogramHovered] = f64.Vec4{1.00, 0.60, 0.00, 1.00}
	col[ColTextSelectedBg] = f64.Vec4{0.26, 0.59, 0.98, 0.35}
	col[ColModalWindowDarkening] = f64.Vec4{0.80, 0.80, 0.80, 0.35}
	col[ColDragDropTarget] = f64.Vec4{1.00, 1.00, 0.00, 0.90}
	col[ColNavHighlight] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	col[ColNavWindowingHighlight] = f64.Vec4{1.00, 1.00, 1.00, 0.70}
}