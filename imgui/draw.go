package imgui

import (
	"image/color"

	"github.com/qeedquan/go-media/math/f64"
)

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

type DrawListFlags uint

const (
	DrawListFlagsAntiAliasedLines DrawListFlags = 1 << iota
	DrawListFlagsAntiAliasedFill
)

type DrawDataBuilder struct {
	Layers [2][]*DrawList // Global layers for: regular, tooltip
}

type DrawListSharedData struct {
	TexUvWhitePixel      f64.Vec2 // UV of white pixel in the atlas
	Font                 *Font    // Current/default font (optional, for simplified AddText overload)
	FontSize             float64  // Current/default font size (optional, for simplified AddText overload)
	CurveTessellationTol float64
	ClipRectFullscreen   f64.Vec4 // Value for PushClipRectFullscreen()
}

type DrawList struct {
	CmdBuffer []DrawCmd // Draw commands. Typically 1 command = 1 GPU draw call, unless the command is a callback.
	IdxBuffer []DrawCmd // Index buffer. Each command consume ImDrawCmd::ElemCount of those
	VtxBuffer []DrawCmd // Vertex buffer.

	VtxWritePtr   []DrawCmd
	VtxCurrentIdx int
	IdxWritePtr   []DrawCmd

	Flags DrawListFlags       // Flags, you may poke into these to adjust anti-aliasing settings per-primitives
	Data  *DrawListSharedData // Pointer to shared draw data (you can use ImGui::GetDrawListSharedData() to get the one from current ImGui context)
	Path  []f64.Vec2
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
func (c *Context) RenderFrame(pmin, pmax f64.Vec2, fillCol color.RGBA, border bool, rounding float64) {
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

func (d *DrawList) AddLine(pmin, pmax f64.Vec2, col color.RGBA, thickness float64) {
	if col.A == 0 {
		return
	}
	half := f64.Vec2{0.5, 0.5}
	d.PathLineTo(pmin.Add(half))
	d.PathLineTo(pmax.Add(half))
	d.PathStroke(col, false, thickness)
}

func (d *DrawList) AddRectFilled(pmin, pmax f64.Vec2, fillCol color.RGBA, rounding float64) {
}

func (d *DrawList) AddRect(pmin, pmax f64.Vec2, col color.RGBA, rounding float64, roundingCornerFlags DrawCornerFlags, thickness float64) {
}

func (d *DrawList) AddPolyline(points []f64.Vec2, col color.RGBA, closed bool, thickness float64) {
	pointsCount := len(points)
	if pointsCount < 2 {
		return
	}

	count := pointsCount
	if !closed {
		count--
	}

	thickLine := thickness > 1
	if d.Flags&DrawListFlagsAntiAliasedLines != 0 {
		// anti-aliased stroke

		idxCount := count * 12
		vtxCount := pointsCount * 3
		thickLineCount := 3
		if thickLine {
			idxCount = count * 18
			vtxCount = pointsCount * 4
			thickLineCount = 5
		}
		_, _ = idxCount, vtxCount

		tempNormals := make([]f64.Vec2, pointsCount*thickLineCount)

		for i1 := 0; i1 < count; i1++ {
			i2 := (i1 + 1) % count
			diff := points[i2].Sub(points[i1])
			diff = diff.Normalize()
			tempNormals[i1].X = diff.Y
			tempNormals[i1].Y = -diff.X
		}
		if !closed {
			tempNormals[pointsCount-1] = tempNormals[pointsCount-2]
		}

		if !thickLine {
		} else {
		}

	} else {
		// non anti-aliased stroke
	}
}

func (d *DrawList) AddText(pos f64.Vec2, col color.RGBA, text string) {
}

func (d *DrawList) AddImage(id TextureID, pmin, pmax, uvmin, uvmax f64.Vec2, col color.RGBA) {
}

func (d *DrawList) PathClear() {
	d.Path = d.Path[:0]
}

func (d *DrawList) PathLineTo(pos f64.Vec2) {
	d.Path = append(d.Path, pos)
}

func (d *DrawList) PathStroke(col color.RGBA, closed bool, thickness float64) {
	d.AddPolyline(d.Path, col, closed, thickness)
}
