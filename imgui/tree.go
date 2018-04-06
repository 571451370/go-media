package imgui

import (
	"math"

	"github.com/qeedquan/go-media/math/f64"
)

type TreeNodeFlags int

const (
	TreeNodeFlagsSelected          TreeNodeFlags = 1 << 0  // Draw as selected
	TreeNodeFlagsFramed            TreeNodeFlags = 1 << 1  // Full colored frame (e.g. for CollapsingHeader)
	TreeNodeFlagsAllowItemOverlap  TreeNodeFlags = 1 << 2  // Hit testing to allow subsequent widgets to overlap this one
	TreeNodeFlagsNoTreePushOnOpen  TreeNodeFlags = 1 << 3  // Don't do a TreePush() when open (e.g. for CollapsingHeader) = no extra indent nor pushing on ID stack
	TreeNodeFlagsNoAutoOpenOnLog   TreeNodeFlags = 1 << 4  // Don't automatically and temporarily open node when Logging is active (by default logging will automatically open tree nodes)
	TreeNodeFlagsDefaultOpen       TreeNodeFlags = 1 << 5  // Default node to be open
	TreeNodeFlagsOpenOnDoubleClick TreeNodeFlags = 1 << 6  // Need double-click to open node
	TreeNodeFlagsOpenOnArrow       TreeNodeFlags = 1 << 7  // Only open when clicking on the arrow part. If TreeNodeFlagsOpenOnDoubleClick is also set single-click arrow or double-click all box to open.
	TreeNodeFlagsLeaf              TreeNodeFlags = 1 << 8  // No collapsing no arrow (use as a convenience for leaf nodes).
	TreeNodeFlagsBullet            TreeNodeFlags = 1 << 9  // Display a bullet instead of arrow
	TreeNodeFlagsFramePadding      TreeNodeFlags = 1 << 10 // Use FramePadding (even for an unframed text node) to vertically align text baseline to regular widget height. Equivalent to calling AlignTextToFramePadding().
	//ImGuITreeNodeFlags_SpanAllAvailWidth TreeNodeFlags = 1 << 11  // FIXME: TODO: Extend hit box horizontally even if not framed
	//TreeNodeFlagsNoScrollOnOpen   TreeNodeFlags  = 1 << 12  // FIXME: TODO: Disable automatic scroll on TreePop() if node got just open and contents is not visible
	TreeNodeFlagsNavLeftJumpsBackHere TreeNodeFlags = 1 << 13 // (WIP) Nav: left direction may move to this TreeNode() from any of its child (items submitted between TreeNode and TreePop)
	TreeNodeFlagsCollapsingHeader     TreeNodeFlags = TreeNodeFlagsFramed | TreeNodeFlagsNoAutoOpenOnLog
)

type ItemHoveredDataBackup struct {
	LastItemId          ID
	LastItemStatusFlags ItemStatusFlags
	LastItemRect        f64.Rectangle
	LastItemDisplayRect f64.Rectangle
}

func (b *ItemHoveredDataBackup) Backup(c *Context) {
	window := c.CurrentWindow
	b.LastItemId = window.DC.LastItemId
	b.LastItemStatusFlags = window.DC.LastItemStatusFlags
	b.LastItemRect = window.DC.LastItemRect
	b.LastItemDisplayRect = window.DC.LastItemDisplayRect
}

func (b *ItemHoveredDataBackup) Restore(c *Context) {
	window := c.CurrentWindow
	window.DC.LastItemId = b.LastItemId
	window.DC.LastItemStatusFlags = b.LastItemStatusFlags
	window.DC.LastItemRect = b.LastItemRect
	window.DC.LastItemDisplayRect = b.LastItemDisplayRect
}

func (c *Context) CollapsingHeader(label string, flags TreeNodeFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}
	return c.TreeNodeBehavior(window.GetID(label), flags|TreeNodeFlagsCollapsingHeader|TreeNodeFlagsNoTreePushOnOpen, label)
}

func (c *Context) CollapsingHeaderEx(label string, p_open *bool, flags TreeNodeFlags) bool {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return false
	}

	if p_open != nil && !*p_open {
		return false
	}

	id := window.GetID(label)
	if p_open != nil {
		flags |= TreeNodeFlagsAllowItemOverlap
	}
	is_open := c.TreeNodeBehavior(id, flags|TreeNodeFlagsCollapsingHeader|TreeNodeFlagsNoTreePushOnOpen, label)
	if p_open != nil {
		// Create a small overlapping close button // FIXME: We can evolve this into user accessible helpers to add extra buttons on title bars, headers, etc.
		button_sz := c.FontSize * 0.5
		pos := f64.Vec2{
			math.Min(window.DC.LastItemRect.Max.X, window.ClipRect.Max.X) - c.Style.FramePadding.X - button_sz,
			window.DC.LastItemRect.Min.Y + c.Style.FramePadding.Y + button_sz,
		}
		var last_item_backup ItemHoveredDataBackup
		last_item_backup.Backup(c)
		// TODO
		if c.CloseButton(id+1, pos, button_sz) {
			*p_open = false
		}
		last_item_backup.Restore(c)
	}
	return is_open
}

func (c *Context) TreeNode(label string) bool {
	return false
}

func (c *Context) TreeNodeBehavior(id ID, flags TreeNodeFlags, label string) bool {
	return false
}