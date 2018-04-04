package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type SelectableFlags int

const (
	SelectableFlagsDontClosePopups  SelectableFlags = 1 << 0 // Clicking this don't close parent popup window
	SelectableFlagsSpanAllColumns   SelectableFlags = 1 << 1 // Selectable frame can span all columns (text will still fit in current column)
	SelectableFlagsAllowDoubleClick SelectableFlags = 1 << 2 // Generate press events on double clicks too

	// NB: need to be in sync with last value of ImGuiSelectableFlags_
	SelectableFlagsMenu               SelectableFlags = 1 << 3 // -> PressedOnClick
	SelectableFlagsMenuItem           SelectableFlags = 1 << 4 // -> PressedOnRelease
	SelectableFlagsDisabled           SelectableFlags = 1 << 5
	SelectableFlagsDrawFillAvailWidth SelectableFlags = 1 << 6
)

func (c *Context) BeginMenuBar() bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	if window.Flags&WindowFlagsMenuBar == 0 {
		return false
	}
	assert(!window.DC.MenuBarAppending)

	// Save position
	c.BeginGroup()
	c.PushID("##menubar")

	// We don't clip with regular window clipping rectangle as it is already set to the area below. However we clip with window full rect.
	// We remove 1 worth of rounding to Max.x to that text in long menus don't tend to display over the lower-right rounded area, which looks particularly glitchy.
	bar_rect := window.MenuBarRect()
	clip_rect := f64.Rect(
		math.Floor(bar_rect.Min.X+0.5),
		math.Floor(bar_rect.Min.Y+window.WindowBorderSize+0.5),
		math.Floor(math.Max(bar_rect.Min.X, bar_rect.Max.X-window.WindowRounding)+0.5),
		math.Floor(bar_rect.Max.Y+0.5),
	)
	clip_rect = clip_rect.Intersect(window.WindowRectClipped)
	c.PushClipRect(clip_rect.Min, clip_rect.Max, false)

	window.DC.CursorPos = f64.Vec2{bar_rect.Min.X + window.DC.MenuBarOffsetX, bar_rect.Min.Y}
	window.DC.LayoutType = LayoutTypeHorizontal
	window.DC.NavLayerCurrent++
	window.DC.NavLayerCurrentMask <<= 1
	window.DC.MenuBarAppending = true
	c.AlignTextToFramePadding()

	return true
}

func (c *Context) EndMenuBar() {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	// Nav: When a move request within one of our child menu failed, capture the request to navigate among our siblings.
	if c.NavMoveRequestButNoResultYet() && (c.NavMoveDir == DirLeft || c.NavMoveDir == DirRight) && c.NavWindow.Flags&WindowFlagsChildMenu == 0 {
		nav_earliest_child := c.NavWindow
		if nav_earliest_child.ParentWindow != nil && (nav_earliest_child.ParentWindow.Flags&WindowFlagsChildMenu) != 0 {
			nav_earliest_child = nav_earliest_child.ParentWindow
		}
		if nav_earliest_child.ParentWindow == window && nav_earliest_child.DC.ParentLayoutType == LayoutTypeHorizontal && c.NavMoveRequestForward == NavForwardNone {
			// To do so we claim focus back, restore NavId and then process the movement request for yet another frame.
			// This involve a one-frame delay which isn't very problematic in this situation. We could remove it by scoring in advance for multiple window (probably not worth the hassle/cost)
			assert(window.DC.NavLayerActiveMaskNext&0x02 != 0) // Sanity Check
			c.FocusWindow(window)
			c.SetNavIDWithRectRel(window.NavLastIds[1], 1, window.NavRectRel[1])
			c.NavLayer = 1
			// Hide highlight for the current frame so we don't see the intermediary selection.
			c.NavDisableHighlight = true
			c.NavMoveRequestForward = NavForwardForwardQueued
			c.NavMoveRequestCancel()
		}
	}

	assert(window.Flags&WindowFlagsMenuBar != 0)
	assert(window.DC.MenuBarAppending)
	c.PopClipRect()
	c.PopID()
	window.DC.MenuBarOffsetX = window.DC.CursorPos.X - window.MenuBarRect().Min.X
	window.DC.GroupStack[len(window.DC.GroupStack)-1].AdvanceCursor = false
	c.EndGroup()
	window.DC.LayoutType = LayoutTypeVertical
	window.DC.NavLayerCurrent--
	window.DC.NavLayerCurrentMask >>= 1
	window.DC.MenuBarAppending = false
}

func (c *Context) BeginMenu(label string, enabled bool) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	style := &c.Style
	id := window.GetID(label)
	label_size := c.CalcTextSizeEx(label, true, -1)

	menu_is_open := c.IsPopupOpen(label)
	menuset_is_open := window.Flags&WindowFlagsPopup == 0 && len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].OpenParentId == window.IDStack[len(window.IDStack)-1]
	backed_nav_window := c.NavWindow
	if menuset_is_open {
		// Odd hack to allow hovering across menus of a same menu-set (otherwise we wouldn't be able to hover parent)
		c.NavWindow = window
	}

	// The reference position stored in popup_pos will be used by Begin() to find a suitable position for the child menu (using FindBestPopupWindowPos).
	pos := window.DC.CursorPos
	if window.DC.LayoutType == LayoutTypeHorizontal {
		// Menu inside an horizontal menu bar
		// Selectable extend their highlight by half ItemSpacing in each direction.
		// For ChildMenu, the popup position will be overwritten by the call to FindBestPopupWindowPos() in Begin()
		popup_pos := f64.Vec2{pos.X - window.WindowPadding.X, pos.Y - style.FramePadding.Y + window.MenuBarHeight()}
		window.DC.CursorPos.X += float64(int(style.ItemSpacing.X * 0.5))
		c.PushStyleVar(StyleVarItemSpacing, style.ItemSpacing.Scale(2.0))
		w := label_size.X
		select_flags := SelectableFlagsMenu | SelectableFlagsDontClosePopups
		if !enabled {
			select_flags |= SelectableFlagsDisabled
		}
		pressed := c.SelectableEx(label, menu_is_open, select_flags, f64.Vec2{0, 0})
		if !enabled {
		}
		c.PopStyleVar()
		// -1 spacing to compensate the spacing added when Selectable() did a SameLine(). It would also work to call SameLine() ourselves after the PopStyleVar().
		window.DC.CursorPos.X += float64(int(style.ItemSpacing.X * (-1.0 + 0.5)))
		_, _, _, _, _ = id, backed_nav_window, popup_pos, w, pressed
	} else {
		// Menu inside a menu
		// TODO
	}

	return false
}

// Tip: pass an empty label (e.g. "##dummy") then you can use the space to draw other text or image.
// But you need to make sure the ID is unique, e.g. enclose calls in PushID/PopID.
func (c *Context) Selectable(label string, selected bool, flags SelectableFlags, size_arg f64.Vec2) bool {
	return c.SelectableEx(label, selected, 0, f64.Vec2{0, 0})
}

func (c *Context) SelectableEx(label string, selected bool, flags SelectableFlags, size_arg f64.Vec2) bool {
	return false
}

func (c *Context) EndMenu() {
	// Nav: When a left move request _within our child menu_ failed, close the menu.
	// A menu doesn't close itself because EndMenuBar() wants the catch the last Left<>Right inputs.
	// However it means that with the current code, a BeginMenu() from outside another menu or a menu-bar won't be closable with the Left direction.
	window := c.CurrentWindow
	if c.NavWindow != nil && c.NavWindow.ParentWindow == window && c.NavMoveDir == DirLeft && c.NavMoveRequestButNoResultYet() && window.DC.LayoutType == LayoutTypeVertical {
		c.ClosePopupToLevel(len(c.OpenPopupStack) - 1)
		c.NavMoveRequestCancel()
	}
	c.EndPopup()
}

func (m *MenuColumns) Update(count int, spacing float64, clear bool) {
	m.Count = count
	m.Width = 0
	m.NextWidth = 0
	m.Spacing = spacing
	if clear {
		for i := range m.NextWidths {
			m.NextWidths[i] = 0
		}
	}
	for i := 0; i < m.Count; i++ {
		if i > 0 && m.NextWidths[i] > 0 {
			m.Width += m.Spacing
		}
		m.Pos[i] = float64(int(m.Width))
		m.Width += m.NextWidths[i]
		m.NextWidths[i] = 0
	}
}