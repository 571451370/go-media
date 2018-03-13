package imgui

import "github.com/qeedquan/go-media/math/f64"

func (c *Context) Image(user_texture_id TextureID, size, uv0, uv1 f64.Vec2, tint_col, border_col f64.Vec4) {
	window := c.GetCurrentWindow()
	if window.SkipItems {
		return
	}

	bb := f64.Rectangle{
		window.DC.CursorPos,
		window.DC.CursorPos.Add(size),
	}
	if border_col.W > 0 {
		bb.Max = bb.Max.Add(f64.Vec2{2, 2})
	}
	c.ItemSizeBB(bb)
	if !c.ItemAdd(bb, 0) {
		return
	}

	if border_col.W > 0 {
		window.DrawList.AddRect(bb.Min, bb.Max, border_col.ToRGBA())
		window.DrawList.AddImageEx(user_texture_id, bb.Min.Add(f64.Vec2{1, 1}), bb.Max.Sub(f64.Vec2{1, 1}), uv0, uv1, tint_col.ToRGBA())
	} else {
		window.DrawList.AddImageEx(user_texture_id, bb.Min, bb.Max, uv0, uv1, tint_col.ToRGBA())
	}
}