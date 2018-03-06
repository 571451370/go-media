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

	NavLastChildNavWindow *Window          // When going to the menu bar, we remember the child window we came from. (This could probably be made implicit if we kept g.Windows sorted by last focused including child window.)
	NavLastIds            [2]ID            // Last known NavId for this window, per layer (0/1)
	NavRectRel            [2]f64.Rectangle // Reference rectangle, in window relative space

	// Navigation / Focus
	// FIXME-NAV: Merge all this with the new Nav system, at least the request variables should be moved to ImGuiContext
	FocusIdxAllCounter        int // Start at -1 and increase as assigned via FocusItemRegister()
	FocusIdxTabCounter        int // (same, but only count widgets which you can Tab through)
	FocusIdxAllRequestCurrent int // Item being requested for focus
	FocusIdxTabRequestCurrent int // Tab-able item being requested for focus
	FocusIdxAllRequestNext    int // Item being requested for focus, for next update (relies on layout to be stable between the frame pressing TAB and the next frame)
	FocusIdxTabRequestNext    int // "
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
	ItemFlags                 ItemFlags
	ChildWindows              []*Window
	LayoutType                LayoutType
	ColumnsSet                *ColumnsSet // Current columns set
	LastItemId                ID
	LastItemStatusFlags       ItemStatusFlags
	LastItemRect              f64.Rectangle // Interaction rect
	LastItemDisplayRect       f64.Rectangle // End-user display rect (only valid if LastItemStatusFlags & ImGuiItemStatusFlags_HasDisplayRect)
	NavHideHighlightOneFrame  bool
	NavHasScroll              bool // Set when scrolling can be used (ScrollMax > 0.0f)
	NavLayerCurrent           int  // Current layer, 0..31 (we currently only use 0..1)
	NavLayerCurrentMask       int  // = (1 << NavLayerCurrent) used by ItemAdd prior to clipping.
	NavLayerActiveMask        int  // Which layer have been written to (result from previous frame)
	NavLayerActiveMaskNext    int  // Which layer have been written to (buffer for current frame)

}

type ItemFlags int

const (
	ItemFlagsAllowKeyboardFocus       ItemFlags = 1 << iota // true
	ItemFlagsButtonRepeat                                   // false    // Button() will return true multiple times based on io.KeyRepeatDelay and io.KeyRepeatRate settings.
	ItemFlagsDisabled                                       // false    // FIXME-WIP: Disable interactions but doesn't affect visuals. Should be: grey out and disable interactions with widgets that affect data + view widgets (WIP)
	ItemFlagsNoNav                                          // false
	ItemFlagsNoNavDefaultFocus                              // false
	ItemFlagsSelectableDontClosePopup                       // false    // MenuItem/Selectable() automatically closes current Popup window
	ItemFlagsDefault                  = ItemFlagsAllowKeyboardFocus
)

type DragDropFlags uint

const (
	// BeginDragDropSource() flags
	DragDropFlagsSourceNoPreviewTooltip   DragDropFlags = 1 << iota // By default, a successful call to BeginDragDropSource opens a tooltip so you can display a preview or description of the source contents. This flag disable this behavior.
	DragDropFlagsSourceNoDisableHover                               // By default, when dragging we clear data so that IsItemHovered() will return true, to avoid subsequent user code submitting tooltips. This flag disable this behavior so you can still call IsItemHovered() on the source item.
	DragDropFlagsSourceNoHoldToOpenOthers                           // Disable the behavior that allows to open tree nodes and collapsing header by holding over them while dragging a source item.
	DragDropFlagsSourceAllowNullID                                  // Allow items such as Text(), Image() that have no unique identifier to be used as drag source, by manufacturing a temporary identifier based on their window-relative position. This is extremely unusual within the dear imgui ecosystem and so we made it explicit.
	DragDropFlagsSourceExtern                                       // External source (from outside of imgui), won't attempt to read current item/window info. Will always return true. Only one Extern source can be active simultaneously.

	// AcceptDragDropPayload() flags
	DragDropFlagsAcceptBeforeDelivery                                                                               // AcceptDragDropPayload() will returns true even before the mouse button is released. You can then call IsDelivery() to test if the payload needs to be delivered.
	DragDropFlagsAcceptNoDrawDefaultRect                                                                            // Do not draw the default highlight rectangle when hovering over target.
	DragDropFlagsAcceptPeekOnly          = DragDropFlagsAcceptBeforeDelivery | DragDropFlagsAcceptNoDrawDefaultRect // For peeking ahead and inspecting the payload before delivery.
)

type Cond int

type MenuColumns int

type ColumnsFlags int

const (
	// Default: 0
	ColumnsFlagsNoBorder               ColumnsFlags = 1 << iota // Disable column dividers
	ColumnsFlagsNoResize                                        // Disable resizing columns when clicking on the dividers
	ColumnsFlagsNoPreserveWidths                                // Disable column width preservation when adjusting columns
	ColumnsFlagsNoForceWithinWindow                             // Disable forcing columns to fit within window
	ColumnsFlagsGrowParentContentsSize                          // (WIP) Restore pre-1.51 behavior of extending the parent window contents size but _without affecting the columns width at all_. Will eventually remove.
)

type ColumnData struct {
	OffsetNorm             float64 // Column start offset, normalized 0.0 (far left) -> 1.0 (far right)
	OffsetNormBeforeResize float64
	Flags                  ColumnsFlags // Not exposed
	ClipRect               f64.Rectangle
}

type ColumnsSet struct {
	ID                 ID
	Flags              ColumnsFlags
	IsFirstFrame       bool
	IsBeingResized     bool
	Current            int
	Count              int
	MinX, MaxX         float64
	StartPosY          float64
	StartMaxPosX       float64 // Backup of CursorMaxPos
	CellMinY, CellMaxY float64
	Columns            []ColumnData
}

type Storage int

type HoveredFlags int

const (
	HoveredFlagsDefault                 HoveredFlags = 0               // Return true if directly over the item/window not obstructed by another window not obstructed by an active popup or modal blocking inputs under them.
	HoveredFlagsChildWindows            HoveredFlags = 1 << (iota - 1) // IsWindowHovered() only: Return true if any children of the window is hovered
	HoveredFlagsRootWindow                                             // IsWindowHovered() only: Test from root window (top most parent of the current hierarchy)
	HoveredFlagsAnyWindow                                              // IsWindowHovered() only: Return true if any window is hovered
	HoveredFlagsAllowWhenBlockedByPopup                                // Return true even if a popup window is normally blocking access to this item/window
	//HoveredFlagsAllowWhenBlockedByModal        // Return true even if a modal popup window is normally blocking access to this item/window. FIXME-TODO: Unavailable yet.
	HoveredFlagsAllowWhenBlockedByActiveItem              // Return true even if an active item is blocking access to this item/window. Useful for Drag and Drop patterns.
	HoveredFlagsAllowWhenOverlapped                       // Return true even if the position is overlapped by another window
	HoveredFlagsRectOnly                     HoveredFlags = HoveredFlagsAllowWhenBlockedByPopup | HoveredFlagsAllowWhenBlockedByActiveItem | HoveredFlagsAllowWhenOverlapped
	HoveredFlagsRootAndChildWindows          HoveredFlags = HoveredFlagsRootWindow | HoveredFlagsChildWindows
)

type ItemStatusFlags uint

const (
	ItemStatusFlagsHoveredRect ItemStatusFlags = 1 << iota
	ItemStatusFlagsHasDisplayRect
)

func (c *Context) ItemAdd(bb f64.Rectangle, id ID) bool {
	return c.ItemAddEx(bb, id, nil)
}

// Declare item bounding box for clipping and interaction.
// Note that the size can be different than the one provided to ItemSize(). Typically, widgets that spread over available surface
// declare their minimum size requirement to ItemSize() and then use a larger region for drawing/interaction, which is passed to ItemAdd().
func (c *Context) ItemAddEx(bb f64.Rectangle, id ID, navBB *f64.Rectangle) bool {
	window := c.CurrentWindow

	if id != 0 {
		// Navigation processing runs prior to clipping early-out
		//  (a) So that NavInitRequest can be honored, for newly opened windows to select a default widget
		//  (b) So that we can scroll up/down past clipped items. This adds a small O(N) cost to regular navigation requests unfortunately, but it is still limited to one window.
		//      it may not scale very well for windows with ten of thousands of item, but at least NavMoveRequest is only set on user interaction, aka maximum once a frame.
		//      We could early out with "if (is_clipped && !g.NavInitRequest) return false;" but when we wouldn't be able to reach unclipped widgets. This would work if user had explicit scrolling control (e.g. mapped on a stick)
		window.DC.NavLayerActiveMaskNext |= window.DC.NavLayerCurrentMask
		if c.NavId == id || c.NavAnyRequest {
			if c.NavWindow.RootWindowForNav == window.RootWindowForNav {
				if window == c.NavWindow || (window.Flags|c.NavWindow.Flags)&WindowFlagsNavFlattened != 0 {
					if navBB != nil {
						c.NavProcessItem(window, *navBB, id)
					} else {
						c.NavProcessItem(window, bb, id)
					}
				}
			}
		}
	}

	window.DC.LastItemId = id
	window.DC.LastItemRect = bb
	window.DC.LastItemStatusFlags = 0

	// Clipping test
	isClipped := c.IsClippedEx(bb, id, false)
	if isClipped {
		return false
	}

	// We need to calculate this now to take account of the current clipping rectangle (as items like Selectable may change them)
	if c.IsMouseHoveringRect(bb.Min, bb.Max) {
		window.DC.LastItemStatusFlags |= ItemStatusFlagsHoveredRect
	}

	return true
}

func (c *Context) ItemSize(size f64.Vec2) {
	c.ItemSizeEx(size, 0)
}

func (c *Context) ItemSizeEx(size f64.Vec2, textOffsetY float64) {
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

func (c *Context) ItemSizeBB(bb f64.Rectangle) {
	c.ItemSizeBBEx(bb, 0)
}

func (c *Context) ItemSizeBBEx(bb f64.Rectangle, textOffsetY float64) {
	c.ItemSizeEx(bb.Size(), textOffsetY)
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
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}
	dc := &window.DC

	backupLayoutType := dc.LayoutType
	dc.LayoutType = LayoutTypeVertical
	// In the event that we are on a line with items that is smaller that FontSize high, we will preserve its height.
	if dc.CurrentLineHeight > 0 {
		c.ItemSize(f64.Vec2{0, 0})
	} else {
		c.ItemSize(f64.Vec2{0, c.FontSize})
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
	window := c.GetCurrentWindowRead()
	dc := &window.DC
	mx := window.ContentsRegionRect.Max
	if dc.ColumnsSet != nil {
		mx.X = c.GetColumnOffset(dc.ColumnsSet.Current+1) - window.WindowPadding.X
	}
	return mx
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

func (c *Context) ItemHoverable(bb f64.Rectangle, id ID) bool {
	if c.HoveredId != 0 && c.HoveredId == id && !c.HoveredIdAllowOverlap {
		return false
	}

	window := c.CurrentWindow
	if c.HoveredWindow != window {
		return false
	}

	if c.ActiveId != 0 && c.ActiveId != id && !c.ActiveIdAllowOverlap {
		return false
	}

	if !c.IsMouseHoveringRect(bb.Min, bb.Max) {
		return false
	}

	if c.NavDisableMouseHover || !c.IsWindowContentHoverable(window, HoveredFlagsDefault) {
		return false
	}

	if window.DC.ItemFlags&ItemFlagsDisabled != 0 {
		return false
	}

	c.SetHoveredID(id)
	return true
}

func (c *Context) IsMouseHoveringRect(rmin, rmax f64.Vec2) bool {
	return c.IsMouseHoveringRectEx(rmin, rmax, true)
}

func (c *Context) IsMouseHoveringRectEx(rmin, rmax f64.Vec2, clip bool) bool {
	window := c.CurrentWindow
	io := &c.IO
	style := &c.Style

	// Clip
	rectClipped := f64.Rectangle{rmin, rmax}
	if clip {
		rectClipped.Intersect(window.ClipRect)
	}

	// Expand for touch input
	rectForTouch := f64.Rectangle{
		rectClipped.Min.Sub(style.TouchExtraPadding),
		rectClipped.Max.Add(style.TouchExtraPadding),
	}

	return io.MousePos.In(rectForTouch)
}

func (c *Context) IsWindowContentHoverable(window *Window, flags HoveredFlags) bool {
	// An active popup disable hovering on other windows (apart from its own children)
	// FIXME-OPT: This could be cached/stored within the window.
	if c.NavWindow != nil {
		focusedRootWindow := c.NavWindow.RootWindow
		if focusedRootWindow != nil {
			if focusedRootWindow.WasActive && focusedRootWindow != window.RootWindow {
				// For the purpose of those flags we differentiate "standard popup" from "modal popup"
				// NB: The order of those two tests is important because Modal windows are also Popups.
				if focusedRootWindow.Flags&WindowFlagsModal != 0 {
					return false
				}

				if focusedRootWindow.Flags&WindowFlagsPopup != 0 && flags&HoveredFlagsAllowWhenBlockedByPopup == 0 {
					return false
				}
			}
		}
	}

	return true
}

func (c *Context) GetColumnOffset(columnIndex int) float64 {
	window := c.GetCurrentWindowRead()
	dc := &window.DC
	columns := dc.ColumnsSet

	if columnIndex < 0 {
		columnIndex = columns.Current
	}

	t := columns.Columns[columnIndex].OffsetNorm
	xOffset := f64.Lerp(columns.MinX, columns.MaxX, t)
	return xOffset
}

func (c *Context) IsClippedEx(bb f64.Rectangle, id ID, clipEvenWhenLogged bool) bool {
	window := c.CurrentWindow
	if !bb.Overlaps(window.ClipRect) {
		if id == 0 || id != c.ActiveId {
			if clipEvenWhenLogged || !c.LogEnabled {
				return true
			}
		}
	}
	return false
}