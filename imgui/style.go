package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/image/chroma"
	"github.com/qeedquan/go-media/math/f64"
)

type Col int

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
	ColCOUNT
)

type ColMod struct {
	Col         Col
	BackupValue f64.Vec4
}

type StyleMod struct {
	VarIdx StyleVar
	Value  interface{}
}

type Style struct {
	Alpha                  float64  // Global alpha applies to everything in ImGui.
	WindowPadding          f64.Vec2 // Padding within a window.
	WindowRounding         float64  // Radius of window corners rounding. Set to 0.0f to have rectangular windows.
	WindowBorderSize       float64  // Thickness of border around windows. Generally set to 0.0f or 1.0f. (Other values are not well tested and more CPU/GPU costly).
	WindowMinSize          f64.Vec2 // Minimum window size. This is a global setting. If you want to constraint individual windows, use SetNextWindowSizeConstraints().
	WindowTitleAlign       f64.Vec2 // Alignment for title bar text. Defaults to (0.0,0.5f) for left-aligned,vertically centered.
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
	ButtonTextAlign        f64.Vec2 // Alignment of button text when button is larger than text. Defaults to (0.5,0.5f) for horizontally+vertically centered.
	DisplayWindowPadding   f64.Vec2 // Window positions are clamped to be visible within the display area by at least this amount. Only covers regular windows.
	DisplaySafeAreaPadding f64.Vec2 // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows.
	MouseCursorScale       float64  // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	AntiAliasedLines       bool     // Enable anti-aliasing on lines/borders. Disable if you are really tight on CPU/GPU.
	AntiAliasedFill        bool     // Enable anti-aliasing on filled shapes (rounded rectangles, circles, etc.)
	CurveTessellationTol   float64  // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
	Colors                 [ColCOUNT]f64.Vec4
}

// Enumeration for PushStyleVar() / PopStyleVar() to temporarily modify the ImGuiStyle structure.
// NB: the enum only refers to fields of ImGuiStyle which makes sense to be pushed/popped inside UI code. During initialization, feel free to just poke into ImGuiStyle directly.
// NB: if changing this enum, you need to update the associated internal table GStyleVarInfo[] accordingly. This is where we link enum values to members offset/type.
type StyleVar int

const (
	// Enum name ......................// Member in ImGuiStyle structure (see ImGuiStyle for descriptions)
	StyleVarAlpha             StyleVar = iota // float     Alpha
	StyleVarWindowPadding                     // ImVec2    WindowPadding
	StyleVarWindowRounding                    // float     WindowRounding
	StyleVarWindowBorderSize                  // float     WindowBorderSize
	StyleVarWindowMinSize                     // ImVec2    WindowMinSize
	StyleVarWindowTitleAlign                  // ImVec2    WindowTitleAlign
	StyleVarChildRounding                     // float     ChildRounding
	StyleVarChildBorderSize                   // float     ChildBorderSize
	StyleVarPopupRounding                     // float     PopupRounding
	StyleVarPopupBorderSize                   // float     PopupBorderSize
	StyleVarFramePadding                      // ImVec2    FramePadding
	StyleVarFrameRounding                     // float     FrameRounding
	StyleVarFrameBorderSize                   // float     FrameBorderSize
	StyleVarItemSpacing                       // ImVec2    ItemSpacing
	StyleVarItemInnerSpacing                  // ImVec2    ItemInnerSpacing
	StyleVarIndentSpacing                     // float     IndentSpacing
	StyleVarScrollbarSize                     // float     ScrollbarSize
	StyleVarScrollbarRounding                 // float     ScrollbarRounding
	StyleVarGrabMinSize                       // float     GrabMinSize
	StyleVarGrabRounding                      // float     GrabRounding
	StyleVarButtonTextAlign                   // ImVec2    ButtonTextAlign
	StyleVarCOUNT
)

func (c *Context) StyleColorsDark(style *Style) {
	if style == nil {
		style = c.GetStyle()
	}
	colors := style.Colors[:]

	colors[ColText] = f64.Vec4{1.00, 1.00, 1.00, 1.00}
	colors[ColTextDisabled] = f64.Vec4{0.50, 0.50, 0.50, 1.00}
	colors[ColWindowBg] = f64.Vec4{0.06, 0.06, 0.06, 0.94}
	colors[ColChildBg] = f64.Vec4{1.00, 1.00, 1.00, 0.00}
	colors[ColPopupBg] = f64.Vec4{0.08, 0.08, 0.08, 0.94}
	colors[ColBorder] = f64.Vec4{0.43, 0.43, 0.50, 0.50}
	colors[ColBorderShadow] = f64.Vec4{0.00, 0.00, 0.00, 0.00}
	colors[ColFrameBg] = f64.Vec4{0.16, 0.29, 0.48, 0.54}
	colors[ColFrameBgHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColFrameBgActive] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColTitleBg] = f64.Vec4{0.04, 0.04, 0.04, 1.00}
	colors[ColTitleBgActive] = f64.Vec4{0.16, 0.29, 0.48, 1.00}
	colors[ColTitleBgCollapsed] = f64.Vec4{0.00, 0.00, 0.00, 0.51}
	colors[ColMenuBarBg] = f64.Vec4{0.14, 0.14, 0.14, 1.00}
	colors[ColScrollbarBg] = f64.Vec4{0.02, 0.02, 0.02, 0.53}
	colors[ColScrollbarGrab] = f64.Vec4{0.31, 0.31, 0.31, 1.00}
	colors[ColScrollbarGrabHovered] = f64.Vec4{0.41, 0.41, 0.41, 1.00}
	colors[ColScrollbarGrabActive] = f64.Vec4{0.51, 0.51, 0.51, 1.00}
	colors[ColCheckMark] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSliderGrab] = f64.Vec4{0.24, 0.52, 0.88, 1.00}
	colors[ColSliderGrabActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColButton] = f64.Vec4{0.26, 0.59, 0.98, 0.40}
	colors[ColButtonHovered] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColButtonActive] = f64.Vec4{0.06, 0.53, 0.98, 1.00}
	colors[ColHeader] = f64.Vec4{0.26, 0.59, 0.98, 0.31}
	colors[ColHeaderHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.80}
	colors[ColHeaderActive] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColSeparator] = colors[ColBorder]
	colors[ColSeparatorHovered] = f64.Vec4{0.10, 0.40, 0.75, 0.78}
	colors[ColSeparatorActive] = f64.Vec4{0.10, 0.40, 0.75, 1.00}
	colors[ColResizeGrip] = f64.Vec4{0.26, 0.59, 0.98, 0.25}
	colors[ColResizeGripHovered] = f64.Vec4{0.26, 0.59, 0.98, 0.67}
	colors[ColResizeGripActive] = f64.Vec4{0.26, 0.59, 0.98, 0.95}
	colors[ColCloseButton] = f64.Vec4{0.41, 0.41, 0.41, 0.50}
	colors[ColCloseButtonHovered] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	colors[ColCloseButtonActive] = f64.Vec4{0.98, 0.39, 0.36, 1.00}
	colors[ColPlotLines] = f64.Vec4{0.61, 0.61, 0.61, 1.00}
	colors[ColPlotLinesHovered] = f64.Vec4{1.00, 0.43, 0.35, 1.00}
	colors[ColPlotHistogram] = f64.Vec4{0.90, 0.70, 0.00, 1.00}
	colors[ColPlotHistogramHovered] = f64.Vec4{1.00, 0.60, 0.00, 1.00}
	colors[ColTextSelectedBg] = f64.Vec4{0.26, 0.59, 0.98, 0.35}
	colors[ColModalWindowDarkening] = f64.Vec4{0.80, 0.80, 0.80, 0.35}
	colors[ColDragDropTarget] = f64.Vec4{1.00, 1.00, 0.00, 0.90}
	colors[ColNavHighlight] = f64.Vec4{0.26, 0.59, 0.98, 1.00}
	colors[ColNavWindowingHighlight] = f64.Vec4{1.00, 1.00, 1.00, 0.70}
}

func (c *Context) GetColorFromStyle(idx Col) color.RGBA {
	return c.GetColorFromStyleWithAlpha(idx, 1)
}

func (c *Context) GetColorFromStyleWithAlpha(idx Col, alpha_mul float64) color.RGBA {
	style := &c.Style
	col := style.Colors[idx]
	col.W *= style.Alpha * alpha_mul
	return col.ToRGBA()
}

func (s *Style) Init() {
	s.Alpha = 1.0                             // Global alpha applies to everything in ImGui
	s.WindowPadding = f64.Vec2{8, 8}          // Padding within a window
	s.WindowRounding = 7.0                    // Radius of window corners rounding. Set to 0.0f to have rectangular windows
	s.WindowBorderSize = 1.0                  // Thickness of border around windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.WindowMinSize = f64.Vec2{32, 32}        // Minimum window size
	s.WindowTitleAlign = f64.Vec2{0.0, 0.5}   // Alignment for title bar text
	s.ChildRounding = 0.0                     // Radius of child window corners rounding. Set to 0.0f to have rectangular child windows
	s.ChildBorderSize = 1.0                   // Thickness of border around child windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.PopupRounding = 0.0                     // Radius of popup window corners rounding. Set to 0.0f to have rectangular child windows
	s.PopupBorderSize = 1.0                   // Thickness of border around popup or tooltip windows. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.FramePadding = f64.Vec2{4, 3}           // Padding within a framed rectangle (used by most widgets)
	s.FrameRounding = 0.0                     // Radius of frame corners rounding. Set to 0.0f to have rectangular frames (used by most widgets).
	s.FrameBorderSize = 0.0                   // Thickness of border around frames. Generally set to 0.0f or 1.0f. Other values not well tested.
	s.ItemSpacing = f64.Vec2{8, 4}            // Horizontal and vertical spacing between widgets/lines
	s.ItemInnerSpacing = f64.Vec2{4, 4}       // Horizontal and vertical spacing between within elements of a composed widget (e.g. a slider and its label)
	s.TouchExtraPadding = f64.Vec2{0, 0}      // Expand reactive bounding box for touch-based system where touch position is not accurate enough. Unfortunately we don't sort widgets so priority on overlap will always be given to the first widget. So don't grow this too much!
	s.IndentSpacing = 21.0                    // Horizontal spacing when e.g. entering a tree node. Generally == (FontSize + FramePadding.x*2).
	s.ColumnsMinSpacing = 6.0                 // Minimum horizontal spacing between two columns
	s.ScrollbarSize = 16.0                    // Width of the vertical scrollbar, Height of the horizontal scrollbar
	s.ScrollbarRounding = 9.0                 // Radius of grab corners rounding for scrollbar
	s.GrabMinSize = 10.0                      // Minimum width/height of a grab box for slider/scrollbar
	s.GrabRounding = 0.0                      // Radius of grabs corners rounding. Set to 0.0f to have rectangular slider grabs.
	s.ButtonTextAlign = f64.Vec2{0.5, 0.5}    // Alignment of button text when button is larger than text.
	s.DisplayWindowPadding = f64.Vec2{22, 22} // Window positions are clamped to be visible within the display area by at least this amount. Only covers regular windows.
	s.DisplaySafeAreaPadding = f64.Vec2{4, 4} // If you cannot see the edge of your screen (e.g. on a TV) increase the safe area padding. Covers popups/tooltips as well regular windows.
	s.MouseCursorScale = 1.0                  // Scale software rendered mouse cursor (when io.MouseDrawCursor is enabled). May be removed later.
	s.AntiAliasedLines = true                 // Enable anti-aliasing on lines/borders. Disable if you are really short on CPU/GPU.
	s.AntiAliasedFill = true                  // Enable anti-aliasing on filled shapes (rounded rectangles, circles, etc.)
	s.CurveTessellationTol = 1.25             // Tessellation tolerance when using PathBezierCurveTo() without a specific number of segments. Decrease for highly tessellated curves (higher quality, more polygons), increase to reduce quality.
}

func (c *Context) PushStyleVar(idx StyleVar, val interface{}) {
	c.StyleModifiers = append(c.StyleModifiers, StyleMod{idx, val})
}

func (c *Context) PopStyleVar() {
	c.PopStyleVarN(1)
}

func (c *Context) PopStyleVarN(count int) {
	style := &c.Style
	for ; count > 0; count-- {
		n := len(c.StyleModifiers) - 1
		m := c.StyleModifiers[n]
		c.StyleModifiers = c.StyleModifiers[:n]
		switch m.VarIdx {
		case StyleVarAlpha:
			style.Alpha = m.Value.(float64)
		case StyleVarWindowPadding:
			style.WindowPadding = m.Value.(f64.Vec2)
		case StyleVarWindowRounding:
			style.WindowRounding = m.Value.(float64)
		case StyleVarWindowBorderSize:
			style.WindowBorderSize = m.Value.(float64)
		case StyleVarWindowMinSize:
			style.WindowMinSize = m.Value.(f64.Vec2)
		case StyleVarWindowTitleAlign:
			style.WindowTitleAlign = m.Value.(f64.Vec2)
		case StyleVarChildRounding:
			style.ChildRounding = m.Value.(float64)
		case StyleVarChildBorderSize:
			style.ChildBorderSize = m.Value.(float64)
		case StyleVarPopupRounding:
			style.PopupRounding = m.Value.(float64)
		case StyleVarPopupBorderSize:
			style.PopupBorderSize = m.Value.(float64)
		case StyleVarFramePadding:
			style.FramePadding = m.Value.(f64.Vec2)
		case StyleVarFrameRounding:
			style.FrameRounding = m.Value.(float64)
		case StyleVarFrameBorderSize:
			style.FrameBorderSize = m.Value.(float64)
		case StyleVarItemSpacing:
			style.ItemSpacing = m.Value.(f64.Vec2)
		case StyleVarItemInnerSpacing:
			style.ItemInnerSpacing = m.Value.(f64.Vec2)
		case StyleVarIndentSpacing:
			style.IndentSpacing = m.Value.(float64)
		case StyleVarScrollbarSize:
			style.ScrollbarSize = m.Value.(float64)
		case StyleVarScrollbarRounding:
			style.ScrollbarRounding = m.Value.(float64)
		case StyleVarGrabMinSize:
			style.GrabMinSize = m.Value.(float64)
		case StyleVarGrabRounding:
			style.GrabRounding = m.Value.(float64)
		case StyleVarButtonTextAlign:
			style.ButtonTextAlign = m.Value.(f64.Vec2)
		}
	}
}

func (c *Context) PushStyleColor(idx Col, col f64.Vec4) {
	c.ColorModifiers = append(c.ColorModifiers, ColMod{idx, c.Style.Colors[idx]})
	c.Style.Colors[idx] = col
}

func (c *Context) PushStyleColorRGBA(idx Col, col color.RGBA) {
	c.ColorModifiers = append(c.ColorModifiers, ColMod{idx, c.Style.Colors[idx]})
	c.Style.Colors[idx] = chroma.RGBA2VEC4(col)
}

func (c *Context) PopStyleColor() {
	c.PopStyleColorN(1)
}

func (c *Context) PopStyleColorN(count int) {
	for ; count > 0; count-- {
		backup := c.ColorModifiers[len(c.ColorModifiers)-1]
		c.Style.Colors[backup.Col] = backup.BackupValue
		c.ColorModifiers = c.ColorModifiers[:len(c.ColorModifiers)-1]
	}
}
