package imgui

import "github.com/qeedquan/go-media/math/f64"

type DrawCornerFlags uint

const (
	DrawCornerFlagsTopLeft  DrawCornerFlags = 1 << 0                                            // 0x1
	DrawCornerFlagsTopRight DrawCornerFlags = 1 << 1                                            // 0x2
	DrawCornerFlagsBotLeft  DrawCornerFlags = 1 << 2                                            // 0x4
	DrawCornerFlagsBotRight DrawCornerFlags = 1 << 3                                            // 0x8
	DrawCornerFlagsTop      DrawCornerFlags = DrawCornerFlagsTopLeft | DrawCornerFlagsTopRight  // 0x3
	DrawCornerFlagsBot      DrawCornerFlags = DrawCornerFlagsBotLeft | DrawCornerFlagsBotRight  // 0xC
	DrawCornerFlagsLeft     DrawCornerFlags = DrawCornerFlagsTopLeft | DrawCornerFlagsBotLeft   // 0x5
	DrawCornerFlagsRight    DrawCornerFlags = DrawCornerFlagsTopRight | DrawCornerFlagsBotRight // 0xA
	DrawCornerFlagsAll      DrawCornerFlags = 0xF                                               // In your function calls you may use ~0 (= all bits sets) instead of DrawCornerFlagsAll, as a convenience
)

type DrawDataBuilder struct {
	Layers [2][]*DrawList // Global layers for: regular, tooltip
}

type DrawList struct {
	CmdBuffer []DrawCmd // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []DrawCmd // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []DrawCmd // Vertex buffer.
}

type DrawCallback func(parentList *DrawList, cmd *DrawCmd)

type DrawCmd struct {
	ElemCount    int          // Number of indices (multiple of 3) to be rendered as triangles. Vertices are stored in the callee ImDrawList's vtx_buffer[] array, indices in idx_buffer[].
	ClipRect     f64.Vec4     // Clipping rectangle (x1, y1, x2, y2)
	TextureId    TextureID    // User-provided texture ID. Set by user in ImfontAtlas::SetTexID() for fonts or passed to Image*() functions. Ignore if never using images or multiple fonts atlas.
	UserCallback DrawCallback // If != nil, call the function instead of rendering the vertices. clip_rect and texture_id will be set normally.
}

type TextureID int

// RenderFrame renders a rectangle shaped with optional rounding and borders
func (c *Context) RenderFrame(pmin, pmax f64.Vec2, fillCol uint32, border bool, rounding float64) {
	window := c.GetCurrentWindow()
	style := c.GetStyle()

	drawList := window.DrawList
	drawList.AddRectFilled(pmin, pmax, fillCol, rounding)

	borderSize := style.FrameBorderSize
	if border && borderSize > 0 {
		one := f64.Vec2{1, 1}
		drawList.AddRect(pmin.Add(one), pmax.Add(one), fillCol, rounding, DrawCornerFlagsAll, 1)
		drawList.AddRect(pmin, pmax, fillCol, rounding, DrawCornerFlagsAll, 1)
	}
}

func (d *DrawList) AddLine(pmin, pmax f64.Vec2, col uint32, thickness float64) {
}

func (d *DrawList) AddRectFilled(pmin, pmax f64.Vec2, fillCol uint32, rounding float64) {
}

func (d *DrawList) AddRect(pmin, pmax f64.Vec2, col uint32, rounding float64, roundingCornerFlags DrawCornerFlags, thickness float64) {
}