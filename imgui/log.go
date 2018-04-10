package imgui

import (
	"os"

	"github.com/qeedquan/go-media/math/f64"
)

func (c *Context) LogText(text string) {
}

func (c *Context) LogRenderedText(ref_pos *f64.Vec2, text string) {
}

func (c *Context) LogToClipboard() {
	c.LogToClipboardEx(-1)
}

func (c *Context) LogToClipboardEx(max_depth int) {
	if c.LogEnabled {
		return
	}
	window := c.CurrentWindow

	assert(c.LogFile == nil)
	c.LogFile = nil
	c.LogEnabled = true
	c.LogStartDepth = window.DC.TreeDepth
	if max_depth >= 0 {
		c.LogAutoExpandMaxDepth = max_depth
	}
}

func (c *Context) LogToTTY() {
	c.LogToTTYEx(-1)
}

func (c *Context) LogToTTYEx(max_depth int) {
	if c.LogEnabled {
		return
	}
	window := c.CurrentWindow

	assert(c.LogFile == nil)
	c.LogFile = os.Stdout
	c.LogEnabled = true
	c.LogStartDepth = window.DC.TreeDepth
	if max_depth >= 0 {
		c.LogAutoExpandMaxDepth = max_depth
	}
}