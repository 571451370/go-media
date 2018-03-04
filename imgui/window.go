package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type LayoutType uint

const (
	LayoutTypeVertical LayoutType = iota
	LayoutTypeHorizontal
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
	WindowFlagsHorizontalScrollbar                               // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing                                // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus                             // Disable bringing window to front when taking focus (e.g. clicking on it or programatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar                           // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar                         // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding                            // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows because more convenient)
	WindowFlagsResizeFromAnySide                                 // (WIP) Enable resize from any corners and borders. Your back-end needs to honor the different values of io.MouseCursor set by imgui.
	WindowFlagsNoNavInputs                                       // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                                        // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	WindowFlagsNoNav                     WindowFlags = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus

	// [Internal]
	WindowFlagsNavFlattened WindowFlags = 1 << iota // (WIP) Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	WindowFlagsChildWindow                          // Don't use! For internal use by BeginChild()
	WindowFlagsTooltip                              // Don't use! For internal use by BeginTooltip()
	WindowFlagsPopup                                // Don't use! For internal use by BeginPopup()
	WindowFlagsModal                                // Don't use! For internal use by BeginPopupModal()
	WindowFlagsChildMenu                            // Don't use! For internal use by BeginMenu()
)

type Window struct {
	Ctx                 *Context
	Name                string
	Id                  ID
	Flags               WindowFlags
	PosFloat            f64.Vec2
	Pos                 f64.Vec2
	Size                f64.Vec2
	SizeFull            f64.Vec2
	SizeFullAtLastBegin f64.Vec2
	SizeContents        f64.Vec2
	Scroll              f64.Vec2
	DC                  DrawContext
	IdStack             []ID // ID stack. ID are hashes seeded with the value at the top of the stack
	Active              bool
	WasActive           bool
	HiddenFrames        int
	SkipItems           bool
	DrawList            *DrawList
}

type DrawContext struct {
	CursorPos                 f64.Vec2
	CursorMaxPos              f64.Vec2
	CursorPosPrevLine         f64.Vec2
	PrevLineHeight            float64
	PrevLineTextBaseOffset    float64
	CurrentLineTextBaseOffset float64
	CurrentLineHeight         float64
	IndentX                   float64
	GroupOffsetX              float64
	ColumnsOffsetX            float64
	ChildWindows              []*Window
	LayoutType                LayoutType
}

func (c *Context) ItemAdd(bb f64.Rectangle, id ID, navBB *f64.Rectangle) bool {
	return false
}

func (c *Context) ItemSize(size f64.Vec2, textOffsetY float64) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	style := c.GetStyle()
	dc := &window.DC

	// always align ourselves on pixel boundaries
	lineHeight := math.Max(dc.CurrentLineHeight, size.Y)
	textBaseOffset := math.Max(dc.CurrentLineTextBaseOffset, textOffsetY)

	dc.CursorPosPrevLine = f64.Vec2{dc.CursorPos.X + size.X, dc.CursorPos.Y}
	dc.CursorPos = f64.Vec2{
		window.Pos.X + dc.IndentX + dc.ColumnsOffsetX,
		dc.CursorPos.Y + lineHeight + style.ItemSpacing.Y,
	}
	dc.CursorMaxPos.X = math.Max(dc.CursorMaxPos.X, dc.CursorPosPrevLine.X)
	dc.CursorMaxPos.Y = math.Max(dc.CursorMaxPos.Y, dc.CursorPos.Y-style.ItemSpacing.Y)

	dc.PrevLineHeight = lineHeight
	dc.PrevLineTextBaseOffset = textBaseOffset
	dc.CurrentLineHeight = 0
	dc.CurrentLineTextBaseOffset = 0

	// horizontal layout mode
	if dc.LayoutType == LayoutTypeHorizontal {
		c.SameLine(0, -1)
	}
}

// Gets back to previous line and continue with horizontal layout
//      pos_x == 0      : follow right after previous item
//      pos_x != 0      : align to specified x position (relative to window/group left)
//      spacing_w < 0   : use default spacing if pos_x == 0, no spacing if pos_x != 0
//      spacing_w >= 0  : enforce spacing amount
func (c *Context) SameLine(posX, spacingW float64) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	style := c.GetStyle()
	dc := &window.DC

	if posX != 0 {
		spacingW = math.Max(spacingW, 0)
		dc.CursorPos.X = window.Pos.X - window.Scroll.X + posX + spacingW + dc.GroupOffsetX + dc.ColumnsOffsetX
		dc.CursorPos.Y = dc.CursorPosPrevLine.Y
	} else {
		spacingW = math.Max(spacingW, style.ItemSpacing.X)
		dc.CursorPos.X = dc.CursorPosPrevLine.X + spacingW
		dc.CursorPos.Y = dc.CursorPosPrevLine.Y
	}
	dc.CurrentLineHeight = dc.PrevLineHeight
	dc.CurrentLineTextBaseOffset = dc.PrevLineTextBaseOffset
}

func (c *Context) NewLine() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	dc := &window.DC

	backupLayoutType := dc.LayoutType
	dc.LayoutType = LayoutTypeVertical
	// In the event that we are on a line with items that is smaller that FontSize high, we will preserve its height.
	if dc.CurrentLineHeight > 0 {
		c.ItemSize(f64.Vec2{}, 0)
	} else {
		c.ItemSize(f64.Vec2{}, c.FontSize)
	}
	dc.LayoutType = backupLayoutType
}

func (w *Window) CalcFontSize() float64 {
	return 0
}

func (w *Window) GetID(str string) ID {
	ctx := w.Ctx
	id := ID(hash(str))
	ctx.KeepAliveID(id)
	return id
}

func (c *Context) CalcItemSize(size f64.Vec2, defaultX, defaultY float64) f64.Vec2 {
	return f64.Vec2{}
}
