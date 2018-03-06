package imgui

import "github.com/qeedquan/go-media/math/f64"

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
	Backup [2]interface{}
}

type StyleVar int

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

type ColorEditFlags int

const (
	ColorEditFlagsNoAlpha        ColorEditFlags = 1 << 1 //              // ColorEdit ColorPicker ColorButton: ignore Alpha component (read 3 components from the input pointer).
	ColorEditFlagsNoPicker       ColorEditFlags = 1 << 2 //              // ColorEdit: disable picker when clicking on colored square.
	ColorEditFlagsNoOptions      ColorEditFlags = 1 << 3 //              // ColorEdit: disable toggling options menu when right-clicking on inputs/small preview.
	ColorEditFlagsNoSmallPreview ColorEditFlags = 1 << 4 //              // ColorEdit ColorPicker: disable colored square preview next to the inputs. (e.g. to show only the inputs)
	ColorEditFlagsNoInputs       ColorEditFlags = 1 << 5 //              // ColorEdit ColorPicker: disable inputs sliders/text widgets (e.g. to show only the small preview colored square).
	ColorEditFlagsNoTooltip      ColorEditFlags = 1 << 6 //              // ColorEdit ColorPicker ColorButton: disable tooltip when hovering the preview.
	ColorEditFlagsNoLabel        ColorEditFlags = 1 << 7 //              // ColorEdit ColorPicker: disable display of inline text label (the label is still forwarded to the tooltip and picker).
	ColorEditFlagsNoSidePreview  ColorEditFlags = 1 << 8 //              // ColorPicker: disable bigger color preview on right side of the picker use small colored square preview instead.
	// User Options (right-click on widget to change some of them). You can set application defaults using SetColorEditOptions(). The idea is that you probably don't want to override them in most of your calls let the user choose and/or call SetColorEditOptions() during startup.
	ColorEditFlagsAlphaBar                       ColorEditFlags = 1 << 9  //              // ColorEdit ColorPicker: show vertical alpha bar/gradient in picker.
	ColorEditFlagsAlphaPreview                   ColorEditFlags = 1 << 10 //              // ColorEdit ColorPicker ColorButton: display preview as a transparent color over a checkerboard instead of opaque.
	ColorEditFlagsAlphaPreviewHalfColorEditFlags                = 1 << 11 //              // ColorEdit ColorPicker ColorButton: display half opaque / half checkerboard instead of opaque.
	ColorEditFlagsHDR                            ColorEditFlags = 1 << 12 //              // (WIP) ColorEdit: Currently only disable 0.0f..1.0f limits in RGBA edition (note: you probably want to use ColorEditFlagsFloat flag as well).
	ColorEditFlagsRGB                            ColorEditFlags = 1 << 13 // [Inputs]     // ColorEdit: choose one among RGB/HSV/HEX. ColorPicker: choose any combination using RGB/HSV/HEX.
	ColorEditFlagsHSV                            ColorEditFlags = 1 << 14 // [Inputs]     // "
	ColorEditFlagsHEX                            ColorEditFlags = 1 << 15 // [Inputs]     // "
	ColorEditFlagsUint8                          ColorEditFlags = 1 << 16 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0..255.
	ColorEditFlagsFloat                          ColorEditFlags = 1 << 17 // [DataType]   // ColorEdit ColorPicker ColorButton: _display_ values formatted as 0.0f..1.0f floats instead of 0..255 integers. No round-trip of value via integers.
	ColorEditFlagsPickerHueBar                   ColorEditFlags = 1 << 18 // [PickerMode] // ColorPicker: bar for Hue rectangle for Sat/Value.
	ColorEditFlagsPickerHueWheel                 ColorEditFlags = 1 << 19 // [PickerMode] // ColorPicker: wheel for Hue triangle for Sat/Value.
	// Internals/Masks
	ColorEditFlags_InputsMask     ColorEditFlags = ColorEditFlagsRGB | ColorEditFlagsHSV | ColorEditFlagsHEX
	ColorEditFlags_DataTypeMask   ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsFloat
	ColorEditFlags_PickerMask     ColorEditFlags = ColorEditFlagsPickerHueWheel | ColorEditFlagsPickerHueBar
	ColorEditFlags_OptionsDefault ColorEditFlags = ColorEditFlagsUint8 | ColorEditFlagsRGB | ColorEditFlagsPickerHueBar // Change application default using SetColorEditOptions()
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
