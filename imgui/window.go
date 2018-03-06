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

type ItemStatusFlags int

const (
	ItemStatusFlags_HoveredRect    ItemStatusFlags = 1 << 0
	ItemStatusFlags_HasDisplayRect ItemStatusFlags = 1 << 1
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

type Payload struct {
	Data           interface{}
	SourceId       ID            // Source item id
	SourceParentId ID            // Source parent id (if available)
	DataFrameCount int           // Data timestamp
	DataType       [12 + 1]uint8 // Data type tag (short user-supplied string, 12 characters max)
	Preview        bool          // Set when AcceptDragDropPayload() was called and mouse has been hovering the target item (nb: handle overlapping drag targets)
	Delivery       bool          // Set when AcceptDragDropPayload() was called and mouse button is released over the target item.
}

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
	return 0
}

func (c *Context) GetColumnWidthEx(columns *ColumnsSet, column_index int, before_resize bool) {
}

func (c *Context) GetContentRegionAvail() f64.Vec2 {
	return f64.Vec2{}
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