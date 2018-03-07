package imgui

import (
	"hash/fnv"
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type WindowFlags int

const (
	WindowFlagsNoTitleBar        WindowFlags = 1 << 0 // Disable title-bar
	WindowFlagsNoResize          WindowFlags = 1 << 1 // Disable user resizing with the lower-right grip
	WindowFlagsNoMove            WindowFlags = 1 << 2 // Disable user moving the window
	WindowFlagsNoScrollbar       WindowFlags = 1 << 3 // Disable scrollbars (window can still scroll with mouse or programatically)
	WindowFlagsNoScrollWithMouse WindowFlags = 1 << 4 // Disable user vertically scrolling with mouse wheel. On child window mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	WindowFlagsNoCollapse        WindowFlags = 1 << 5 // Disable user collapsing window by double-clicking on it
	WindowFlagsAlwaysAutoResize  WindowFlags = 1 << 6 // Resize every window to its content every frame
	//WindowFlagsShowBorders   WindowFlags       = 1 << 7   // Show borders around windows and items (OBSOLETE! Use e.g. style.FrameBorderSize=1.0f to enable borders).
	WindowFlagsNoSavedSettings           WindowFlags = 1 << 8  // Never load/save settings in .ini file
	WindowFlagsNoInputs                  WindowFlags = 1 << 9  // Disable catching mouse or keyboard inputs hovering test with pass through.
	WindowFlagsMenuBar                   WindowFlags = 1 << 10 // Has a menu-bar
	WindowFlagsHorizontalScrollbar       WindowFlags = 1 << 11 // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	WindowFlagsNoFocusOnAppearing        WindowFlags = 1 << 12 // Disable taking focus when transitioning from hidden to visible state
	WindowFlagsNoBringToFrontOnFocus     WindowFlags = 1 << 13 // Disable bringing window to front when taking focus (e.g. clicking on it or programatically giving it focus)
	WindowFlagsAlwaysVerticalScrollbar   WindowFlags = 1 << 14 // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	WindowFlagsAlwaysHorizontalScrollbar WindowFlags = 1 << 15 // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	WindowFlagsAlwaysUseWindowPadding    WindowFlags = 1 << 16 // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows because more convenient)
	WindowFlagsResizeFromAnySide         WindowFlags = 1 << 17 // (WIP) Enable resize from any corners and borders. Your back-end needs to honor the different values of io.MouseCursor set by imgui.
	WindowFlagsNoNavInputs               WindowFlags = 1 << 18 // No gamepad/keyboard navigation within the window
	WindowFlagsNoNavFocus                WindowFlags = 1 << 19 // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	WindowFlagsNoNav                     WindowFlags = WindowFlagsNoNavInputs | WindowFlagsNoNavFocus

	// [Internal]
	WindowFlagsNavFlattened WindowFlags = 1 << 23 // (WIP) Allow gamepad/keyboard navigation to cross over parent border to this child (only use on child that have no scrolling!)
	WindowFlagsChildWindow  WindowFlags = 1 << 24 // Don't use! For internal use by BeginChild()
	WindowFlagsTooltip      WindowFlags = 1 << 25 // Don't use! For internal use by BeginTooltip()
	WindowFlagsPopup        WindowFlags = 1 << 26 // Don't use! For internal use by BeginPopup()
	WindowFlagsModal        WindowFlags = 1 << 27 // Don't use! For internal use by BeginPopupModal()
	WindowFlagsChildMenu    WindowFlags = 1 << 28 // Don't use! For internal use by BeginMenu()
)

type WindowSettings struct {
	Name      string
	Id        ID
	Pos       f64.Vec2
	Size      f64.Vec2
	Collapsed bool
}

type Window struct {
	Ctx                            *Context
	Name                           string
	ID                             ID          // == ImHash(Name)
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
	SetWindowPosAllowFlags         Cond          // store condition flags for next SetWindowPos() call.
	SetWindowSizeAllowFlags        Cond          // store condition flags for next SetWindowSize() call.
	SetWindowCollapsedAllowFlags   Cond          // store condition flags for next SetWindowCollapsed() call.
	SetWindowPosVal                f64.Vec2      // store window position when using a non-zero Pivot (position set needs to be processed when we know the window size)
	SetWindowPosPivot              f64.Vec2      // store window pivot for positioning. ImVec2(0,0) when positioning from top-left corner; ImVec2(0.5f,0.5f) for centering; ImVec2(1,1) for bottom right.
	DC                             DrawContext   // Temporary per-window data, reset at the beginning of the frame
	IDStack                        []ID          // ID stack. ID are hashes seeded with the value at the top of the stack
	ClipRect                       f64.Rectangle // = DrawList->clip_rect_stack.back(). Scissoring / clipping rectangle. x1, y1, x2, y2.
	WindowRectClipped              f64.Rectangle // = WindowRect just after setup in Begin(). == window->Rect() for root window.
	InnerRect, InnerClipRect       f64.Rectangle
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
	CursorPos                     f64.Vec2
	CursorPosPrevLine             f64.Vec2
	CursorStartPos                f64.Vec2
	CursorMaxPos                  f64.Vec2 // Used to implicitly calculate the size of our contents, always growing during the frame. Turned into window->SizeContents at the beginning of next frame
	CurrentLineHeight             float64
	CurrentLineTextBaseOffset     float64
	PrevLineHeight                float64
	PrevLineTextBaseOffset        float64
	LogLinePosY                   float64
	TreeDepth                     int
	TreeDepthMayJumpToParentOnPop uint32 // Store a copy of !g.NavIdIsAlive for TreeDepth 0..31
	LastItemId                    ID
	LastItemStatusFlags           ItemStatusFlags
	LastItemRect                  f64.Rectangle // Interaction rect
	LastItemDisplayRect           f64.Rectangle // End-user display rect (only valid if LastItemStatusFlags & ImGuiItemStatusFlags_HasDisplayRect)
	NavHideHighlightOneFrame      bool
	NavHasScroll                  bool // Set when scrolling can be used (ScrollMax > 0.0f)
	NavLayerCurrent               int  // Current layer, 0..31 (we currently only use 0..1)
	NavLayerCurrentMask           int  // = (1 << NavLayerCurrent) used by ItemAdd prior to clipping.
	NavLayerActiveMask            int  // Which layer have been written to (result from previous frame)
	NavLayerActiveMaskNext        int  // Which layer have been written to (buffer for current frame)
	MenuBarAppending              bool // FIXME: Remove this
	MenuBarOffsetX                float64
	ChildWindows                  []*Window
	StateStorage                  *Storage
	LayoutType                    LayoutType
	ParentLayoutType              LayoutType // Layout type of parent window at the time of Begin()

	// We store the current settings outside of the vectors to increase memory locality (reduce cache misses). The vectors are rarely modified. Also it allows us to not heap allocate for short-lived windows which are not using those settings.
	ItemFlags        ItemFlags // == ItemFlagsStack.back() [empty == ImGuiItemFlags_Default]
	ItemWidth        float64   // == ItemWidthStack.back(). 0.0: default, >0.0: width in pixels, <0.0: align xx pixels to the right of window
	TextWrapPos      float64   // == TextWrapPosStack.back() [empty == -1.0f]
	ItemFlagsStack   []ItemFlags
	ItemWidthStack   []float64
	TextWrapPosStack []float64
	GroupStack       []GroupData
	StackSizesBackup [6]int // Store size of various stacks for asserting

	IndentX        float64 // Indentation / start position from left of window (increased by TreePush/TreePop, etc.)
	GroupOffsetX   float64
	ColumnsOffsetX float64     // Offset to the current column (if ColumnsCurrent > 0). FIXME: This and the above should be a stack to allow use cases like Tree->Column->Tree. Need revamp columns API.
	ColumnsSet     *ColumnsSet // Current columns set
}

type GroupData struct {
	BackupCursorPos                 f64.Vec2
	BackupCursorMaxPos              f64.Vec2
	BackupIndentX                   float64
	BackupGroupOffsetX              float64
	BackupCurrentLineHeight         float64
	BackupCurrentLineTextBaseOffset float64
	BackupLogLinePosY               float64
	BackupActiveIdIsAlive           bool
	AdvanceCursor                   bool
}

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

type ColumnsFlags int

const (
	// Default: 0
	ColumnsFlagsNoBorder               ColumnsFlags = 1 << 0 // Disable column dividers
	ColumnsFlagsNoResize               ColumnsFlags = 1 << 1 // Disable resizing columns when clicking on the dividers
	ColumnsFlagsNoPreserveWidths       ColumnsFlags = 1 << 2 // Disable column width preservation when adjusting columns
	ColumnsFlagsNoForceWithinWindow    ColumnsFlags = 1 << 3 // Disable forcing columns to fit within window
	ColumnsFlagsGrowParentContentsSize ColumnsFlags = 1 << 4 // (WIP) Restore pre-1.51 behavior of extending the parent window contents size but _without affecting the columns width at all_. Will eventually remove.
)

type MenuColumns struct {
	Count            int
	Spacing          float64
	Width, NextWidth float64
	Pos, NextWidths  [4]float64
}

type ItemFlags int

const (
	ItemFlagsAllowKeyboardFocus       ItemFlags = 1 << 0 // true
	ItemFlagsButtonRepeat             ItemFlags = 1 << 1 // false    // Button() will return true multiple times based on io.KeyRepeatDelay and io.KeyRepeatRate settings.
	ItemFlagsDisabled                 ItemFlags = 1 << 2 // false    // FIXME-WIP: Disable interactions but doesn't affect visuals. Should be: grey out and disable interactions with widgets that affect data + view widgets (WIP)
	ItemFlagsNoNav                    ItemFlags = 1 << 3 // false
	ItemFlagsNoNavDefaultFocus        ItemFlags = 1 << 4 // false
	ItemFlagsSelectableDontClosePopup ItemFlags = 1 << 5 // false    // MenuItem/Selectable() automatically closes current Popup window
	ItemFlagsDefault_                 ItemFlags = ItemFlagsAllowKeyboardFocus
)

type LayoutType int

const (
	LayoutTypeVertical LayoutType = iota
	LayoutTypeHorizontal
)

type Axis int

const (
	AxisNone Axis = -1
	AxisX    Axis = 0
	AxisY    Axis = 1
)

type PopupRef struct {
	PopupId        ID       // Set on OpenPopup()
	Window         *Window  // Resolved on BeginPopup() - may stay unresolved if user never calls OpenPopup()
	ParentWindow   *Window  // Set on OpenPopup()
	OpenFrameCount int      // Set on OpenPopup()
	OpenParentId   ID       // Set on OpenPopup(), we need this to differenciate multiple menu sets from each others (e.g. inside menu bar vs loose menu items)
	OpenPopupPos   f64.Vec2 // Set on OpenPopup(), preferred popup position (typically == OpenMousePos when using mouse)
	OpenMousePos   f64.Vec2 // Set on OpenPopup(), copy of mouse position at the time of opening popup
}

type NextWindowData struct {
	PosCond              Cond
	SizeCond             Cond
	ContentSizeCond      Cond
	CollapsedCond        Cond
	SizeConstraintCond   Cond
	FocusCond            Cond
	BgAlphaCond          Cond
	PosVal               f64.Vec2
	PosPivotVal          f64.Vec2
	SizeVal              f64.Vec2
	ContentSizeVal       f64.Vec2
	CollapsedVal         bool
	SizeConstraintRect   f64.Rectangle // Valid if 'SetNextWindowSizeConstraint' is true
	SizeCallback         func()
	SizeCallbackUserData interface{}
	BgAlphaVal           float64
}

type DragDropFlags int

const (
	// BeginDragDropSource() flags
	DragDropFlagsSourceNoPreviewTooltip   DragDropFlags = 1 << 0 // By default a successful call to BeginDragDropSource opens a tooltip so you can display a preview or description of the source contents. This flag disable this behavior.
	DragDropFlagsSourceNoDisableHover     DragDropFlags = 1 << 1 // By default when dragging we clear data so that IsItemHovered() will return true to avoid subsequent user code submitting tooltips. This flag disable this behavior so you can still call IsItemHovered() on the source item.
	DragDropFlagsSourceNoHoldToOpenOthers DragDropFlags = 1 << 2 // Disable the behavior that allows to open tree nodes and collapsing header by holding over them while dragging a source item.
	DragDropFlagsSourceAllowNullID        DragDropFlags = 1 << 3 // Allow items such as Text() Image() that have no unique identifier to be used as drag source by manufacturing a temporary identifier based on their window-relative position. This is extremely unusual within the dear imgui ecosystem and so we made it explicit.
	DragDropFlagsSourceExtern             DragDropFlags = 1 << 4 // External source (from outside of imgui) won't attempt to read current item/window info. Will always return true. Only one Extern source can be active simultaneously.
	// AcceptDragDropPayload() flags
	DragDropFlagsAcceptBeforeDelivery    DragDropFlags = 1 << 10                                                                  // AcceptDragDropPayload() will returns true even before the mouse button is released. You can then call IsDelivery() to test if the payload needs to be delivered.
	DragDropFlagsAcceptNoDrawDefaultRect DragDropFlags = 1 << 11                                                                  // Do not draw the default highlight rectangle when hovering over target.
	DragDropFlagsAcceptPeekOnly          DragDropFlags = DragDropFlagsAcceptBeforeDelivery | DragDropFlagsAcceptNoDrawDefaultRect // For peeking ahead and inspecting the payload before delivery.
)

type HoveredFlags int

const (
	HoveredFlagsDefault                 HoveredFlags = 0      // Return true if directly over the item/window not obstructed by another window not obstructed by an active popup or modal blocking inputs under them.
	HoveredFlagsChildWindows            HoveredFlags = 1 << 0 // IsWindowHovered() only: Return true if any children of the window is hovered
	HoveredFlagsRootWindow              HoveredFlags = 1 << 1 // IsWindowHovered() only: Test from root window (top most parent of the current hierarchy)
	HoveredFlagsAnyWindow               HoveredFlags = 1 << 2 // IsWindowHovered() only: Return true if any window is hovered
	HoveredFlagsAllowWhenBlockedByPopup HoveredFlags = 1 << 3 // Return true even if a popup window is normally blocking access to this item/window
	//HoveredFlagsAllowWhenBlockedByModal     HoveredFlags= 1 << 4   // Return true even if a modal popup window is normally blocking access to this item/window. FIXME-TODO: Unavailable yet.
	HoveredFlagsAllowWhenBlockedByActiveItem HoveredFlags = 1 << 5 // Return true even if an active item is blocking access to this item/window. Useful for Drag and Drop patterns.
	HoveredFlagsAllowWhenOverlapped          HoveredFlags = 1 << 6 // Return true even if the position is overlapped by another window
	HoveredFlagsRectOnly                     HoveredFlags = HoveredFlagsAllowWhenBlockedByPopup | HoveredFlagsAllowWhenBlockedByActiveItem | HoveredFlagsAllowWhenOverlapped
	HoveredFlagsRootAndChildWindows          HoveredFlags = HoveredFlagsRootWindow | HoveredFlagsChildWindows
)

type Payload struct {
	Data           interface{}
	SourceId       ID            // Source item id
	SourceParentId ID            // Source parent id (if available)
	DataFrameCount int           // Data timestamp
	DataType       [12 + 1]uint8 // Data type tag (short user-supplied string, 12 characters max)
	Preview        bool          // Set when AcceptDragDropPayload() was called and mouse has been hovering the target item (nb: handle overlapping drag targets)
	Delivery       bool          // Set when AcceptDragDropPayload() was called and mouse button is released over the target item.
}

type ItemStatusFlags int

const (
	ItemStatusFlagsHoveredRect    ItemStatusFlags = 1 << 0
	ItemStatusFlagsHasDisplayRect ItemStatusFlags = 1 << 1
)

func (w *Window) GetID(str string) ID {
	h := fnv.New32()
	h.Write([]byte(str))
	id := ID(h.Sum32())
	w.Ctx.KeepAliveID(id)
	return id
}

func (c *Context) ItemSize(size f64.Vec2) {
	c.ItemSizeEx(size, 0)
}

func (c *Context) ItemSizeEx(size f64.Vec2, text_offset_y float64) {
	window := c.CurrentWindow
	if window.SkipItems {
		return
	}

	// Always align ourselves on pixel boundaries
	line_height := math.Max(window.DC.CurrentLineHeight, size.Y)
	text_base_offset := math.Max(window.DC.CurrentLineTextBaseOffset, text_offset_y)
	window.DC.CursorPosPrevLine = f64.Vec2{window.DC.CursorPos.X + size.X, window.DC.CursorPos.Y}
	window.DC.CursorPos = f64.Vec2{
		float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX)),
		float64(int(window.DC.CursorPos.Y + line_height + c.Style.ItemSpacing.Y)),
	}
	window.DC.CursorMaxPos.X = math.Max(window.DC.CursorMaxPos.X, window.DC.CursorPosPrevLine.X)
	window.DC.CursorMaxPos.Y = math.Max(window.DC.CursorMaxPos.Y, window.DC.CursorPos.Y-c.Style.ItemSpacing.Y)

	window.DC.PrevLineHeight = line_height
	window.DC.PrevLineTextBaseOffset = text_base_offset
	window.DC.CurrentLineHeight = 0
	window.DC.CurrentLineTextBaseOffset = 0

	// Horizontal layout mode
	if window.DC.LayoutType == LayoutTypeHorizontal {
		c.SameLine()
	}
}

func (c *Context) ItemSizeBB(bb f64.Rectangle) {
	c.ItemSizeBBEx(bb, 0)
}

func (c *Context) ItemSizeBBEx(bb f64.Rectangle, text_offset_y float64) {
	c.ItemSizeEx(bb.Size(), text_offset_y)
}

func (c *Context) ItemAdd(bb f64.Rectangle, id ID) bool {
	return c.ItemAddEx(bb, id, nil)
}

func (c *Context) ItemAddEx(bb f64.Rectangle, id ID, nav_bb *f64.Rectangle) bool {
	return false
}

func (c *Context) SameLine() {
	c.SameLineEx(0, -1)
}

// Gets back to previous line and continue with horizontal layout
//      pos_x == 0      : follow right after previous item
//      pos_x != 0      : align to specified x position (relative to window/group left)
//      spacing_w < 0   : use default spacing if pos_x == 0, no spacing if pos_x != 0
//      spacing_w >= 0  : enforce spacing amount
func (c *Context) SameLineEx(pos_x, spacing_w float64) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	if pos_x != 0 {
		if spacing_w < 0 {
			spacing_w = 0
		}
		window.DC.CursorPos.X = window.Pos.X - window.Scroll.X + pos_x + spacing_w + window.DC.GroupOffsetX + window.DC.ColumnsOffsetX
		window.DC.CursorPos.Y = window.DC.CursorPosPrevLine.Y
	} else {
		if spacing_w < 0 {
			spacing_w = c.Style.ItemSpacing.X
		}
		window.DC.CursorPos.X = window.DC.CursorPosPrevLine.X + spacing_w
		window.DC.CursorPos.Y = window.DC.CursorPosPrevLine.Y
	}
	window.DC.CurrentLineHeight = window.DC.PrevLineHeight
	window.DC.CurrentLineTextBaseOffset = window.DC.PrevLineTextBaseOffset
}

func (c *Context) NewLine() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	backup_layout_type := window.DC.LayoutType
	window.DC.LayoutType = LayoutTypeVertical
	if window.DC.CurrentLineHeight > 0.0 { // In the event that we are on a line with items that is smaller that FontSize high, we will preserve its height.
		c.ItemSizeEx(f64.Vec2{0, 0}, -1)
	} else {
		c.ItemSizeEx(f64.Vec2{0, c.FontSize}, -1)
	}
	window.DC.LayoutType = backup_layout_type
}

func (c *Context) NextColumn() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	c.PopItemWidth()
	c.PopClipRect()

	columns := window.DC.ColumnsSet
	columns.CellMaxY = math.Max(columns.CellMaxY, window.DC.CursorPos.Y)
	if columns.Current++; columns.Current < columns.Count {
		// Columns 1+ cancel out IndentX
		window.DC.ColumnsOffsetX = c.GetColumnOffset(columns.Current) - window.DC.IndentX + c.Style.ItemSpacing.X
		window.DrawList.ChannelsSetCurrent(columns.Current)
	} else {
		window.DC.ColumnsOffsetX = 0
		window.DrawList.ChannelsSetCurrent(0)
		columns.Current = 0
		columns.CellMinY = columns.CellMaxY
	}

	window.DC.CursorPos.X = float64(int(window.Pos.X + window.DC.IndentX + window.DC.ColumnsOffsetX))
	window.DC.CursorPos.Y = columns.CellMinY
	window.DC.CurrentLineHeight = 0
	window.DC.CurrentLineTextBaseOffset = 0

	c.PushColumnClipRect()
	c.PushItemWidth(c.GetColumnWidth() * 0.65) // FIXME: Move on columns setup
}

func (c *Context) PushColumnClipRect() {
	c.PushColumnClipRectEx(-1)
}

func (c *Context) PushColumnClipRectEx(column_index int) {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet
	if column_index < 0 {
		column_index = columns.Current
	}

	c.PushClipRect(
		columns.Columns[column_index].ClipRect.Min,
		columns.Columns[column_index].ClipRect.Max,
		false,
	)
}

// When using this function it is sane to ensure that float are perfectly rounded to integer values, to that e.g. (int)(max.x-min.x) in user's render produce correct result.
func (c *Context) PushClipRect(clip_rect_min, clip_rect_max f64.Vec2, intersect_with_current_clip_rect bool) {
	window := c.GetCurrentWindow()
	window.DrawList.PushClipRectEx(clip_rect_min, clip_rect_max, intersect_with_current_clip_rect)
	length := len(window.DrawList._ClipRectStack)
	clipRect := window.DrawList._ClipRectStack[length-1]
	window.ClipRect = f64.Rectangle{f64.Vec2{clipRect.X, clipRect.Y}, f64.Vec2{clipRect.Z, clipRect.W}}
}

func (c *Context) PushItemWidth(item_width float64) {
	window := c.GetCurrentWindow()
	window.DC.ItemWidth = item_width
	if window.DC.ItemWidth == 0 {
		window.DC.ItemWidth = window.ItemWidthDefault
	}
	window.DC.ItemWidthStack = append(window.DC.ItemWidthStack, window.DC.ItemWidth)
}

func (c *Context) GetColumnWidth() float64 {
	return c.GetColumnWidthDx(-1)
}

func (c *Context) GetColumnWidthDx(column_index int) float64 {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet
	if column_index < 0 {
		column_index = columns.Current
	}
	return c.OffsetNormToPixels(
		columns,
		columns.Columns[column_index+1].OffsetNorm-columns.Columns[column_index].OffsetNorm,
	)
}

func (c *Context) OffsetNormToPixels(columns *ColumnsSet, offset_norm float64) float64 {
	return offset_norm * (columns.MaxX - columns.MinX)
}

func (c *Context) GetColumnWidthEx(columns *ColumnsSet, column_index int, before_resize bool) float64 {
	if column_index < 0 {
		column_index = columns.Current
	}

	var offset_norm float64
	if before_resize {
		offset_norm = columns.Columns[column_index+1].OffsetNormBeforeResize - columns.Columns[column_index].OffsetNormBeforeResize
	} else {
		offset_norm = columns.Columns[column_index+1].OffsetNorm - columns.Columns[column_index].OffsetNorm
	}

	return c.OffsetNormToPixels(columns, offset_norm)
}

func (c *Context) GetContentRegionAvail() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	regionMax := c.GetContentRegionMax()
	windowRegion := window.DC.CursorPos.Sub(window.Pos)
	return regionMax.Sub(windowRegion)
}

func (c *Context) GetContentRegionMax() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	mx := window.ContentsRegionRect.Max
	if window.DC.ColumnsSet != nil {
		mx.X = c.GetColumnOffset(window.DC.ColumnsSet.Current+1) - window.WindowPadding.X
	}
	return mx
}

func (c *Context) GetContentRegionAvailWidth() float64 {
	return c.GetContentRegionAvail().X
}

func (c *Context) GetWindowContentRegionMin() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Min
}

func (c *Context) GetWindowContentRegionMax() f64.Vec2 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Max
}

func (c *Context) GetWindowContentRegionWidth() float64 {
	window := c.GetCurrentWindowRead()
	return window.ContentsRegionRect.Max.X - window.ContentsRegionRect.Min.X
}

func (c *Context) GetColumnOffset(column_index int) float64 {
	window := c.GetCurrentWindowRead()
	columns := window.DC.ColumnsSet

	if column_index < 0 {
		column_index = columns.Current
	}

	t := columns.Columns[column_index].OffsetNorm
	x_offset := f64.Lerp(t, columns.MinX, columns.MaxX)
	return x_offset
}

func (c *Context) GetColumnIndex() int {
	window := c.GetCurrentWindowRead()
	if window.DC.ColumnsSet != nil {
		return window.DC.ColumnsSet.Current
	}
	return 0
}

func (c *Context) GetColumnsCount() int {
	window := c.GetCurrentWindowRead()
	if window.DC.ColumnsSet != nil {
		return window.DC.ColumnsSet.Count
	}
	return 0
}

func (c *Context) PopItemWidth() {
	window := c.GetCurrentWindow()
	length := len(window.DC.ItemWidthStack)
	window.DC.ItemWidthStack = window.DC.ItemWidthStack[:length-1]
	if length--; length == 0 {
		window.DC.ItemWidth = window.ItemWidthDefault
	} else {
		window.DC.ItemWidth = window.DC.ItemWidthStack[length-1]
	}
}

func (c *Context) CalcItemWidth() float64 {
	window := c.GetCurrentWindowRead()
	w := window.DC.ItemWidth
	if w < 0 {
		// Align to a right-side limit. We include 1 frame padding in the calculation because this is how the width is always used (we add 2 frame padding to it), but we could move that responsibility to the widget as well.
		width_to_right_edge := c.GetContentRegionAvail().X
		w = math.Max(1, width_to_right_edge+w)
	}
	w = float64(int(w))
	return w
}

func (c *Context) PopClipRect() {
	window := c.GetCurrentWindow()
	window.DrawList.PopClipRect()
	length := len(window.DrawList._ClipRectStack)
	clipRect := window.DrawList._ClipRectStack[length-1]
	window.ClipRect = f64.Rectangle{f64.Vec2{clipRect.X, clipRect.Y}, f64.Vec2{clipRect.Z, clipRect.W}}
}

func (c *Context) CalcItemSize(size f64.Vec2, default_x, default_y float64) f64.Vec2 {
	var content_max f64.Vec2
	if size.X < 0 || size.Y < 0 {
		content_max = c.CurrentWindow.Pos.Add(c.GetContentRegionMax())
	}
	if size.X == 0 {
		size.X = default_x
	} else {
		size.X += math.Max(content_max.X-c.CurrentWindow.DC.CursorPos.X, 4)
	}

	if size.Y == 0 {
		size.Y = default_y
	} else {
		size.Y += math.Max(content_max.Y-c.CurrentWindow.DC.CursorPos.Y, 4)
	}

	return size
}

func (c *Context) CalcWrapWidthForPos(pos f64.Vec2, wrap_pos_x float64) float64 {
	if wrap_pos_x < 0 {
		return 0
	}

	window := c.GetCurrentWindowRead()
	if wrap_pos_x == 0 {
		wrap_pos_x = c.GetContentRegionMax().X + window.Pos.X
	} else if wrap_pos_x > 0 {
		wrap_pos_x += window.Pos.X - window.Scroll.X // wrap_pos_x is provided is window local space
	}

	return math.Max(wrap_pos_x-pos.X, 1)
}

func (w *Window) CalcFontSize() float64 {
	return w.Ctx.FontBaseSize * w.FontWindowScale
}

func (c *Context) ItemHoverable(bb f64.Rectangle, id ID) bool {
	if c.HoveredId != 0 && c.HoveredId != id && !c.HoveredIdAllowOverlap {
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

func (c *Context) IsItemHovered(flags HoveredFlags) bool {
	window := c.CurrentWindow
	if c.NavDisableMouseHover && !c.NavDisableHighlight {
		return c.IsItemFocused()
	}

	// Test for bounding box overlap, as updated as ItemAdd()
	if window.DC.LastItemStatusFlags&ItemStatusFlagsHoveredRect != 0 {
		return false
	}

	if c.HoveredRootWindow != window.RootWindow && flags&HoveredFlagsAllowWhenOverlapped == 0 {
		return false
	}

	// Test if another item is active (e.g. being dragged)
	if flags&HoveredFlagsAllowWhenBlockedByActiveItem == 0 {
		if c.ActiveId != 0 && c.ActiveId != window.DC.LastItemId && !c.ActiveIdAllowOverlap && c.ActiveId != window.MoveId {
			return false
		}
	}

	// Test if interactions on this window are blocked by an active popup or modal
	if !c.IsWindowContentHoverable(window, flags) {
		return false
	}

	// Test if the item is disabled
	if window.DC.ItemFlags&ItemFlagsDisabled != 0 {
		return false
	}

	// Special handling for the 1st item after Begin() which represent the title bar. When the window is collapsed (SkipItems==true) that last item will never be overwritten so we need to detect tht case.
	if window.DC.LastItemId == window.MoveId && window.WriteAccessed {
		return false
	}

	return true
}

func (c *Context) IsMouseHoveringRect(r_min, r_max f64.Vec2) bool {
	return c.IsMouseHoveringRectEx(r_min, r_max, true)
}

// Test if mouse cursor is hovering given rectangle
// NB- Rectangle is clipped by our current clip setting
// NB- Expand the rectangle to be generous on imprecise inputs systems (g.Style.TouchExtraPadding)
func (c *Context) IsMouseHoveringRectEx(r_min, r_max f64.Vec2, clip bool) bool {
	window := c.CurrentWindow

	rect_clipped := f64.Rectangle{r_min, r_max}
	if clip {
		rect_clipped.Intersect(window.ClipRect)
	}

	// Expand for touch input
	rect_for_touch := f64.Rectangle{
		rect_clipped.Min.Sub(c.Style.TouchExtraPadding),
		rect_clipped.Max.Add(c.Style.TouchExtraPadding),
	}
	return c.IO.MousePos.In(rect_for_touch)
}

func (c *Context) IsWindowContentHoverable(window *Window, flags HoveredFlags) bool {
	if c.NavWindow != nil {
		focused_root_window := c.NavWindow.RootWindow
		if focused_root_window != nil {
			if focused_root_window.WasActive && focused_root_window != window.RootWindow {
				// For the purpose of those flags we differentiate "standard popup" from "modal popup"
				// NB: The order of those two tests is important because Modal windows are also Popups.
				if focused_root_window.Flags&WindowFlagsModal != 0 {
					return false
				}
				if focused_root_window.Flags&WindowFlagsPopup != 0 && flags&HoveredFlagsAllowWhenBlockedByPopup == 0 {
					return false
				}
			}
		}
	}
	return true
}

func (c *Context) IsItemFocused() bool {
	return c.NavId != 0 && !c.NavDisableHighlight && c.NavId == c.CurrentWindow.DC.LastItemId
}

func (c *Context) IsItemClicked(mouse_button int) bool {
	return c.IsMouseClicked(mouse_button, false) && c.IsItemHovered(HoveredFlagsDefault)
}

func (c *Context) IsAnyItemHovered() bool {
	return c.HoveredId != 0 || c.HoveredIdPreviousFrame != 0
}

func (c *Context) IsMouseClicked(button int, repeat bool) bool {
	t := c.IO.MouseDownDuration[button]
	if t == 0 {
		return true
	}

	if repeat && t > c.IO.KeyRepeatDelay {
		delay := c.IO.KeyRepeatDelay
		rate := c.IO.KeyRepeatRate
		mod1 := math.Mod(t-delay, rate) > rate*0.5
		mod2 := math.Mod(t-delay-c.IO.DeltaTime, rate) > rate*0.5
		if mod1 != mod2 {
			return true
		}
	}

	return false
}

func (c *Context) IsMouseReleased(button int) bool {
	return c.IO.MouseReleased[button]
}

func (c *Context) IsMouseDoubleClicked(button int) bool {
	return c.IO.MouseDoubleClicked[button]
}

func (c *Context) IsMouseDragging(button int, lock_threshold float64) bool {
	if !c.IO.MouseDown[button] {
		return false
	}
	if lock_threshold < 0.0 {
		lock_threshold = c.IO.MouseDragThreshold
	}
	return c.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold
}

func (c *Context) GetMousePos() f64.Vec2 {
	return c.IO.MousePos
}

func (c *Context) IsMousePosValid(mouse_pos *f64.Vec2) bool {
	if mouse_pos == nil {
		mouse_pos = &c.IO.MousePos
	}
	const MOUSE_INVALID = -256000.0
	return mouse_pos.X >= MOUSE_INVALID && mouse_pos.Y >= MOUSE_INVALID
}

func (c *Context) GetMouseDragDelta(button int, lock_threshold float64) f64.Vec2 {
	if lock_threshold < 0 {
		lock_threshold = c.IO.MouseDragThreshold
	}
	if c.IO.MouseDown[button] {
		if c.IO.MouseDragMaxDistanceSqr[button] >= lock_threshold*lock_threshold {
			return c.IO.MousePos.Sub(c.IO.MouseClickedPos[button]) // Assume we can only get active with left-mouse button (at the moment).
		}
	}
	return f64.Vec2{0, 0}
}

func (c *Context) ResetMouseDragDelta(button int) {
	c.IO.MouseClickedPos[button] = c.IO.MousePos
}

func (c *Context) GetMouseCursor() MouseCursor {
	return c.MouseCursor
}

func (c *Context) SetMouseCursor(cursor_type MouseCursor) {
	c.MouseCursor = cursor_type
}

func (c *Context) IsItemActive() bool {
	if c.ActiveId != 0 {
		window := c.CurrentWindow
		return c.ActiveId == window.DC.LastItemId
	}
	return false
}

func (c *Context) IsAnyItemFocused() bool {
	return c.NavId != 0 && !c.NavDisableHighlight
}

func (c *Context) IsItemVisible() bool {
	window := c.GetCurrentWindowRead()
	return window.ClipRect.Overlaps(window.DC.LastItemRect)
}

// Allow last item to be overlapped by a subsequent item. Both may be activated during the same frame before the later one takes priority.
func (c *Context) SetItemAllowOverlap() {
	if c.HoveredId == c.CurrentWindow.DC.LastItemId {
		c.HoveredIdAllowOverlap = true
	}
	if c.ActiveId == c.CurrentWindow.DC.LastItemId {
		c.ActiveIdAllowOverlap = true
	}
}

func (c *Context) GetViewportRect() f64.Rectangle {
	if c.IO.DisplayVisibleMin.X != c.IO.DisplayVisibleMax.X && c.IO.DisplayVisibleMin.Y != c.IO.DisplayVisibleMax.Y {
		return f64.Rectangle{c.IO.DisplayVisibleMin, c.IO.DisplayVisibleMax}
	}
	return f64.Rect(0, 0, c.IO.DisplaySize.X, c.IO.DisplaySize.Y)
}

// Moving window to front of display and set focus (which happens to be back of our sorted list)
func (c *Context) FocusWindow(window *Window) {
	if c.NavWindow != window {
		c.NavWindow = window
		if window != nil && c.NavDisableMouseHover {
			c.NavMousePosDirty = true
		}
		c.NavInitRequest = false
		c.NavId = 0
		if window != nil {
			c.NavId = window.NavLastIds[0]
		}
		c.NavIdIsAlive = false
		c.NavLayer = 0
	}

	// Passing NULL allow to disable keyboard focus
	if window == nil {
		return
	}

	// Move the root window to the top of the pile
	if window.RootWindow != nil {
		window = window.RootWindow
	}

	// Steal focus on active widgets
	// FIXME: This statement should be unnecessary. Need further testing before removing it..
	if window.Flags&WindowFlagsPopup != 0 {
		if c.ActiveId != 0 && c.ActiveIdWindow != nil && c.ActiveIdWindow.RootWindow != window {
			c.ClearActiveID()
		}
	}

	// Bring to front
	if window.Flags&WindowFlagsNoBringToFrontOnFocus == 0 {
		c.BringWindowToFront(window)
	}
}

func (c *Context) BringWindowToFront(window *Window) {
	current_front_window := c.Windows[len(c.Windows)-1]
	if current_front_window == window || current_front_window.RootWindow == window {
		return
	}
	// We can ignore the front most window
	for i := len(c.Windows) - 2; i >= 0; i-- {
		if c.Windows[i] == window {
			c.Windows = append(c.Windows[:i], c.Windows[i+1:]...)
			c.Windows = append(c.Windows, window)
			break
		}
	}
}

func (c *Context) EndColumns() {
}

func (c *Context) SetCurrentWindow(window *Window) {
	c.CurrentWindow = window
	if window != nil {
		c.FontSize = window.CalcFontSize()
		c.DrawListSharedData.FontSize = c.FontSize
	}
}