package imgui

import "github.com/qeedquan/go-media/math/f64"

type Context struct {
}

func New() *Context {
	return &Context{}
}

func (c *Context) NewFrame() {
}

func (c *Context) Render() {
}

func (c *Context) EndFrame() {
}

// Begin push window to the stack and start appending to it. return false when window is collapsed
// (so you can early out in your code) but you always need to call End() regardless.
// 'open' creates a widget on the upper-right to close the window (which sets your bool to false).
func (c *Context) Begin(name string, open bool, flags int) bool {
	return true
}

// End pops window off the stack, always call even if Begin() return false (which indicates a collapsed window)!
func (c *Context) End() {
}

// GetWindowSize gets the current window size
func (c *Context) GetWindowSize() f64.Vec2 {
	return f64.Vec2{}
}

// GetWindowPos gets the current window position in screen space (useful if you want to do your own drawing via the DrawList api)
func (c *Context) GetWindowPos() f64.Vec2 {
	return f64.Vec2{}
}

// GetScroll gets  the scrolling amount [0..GetScrollMax()]
func (c *Context) GetScroll() f64.Vec2 {
	return f64.Vec2{}
}

// GetScrollMax get the scrolling max
func (c *Context) GetScrollMax() f64.Vec2 {
	return f64.Vec2{}
}

// Button draws a button
func (c *Context) Button(label string, size f64.Vec2) bool {
	return true
}

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

type WindowFlags uint

const (
	WindowFlagsNoTitleBar                WindowFlags = 1 << iota // Disable title-bar
	WindowFlagsNoResize                                          // Disable user resizing with the lower-right grip
	WindowFlagsNoMove                                            // Disable user moving the window
	WindowFlagsNoScrollbar                                       // Disable scrollbars (window can still scroll with mouse or programatically)
	WindowFlagsNoScrollWithMouse                                 // Disable user vertically scrolling with mouse wheel. On child window mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	WindowFlagsNoCollapse                                        // Disable user collapsing window by double-clicking on it
	WindowFlagsAlwaysAutoResize                                  // Resize every window to its content every frame
	WindowFlagsNoSavedSettings                                   // Never load/save settings in .ini file
	WindowFlagsNoInputs                                          // Disable catching mouse or keyboard inputs hovering test with pass through.
	WindowFlagsMenuBar                                           // Has a menu-bar
	WindowFlagsHorizontalScrollbar                               // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width0.0f)); prior to calling Begin() to specify width. Read code in imguidemo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing                                // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus                             // Disable bringing window to front when taking focus (e.g. clicking on it or programatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar                           // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar                         // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding                            // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows because more convenient)
	WindowFlagsResizeFromAnySide                                 // (WIP) Enable resize from any corners and borders. Your back-end needs to honor the different values of io.MouseCursor set by imgui.
	WindowFlagsNoNavInputs                                       // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                                        // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)

	// [Internal]
	windowFlagsNavFlattened // (WIP) Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	windowFlagsChildWindow  // Don't use! For internal use by BeginChild()
	windowFlagsTooltip      // Don't use! For internal use by BeginTooltip()
	windowFlagsPopup        // Don't use! For internal use by BeginPopup()
	windowFlagsModal        // Don't use! For internal use by BeginPopupModal()
	windowFlagsChildMenu    // Don't use! For internal use by BeginMenu()

	WindowFlagsNoNav WindowFlags = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus
)

type InputTextFlags uint

const (
	InputTextFlagsCharsDecimal        InputTextFlags = 1 << iota // Allow 0123456789.+-*/
	InputTextFlagsCharsHexadecimal                               // Allow 0123456789ABCDEFabcdef
	InputTextFlagsCharsUppercase                                 // Turn a..z into A..Z
	InputTextFlagsCharsNoBlank                                   // Filter out spaces, tabs
	InputTextFlagsAutoSelectAll                                  // Select entire text when first taking mouse focus
	InputTextFlagsEnterReturnsTrue                               // Return 'true' when Enter is pressed (as opposed to when the value was modified)
	InputTextFlagsCallbackCompletion                             // Call user function on pressing TAB (for completion handling)
	InputTextFlagsCallbackHistory                                // Call user function on pressing Up/Down arrows (for history handling)
	InputTextFlagsCallbackAlways                                 // Call user function every time. User code may query cursor position, modify text buffer.
	InputTextFlagsCallbackCharFilter                             // Call user function to filter character. Modify data->EventChar to replace/filter input, or return 1 to discard character.
	InputTextFlagsAllowTabInput                                  // Pressing TAB input a '\t' character into the text field
	InputTextFlagsCtrlEnterForNewLine                            // In multi-line mode, unfocus with Enter, add new line with Ctrl+Enter (default is opposite: unfocus with Ctrl+Enter, add line with Enter).
	InputTextFlagsNoHorizontalScroll                             // Disable following the cursor horizontally
	InputTextFlagsAlwaysInsertMode                               // Insert mode
	InputTextFlagsReadOnly                                       // Read-only mode
	InputTextFlagsPassword                                       // Password mode, display all characters as '*'
	InputTextFlagsNoUndoRedo                                     // Disable undo/redo. Note that input text owns the text data while active, if you want to provide your own undo/redo stack you need e.g. to call ClearActiveID().

	//	[Internal]
	inputTextFlagsMultiline // For internal use by InputTextMultiline()

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
