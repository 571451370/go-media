// +build ignore
package main

import (
	"log"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/qeedquan/go-media/imgui"
	"github.com/qeedquan/go-media/math/f64"
	"github.com/qeedquan/go-media/sdl"
)

type UI struct {
	MousePressed [3]bool
	FontTexture  uint32
	Time         float64
	MouseCursors [imgui.MouseCursorCOUNT]*sdl.Cursor

	ShaderHandle uint32
	VertHandle   uint32
	FragHandle   uint32

	AttribLocationTex      int32
	AttribLocationProjMtx  int32
	AttribLocationPosition int32
	AttribLocationUV       int32
	AttribLocationColor    int32

	VboHandle      uint32
	ElementsHandle uint32

	SliderValue       float64
	ShowSimpleWindow  bool
	ShowDemoWindow    bool
	ShowAnotherWindow bool
	Counter           int

	ShowAppMainMenuBar       bool
	ShowAppConsole           bool
	ShowAppLog               bool
	ShowAppLayout            bool
	ShowAppPropertyEditor    bool
	ShowAppLongText          bool
	ShowAppAutoResize        bool
	ShowAppConstrainedResize bool
	ShowAppFixedOverlay      bool
	ShowAppWindowTitles      bool
	ShowAppCustomRendering   bool
	ShowAppStyleEditor       bool
	ShowAppMetrics           bool
	ShowAppAbout             bool

	NoTitlebar  bool
	NoScrollbar bool
	NoMenu      bool
	NoMove      bool
	NoResize    bool
	NoCollapse  bool
	NoNav       bool
	NoClose     bool

	MenuOptions struct {
		Enabled    bool
		Float      float64
		Check      bool
		ComboItems int
	}

	Widgets struct {
		Clicked int
		Check   bool
	}

	Colors struct {
		OutputOnlyModified bool
		OutputDest         int
		AlphaFlags         int
	}
}

var (
	im     *imgui.Context
	ui     UI
	window *sdl.Window
)

func main() {
	runtime.LockOSThread()
	log.SetFlags(0)
	log.SetPrefix("")
	initSDL()
	initGL()
	initIM()
	for {
		event()
		newFrame()
		render()
	}
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initSDL() {
	// Setup SDL
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_TIMER)
	ck(err)

	// Setup window
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	_, err = sdl.GetCurrentDisplayMode(0)
	ck(err)

	window, err = sdl.CreateWindow(
		"ImGui SDL2+OpenGL3 example",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		1280, 720,
		sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
	)
	ck(err)

	window.CreateContextGL()
	sdl.GLSetSwapInterval(1)
}

func initGL() {
	err := gl.Init()
	ck(err)
}

func initIM() {
	// Setup ImGui binding
	im = imgui.CreateContext()

	// Setup style
	im.StyleColorsDark(nil)

	// Setup back-end capabilities flags
	io := im.GetIO()
	// We can honor GetMouseCursor() values (optional)
	io.BackendFlags |= imgui.BackendFlagsHasMouseCursors

	// Keyboard mapping. ImGui will use those indices to peek into the io.KeyDown[] array.
	io.KeyMap[imgui.KeyTab] = sdl.SCANCODE_TAB
	io.KeyMap[imgui.KeyLeftArrow] = sdl.SCANCODE_LEFT
	io.KeyMap[imgui.KeyRightArrow] = sdl.SCANCODE_RIGHT
	io.KeyMap[imgui.KeyUpArrow] = sdl.SCANCODE_UP
	io.KeyMap[imgui.KeyDownArrow] = sdl.SCANCODE_DOWN
	io.KeyMap[imgui.KeyPageUp] = sdl.SCANCODE_PAGEUP
	io.KeyMap[imgui.KeyPageDown] = sdl.SCANCODE_PAGEDOWN
	io.KeyMap[imgui.KeyHome] = sdl.SCANCODE_HOME
	io.KeyMap[imgui.KeyEnd] = sdl.SCANCODE_END
	io.KeyMap[imgui.KeyInsert] = sdl.SCANCODE_INSERT
	io.KeyMap[imgui.KeyDelete] = sdl.SCANCODE_DELETE
	io.KeyMap[imgui.KeyBackspace] = sdl.SCANCODE_BACKSPACE
	io.KeyMap[imgui.KeySpace] = sdl.SCANCODE_SPACE
	io.KeyMap[imgui.KeyEnter] = sdl.SCANCODE_RETURN
	io.KeyMap[imgui.KeyEscape] = sdl.SCANCODE_ESCAPE
	io.KeyMap[imgui.KeyA] = sdl.SCANCODE_A
	io.KeyMap[imgui.KeyC] = sdl.SCANCODE_C
	io.KeyMap[imgui.KeyV] = sdl.SCANCODE_V
	io.KeyMap[imgui.KeyX] = sdl.SCANCODE_X
	io.KeyMap[imgui.KeyY] = sdl.SCANCODE_Y
	io.KeyMap[imgui.KeyZ] = sdl.SCANCODE_Z

	ui.MouseCursors[imgui.MouseCursorArrow] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)
	ui.MouseCursors[imgui.MouseCursorTextInput] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
	ui.MouseCursors[imgui.MouseCursorResizeAll] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEALL)
	ui.MouseCursors[imgui.MouseCursorResizeNS] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENS)
	ui.MouseCursors[imgui.MouseCursorResizeEW] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZEWE)
	ui.MouseCursors[imgui.MouseCursorResizeNESW] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENESW)
	ui.MouseCursors[imgui.MouseCursorResizeNWSE] = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_SIZENWSE)

	ui.ShowSimpleWindow = true
	ui.ShowDemoWindow = true
	ui.MenuOptions.Enabled = true
	ui.MenuOptions.Float = 0.5
}

func evKey(key int, down bool) {
	mod := sdl.GetModState()
	io := im.GetIO()
	io.KeysDown[key] = down
	io.KeyShift = mod&sdl.KMOD_SHIFT != 0
	io.KeyCtrl = mod&sdl.KMOD_CTRL != 0
	io.KeyAlt = mod&sdl.KMOD_ALT != 0
	io.KeySuper = mod&sdl.KMOD_GUI != 0
}

func event() {
	io := im.GetIO()
	for {
		ev := sdl.PollEvent()
		if ev == nil {
			break
		}
		switch ev := ev.(type) {
		case sdl.QuitEvent:
			os.Exit(0)
		case sdl.KeyDownEvent:
			switch ev.Sym {
			case sdl.K_ESCAPE:
				os.Exit(0)
			}
			evKey(int(ev.Scancode), true)
		case sdl.KeyUpEvent:
			evKey(int(ev.Scancode), false)
		case sdl.MouseWheelEvent:
			if ev.X > 0 {
				io.MouseWheelH += 1
			} else if ev.X < 0 {
				io.MouseWheelH -= 1
			}
			if ev.Y > 0 {
				io.MouseWheel += 1
			} else if ev.Y < 0 {
				io.MouseWheel -= 1
			}
		case sdl.MouseButtonDownEvent:
			switch ev.Button {
			case sdl.BUTTON_LEFT:
				ui.MousePressed[0] = true
			case sdl.BUTTON_RIGHT:
				ui.MousePressed[1] = true
			case sdl.BUTTON_MIDDLE:
				ui.MousePressed[2] = true
			}
		}
	}
}

func render() {
	// 1. Show a simple window.
	// Tip: if we don't call ImGui::Begin()/ImGui::End() the widgets automatically appears in a window called "Debug".
	showSimpleWindow()
	showAnotherWindow()
	showDemoWindow()

	// Rendering
	clearColor := f64.Vec4{0.45, 0.55, 0.60, 1.00}
	io := im.GetIO()
	gl.Viewport(0, 0, int32(io.DisplaySize.X), int32(io.DisplaySize.Y))
	gl.ClearColor(float32(clearColor.X), float32(clearColor.Y), float32(clearColor.Z), float32(clearColor.W))
	gl.Clear(gl.COLOR_BUFFER_BIT)
	im.Render()
	renderDrawData(im.GetDrawData())
	window.SwapGL()
}

func showSimpleWindow() {
	if !ui.ShowSimpleWindow {
		return
	}
	im.Text("Hello World!")
	im.SliderFloat("float", &ui.SliderValue, 0, 1)
	im.Checkbox("Demo Window", &ui.ShowDemoWindow)
	im.Checkbox("Another Window", &ui.ShowAnotherWindow)
	if im.Button("Button") {
		ui.Counter++
	}
	im.SameLine()
	im.Text("counter = %d", ui.Counter)
	im.Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0/im.GetIO().Framerate, im.GetIO().Framerate)
}

func showAnotherWindow() {
	if !ui.ShowAnotherWindow {
		return
	}

	im.BeginEx("Another Window", &ui.ShowAnotherWindow, 0)
	im.Text("Hello from another window!")
	if im.Button("Close Me") {
		ui.ShowAnotherWindow = false
	}
	im.End()
}

func newFrame() {
	if ui.FontTexture == 0 {
		createDeviceObjects()
	}

	io := im.GetIO()

	// Setup display size (every frame to accommodate for window resizing)
	w, h := window.Size()
	dw, dh := window.DrawableSizeGL()
	io.DisplaySize = f64.Vec2{float64(w), float64(h)}
	fw, fh := 0.0, 0.0
	if w > 0 {
		fw = float64(dw) / float64(w)
	}
	if h > 0 {
		fh = float64(dh) / float64(h)
	}
	io.DisplayFramebufferScale = f64.Vec2{float64(fw), float64(fh)}

	// Setup time step (we don't use SDL_GetTicks() because it is using millisecond resolution)
	frequency := sdl.GetPerformanceFrequency()
	currentTime := sdl.GetPerformanceCounter()
	if ui.Time > 0 {
		io.DeltaTime = (float64(currentTime) - ui.Time) / float64(frequency)
	} else {
		io.DeltaTime = 1.0 / 60
	}
	ui.Time = float64(currentTime)

	// Setup mouse inputs (we already got mouse wheel, keyboard keys & characters from our event handler)
	mx, my, button := sdl.GetMouseState()

	// If a mouse press event came, always pass it as "mouse held this frame", so we don't miss click-release events that are shorter than 1 frame.
	io.MouseDown[0] = ui.MousePressed[0] || (button&sdl.BUTTON(sdl.BUTTON_LEFT)) != 0
	io.MouseDown[1] = ui.MousePressed[1] || (button&sdl.BUTTON(sdl.BUTTON_RIGHT)) != 0
	io.MouseDown[2] = ui.MousePressed[2] || (button&sdl.BUTTON(sdl.BUTTON_MIDDLE)) != 0
	ui.MousePressed[0] = false
	ui.MousePressed[1] = false
	ui.MousePressed[2] = false

	// We need to use SDL_CaptureMouse() to easily retrieve mouse coordinates outside of the client area.
	if window.Flags()&(sdl.WINDOW_MOUSE_FOCUS|sdl.WINDOW_MOUSE_CAPTURE) != 0 {
		io.MousePos = f64.Vec2{float64(mx), float64(my)}
	}
	anyMouseButtonDown := false
	for i := range io.MouseDown {
		if io.MouseDown[i] {
			anyMouseButtonDown = true
		}
	}
	if anyMouseButtonDown && window.Flags()&sdl.WINDOW_MOUSE_CAPTURE == 0 {
		sdl.CaptureMouse(true)
	}
	if !anyMouseButtonDown && window.Flags()&sdl.WINDOW_MOUSE_CAPTURE != 0 {
		sdl.CaptureMouse(false)
	}

	// Update OS/hardware mouse cursor if imgui isn't drawing a software cursor
	if io.ConfigFlags&imgui.ConfigFlagsNoMouseCursorChange == 0 {
		cursor := im.GetMouseCursor()
		if io.MouseDrawCursor || cursor == imgui.MouseCursorNone {
			sdl.ShowCursor(0)
		} else {
			if ui.MouseCursors[cursor] != nil {
				sdl.SetCursor(ui.MouseCursors[cursor])
			} else {
				sdl.SetCursor(ui.MouseCursors[imgui.MouseCursorArrow])
			}
			sdl.ShowCursor(1)
		}
	}

	// Start the frame. This call will update the io.WantCaptureMouse, io.WantCaptureKeyboard flag that you can use to dispatch inputs (or not) to your application.
	im.NewFrame()
}

func createDeviceObjects() {
	var lastTexture, lastArrayBuffer, lastVertexArray int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)

	ui.ShaderHandle = gl.CreateProgram()

	vertexShader := `
	#version 330
	
	uniform mat4 ProjMtx;
	in vec2 Position;
	in vec2 UV;
	in vec4 Color;
	out vec2 Frag_UV;
	out vec4 Frag_Color;
	
	void main() {
		Frag_UV = UV;
		Frag_Color = Color;
		gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
	}
	` + "\x00"

	fragmentShader := `
	#version 330
	
	uniform sampler2D Texture;
	in vec2 Frag_UV;
	in vec4 Frag_Color;
	out vec4 Out_Color;
	
	void main() {
		Out_Color = Frag_Color * texture(Texture, Frag_UV.st);
	}
	` + "\x00"

	ui.VertHandle = gl.CreateShader(gl.VERTEX_SHADER)
	vsrc, free := gl.Strs(vertexShader)
	gl.ShaderSource(ui.VertHandle, 1, vsrc, nil)
	gl.CompileShader(ui.VertHandle)
	checkShaderCompileError(ui.VertHandle, gl.COMPILE_STATUS)
	free()

	ui.FragHandle = gl.CreateShader(gl.FRAGMENT_SHADER)
	fsrc, free := gl.Strs(fragmentShader)
	gl.ShaderSource(ui.FragHandle, 1, fsrc, nil)
	gl.CompileShader(ui.FragHandle)
	checkShaderCompileError(ui.FragHandle, gl.COMPILE_STATUS)
	free()

	gl.AttachShader(ui.ShaderHandle, ui.VertHandle)
	gl.AttachShader(ui.ShaderHandle, ui.FragHandle)
	gl.LinkProgram(ui.ShaderHandle)
	checkShaderLinkError(ui.ShaderHandle)

	ui.AttribLocationTex = gl.GetUniformLocation(ui.ShaderHandle, gl.Str("Texture\x00"))
	ui.AttribLocationProjMtx = gl.GetUniformLocation(ui.ShaderHandle, gl.Str("ProjMtx\x00"))
	ui.AttribLocationPosition = gl.GetAttribLocation(ui.ShaderHandle, gl.Str("Position\x00"))
	ui.AttribLocationUV = gl.GetAttribLocation(ui.ShaderHandle, gl.Str("UV\x00"))
	ui.AttribLocationColor = gl.GetAttribLocation(ui.ShaderHandle, gl.Str("Color\x00"))

	gl.GenBuffers(1, &ui.VboHandle)
	gl.GenBuffers(1, &ui.ElementsHandle)

	createFontsTexture()

	// Restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindVertexArray(uint32(lastVertexArray))
}

func createFontsTexture() {
	// Build texture atlas
	io := im.GetIO()
	// Load as RGBA 32-bits for OpenGL3 demo because it is more likely to be compatible with user's existing shader.
	pixels, width, height, _ := io.Fonts.GetTexDataAsRGBA32()

	// Upload texture to graphics system
	var last_texture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &last_texture)
	gl.GenTextures(1, &ui.FontTexture)
	gl.BindTexture(gl.TEXTURE_2D, ui.FontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, unsafe.Pointer(&pixels[0]))

	// Store our identifier
	io.Fonts.TexID = ui.FontTexture

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(last_texture))
}

func checkShaderCompileError(shader, typ uint32) {
	var status int32
	gl.GetShaderiv(shader, typ, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		str := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(str))
		panic(str)
	}
}

func checkShaderLinkError(program uint32) {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		str := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(program, logLength, nil, gl.Str(str))
		panic(str)
	}
}

func renderDrawData(draw_data *imgui.DrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	io := im.GetIO()
	fb_width := int(io.DisplaySize.X * io.DisplayFramebufferScale.X)
	fb_height := int(io.DisplaySize.Y * io.DisplayFramebufferScale.Y)
	if fb_width == 0 || fb_height == 0 {
		return
	}
	draw_data.ScaleClipRects(io.DisplayFramebufferScale)

	// Backup GL state
	var (
		last_program              int32
		last_active_texture       int32
		last_texture              int32
		last_sampler              int32
		last_array_buffer         int32
		last_element_array_buffer int32
		last_vertex_array         int32
		last_polygon_mode         [2]int32
		last_viewport             [4]int32
		last_scissor_box          [4]int32
		last_blend_src_rgb        int32
		last_blend_dst_rgb        int32
		last_blend_src_alpha      int32
		last_blend_dst_alpha      int32
		last_blend_equation_rgb   int32
		last_blend_equation_alpha int32
	)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &last_program)
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &last_active_texture)
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &last_texture)
	gl.GetIntegerv(gl.SAMPLER_BINDING, &last_sampler)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &last_array_buffer)
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &last_element_array_buffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &last_vertex_array)
	gl.GetIntegerv(gl.POLYGON_MODE, &last_polygon_mode[0])
	gl.GetIntegerv(gl.VIEWPORT, &last_viewport[0])
	gl.GetIntegerv(gl.SCISSOR_BOX, &last_scissor_box[0])
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &last_blend_src_rgb)
	gl.GetIntegerv(gl.BLEND_DST_RGB, &last_blend_dst_rgb)
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &last_blend_src_alpha)
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &last_blend_dst_alpha)
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &last_blend_equation_rgb)
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &last_blend_equation_alpha)
	last_enable_blend := gl.IsEnabled(gl.BLEND)
	last_enable_cull_face := gl.IsEnabled(gl.CULL_FACE)
	last_enable_depth_test := gl.IsEnabled(gl.DEPTH_TEST)
	last_enable_scissor_test := gl.IsEnabled(gl.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	ortho_projection := [4][4]float32{
		{2.0 / float32(io.DisplaySize.X), 0.0, 0.0, 0.0},
		{0.0, 2.0 / float32(-io.DisplaySize.Y), 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	gl.Viewport(0, 0, int32(fb_width), int32(fb_height))
	gl.UseProgram(ui.ShaderHandle)
	gl.Uniform1i(ui.AttribLocationTex, 0)
	gl.UniformMatrix4fv(ui.AttribLocationProjMtx, 1, false, &ortho_projection[0][0])

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	var vao_handle uint32
	gl.GenVertexArrays(1, &vao_handle)
	gl.BindVertexArray(vao_handle)
	gl.BindBuffer(gl.ARRAY_BUFFER, ui.VboHandle)
	gl.EnableVertexAttribArray(uint32(ui.AttribLocationPosition))
	gl.EnableVertexAttribArray(uint32(ui.AttribLocationUV))
	gl.EnableVertexAttribArray(uint32(ui.AttribLocationColor))
	sizeofDrawVert := 20
	sizeofDrawIdx := 4
	gl.VertexAttribPointer(uint32(ui.AttribLocationPosition), 2, gl.FLOAT, false, int32(sizeofDrawVert), unsafe.Pointer(uintptr(0)))
	gl.VertexAttribPointer(uint32(ui.AttribLocationUV), 2, gl.FLOAT, false, int32(sizeofDrawVert), unsafe.Pointer(uintptr(8)))
	gl.VertexAttribPointer(uint32(ui.AttribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(sizeofDrawVert), unsafe.Pointer(uintptr(16)))

	// Draw
	for n := range draw_data.CmdLists {
		cmd_list := draw_data.CmdLists[n]
		idx_buffer_offset := 0
		gl.BindBuffer(gl.ARRAY_BUFFER, ui.VboHandle)
		gl.BufferData(gl.ARRAY_BUFFER, len(cmd_list.VtxBuffer)*sizeofDrawVert, unsafe.Pointer(&cmd_list.VtxBuffer[0]), gl.STREAM_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ui.ElementsHandle)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(cmd_list.IdxBuffer)*sizeofDrawIdx, unsafe.Pointer(&cmd_list.IdxBuffer[0]), gl.STREAM_DRAW)

		for cmd_i := 0; cmd_i < len(cmd_list.CmdBuffer); cmd_i++ {
			pcmd := &cmd_list.CmdBuffer[cmd_i]
			if pcmd.UserCallback != nil {
				pcmd.UserCallback(cmd_list)
			} else {
				gl.BindTexture(gl.TEXTURE_2D, pcmd.TextureId.(uint32))
				gl.Scissor(
					int32(pcmd.ClipRect.X),
					int32(float64(fb_height)-pcmd.ClipRect.W),
					int32(pcmd.ClipRect.Z-pcmd.ClipRect.X),
					int32(pcmd.ClipRect.W-pcmd.ClipRect.Y),
				)
				gl.DrawElements(gl.TRIANGLES, int32(pcmd.ElemCount), gl.UNSIGNED_INT, unsafe.Pointer(uintptr(idx_buffer_offset)))
			}
			idx_buffer_offset += pcmd.ElemCount
		}
	}

	gl.DeleteVertexArrays(1, &vao_handle)

	// Restore modified GL state
	gl.UseProgram(uint32(last_program))
	gl.BindTexture(gl.TEXTURE_2D, uint32(last_texture))
	gl.BindSampler(0, uint32(last_sampler))
	gl.ActiveTexture(uint32(last_active_texture))
	gl.BindVertexArray(uint32(last_vertex_array))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(last_array_buffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(last_element_array_buffer))
	gl.BlendEquationSeparate(uint32(last_blend_equation_rgb), uint32(last_blend_equation_alpha))
	gl.BlendFuncSeparate(uint32(last_blend_src_rgb), uint32(last_blend_dst_rgb), uint32(last_blend_src_alpha), uint32(last_blend_dst_alpha))
	if last_enable_blend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	if last_enable_cull_face {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	if last_enable_depth_test {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	if last_enable_scissor_test {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(last_polygon_mode[0]))
	gl.Viewport(last_viewport[0], last_viewport[1], last_viewport[2], last_viewport[3])
	gl.Scissor(last_scissor_box[0], last_scissor_box[1], last_scissor_box[2], last_scissor_box[3])
}

func showDemoWindow() {
	if !ui.ShowDemoWindow {
		return
	}
	im.SetNextWindowPos(f64.Vec2{650, 20}, imgui.CondFirstUseEver, f64.Vec2{0, 0})
	p_open := &ui.ShowDemoWindow

	// Demonstrate the various window flags. Typically you would just use the default.
	if ui.ShowAppMainMenuBar {
		showExampleAppMainMenuBar()
	}
	if ui.ShowAppConsole {
		showExampleAppConsole()
	}
	if ui.ShowAppLog {
		showExampleAppLog()
	}
	if ui.ShowAppLayout {
		showExampleAppLayout()
	}
	if ui.ShowAppPropertyEditor {
		showExampleAppPropertyEditor()
	}
	if ui.ShowAppLongText {
		showExampleAppLongText()
	}
	if ui.ShowAppAutoResize {
		showExampleAppAutoResize()
	}
	if ui.ShowAppConstrainedResize {
		showExampleAppConstrainedResize()
	}
	if ui.ShowAppFixedOverlay {
		showExampleAppFixedOverlay()
	}
	if ui.ShowAppWindowTitles {
		showExampleAppWindowTitles()
	}
	if ui.ShowAppCustomRendering {
		showExampleAppCustomRendering()
	}
	if ui.ShowAppAbout {
		showAppAbout()
	}
	// Demonstrate the various window flags. Typically you would just use the default.
	var window_flags imgui.WindowFlags
	if ui.NoTitlebar {
		window_flags |= imgui.WindowFlagsNoTitleBar
	}
	if ui.NoScrollbar {
		window_flags |= imgui.WindowFlagsNoScrollbar
	}
	if !ui.NoMenu {
		window_flags |= imgui.WindowFlagsMenuBar
	}
	if ui.NoMove {
		window_flags |= imgui.WindowFlagsNoMove
	}
	if ui.NoResize {
		window_flags |= imgui.WindowFlagsNoResize
	}
	if ui.NoCollapse {
		window_flags |= imgui.WindowFlagsNoCollapse
	}
	if ui.NoNav {
		window_flags |= imgui.WindowFlagsNoNav
	}
	if ui.NoClose {
		p_open = nil
	}

	im.SetNextWindowSize(f64.Vec2{550, 680}, imgui.CondFirstUseEver)
	if !im.BeginEx("ImGui Demo", p_open, window_flags) {
		// Early out if the window is collapsed, as an optimization.
		im.End()
		return
	}

	im.PushItemWidth(-140)
	im.Text("dear imgui says hello (%s)", im.GetVersion())

	// Menu
	if im.BeginMenuBar() {
		if im.BeginMenu("Menu") {
			showExampleMenuFile()
			im.EndMenu()
		}
		if im.BeginMenu("Examples") {
			im.MenuItemSelect("Main menu bar", "", &ui.ShowAppMainMenuBar)
			im.MenuItemSelect("Console", "", &ui.ShowAppConsole)
			im.MenuItemSelect("Log", "", &ui.ShowAppLog)
			im.MenuItemSelect("Simple layout", "", &ui.ShowAppLayout)
			im.MenuItemSelect("Property editor", "", &ui.ShowAppPropertyEditor)
			im.MenuItemSelect("Long text display", "", &ui.ShowAppLongText)
			im.MenuItemSelect("Auto-resizing window", "", &ui.ShowAppAutoResize)
			im.MenuItemSelect("Constrained-resizing window", "", &ui.ShowAppConstrainedResize)
			im.MenuItemSelect("Simple overlay", "", &ui.ShowAppFixedOverlay)
			im.MenuItemSelect("Manipulating window titles", "", &ui.ShowAppWindowTitles)
			im.MenuItemSelect("Custom rendering", "", &ui.ShowAppCustomRendering)
			im.EndMenu()
		}
		if im.BeginMenu("Help") {
			im.MenuItemSelect("Metrics", "", &ui.ShowAppMetrics)
			im.MenuItemSelect("Style Editor", "", &ui.ShowAppStyleEditor)
			im.MenuItemSelect("About Dear ImGui", "", &ui.ShowAppAbout)
			im.EndMenu()
		}

		im.EndMenuBar()
	}

	im.Spacing()
	if im.CollapsingHeader("Help") {
		im.TextWrapped("This window is being created by the ShowDemoWindow() function. Please refer to the code in imgui_demo.cpp for reference.\n\n")
		im.Text("USER GUIDE:")
		showUserGuide()
	}

	if im.CollapsingHeader("Window options") {
		im.Checkbox("No titlebar", &ui.NoTitlebar)
		im.SameLineEx(150, -1)
		im.Checkbox("No scrollbar", &ui.NoScrollbar)
		im.SameLineEx(300, -1)
		im.Checkbox("No menu", &ui.NoMenu)
		im.Checkbox("No move", &ui.NoMove)
		im.SameLineEx(150, -1)
		im.Checkbox("No resize", &ui.NoResize)
		im.SameLineEx(300, -1)
		im.Checkbox("No collapse", &ui.NoCollapse)
		im.Checkbox("No close", &ui.NoClose)
		im.SameLineEx(150, -1)
		im.Checkbox("No nav", &ui.NoNav)
		if im.TreeNode("Style") {
			im.TreePop()
		}
		if im.TreeNode("Captured/Logging") {
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Widgets") {
		if im.TreeNode("Basic") {
			if im.Button("Button") {
				ui.Widgets.Clicked++
			}
			if ui.Widgets.Clicked != 0 {
				im.SameLine()
				im.Text("Thanks for clicking me!")
			}

			im.Checkbox("checkbox", &ui.Widgets.Check)
		}
	}
	if im.CollapsingHeader("Layout") {
	}
	if im.CollapsingHeader("Popups & Modal windows") {
	}
	if im.CollapsingHeader("Columns") {
	}
	if im.CollapsingHeader("Filtering") {
	}
	if im.CollapsingHeader("Inputs, Navigation & Focus") {
	}

	im.End()
}

func showUserGuide() {
	im.BulletText("Double-click on title bar to collapse window.")
	im.BulletText("Click and drag on lower right corner to resize window\n(double-click to auto fit window to its contents).")
	im.BulletText("Click and drag on any empty space to move window.")
	im.BulletText("TAB/SHIFT+TAB to cycle through keyboard editable fields.")
	im.BulletText("CTRL+Click on a slider or drag box to input value as text.")
	if im.GetIO().FontAllowUserScaling {
		im.BulletText("CTRL+Mouse Wheel to zoom window contents.")
	}
	im.BulletText("Mouse Wheel to scroll.")
	im.BulletText("While editing text:\n")
	im.Indent()
	im.BulletText("Hold SHIFT or use mouse to select text.")
	im.BulletText("CTRL+Left/Right to word jump.")
	im.BulletText("CTRL+A or double-click to select all.")
	im.BulletText("CTRL+X,CTRL+C,CTRL+V to use clipboard.")
	im.BulletText("CTRL+Z,CTRL+Y to undo/redo.")
	im.BulletText("ESCAPE to revert.")
	im.BulletText("You can apply arithmetic operators +,*,/ on numerical values.\nUse +- to subtract.")
	im.Unindent()
}

func showExampleAppMainMenuBar() {
}

func showExampleAppConsole() {
}

func showExampleAppLog() {
}

func showExampleAppLayout() {
}

func showExampleAppPropertyEditor() {
}

func showExampleAppLongText() {
}

func showExampleAppAutoResize() {
}

func showExampleAppConstrainedResize() {
}

func showExampleAppFixedOverlay() {
}

func showExampleAppWindowTitles() {
}

func showExampleAppCustomRendering() {
}

func showAppAbout() {
	im.BeginEx("About Dear ImGui", &ui.ShowAppAbout, imgui.WindowFlagsAlwaysAutoResize)
	im.Text("Dear ImGui, %s", im.GetVersion())
	im.Separator()
	im.Text("By Omar Cornut and all dear imgui contributors.")
	im.Text("Dear ImGui is licensed under the MIT License, see LICENSE for more information.")
	im.End()
}

func showExampleMenuFile() {
	im.MenuItemEx("(dummy menu)", "", false, false)
	im.MenuItem("New")
	im.MenuItemEx("Open", "Ctrl+O", false, true)
	if im.BeginMenu("Open Recent") {
		im.MenuItem("fish_hat.c")
		im.MenuItem("fish_hat.inl")
		im.MenuItem("fish_hat.h")
		if im.BeginMenu("More..") {
			im.MenuItem("Hello")
			im.MenuItem("Sailor")
			if im.BeginMenu("Recurse..") {
				showExampleMenuFile()
				im.EndMenu()
			}
			im.EndMenu()
		}
		im.EndMenu()
	}
	im.MenuItemEx("Save", "Ctrl+S", false, true)
	im.MenuItem("Save As..")
	im.Separator()
	if im.BeginMenu("Options") {
		im.MenuItemSelect("Enabled", "", &ui.MenuOptions.Enabled)
		im.BeginChildEx("child", f64.Vec2{0, 60}, true, 0)
		for i := 0; i < 10; i++ {
			im.Text("Scrolling Text %d", i)
		}
		im.EndChild()
		im.SliderFloat("Value", &ui.MenuOptions.Float, 0, 1)
		im.InputFloat("Input", &ui.MenuOptions.Float, 0.1)
		im.ComboString("Combo", &ui.MenuOptions.ComboItems, []string{"Yes", "No", "Maybe"})
		im.Checkbox("Check", &ui.MenuOptions.Check)
		im.EndMenu()
	}
	if im.BeginMenu("Colors") {
		sz := im.GetTextLineHeight()
		for i := imgui.Col(0); i < imgui.ColCOUNT; i++ {
			name := im.GetStyleColorName(i)
			p := im.GetCursorScreenPos()
			im.GetWindowDrawList().AddRectFilled(p, f64.Vec2{p.X + sz, p.Y + sz}, im.GetColorFromStyle(i))
			im.Dummy(f64.Vec2{sz, sz})
			im.SameLine()
			im.MenuItem(name)
		}
		im.EndMenu()
	}
	if im.BeginMenuEx("Disabled", false) {
		assert(false)
	}
	im.MenuItemEx("Checked", "Checked", true, true)
	im.MenuItemEx("Quit", "Alt+F4", false, true)
}

func assert(x bool) {
	if !x {
		panic("assert failed")
	}
}
