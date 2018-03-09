package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type NavInput int

const (
	// Gamepad Mapping
	NavInputActivate    NavInput = iota // activate / open / toggle / tweak value       // e.g. Circle (PS4) A (Xbox) A (Switch) Space (Keyboard)
	NavInputCancel                      // cancel / close / exit                        // e.g. Cross  (PS4) B (Xbox) B (Switch) Escape (Keyboard)
	NavInputInput                       // text input / on-screen keyboard              // e.g. Triang.(PS4) Y (Xbox) X (Switch) Return (Keyboard)
	NavInputMenu                        // tap: toggle menu / hold: focus move resize // e.g. Square (PS4) X (Xbox) Y (Switch) Alt (Keyboard)
	NavInputDpadLeft                    // move / tweak / resize window (w/ PadMenu)    // e.g. D-pad Left/Right/Up/Down (Gamepads) Arrow keys (Keyboard)
	NavInputDpadRight                   //
	NavInputDpadUp                      //
	NavInputDpadDown                    //
	NavInputLStickLeft                  // scroll / move window (w/ PadMenu)            // e.g. Left Analog Stick Left/Right/Up/Down
	NavInputLStickRight                 //
	NavInputLStickUp                    //
	NavInputLStickDown                  //
	NavInputFocusPrev                   // next window (w/ PadMenu)                     // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputFocusNext                   // prev window (w/ PadMenu)                     // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)
	NavInputTweakSlow                   // slower tweaks                                // e.g. L1 or L2 (PS4) LB or LT (Xbox) L or ZL (Switch)
	NavInputTweakFast                   // faster tweaks                                // e.g. R1 or R2 (PS4) RB or RT (Xbox) R or ZL (Switch)

	// [Internal] Don't use directly! This is used internally to differentiate keyboard from gamepad inputs for behaviors that require to differentiate them.
	// Keyboard behavior that have no corresponding gamepad mapping (e.g. CTRL+TAB) may be directly reading from io.KeyDown[] instead of io.NavInputs[].
	NavInputKeyMenu_  // toggle menu                                  // = io.KeyAlt
	NavInputKeyLeft_  // move left                                    // = Arrow keys
	NavInputKeyRight_ // move right
	NavInputKeyUp_    // move up
	NavInputKeyDown_  // move down
	NavInputCOUNT
	NavInputInternalStart_ = NavInputKeyMenu_
)

type Dir int

const (
	DirNone  Dir = -1
	DirLeft  Dir = 0
	DirRight Dir = 1
	DirUp    Dir = 2
	DirDown  Dir = 3
	DirCOUNT Dir = 4
)

type Cond int

const (
	CondAlways       Cond = 1 << 0 // Set the variable
	CondOnce         Cond = 1 << 1 // Set the variable once per runtime session (only the first call with succeed)
	CondFirstUseEver Cond = 1 << 2 // Set the variable if the window has no saved data (if doesn't exist in the .ini file)
	CondAppearing    Cond = 1 << 3 // Set the variable if the window is appearing after being hidden/inactive (or the first time)
)

type NavHighlightFlags int

const (
	NavHighlightFlagsTypeDefault NavHighlightFlags = 1 << 0
	NavHighlightFlagsTypeThin    NavHighlightFlags = 1 << 1
	NavHighlightFlagsAlwaysDraw  NavHighlightFlags = 1 << 2
	NavHighlightFlagsNoRounding  NavHighlightFlags = 1 << 3
)

type NavDirSourceFlags int

const (
	NavDirSourceFlagsKeyboard  NavDirSourceFlags = 1 << 0
	NavDirSourceFlagsPadDPad   NavDirSourceFlags = 1 << 1
	NavDirSourceFlagsPadLStick NavDirSourceFlags = 1 << 2
)

type NavForward int

const (
	NavForwardNone NavForward = iota
	NavForwardForwardQueued
	NavForwardForwardActive
)

type NavMoveResult struct {
	ID         ID      // Best candidate
	ParentID   ID      // Best candidate window->IDStack.back() - to compare context
	Window     *Window // Best candidate window
	DistBox    float64 // Best candidate box distance to current NavId
	DistCenter float64 // Best candidate center distance to current NavId
	DistAxial  float64
	RectRel    f64.Rectangle // Best candidate bounding box in window relative space
}

func (c *Context) IsNavInputPressed(n NavInput, mode InputReadMode) bool {
	return c.GetNavInputAmount(n, mode) > 0
}

func (c *Context) GetNavInputAmount(n NavInput, mode InputReadMode) float64 {
	// Instant, read analog input (0.0f..1.0f, as provided by user)
	if mode == InputReadModeDown {
		return c.IO.NavInputs[n]
	}

	t := c.IO.NavInputsDownDuration[n]
	// Return 1.0f when just released, no repeat, ignore analog input.
	if t < 0 && mode == InputReadModeReleased {
		if c.IO.NavInputsDownDurationPrev[n] >= 0 {
			return 1
		}
		return 0
	}

	if t < 0 {
		return 0
	}

	// Return 1.0f when just pressed, no repeat, ignore analog input.
	if mode == InputReadModePressed {
		if t == 0 {
			return 1
		}
		return 0
	}

	if mode == InputReadModeRepeat {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.80,
		))
	}

	if mode == InputReadModeRepeatSlow {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*1,
			c.IO.KeyRepeatRate*2,
		))
	}

	if mode == InputReadModeRepeatFast {
		return float64(c.CalcTypematicPressedRepeatAmount(
			t,
			t-c.IO.DeltaTime,
			c.IO.KeyRepeatDelay*0.80,
			c.IO.KeyRepeatRate*0.30,
		))
	}

	return 0
}

func (c *Context) GetNavInputAmount2d(dir_sources NavDirSourceFlags, mode InputReadMode, slow_factor, fast_factor float64) f64.Vec2 {
	delta := f64.Vec2{}
	if dir_sources&NavDirSourceFlagsKeyboard != 0 {
		right := c.GetNavInputAmount(NavInputKeyRight_, mode)
		left := c.GetNavInputAmount(NavInputKeyLeft_, mode)
		down := c.GetNavInputAmount(NavInputKeyRight_, mode)
		up := c.GetNavInputAmount(NavInputKeyRight_, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadDPad != 0 {
		right := c.GetNavInputAmount(NavInputDpadRight, mode)
		left := c.GetNavInputAmount(NavInputDpadLeft, mode)
		down := c.GetNavInputAmount(NavInputDpadDown, mode)
		up := c.GetNavInputAmount(NavInputDpadUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if dir_sources&NavDirSourceFlagsPadLStick != 0 {
		right := c.GetNavInputAmount(NavInputLStickRight, mode)
		left := c.GetNavInputAmount(NavInputLStickLeft, mode)
		down := c.GetNavInputAmount(NavInputLStickDown, mode)
		up := c.GetNavInputAmount(NavInputLStickUp, mode)
		dir := f64.Vec2{right - left, down - up}
		delta = delta.Add(dir)
	}
	if slow_factor != 0.0 && c.IsNavInputDown(NavInputTweakSlow) {
		delta = delta.Scale(slow_factor)
	}
	if fast_factor != 0.0 && c.IsNavInputDown(NavInputTweakFast) {
		delta = delta.Scale(fast_factor)
	}
	return delta
}

// FIXME-OPT O(N)
func (c *Context) FindWindowIndex(window *Window) int {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		if c.Windows[i] == window {
			return i
		}
	}
	return -1
}

// FIXME-OPT O(N)
func (c *Context) FindWindowNavigable(i_start, i_stop, dir int) *Window {
	for i := i_start; i >= 0 && i < len(c.Windows) && i != i_stop; i += dir {
		if c.IsWindowNavFocusable(c.Windows[i]) {
			return c.Windows[i]
		}
	}
	return nil
}

func (c *Context) IsWindowNavFocusable(window *Window) bool {
	return window.Active && window == window.RootWindowForTabbing && window.Flags&WindowFlagsNoNavFocus == 0 || window == c.NavWindow
}

func (c *Context) NavUpdateWindowingHighlightWindow(focus_change_dir int) {
	if c.NavWindowingTarget.Flags&WindowFlagsModal != 0 {
		return
	}

	i_current := c.FindWindowIndex(c.NavWindowingTarget)
	window_target := c.FindWindowNavigable(i_current+focus_change_dir, -math.MinInt32, focus_change_dir)
	if window_target == nil {
		if focus_change_dir < 0 {
			window_target = c.FindWindowNavigable(len(c.Windows)-1, i_current, focus_change_dir)
		} else {
			window_target = c.FindWindowNavigable(0, i_current, focus_change_dir)
		}
	}
	c.NavWindowingTarget = window_target
	c.NavWindowingToggleLayer = false
}

// Equivalent of IsKeyDown() for NavInputs[]
func (c *Context) IsNavInputDown(n NavInput) bool {
	return c.IO.NavInputs[n] > 0.0
}

func (n *NavMoveResult) Clear() {
	n.ID = 0
	n.ParentID = 0
	n.Window = nil
	n.DistBox = math.MaxFloat32
	n.DistCenter = math.MaxFloat32
	n.DistAxial = math.MaxFloat32
	n.RectRel = f64.Rectangle{}
}

func (c *Context) NavUpdate() {
	c.IO.WantMoveMouse = false

	// Update Keyboard->Nav inputs mapping
	for i := int(NavInputInternalStart_); i < len(c.IO.NavInputs); i++ {
		c.IO.NavInputs[i] = 0
	}
	if c.IO.ConfigFlags&ConfigFlagsNavEnableKeyboard != 0 {
		c.navMapKey(KeySpace, NavInputActivate)
		c.navMapKey(KeyEnter, NavInputInput)
		c.navMapKey(KeyEscape, NavInputCancel)
		c.navMapKey(KeyLeftArrow, NavInputKeyLeft_)
		c.navMapKey(KeyRightArrow, NavInputKeyRight_)
		c.navMapKey(KeyUpArrow, NavInputKeyUp_)
		c.navMapKey(KeyDownArrow, NavInputKeyDown_)
		if c.IO.KeyCtrl {
			c.IO.NavInputs[NavInputTweakSlow] = 1
		}
		if c.IO.KeyShift {
			c.IO.NavInputs[NavInputTweakFast] = 1
		}
		if c.IO.KeyAlt {
			c.IO.NavInputs[NavInputKeyMenu_] = 1
		}
	}
}

func (c *Context) navMapKey(key Key, nav_input NavInput) {
	if c.IO.KeyMap[key] != -1 && c.IsKeyDown(c.IO.KeyMap[key]) {
		c.IO.NavInputs[nav_input] = 1
	}
}

// We get there when either NavId == id, or when g.NavAnyRequest is set (which is updated by NavUpdateAnyRequestFlag above)
func (c *Context) NavProcessItem(window *Window, nav_bb f64.Rectangle, id ID) {
	item_flags := window.DC.ItemFlags
	nav_bb_rel := f64.Rectangle{
		nav_bb.Min.Sub(window.Pos),
		nav_bb.Max.Sub(window.Pos),
	}
	if c.NavInitRequest && c.NavLayer == window.DC.NavLayerCurrent {
		// Even if 'ImGuiItemFlags_NoNavDefaultFocus' is on (typically collapse/close button) we record the first ResultId so they can be used as a fallback
		if item_flags&ItemFlagsNoNavDefaultFocus == 0 || c.NavInitResultId == 0 {
			c.NavInitResultId = id
			c.NavInitResultRectRel = nav_bb_rel
		}

		if item_flags&ItemFlagsNoNavDefaultFocus == 0 {
			c.NavInitRequest = false // Found a match, clear request
			c.NavUpdateAnyRequestFlag()
		}
	}

	// Scoring for navigation
	if c.NavId != id && item_flags&ItemFlagsNoNav == 0 {
		var result *NavMoveResult
		if window == c.NavWindow {
			result = &c.NavMoveResultLocal
		} else {
			result = &c.NavMoveResultOther
		}

		new_best := c.NavMoveRequest && c.NavScoreItem(result, nav_bb)
		if new_best {
			result.ID = id
			result.ParentID = window.IDStack[len(window.IDStack)-1]
			result.Window = window
			result.RectRel = nav_bb_rel
		}
	}

	// Update window-relative bounding box of navigated item
	if c.NavId == id {
		// Always refresh g.NavWindow, because some operations such as FocusItem() don't have a window.
		c.NavWindow = window
		c.NavLayer = window.DC.NavLayerCurrent
		c.NavIdIsAlive = true
		c.NavIdTabCounter = window.FocusIdxTabCounter
		// Store item bounding box (relative to window position)
		window.NavRectRel[window.DC.NavLayerCurrent] = nav_bb_rel
	}
}

func (c *Context) NavUpdateAnyRequestFlag() {
	c.NavAnyRequest = c.NavMoveRequest || c.NavInitRequest
}

// Scoring function for directional navigation. Based on https://gist.github.com/rygorous/6981057
func (c *Context) NavScoreItem(result *NavMoveResult, cand f64.Rectangle) bool {
	window := c.CurrentWindow
	if c.NavLayer != window.DC.NavLayerCurrent {
		return false
	}

	// Current modified source rect (NB: we've applied Max.x = Min.x in NavUpdate() to inhibit the effect of having varied item width)
	curr := &c.NavScoringRectScreen
	c.NavScoringCount++

	// We perform scoring on items bounding box clipped by their parent window on the other axis (clipping on our movement axis would give us equal scores for all clipped items)
	if c.NavMoveDir == DirLeft || c.NavMoveDir == DirRight {
		cand.Min.Y = f64.Clamp(cand.Min.Y, window.ClipRect.Min.Y, window.ClipRect.Max.Y)
		cand.Max.Y = f64.Clamp(cand.Max.Y, window.ClipRect.Min.Y, window.ClipRect.Max.Y)
	} else {
		cand.Min.X = f64.Clamp(cand.Min.X, window.ClipRect.Min.X, window.ClipRect.Max.X)
		cand.Max.X = f64.Clamp(cand.Max.X, window.ClipRect.Min.X, window.ClipRect.Max.X)
	}

	// Compute distance between boxes
	// FIXME-NAV: Introducing biases for vertical navigation, needs to be removed.
	dbx := c.NavScoreItemDistInterval(cand.Min.X, cand.Max.X, curr.Min.X, curr.Max.X)
	// Scale down on Y to keep using box-distance for vertically touching items
	dby := c.NavScoreItemDistInterval(
		f64.Lerp(0.2, cand.Min.Y, cand.Max.Y),
		f64.Lerp(0.8, cand.Min.Y, cand.Max.Y),
		f64.Lerp(0.2, curr.Min.Y, curr.Max.Y),
		f64.Lerp(0.8, curr.Min.Y, curr.Max.Y),
	)
	if dby != 0 && dbx != 0 {
		if dbx > 0 {
			dbx = dbx/1000 + 1
		} else {
			dbx = dbx/1000 - 1
		}
	}
	dist_box := math.Abs(dbx) + math.Abs(dby)

	// Compute distance between centers (this is off by a factor of 2, but we only compare center distances with each other so it doesn't matter)
	dcx := (cand.Min.X + cand.Max.X) - (curr.Min.X + curr.Max.X)
	dcy := (cand.Min.Y + cand.Max.Y) - (curr.Min.Y + curr.Max.Y)
	dist_center := math.Abs(dcx) + math.Abs(dcy) // L1 metric (need this for our connectedness guarantee)

	// Determine which quadrant of 'curr' our candidate item 'cand' lies in based on distance
	var quadrant Dir
	var dax, day, dist_axial float64
	if dbx != 0 || dby != 0 {
		// For non-overlapping boxes, use distance between boxes
		dax = dbx
		day = dby
		dist_axial = dist_box
		quadrant = c.NavScoreItemGetQuadrant(dbx, dby)
	} else if dcx != 0 || dcy != 0 {
		// For overlapping boxes with different centers, use distance between centers
		dax = dcx
		day = dcy
		dist_axial = dist_center
		quadrant = c.NavScoreItemGetQuadrant(dcx, dcy)
	} else {
		// Degenerate case: two overlapping buttons with same center, break ties arbitrarily (note that LastItemId here is really the _previous_ item order, but it doesn't matter)
		if window.DC.LastItemId < c.NavId {
			quadrant = DirLeft
		} else {
			quadrant = DirRight
		}
	}

	// Is it in the quadrant we're interesting in moving to?
	new_best := false
	if quadrant == c.NavMoveDir {
		// Does it beat the current best candidate?
		if dist_box < result.DistBox {
			result.DistBox = dist_box
			result.DistCenter = dist_center
			return true
		}

		if dist_box == result.DistBox {
			// Try using distance between center points to break ties
			if dist_center < result.DistCenter {
				result.DistCenter = dist_center
				new_best = true
			} else if dist_center == result.DistCenter {
				// Still tied! we need to be extra-careful to make sure everything gets linked properly. We consistently break ties by symbolically moving "later" items
				// (with higher index) to the right/downwards by an infinitesimal amount since we the current "best" button already (so it must have a lower index),
				// this is fairly easy. This rule ensures that all buttons with dx==dy==0 will end up being linked in order of appearance along the x axis.

				// moving bj to the right/down decreases distance
				if c.NavMoveDir == DirUp || c.NavMoveDir == DirDown {
					new_best = dby < 0
				} else {
					new_best = dbx < 0
				}
			}
		}
	}

	// Axial check: if 'curr' has no link at all in some direction and 'cand' lies roughly in that direction, add a tentative link. This will only be kept if no "real" matches
	// are found, so it only augments the graph produced by the above method using extra links. (important, since it doesn't guarantee strong connectedness)
	// This is just to avoid buttons having no links in a particular direction when there's a suitable neighbor. you get good graphs without this too.
	// 2017/09/29: FIXME: This now currently only enabled inside menu bars, ideally we'd disable it everywhere. Menus in particular need to catch failure. For general navigation it feels awkward.
	// Disabling it may however lead to disconnected graphs when nodes are very spaced out on different axis. Perhaps consider offering this as an option?

	// Check axial match
	if result.DistBox == math.MaxFloat32 && dist_axial < result.DistAxial {
		if c.NavLayer == 1 && c.NavWindow.Flags&WindowFlagsChildMenu == 0 {
			if (c.NavMoveDir == DirRight && dax > 0) || (c.NavMoveDir == DirRight && dax > 0) ||
				(c.NavMoveDir == DirUp && day < 0) || (c.NavMoveDir == DirDown && day > 0) {
				result.DistAxial = dist_axial
				new_best = true
			}
		}
	}

	return new_best
}

func (c *Context) NavScoreItemDistInterval(a0, a1, b0, b1 float64) float64 {
	if a1 < b0 {
		return a1 - b0
	}
	if b1 < a0 {
		return a0 - b1
	}
	return 0
}

func (c *Context) NavScoreItemGetQuadrant(dx, dy float64) Dir {
	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			return DirRight
		}
		return DirLeft
	}

	if dy > 0 {
		return DirDown
	}
	return DirUp
}

func (c *Context) NavProcessMoveRequestWrapAround(window *Window) {
	if c.NavWindow == window && c.NavMoveRequestButNoResultYet() {
		if (c.NavMoveDir == DirUp || c.NavMoveDir == DirDown) &&
			c.NavMoveRequestForward == NavForwardNone && c.NavLayer == 0 {
			c.NavMoveRequestForward = NavForwardForwardQueued
			c.NavMoveRequestCancel()

			c.NavWindow.NavRectRel[0].Min.Y = 0
			if c.NavMoveDir == DirUp {
				c.NavWindow.NavRectRel[0].Min.Y = math.Max(window.SizeFull.Y, window.SizeContents.Y) - window.Scroll.Y
			} else {
				c.NavWindow.NavRectRel[0].Min.Y = -window.Scroll.Y
			}

			c.NavWindow.NavRectRel[0].Max.Y = c.NavWindow.NavRectRel[0].Min.Y
		}
	}
}

func (c *Context) NavMoveRequestButNoResultYet() bool {
	return c.NavMoveRequest && c.NavMoveResultLocal.ID == 0 && c.NavMoveResultOther.ID == 0
}

func (c *Context) NavMoveRequestCancel() {
	c.NavMoveRequest = false
	c.NavUpdateAnyRequestFlag()
}

func (c *Context) NavRestoreLayer(layer int) {
	c.NavLayer = layer
	if layer == 0 {
		c.NavWindow = c.NavRestoreLastChildNavWindow(c.NavWindow)
	}
	if layer == 0 && c.NavWindow.NavLastIds[0] != 0 {
		c.SetNavIDAndMoveMouse(c.NavWindow.NavLastIds[0], layer, c.NavWindow.NavRectRel[0])
	} else {
		c.NavInitWindow(c.NavWindow, true)
	}
}

func (c *Context) SetNavIDAndMoveMouse(id ID, nav_layer int, rect_rel f64.Rectangle) {
	c.SetNavID(id, nav_layer)
	c.NavWindow.NavRectRel[nav_layer] = rect_rel
	c.NavMousePosDirty = true
	c.NavDisableHighlight = false
	c.NavDisableMouseHover = true
}

func (c *Context) SetNavID(id ID, nav_layer int) {
	c.NavId = id
	c.NavWindow.NavLastIds[nav_layer] = id
}

func (c *Context) NavInitWindow(window *Window, force_reinit bool) {
	var init_for_nav bool
	if window.Flags&WindowFlagsNoNavInputs == 0 {
		if window.Flags&WindowFlagsChildWindow == 0 || window.Flags&WindowFlagsPopup != 0 ||
			window.NavLastIds[0] == 0 || force_reinit {
			init_for_nav = true
		}
	}

	if init_for_nav {
		c.SetNavID(0, c.NavLayer)
		c.NavInitRequest = true
		c.NavInitRequestFromMove = false
		c.NavInitResultId = 0
		c.NavInitResultRectRel = f64.Rectangle{}
		c.NavUpdateAnyRequestFlag()
	} else {
		c.NavId = window.NavLastIds[0]
	}
}

func (c *Context) NavCalcPreferredMousePos() f64.Vec2 {
	window := c.NavWindow
	if window == nil {
		return c.IO.MousePos
	}
	rect_rel := window.NavRectRel[c.NavLayer]
	pos := f64.Vec2{
		rect_rel.Min.X + math.Min(c.Style.FramePadding.X*4, rect_rel.Dx()),
		rect_rel.Max.Y - math.Min(c.Style.FramePadding.Y, rect_rel.Dy()),
	}
	visible_rect := c.GetViewportRect()

	// ImFloor() is important because non-integer mouse position application in back-end might be lossy and result in undesirable non-zero delta.
	pos.X = f64.Clamp(pos.X, visible_rect.Min.X, visible_rect.Max.X)
	pos.Y = f64.Clamp(pos.Y, visible_rect.Min.Y, visible_rect.Max.Y)

	pos.X = math.Floor(pos.X)
	pos.Y = math.Floor(pos.Y)
	return pos
}