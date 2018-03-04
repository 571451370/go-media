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
	Ctx                            *Context
	Name                           string
	Id                             ID          // == ImHash(Name)
	Flags                          WindowFlags // See enum ImGuiWindowFlags_
	PosFloat                       f64.Vec2
	Pos                            f64.Vec2      // Position rounded-up to nearest pixel
	Size                           f64.Vec2      // Current size (==SizeFull or collapsed title bar size)
	SizeFull                       f64.Vec2      // Size when non collapsed
	SizeFullAtLastBegin            f64.Vec2      // Copy of SizeFull at the end of Begin. This is the reference value we'll use on the next frame to decide if we need scrollbars.
	SizeContents                   f64.Vec2      // Size of contents (== extents reach of the drawing cursor) from previous frame. Include decoration, window title, border, menu, etc.
	SizeContentsExplicit           f64.Vec2      // Size of contents explicitly set by the user via SetNextWindowContentSize()
	ContentsRegionRect             f64.Rectangle // Maximum visible content position in window coordinates. ~~ (SizeContentsExplicit ? SizeContentsExplicit : Size - ScrollbarSizes) - CursorStartPos, per axis
	WindowPadding                  f64.Vec2      // Window padding at the time of begin.
	WindowRounding                 float64       // Window rounding at the time of begin.
	WindowBorderSize               float64       // Window border size at the time of begin.
	MoveId                         ID            // == window->GetID("#MOVE")
	ChildId                        ID            // Id of corresponding item in parent window (for child windows)
	Scroll                         f64.Vec2
	ScrollTarget                   f64.Vec2 // target scroll position. stored as cursor position with scrolling canceled out, so the highest point is always 0.0f. (FLT_MAX for no change)
	ScrollTargetCenterRatio        f64.Vec2 // 0.0f = scroll so that target position is at top, 0.5f = scroll so that target position is centered
	ScrollbarX, ScrollbarY         bool
	ScrollbarSizes                 f64.Vec2
	Active                         bool // Set to true on Begin(), unless Collapsed
	WasActive                      bool
	WriteAccessed                  bool // Set to true when any widget access the current window
	Collapsed                      bool // Set when collapsing window to become only title-bar
	CollapseToggleWanted           bool
	SkipItems                      bool // Set when items can safely be all clipped (e.g. window not visible or collapsed)
	Appearing                      bool // Set during the frame where the window is appearing (or re-appearing)
	CloseButton                    bool // Set when the window has a close button (p_open != NULL)
	BeginOrderWithinParent         int  // Order within immediate parent window, if we are a child window. Otherwise 0.
	BeginOrderWithinContext        int  // Order within entire imgui context. This is mostly used for debugging submission order related issues.
	BeginCount                     int  // Number of Begin() during the current frame (generally 0 or 1, 1+ if appending via multiple Begin/End pairs)
	PopupId                        ID   // ID in the popup stack when this window is used as a popup/menu (because we use generic Name/ID for recycling)
	AutoFitFramesX, AutoFitFramesY int
	AutoFitOnlyGrows               bool
	AutoFitChildAxises             int
	AutoPosLastDirection           Dir
	HiddenFrames                   int
	SetWindowPosAllowFlags         Cond     // store condition flags for next SetWindowPos() call.
	SetWindowSizeAllowFlags        Cond     // store condition flags for next SetWindowSize() call.
	SetWindowCollapsedAllowFlags   Cond     // store condition flags for next SetWindowCollapsed() call.
	SetWindowPosVal                f64.Vec2 // store window position when using a non-zero Pivot (position set needs to be processed when we know the window size)
	SetWindowPosPivot              f64.Vec2 // store window pivot for positioning. ImVec2(0,0) when positioning from top-left corner; ImVec2(0.5f,0.5f) for centering; ImVec2(1,1) for bottom right.

	DC                             DrawContext
	IdStack                        []ID          // ID stack. ID are hashes seeded with the value at the top of the stack
	ClipRect                       f64.Rectangle // = DrawList->clip_rect_stack.back(). Scissoring / clipping rectangle. x1, y1, x2, y2.
	WindowRectClipped              f64.Rectangle // = WindowRect just after setup in Begin(). == window->Rect() for root window.
	InnerRect                      f64.Rectangle
	LastFrameActive                int
	ItemWidthDefault               float64
	MenuColumns                    MenuColumns // Simplified columns storage for menu items
	StateStorage                   Storage
	ColumnsStorage                 []ColumnsSet
	FontWindowScale                float64 // Scale multiplier per-window
	DrawList                       *DrawList
	ParentWindow                   *Window // If we are a child _or_ popup window, this is pointing to our parent. Otherwise NULL.
	RootWindow                     *Window // Point to ourself or first ancestor that is not a child window.
	RootWindowForTitleBarHighlight *Window // Point to ourself or first ancestor which will display TitleBgActive color when this window is active.
	RootWindowForTabbing           *Window // Point to ourself or first ancestor which can be CTRL-Tabbed into.
	RootWindowForNav               *Window // Point to ourself or first ancestor which doesn't have the NavFlattened flag.
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

type Cond int

type MenuColumns int

type ColumnsSet int

type Storage int

func (c *Context) ItemAdd(bb f64.Rectangle, id ID, navBB *f64.Rectangle) bool {
	return false
}

func (c *Context) ItemSize(size f64.Vec2, textOffsetY float64) {
	window := c.CurrentWindow
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
	window := c.CurrentWindow
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
	window := c.CurrentWindow
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

// In window space (not screen space!)
func (c *Context) GetContentRegionMax() f64.Vec2 {
	return f64.Vec2{}
}

func (c *Context) CalcItemSize(size f64.Vec2, defaultX, defaultY float64) f64.Vec2 {
	var contentMax f64.Vec2
	window := c.CurrentWindow
	dc := &window.DC

	if size.X < 0 || size.Y < 0 {
		contentMax = window.Pos.Add(c.GetContentRegionMax())
	}

	if size.X == 0 {
		size.X = defaultX
	} else {
		size.X += math.Max(contentMax.X-dc.CursorPos.X, 4)
	}

	if size.Y == 0 {
		size.Y = defaultY
	} else {
		size.Y += math.Max(contentMax.Y-dc.CursorPos.Y, 4)
	}

	return size
}
