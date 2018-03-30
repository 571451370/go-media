package imgui

type PopupPositionPolicy int

const (
	PopupPositionPolicyDefault PopupPositionPolicy = iota
	PopupPositionPolicyComboBox
)

func (c *Context) OpenPopup(str_id string) {
	c.OpenPopupEx(c.CurrentWindow.GetID(str_id))
}

// Mark popup as open (toggle toward open state).
// Popups are closed when user click outside, or activate a pressable item, or CloseCurrentPopup() is called within a BeginPopup()/EndPopup() block.
// Popup identifiers are relative to the current ID-stack (so OpenPopup and BeginPopup needs to be at the same level).
// One open popup per level of the popup hierarchy (NB: when assigning we reset the Window member of ImGuiPopupRef to NULL)
func (c *Context) OpenPopupEx(id ID) {
	parent_window := c.CurrentWindow
	current_stack_size := len(c.CurrentPopupStack)
	// Tagged as new ref as Window will be set back to NULL if we write this into OpenPopupStack.
	var popup_ref PopupRef
	popup_ref.PopupId = id
	popup_ref.Window = nil
	popup_ref.ParentWindow = parent_window
	popup_ref.OpenFrameCount = c.FrameCount
	popup_ref.OpenParentId = parent_window.IDStack[len(parent_window.IDStack)-1]
	popup_ref.OpenMousePos = c.IO.MousePos
	popup_ref.OpenPopupPos = c.IO.MousePos
	if !c.NavDisableHighlight && c.NavDisableMouseHover {
		c.NavCalcPreferredMousePos()
	}

	if len(c.OpenPopupStack) < current_stack_size+1 {
		c.OpenPopupStack = append(c.OpenPopupStack, popup_ref)
	} else {
		// Close child popups if any
		c.OpenPopupStack = c.OpenPopupStack[:current_stack_size+1]

		// Gently handle the user mistakenly calling OpenPopup() every frame. It is a programming mistake! However, if we were to run the regular code path, the ui
		// would become completely unusable because the popup will always be in hidden-while-calculating-size state _while_ claiming focus. Which would be a very confusing
		// situation for the programmer. Instead, we silently allow the popup to proceed, it will keep reappearing and the programming error will be more obvious to understand.
		if c.OpenPopupStack[current_stack_size].PopupId == id && c.OpenPopupStack[current_stack_size].OpenFrameCount == c.FrameCount-1 {
			c.OpenPopupStack[current_stack_size].OpenFrameCount = popup_ref.OpenFrameCount
		} else {
			c.OpenPopupStack[current_stack_size] = popup_ref
		}

		// When reopening a popup we first refocus its parent, otherwise if its parent is itself a popup it would get closed by ClosePopupsOverWindow().
		// This is equivalent to what ClosePopupToLevel() does.
		//if (g.OpenPopupStack[current_stack_size].PopupId == id)
		//    FocusWindow(parent_window);
	}
}

func (c *Context) IsPopupOpen(str_id string) bool {
	return len(c.OpenPopupStack) > len(c.CurrentPopupStack) && c.OpenPopupStack[len(c.CurrentPopupStack)].PopupId == c.CurrentWindow.GetID(str_id)
}

func (c *Context) EndPopup() {
	// Make all menus and popups wrap around for now, may need to expose that policy.
	c.NavProcessMoveRequestWrapAround(c.CurrentWindow)
	c.End()
}

func (c *Context) ClosePopupsOverWindow(ref_window *Window) {
	if len(c.OpenPopupStack) == 0 {
		return
	}

	// When popups are stacked, clicking on a lower level popups puts focus back to it and close popups above it.
	// Don't close our own child popup windows.
	var n int
	if ref_window != nil {
		for n = range c.OpenPopupStack {
			popup := &c.OpenPopupStack[n]
			if popup.Window == nil {
				continue
			}
			if popup.Window.Flags&WindowFlagsChildWindow != 0 {
				continue
			}

			// Trim the stack if popups are not direct descendant of the reference window (which is often the NavWindow)
			has_focus := false
			for m := n; m < len(c.OpenPopupStack) && !has_focus; m++ {
				has_focus = c.OpenPopupStack[m].Window != nil && c.OpenPopupStack[m].Window.RootWindow == ref_window.RootWindow
			}
			if !has_focus {
				break
			}
		}
	}

	// This test is not required but it allows to set a convenient breakpoint on the block below
	if n < len(c.OpenPopupStack) {
		c.ClosePopupToLevel(n)
	}
}

func (c *Context) ClosePopupToLevel(remaining int) {
	var focus_window *Window
	if remaining > 0 {
		focus_window = c.OpenPopupStack[remaining-1].Window
	} else {
		focus_window = c.OpenPopupStack[0].ParentWindow
	}

	if c.NavLayer == 0 {
		focus_window = c.NavRestoreLastChildNavWindow(focus_window)
	}
	c.FocusWindow(focus_window)
	focus_window.DC.NavHideHighlightOneFrame = true
	c.OpenPopupStack = c.OpenPopupStack[:remaining]
}