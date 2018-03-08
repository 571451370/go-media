package imgui

import (
	"os"

	"github.com/qeedquan/go-media/math/f64"
)

type (
	ID        uint
	TextureID interface{}
)

type Context struct {
	Initialized             bool
	FontAtlasOwnedByContext bool // Io.Fonts-> is owned by the ImGuiContext and will be destructed along with it.
	IO                      IO
	Style                   Style
	Font                    *Font   // (Shortcut) == FontStack.empty() ? IO.Font : FontStack.back()
	FontSize                float64 // (Shortcut) == FontBaseSize * g.CurrentWindow->FontWindowScale == window->FontSize(). Text height for current window.
	FontBaseSize            float64 // (Shortcut) == IO.FontGlobalScale * Font->Scale * Font->FontSize. Base text height.
	DrawListSharedData      DrawListSharedData

	Time                     float64
	FrameCount               int
	FrameCountEnded          int
	FrameCountRendered       int
	Windows                  []*Window
	WindowsSortBuffer        []*Window
	CurrentWindowStack       []*Window
	WindowsById              Storage
	WindowsActiveCount       int
	CurrentWindow            *Window // Being drawn into
	HoveredWindow            *Window // Will catch mouse inputs
	HoveredRootWindow        *Window // Will catch mouse inputs (for focus/move only)
	HoveredId                ID      // Hovered widget
	HoveredIdAllowOverlap    bool
	HoveredIdPreviousFrame   ID
	HoveredIdTimer           float64
	ActiveId                 ID // Active widget
	ActiveIdPreviousFrame    ID
	ActiveIdTimer            float64
	ActiveIdIsAlive          bool     // Active widget has been seen this frame
	ActiveIdIsJustActivated  bool     // Set at the time of activation for one frame
	ActiveIdAllowOverlap     bool     // Active widget allows another widget to steal active id (generally for overlapping widgets, but not always)
	ActiveIdAllowNavDirFlags int      // Active widget allows using directional navigation (e.g. can activate a button and move away from it)
	ActiveIdClickOffset      f64.Vec2 // Clicked offset from upper-left corner, if applicable (currently only set by ButtonBehavior)
	ActiveIdWindow           *Window
	ActiveIdSource           InputSource    // Activating with mouse or nav (gamepad/keyboard)
	MovingWindow             *Window        // Track the window we clicked on (in order to preserve focus). The actually window that is moved is generally MovingWindow->RootWindow.
	ColorModifiers           []ColMod       // Stack for PushStyleColor()/PopStyleColor()
	StyleModifiers           []StyleMod     // Stack for PushStyleVar()/PopStyleVar()
	FontStack                []*Font        // Stack for PushFont()/PopFont()
	OpenPopupStack           []PopupRef     // Which popups are open (persistent)
	CurrentPopupStack        []PopupRef     // Which level of BeginPopup() we are in (reset every frame)
	NextWindowData           NextWindowData // Storage for SetNextWindow** functions
	NextTreeNodeOpenVal      bool           // Storage for SetNextTreeNode** functions
	NextTreeNodeOpenCond     Cond

	// Navigation data (for gamepad/keyboard)
	NavWindow                  *Window       // Focused window for navigation. Could be called 'FocusWindow'
	NavId                      ID            // Focused item for navigation
	NavActivateId              ID            // ~~ (g.ActiveId == 0) && IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0, also set when calling ActivateItem()
	NavActivateDownId          ID            // ~~ IsNavInputDown(ImGuiNavInput_Activate) ? NavId : 0
	NavActivatePressedId       ID            // ~~ IsNavInputPressed(ImGuiNavInput_Activate) ? NavId : 0
	NavInputId                 ID            // ~~ IsNavInputPressed(ImGuiNavInput_Input) ? NavId : 0
	NavJustTabbedId            ID            // Just tabbed to this id.
	NavNextActivateId          ID            // Set by ActivateItem(), queued until next frame
	NavJustMovedToId           ID            // Just navigated to this id (result of a successfully MoveRequest)
	NavScoringRectScreen       f64.Rectangle // Rectangle used for scoring, in screen space. Based of window->DC.NavRefRectRel[], modified for directional navigation scoring.
	NavScoringCount            int           // Metrics for debugging
	NavWindowingTarget         *Window       // When selecting a window (holding Menu+FocusPrev/Next, or equivalent of CTRL-TAB) this window is temporarily displayed front-most.
	NavWindowingHighlightTimer float64
	NavWindowingHighlightAlpha float64
	NavWindowingToggleLayer    bool
	NavWindowingInputSource    InputSource // Gamepad or keyboard mode
	NavLayer                   int         // Layer we are navigating on. For now the system is hard-coded for 0=main contents and 1=menu/title bar, may expose layers later.
	NavIdTabCounter            int         // == NavWindow->DC.FocusIdxTabCounter at time of NavId processing
	NavIdIsAlive               bool        // Nav widget has been seen this frame ~~ NavRefRectRel is valid
	NavMousePosDirty           bool        // When set we will update mouse position if (io.ConfigFlags & ImGuiConfigFlags_NavMoveMouse) if set (NB: this not enabled by default)
	NavDisableHighlight        bool        // When user starts using mouse, we hide gamepad/keyboard highlight (NB: but they are still available, which is why NavDisableHighlight isn't always != NavDisableMouseHover)
	NavDisableMouseHover       bool        // When user starts using gamepad/keyboard, we hide mouse hovering highlight until mouse is touched again.
	NavAnyRequest              bool        // ~~ NavMoveRequest || NavInitRequest
	NavInitRequest             bool        // Init request for appearing window to select first item
	NavInitRequestFromMove     bool
	NavInitResultId            ID
	NavInitResultRectRel       f64.Rectangle
	NavMoveFromClampedRefRect  bool          // Set by manual scrolling, if we scroll to a point where NavId isn't visible we reset navigation from visible items
	NavMoveRequest             bool          // Move request for this frame
	NavMoveRequestForward      NavForward    // None / ForwardQueued / ForwardActive (this is used to navigate sibling parent menus from a child menu)
	NavMoveDir, NavMoveDirLast Dir           // Direction of the move request (left/right/up/down), direction of the previous move request
	NavMoveResultLocal         NavMoveResult // Best move request candidate within NavWindow
	NavMoveResultOther         NavMoveResult // Best move request candidate within NavWindow's flattened hierarchy (when using the NavFlattened flag)

	// Render
	DrawData                  DrawData // Main ImDrawData instance to pass render information to the user
	DrawDataBuilder           DrawDataBuilder
	ModalWindowDarkeningRatio float64
	OverlayDrawList           DrawList // Optional software render of mouse cursors, if io.MouseDrawCursor is set + a few debug overlays
	MouseCursor               MouseCursor

	// Drag and Drop
	DragDropActive                  bool
	DragDropSourceFlags             DragDropFlags
	DragDropMouseButton             int
	DragDropPayload                 Payload
	DragDropTargetRect              f64.Rectangle
	DragDropTargetId                ID
	DragDropAcceptIdCurrRectSurface float64
	DragDropAcceptIdCurr            ID      // Target item id (set at the time of accepting the payload)
	DragDropAcceptIdPrev            ID      // Target item id from previous frame (we need to store this to allow for overlapping drag and drop targets)
	DragDropAcceptFrameCount        int     // Last time a target expressed a desire to accept the source
	DragDropPayloadBufHeap          []uint8 // We don't expose the ImVector<> directly
	DragDropPayloadBufLocal         [8]uint8

	// Widget state
	InputTextState                  TextEditState
	InputTextPasswordFont           Font
	ScalarAsInputTextId             ID             // Temporary text input when CTRL+clicking on a slider, etc.
	ColorEditOptions                ColorEditFlags // Store user options for color edit widgets
	ColorPickerRef                  f64.Vec4
	DragCurrentValue                float64 // Currently dragged value, always float, not rounded by end-user precision settings
	DragLastMouseDelta              f64.Vec2
	DragSpeedDefaultRatio           float64 // If speed == 0.0f, uses (max-min) * DragSpeedDefaultRatio
	DragSpeedScaleSlow              float64
	DragSpeedScaleFast              float64
	ScrollbarClickDeltaToGrabCenter f64.Vec2 // Distance between mouse and center of grab box, normalized in parent space. Use storage?
	TooltipOverrideCount            int
	PrivateClipboard                []int8   // If no custom clipboard handler is defined
	OsImePosRequest, OsImePosSet    f64.Vec2 // Cursor position request & last passed to the OS Input Method Editor

	// Settings
	SettingsLoaded     bool
	SettingsDirtyTimer float64          // Save .ini Settings on disk when time reaches zero
	SettingsWindows    []WindowSettings // .ini settings for ImGuiWindow
	SettingsHandlers   []func()         // List of .ini settings handlers

	// Logging
	LogEnabled            bool
	LogFile               *os.File // If != NULL log to stdout/ file
	LogClipboard          []rune   // Else log to clipboard. This is pointer so our GImGui static constructor doesn't call heap allocators.
	LogStartDepth         int
	LogAutoExpandMaxDepth int

	// Misc
	FramerateSecPerFrame         [120]float64 // calculate estimate of framerate for user
	FramerateSecPerFrameIdx      int
	FramerateSecPerFrameAccum    float64
	WantCaptureMouseNextFrame    int // explicit capture via CaptureInputs() sets those flags
	WantCaptureKeyboardNextFrame int
	WantTextInputNextFrame       int
	TempBuffer                   [1024*3 + 1]uint8 // temporary text buffer
}

type ConfigFlags int

const (
	ConfigFlagsNavEnableKeyboard    ConfigFlags = 1 << 0 // Master keyboard navigation enable flag. NewFrame() will automatically fill io.NavInputs[] based on io.KeyDown[].
	ConfigFlagsNavEnableGamepad     ConfigFlags = 1 << 1 // Master gamepad navigation enable flag. This is mostly to instruct your imgui back-end to fill io.NavInputs[].
	ConfigFlagsNavMoveMouse         ConfigFlags = 1 << 2 // Request navigation to allow moving the mouse cursor. May be useful on TV/console systems where moving a virtual mouse is awkward. Will update io.MousePos and set io.WantMoveMouseConfigFlags=true. If enabled you MUST honor io.WantMoveMouse requests in your binding otherwise ImGui will react as if the mouse is jumping around back and forth.
	ConfigFlagsNavNoCaptureKeyboard ConfigFlags = 1 << 3 // Do not set the io.WantCaptureKeyboard flag with io.NavActive is set.

	// User storage (to allow your back-end/engine to communicate to code that may be shared between multiple projects. Those flags are not used by core ImGui)
	ConfigFlagsIsSRGB        ConfigFlags = 1 << 20 // Back-end is SRGB-aware.
	ConfigFlagsIsTouchScreen ConfigFlags = 1 << 21 // Back-end is using a touch screen instead of a mouse.
)

type Storage struct {
}

func CreateContext() *Context {
	return CreateContextEx(nil)
}

func CreateContextEx(shared_font_atlas *FontAtlas) *Context {
	return &Context{}
}

func (c *Context) GetVersion() string {
	return "1.6.0 WIP"
}

func (c *Context) GetIO() *IO {
	return &c.IO
}

func (c *Context) GetStyle() *Style {
	return &c.Style
}

func (c *Context) GetCurrentWindowRead() *Window {
	c.CurrentWindow.WriteAccessed = true
	return c.CurrentWindow
}

func (c *Context) GetCurrentWindow() *Window {
	c.CurrentWindow.WriteAccessed = true
	return c.CurrentWindow
}

func (c *Context) SetFocusID(id ID, window *Window) {
	// Assume that SetFocusID() is called in the context where its NavLayer is the current layer, which is the case everywhere we call it.
	nav_layer := window.DC.NavLayerCurrent
	if c.NavWindow != window {
		c.NavInitRequest = false
	}
	c.NavId = id
	c.NavWindow = window
	c.NavLayer = nav_layer
	window.NavLastIds[nav_layer] = id
	if window.DC.LastItemId == id {
		window.NavRectRel[nav_layer] = f64.Rectangle{
			window.DC.LastItemRect.Min.Sub(window.Pos),
			window.DC.LastItemRect.Max.Sub(window.Pos),
		}
	}

	if c.ActiveIdSource == InputSourceNav {
		c.NavDisableMouseHover = true
	} else {
		c.NavDisableHighlight = true
	}
}

func (c *Context) KeepAliveID(id ID) {
	if c.ActiveId == id {
		c.ActiveIdIsAlive = true
	}
}

func (c *Context) ClearActiveID() {
	c.SetActiveID(0, nil)
}

func (c *Context) SetActiveID(id ID, window *Window) {
	c.ActiveIdIsJustActivated = c.ActiveId != id
	if c.ActiveIdIsJustActivated {
		c.ActiveIdTimer = 0
	}
	c.ActiveId = id
	c.ActiveIdAllowNavDirFlags = 0
	c.ActiveIdAllowOverlap = false
	c.ActiveIdWindow = window
	if id != 0 {
		c.ActiveIdIsAlive = true
		c.ActiveIdSource = InputSourceMouse
		if c.NavActivateId == id || c.NavInputId == id || c.NavJustTabbedId == id || c.NavJustMovedToId == id {
			c.ActiveIdSource = InputSourceNav
		}
	}
}

func (c *Context) GetFrameHeight() float64 {
	return c.FontSize + c.Style.FramePadding.Y*2
}

func (c *Context) GetFrameHeightWithSpacing() float64 {
	return c.FontSize + c.Style.FramePadding.Y*2 + c.Style.ItemSpacing.Y
}

func (c *Context) GetTime() float64 {
	return c.Time
}

func (c *Context) GetFrameCount() int {
	return c.FrameCount
}

func (c *Context) GetOverlayDrawList() *DrawList {
	return &c.OverlayDrawList
}

func (c *Context) GetDrawListSharedData() *DrawListSharedData {
	return &c.DrawListSharedData
}

func (c *Context) SetHoveredID(id ID) {
	c.HoveredId = id
	c.HoveredIdAllowOverlap = false
	if id != 0 && c.HoveredIdPreviousFrame == id {
		c.HoveredIdTimer = c.IO.DeltaTime
	} else {
		c.HoveredIdTimer = 0
	}
}

func (c *Context) LogFinish() {
	if !c.LogEnabled {
		return
	}
	c.LogEnabled = false
}

func (c *Context) PushID(str_id string) {
	window := c.GetCurrentWindowRead()
	window.IDStack = append(window.IDStack, window.GetID(str_id))

}