package imgui

import "github.com/qeedquan/go-media/math/f64"

func (c *Context) NavProcessItem(window *Window, navBB f64.Rectangle, id ID) {
	itemFlags := window.DC.ItemFlags
	navBBRel := f64.Rectangle{navBB.Min.Sub(window.Pos), navBB.Max.Sub(window.Pos)}
	if c.NavInitRequest && c.NavLayer == window.DC.NavLayerCurrent {
		// Even if 'ImGuiItemFlags_NoNavDefaultFocus' is on (typically collapse/close button) we record the first ResultId so they can be used as a fallback
		if itemFlags&ItemFlagsNoNavDefaultFocus == 0 || c.NavInitResultId == 0 {
			c.NavInitResultId = id
			c.NavInitResultRectRel = navBBRel
		}
		if itemFlags&ItemFlagsNoNavDefaultFocus == 0 {
			c.NavInitRequest = false // Found a match, clear request
			c.NavUpdateAnyRequestFlag()
		}
	}

	// Scoring for navigation
	if c.NavId != id && itemFlags&ItemFlagsNoNav == 0 {
		result := &c.NavMoveResultOther
		if window == c.NavWindow {
			result = &c.NavMoveResultLocal
		}
		newBest := c.NavMoveRequest && c.NavScoreItem(result, navBB)
		if newBest {
			result.Id = id
			result.ParentId = window.IdStack[len(window.IdStack)-1]
			result.Window = window
			result.RectRel = navBBRel
		}
	}

	// Update window-relative bounding box of navigated item
	if c.NavId == id {
		c.NavWindow = window // Always refresh g.NavWindow, because some operations such as FocusItem() don't have a window.
		c.NavLayer = window.DC.NavLayerCurrent
		c.NavIdIsAlive = true
		c.NavIdTabCounter = window.FocusIdxTabCounter
		window.NavRectRel[window.DC.NavLayerCurrent] = navBBRel // Store item bounding box (relative to window position)
	}
}

func (c *Context) NavUpdateAnyRequestFlag() {
	c.NavAnyRequest = c.NavMoveRequest || c.NavInitRequest
}

func (c *Context) NavScoreItem(result *NavMoveResult, cand f64.Rectangle) bool {
	return false
}