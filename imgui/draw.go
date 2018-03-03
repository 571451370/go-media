package imgui

import "github.com/qeedquan/go-media/math/f64"

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