// +build ignore
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/qeedquan/go-media/image/chroma"
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

	ClearColor color.RGBA

	AppAutoResize struct {
		Lines int
	}

	AppConstrainedResize struct {
		AutoResize   bool
		Type         int
		DisplayLines int
	}

	AppFixedOverlay struct {
		Corner int
	}

	MenuOptions struct {
		Enabled    bool
		Float      float64
		Check      bool
		ComboItems int
	}

	AppLog struct {
		LastTime float64
	}

	Widgets struct {
		Clicked                        int
		Check                          bool
		AlignLabelWithCurrentXPosition bool
		ClosableGroup                  bool
		WrapWidth                      float64
		UTF8_Buf                       [32]byte
		Layout                         struct {
			Line              int
			DisableMouseWheel bool
			DisableMenu       bool
			WidgetsWidth      struct {
				F float64
			}
		}
	}

	Colors struct {
		OutputOnlyModified bool
		OutputDest         int
		AlphaFlags         int
	}

	MetricsWindow struct {
		ShowClipRects bool
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

	ui.MetricsWindow.ShowClipRects = true

	ui.ClearColor = f64.Vec4{0.45, 0.55, 0.60, 1.00}.ToRGBA()
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
	clearColor := chroma.RGBA2VEC4(ui.ClearColor)
	// Rendering
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
	im.Text("Hello, world!")
	im.SliderFloat("float", &ui.SliderValue, 0, 1)

	im.ColorEdit3("clear color", &ui.ClearColor)
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
			idx_buffer_offset += pcmd.ElemCount * 4
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

	if ui.ShowAppMetrics {
		showMetricsWindow()
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
	im.Text("dear imgui says hello. (%s)", im.GetVersion())

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
			showStyleEditor()
			im.TreePop()
		}

		if im.TreeNode("Captured/Logging") {
			im.TextWrapped("The logging API redirects all text output so you can easily capture the content of a window or a block. Tree nodes can be automatically expanded. You can also call ImGui::LogText() to output directly to the log without a visual output.")
			im.LogButtons()
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
			im.TreePop()
		}
		if im.TreeNode("Trees") {
			if im.TreeNode("Basic trees") {
				for i := imgui.ID(0); i < 5; i++ {
					if im.TreeNodeIDEx(i, 0, "Child %d", i) {
						im.Text("blah blah")
						im.SameLine()
						if im.SmallButton("button") {
						}
						im.TreePop()
					}
				}
				im.TreePop()
			}

			if im.TreeNode("Advanced, with Selectable nodes") {
				showHelpMarker("This is a more standard looking tree with selectable nodes.\nClick to select, CTRL+Click to toggle, click on arrows or double-click to open.")
				im.Checkbox("Align label with current X position)", &ui.Widgets.AlignLabelWithCurrentXPosition)
				im.Text("Hello!")
				if ui.Widgets.AlignLabelWithCurrentXPosition {
					im.UnindentEx(im.GetTreeNodeToLabelSpacing())
				}
				im.TreePop()
			}

			im.TreePop()
		}

		if im.TreeNode("Collapsing Headers") {
			im.Checkbox("Enable extra group", &ui.Widgets.ClosableGroup)
			if im.CollapsingHeader("Header") {
				im.Text("IsItemHovered: %v", im.IsItemHovered())
				for i := 0; i < 5; i++ {
					im.Text("Some content %d", i)
				}
			}
			if im.CollapsingHeaderOpen("Header with a close button", &ui.Widgets.ClosableGroup) {
				im.Text("IsItemHovered: %v", im.IsItemHovered())
				for i := 0; i < 5; i++ {
					im.Text("More content %d", i)
				}
			}
			im.TreePop()
		}

		if im.TreeNode("Bullets") {
			im.BulletText("Bullet point 1")
			im.BulletText("Bullet point 2\nOn multiple lines")
			im.Bullet()
			im.Text("Bullet point 3 (two calls)")
			im.Bullet()
			im.SmallButton("Button")
			im.TreePop()
		}

		if im.TreeNode("Text") {
			if im.TreeNode("Colored Text") {
				// Using shortcut. You can use PushStyleColor()/PopStyleColor() for more flexibility.
				im.TextColored(color.RGBA{255, 0, 255, 255}, "Pink")
				im.TextColored(color.RGBA{255, 255, 0, 255}, "Yellow")
				im.TextDisabled("Disabled")
				im.SameLine()
				showHelpMarker("The TextDisabled color is stored in ImGuiStyle.")
				im.TreePop()
			}
			im.TreePop()
		}

		if im.TreeNode("Word Wrapping") {
			// Using shortcut. You can use PushTextWrapPos()/PopTextWrapPos() for more flexibility.
			im.TextWrapped("This text should automatically wrap on the edge of the window. The current implementation for text wrapping follows simple rules suitable for English and possibly other languages.")
			im.Spacing()

			im.SliderFloatEx("Wrap width", &ui.Widgets.WrapWidth, -20, 600, "%.0f", 1)

			im.Text("Test paragraph 1:")
			pos := im.GetCursorScreenPos()
			im.GetWindowDrawList().AddRectFilled(
				f64.Vec2{pos.X + ui.Widgets.WrapWidth, pos.Y},
				f64.Vec2{pos.X + ui.Widgets.WrapWidth + 10, pos.Y + im.GetTextLineHeight()},
				color.RGBA{255, 0, 255, 255},
			)
			im.PushTextWrapPos(im.GetCursorPos().X + ui.Widgets.WrapWidth)
			im.Text("The lazy dog is a good dog. This paragraph is made to fit within %.0f pixels. Testing a 1 character word. The quick brown fox jumps over the lazy dog.", ui.Widgets.WrapWidth)
			im.GetWindowDrawList().AddRect(im.GetItemRectMin(), im.GetItemRectMax(), color.RGBA{255, 255, 0, 255})
			im.PopTextWrapPos()

			im.Text("Test paragraph 2:")
			pos = im.GetCursorScreenPos()
			im.GetWindowDrawList().AddRectFilled(
				f64.Vec2{pos.X + ui.Widgets.WrapWidth, pos.Y},
				f64.Vec2{pos.X + ui.Widgets.WrapWidth + 10, pos.Y + im.GetTextLineHeight()},
				color.RGBA{255, 0, 255, 255},
			)
			im.PushTextWrapPos(im.GetCursorPos().X + ui.Widgets.WrapWidth)
			im.Text("aaaaaaaa bbbbbbbb, c cccccccc,dddddddd. d eeeeeeee   ffffffff. gggggggg!hhhhhhhh")
			im.GetWindowDrawList().AddRect(im.GetItemRectMin(), im.GetItemRectMax(), color.RGBA{255, 255, 0, 255})
			im.PopTextWrapPos()

			im.TreePop()
		}
		if im.TreeNode("UTF-8 Text") {
			// UTF-8 test with Japanese characters
			// (needs a suitable font, try Arial Unicode or M+ fonts http://mplus-fonts.sourceforge.jp/mplus-outline-fonts/index-en.html)
			// - From C++11 you can use the u8"my text" syntax to encode literal strings as UTF-8
			// - For earlier compiler, you may be able to encode your sources as UTF-8 (e.g. Visual Studio save your file as 'UTF-8 without signature')
			// - HOWEVER, FOR THIS DEMO FILE, BECAUSE WE WANT TO SUPPORT COMPILER, WE ARE *NOT* INCLUDING RAW UTF-8 CHARACTERS IN THIS SOURCE FILE.
			//   Instead we are encoding a few string with hexadecimal constants. Don't do this in your application!
			// Note that characters values are preserved even by InputText() if the font cannot be displayed, so you can safely copy & paste garbled characters into another application.
			im.TextWrapped("CJK text will only appears if the font was loaded with the appropriate CJK character ranges. Call io.Font->LoadFromFileTTF() manually to load extra character ranges.")
			im.Text("Hiragana: \xe3\x81\x8b\xe3\x81\x8d\xe3\x81\x8f\xe3\x81\x91\xe3\x81\x93 (kakikukeko)")
			im.Text("Kanjis: \xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e (nihongo)")
			if ui.Widgets.UTF8_Buf == [32]byte{} {
				copy(ui.Widgets.UTF8_Buf[:], "\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e") // "nihongo"
			}
			im.InputText("UTF-8 input", ui.Widgets.UTF8_Buf[:])
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Layout") {
		if im.TreeNode("Child regions") {
			im.Checkbox("Disable Mouse Wheel", &ui.Widgets.Layout.DisableMouseWheel)
			im.Checkbox("Disable Menu", &ui.Widgets.Layout.DisableMenu)
			im.Button("Goto")
			im.SameLine()
			im.PushItemWidth(100)
			im.InputIntEx("##Line", &ui.Widgets.Layout.Line, 0, 0, imgui.InputTextFlagsEnterReturnsTrue)
			im.PopItemWidth()

			// Child 1: no border, enable horizontal scrollbar
			{

			}

			// Child 2: rounded border
			{
			}
			im.TreePop()
		}

		if im.TreeNode("Widgets Width") {
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Popups & Modal windows") {
		if im.TreeNode("Popups") {
			im.TreePop()
		}
		if im.TreeNode("Context menus") {
			im.TreePop()
		}
		if im.TreeNode("Modals") {
			im.TreePop()
		}
		if im.TreeNode("Menus inside a regular window") {
			im.TreePop()
		}
	}
	if im.CollapsingHeader("Columns") {
		im.PushStringID("Columns")

		// Basic columns
		if im.TreeNode("Basic") {
			im.Text("Without border:")
			im.Columns(3, "mycolumns3", false) // 3-ways, no border
			im.Separator()
			for n := 0; n < 14; n++ {
				label := fmt.Sprintf("Item %d", n)
				if im.Selectable(label) {
				}
				im.NextColumn()
			}
			im.Columns(1, "", true)
			im.Separator()

			im.Text("With border:")
			im.Columns(4, "mycolumns", true) // 4-ways, with border
			im.Separator()

			im.Text("ID")
			im.NextColumn()
			im.Text("Name")
			im.NextColumn()
			im.Text("Path")
			im.NextColumn()
			im.Text("Hovered")
			im.NextColumn()
			im.Separator()

			names := []string{"One", "Two", "Three"}
			paths := []string{"/path/one", "/path/two", "/path/three"}
			selected := -1
			for i := range names {
				label := fmt.Sprintf("%04d", i)
				if im.SelectableEx(label, selected == i, imgui.SelectableFlagsSpanAllColumns, f64.Vec2{0, 0}) {
					selected = i
				}
				hovered := im.IsItemHovered()
				im.NextColumn()
				im.Text(names[i])
				im.NextColumn()
				im.Text(paths[i])
				im.NextColumn()
				im.Text("%v", hovered)
				im.NextColumn()
			}
			im.Columns(1, "", true)
			im.Separator()
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Filtering") {
		im.Text("Filter usage:\n" +
			"  \"\"         display all lines\n" +
			"  \"xxx\"      display lines containing \"xxx\"\n" +
			"  \"xxx,yyy\"  display lines containing \"xxx\" or \"yyy\"\n" +
			"  \"-xxx\"     hide lines containing \"xxx\"",
		)
		lines := []string{"aaa1.c", "bbb1.c", "ccc1.c", "aaa2.cpp", "bbb2.cpp", "ccc2.cpp", "abc.h", "hello, world"}
		for i := range lines {
			im.BulletText(lines[i])
		}
	}
	if im.CollapsingHeader("Inputs, Navigation & Focus") {
		io := im.GetIO()

		im.Text("WantCaptureMouse: %v", io.WantCaptureMouse)
		im.Text("WantCaptureKeyboard: %v", io.WantCaptureKeyboard)
		im.Text("WantTextInput: %v", io.WantTextInput)
		im.Text("WantSetMousePos: %v", io.WantSetMousePos)
		im.Text("NavActive: %v, NavVisible: %v", io.NavActive, io.NavVisible)

		im.Checkbox("io.MouseDrawCursor", &io.MouseDrawCursor)
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
	if im.BeginMainMenuBar() {
		if im.BeginMenu("File") {
			showExampleMenuFile()
			im.EndMenu()
		}
		if im.BeginMenu("Edit") {
			if im.MenuItemEx("Undo", "Ctrl+Z", false, true) {
			}
			if im.MenuItemEx("Redo", "Ctrl+Y", false, false) {
			}
			im.Separator()
			if im.MenuItemEx("Cut", "Ctrl+X", false, true) {
			}
			if im.MenuItemEx("Copy", "Ctrl+C", false, true) {
			}
			if im.MenuItemEx("Paste", "Ctrl+V", false, true) {
			}
			im.EndMenu()
		}
		im.EndMainMenuBar()
	}
}

func showExampleAppConsole() {
}

// Demonstrate creating a simple log window with basic filtering.
func showExampleAppLog() {
	// Demo: add random items (unless Ctrl is held)
}

func showExampleAppLayout() {
}

func showExampleAppPropertyEditor() {
}

func showExampleAppLongText() {
}

// Demonstrate creating a window which gets auto-resized according to its content.
func showExampleAppAutoResize() {
	if !im.BeginEx("Example: Auto-resizing window", &ui.ShowAppAutoResize, imgui.WindowFlagsAlwaysAutoResize) {
		im.End()
		return
	}

	im.Text("Window will resize every-frame to the size of its content.\nNote that you probably don't want to query the window size to\noutput your content because that would create a feedback loop.")
	im.SliderIntEx("Number of lines", &ui.AppAutoResize.Lines, 1, 20, "")
	for i := 0; i < ui.AppAutoResize.Lines; i++ {
		// Pad with space to extend size horizontally
		im.Text("%*sThis is line %d", i*4, "", i)
	}
	im.End()
}

// Demonstrate creating a window with custom resize constraints.
func showExampleAppConstrainedResize() {
	Square := func(data *imgui.SizeCallbackData) {
		data.DesiredSize = f64.Vec2{math.Max(data.DesiredSize.X, data.DesiredSize.Y), math.Max(data.DesiredSize.X, data.DesiredSize.Y)}
	}

	Step := func(data *imgui.SizeCallbackData) {
		step := 10.0
		data.DesiredSize = f64.Vec2{float64(int(data.DesiredSize.X/step+0.5)) * step, float64(int(data.DesiredSize.Y/step+0.5)) * step}
	}

	switch ui.AppConstrainedResize.Type {
	case 0:
		// Vertical only
		im.SetNextWindowSizeConstraints(f64.Vec2{-1, 0}, f64.Vec2{-1, math.MaxFloat32}, nil)
	case 1:
		// Horizontal only
		im.SetNextWindowSizeConstraints(f64.Vec2{0, -1}, f64.Vec2{math.MaxFloat32, -1}, nil)
	case 2:
		// Width > 100, Height > 100
		im.SetNextWindowSizeConstraints(f64.Vec2{100, 100}, f64.Vec2{math.MaxFloat32, math.MaxFloat32}, nil)
	case 3:
		// Width 400-500
		im.SetNextWindowSizeConstraints(f64.Vec2{400, -1}, f64.Vec2{500, -1}, nil)
	case 4:
		// Height 400-500
		im.SetNextWindowSizeConstraints(f64.Vec2{-1, 400}, f64.Vec2{-1, 500}, nil)
	case 5:
		// Always Square
		im.SetNextWindowSizeConstraints(f64.Vec2{0, 0}, f64.Vec2{math.MaxFloat32, math.MaxFloat32}, Square)
	case 6:
		// Fixed Step
		im.SetNextWindowSizeConstraints(f64.Vec2{0, 0}, f64.Vec2{math.MaxFloat32, math.MaxFloat32}, Step)
	}

	var flags imgui.WindowFlags
	if ui.AppConstrainedResize.AutoResize {
		flags = imgui.WindowFlagsAlwaysAutoResize
	}
	if im.BeginEx("Example: Constrained Resize", &ui.ShowAppAutoResize, flags) {
		desc := []string{
			"Resize vertical only",
			"Resize horizontal only",
			"Width > 100, Height > 100",
			"Width 400-500",
			"Height 400-500",
			"Custom: Always Square",
			"Custom: Fixed Steps (100)",
		}
		if im.Button("200x200") {
			im.SetCurrentWindowSize(f64.Vec2{200, 200}, 0)
		}
		im.SameLine()
		if im.Button("500x500") {
			im.SetCurrentWindowSize(f64.Vec2{500, 500}, 0)
		}
		im.SameLine()
		if im.Button("800x500") {
			im.SetCurrentWindowSize(f64.Vec2{800, 200}, 0)
		}
		im.PushItemWidth(200)
		im.ComboString("Constraint", &ui.AppConstrainedResize.Type, desc)
		im.DragIntEx("Lines", &ui.AppConstrainedResize.DisplayLines, 0.2, 1, 100, "")
		im.PopItemWidth()
		im.Checkbox("Auto-resize", &ui.AppConstrainedResize.AutoResize)
		for i := 0; i < ui.AppConstrainedResize.DisplayLines; i++ {
			im.Text("%*sHello, sailor! Making this line long enough for the example.", i*4, "")
		}
	}
	im.End()
}

// Demonstrate creating a simple static window with no decoration + a context-menu to choose which corner of the screen to use.
func showExampleAppFixedOverlay() {
	const DISTANCE = 10.0
	window_pos := f64.Vec2{DISTANCE, DISTANCE}
	window_pos_pivot := f64.Vec2{0, 0}
	if ui.AppFixedOverlay.Corner&1 != 0 {
		window_pos.X = im.GetIO().DisplaySize.X - DISTANCE
		window_pos_pivot.X = 1
	}
	if ui.AppFixedOverlay.Corner&2 != 0 {
		window_pos.Y = im.GetIO().DisplaySize.Y - DISTANCE
		window_pos_pivot.Y = 1
	}
	im.SetNextWindowPos(window_pos, imgui.CondAlways, window_pos_pivot)
	// Transparent background
	im.SetNextWindowBgAlpha(0.3)
	window_flags := imgui.WindowFlagsNoTitleBar | imgui.WindowFlagsNoResize | imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings | imgui.WindowFlagsNoFocusOnAppearing | imgui.WindowFlagsNoNav
	if im.BeginEx("Example: Fixed Overlay", &ui.ShowAppFixedOverlay, window_flags) {
		im.Text("Simple overlay\nin the corner of the screen.\n(right-click to change position)")
		im.Separator()
		im.Text("Mouse Position: (%.1f,%.1f)", im.GetIO().MousePos.X, im.GetIO().MousePos.Y)
		if im.BeginPopupContextWindow() {
			if im.MenuItemEx("Top-left", "", ui.AppFixedOverlay.Corner == 0, true) {
				ui.AppFixedOverlay.Corner = 0
			}
			if im.MenuItemEx("Top-right", "", ui.AppFixedOverlay.Corner == 1, true) {
				ui.AppFixedOverlay.Corner = 1
			}
			if im.MenuItemEx("Bottom-left", "", ui.AppFixedOverlay.Corner == 2, true) {
				ui.AppFixedOverlay.Corner = 2
			}
			if im.MenuItemEx("Bottom-right", "", ui.AppFixedOverlay.Corner == 3, true) {
				ui.AppFixedOverlay.Corner = 3
			}
			if im.MenuItem("Close") {
				ui.ShowAppFixedOverlay = false
			}
			im.EndPopup()
		}
		im.End()
	}
}

// Demonstrate using "##" and "###" in identifiers to manipulate ID generation.
// This apply to regular items as well. Read FAQ section "How can I have multiple widgets with the same label? Can I have widget without a label? (Yes). A primer on the purpose of labels/IDs." for details.
func showExampleAppWindowTitles() {
	// By default, Windows are uniquely identified by their title.
	// You can use the "##" and "###" markers to manipulate the display/ID.
	// Using "##" to display same title but have unique identifier.
	im.SetNextWindowPos(f64.Vec2{100, 100}, imgui.CondFirstUseEver, f64.Vec2{0, 0})
	im.Begin("Same title as another window##1")
	im.Text("This is window 1.\nMy title is the same as window 2, but my identifier is unique.")
	im.End()

	im.SetNextWindowPos(f64.Vec2{100, 200}, imgui.CondFirstUseEver, f64.Vec2{0, 0})
	im.Begin("Same title as another window##2")
	im.Text("This is window 2.\nMy title is the same as window 1, but my identifier is unique.")
	im.End()

	// Using "###" to display a changing title but keep a static identifier "AnimatedTitle"
	buf := fmt.Sprintf("Animated title %c %d###AnimatedTitle", "|/-\\"[int(im.GetTime()/0.25)&3], im.GetFrameCount())
	im.SetNextWindowPos(f64.Vec2{100, 300}, imgui.CondFirstUseEver, f64.Vec2{0, 0})
	im.Begin(buf)
	im.Text("This window has a changing title.")
	im.End()
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

func showStyleEditor() {
	style := im.GetStyle()

	im.PushItemWidth(im.GetWindowWidth() * 0.50)

	// Simplified Settings
	if im.SliderFloatEx("FrameRounding", &style.FrameRounding, 0.0, 12.0, "%.0f", 1.0) {
		// Make GrabRounding always the same as FrameRounding
		style.GrabRounding = style.FrameRounding
	}

	window_border := style.WindowBorderSize > 0
	if im.Checkbox("WindowBorder", &window_border) {
		style.WindowBorderSize = 0
		if window_border {
			style.WindowBorderSize = 1
		}
	}
	im.SameLine()

	frame_border := style.FrameBorderSize > 0
	if im.Checkbox("FrameBorder", &frame_border) {
		style.FrameBorderSize = 0
		if frame_border {
			style.FrameBorderSize = 1
		}
	}
	im.SameLine()

	popup_border := style.PopupBorderSize > 0
	if im.Checkbox("PopupBorder", &popup_border) {
		style.PopupBorderSize = 0
		if popup_border {
			style.PopupBorderSize = 1
		}
	}

	// Save/Revert button
	if im.Button("Save Ref") {
	}
	im.SameLine()
	if im.Button("Revert Ref") {
	}
	im.SameLine()
	showHelpMarker("Save/Revert in local non-persistent storage. Default Colors definition are not affected. Use \"Export Colors\" below to save them somewhere.")

	if im.TreeNode("Rendering") {
		im.Checkbox("Anti-aliased lines", &style.AntiAliasedLines)
		im.SameLine()
		showHelpMarker("When disabling anti-aliasing lines, you'll probably want to disable borders in your style as well.")
		im.PushItemWidth(100)
		im.DragFloatEx("Curve Tessellation Tolerance", &style.CurveTessellationTol, 0.02, 0.10, math.MaxFloat32, "", 2.0)
		if style.CurveTessellationTol < 0.0 {
			style.CurveTessellationTol = 0.10
		}
		// Not exposing zero here so user doesn't "lose" the UI (zero alpha clips all widgets). But application code could have a toggle to switch between zero and non-zero.
		im.DragFloatEx("Global Alpha", &style.Alpha, 0.005, 0.20, 1.0, "%.2f", 1.0)
		im.PopItemWidth()
		im.TreePop()
	}

	if im.TreeNode("Settings") {
		im.TreePop()
	}

	if im.TreeNode("Colors") {
		im.TreePop()
	}

	im.PopItemWidth()
}

func showHelpMarker(desc string) {
	im.TextDisabled("(?)")
	if im.IsItemHovered() {
		im.BeginTooltip()
		im.PushTextWrapPos(im.GetFontSize() * 35.0)
		im.TextUnformatted(desc)
		im.PopTextWrapPos()
		im.EndTooltip()
	}
}

func showMetricsWindow() {
	p_open := &ui.ShowAppMetrics
	if im.BeginEx("ImGui Metrics", p_open, 0) {
		im.Text("Dear ImGui %s", im.GetVersion())
		im.Text("Application average %.3f ms/frame (%.1f FPS)", 1000.0/im.GetIO().Framerate, im.GetIO().Framerate)
		im.Text("%d vertices, %d indices (%d triangles)", im.GetIO().MetricsRenderVertices, im.GetIO().MetricsRenderIndices, im.GetIO().MetricsRenderIndices/3)
		show_clip_rects := &ui.MetricsWindow.ShowClipRects
		im.Checkbox("Show clipping rectangles when hovering draw commands", show_clip_rects)
		im.Separator()

		// Access private state, we are going to display the draw lists from last frame
		if im.TreeNodeStringIDEx("DrawList", 0, "Active DrawLists (%d)", len(im.DrawDataBuilder.Layers[0])) {
			for i := range im.DrawDataBuilder.Layers[0] {
				_ = im.DrawDataBuilder.Layers[0][i]
			}
			im.TreePop()
		}
		if im.TreeNodeStringIDEx("Popups", 0, "Open Popups Stack (%d)", len(im.OpenPopupStack)) {
			for i := range im.OpenPopupStack {
				window := im.OpenPopupStack[i].Window
				window_name := "NULL"
				child_window := ""
				child_menu := ""
				if window != nil {
					window_name = window.Name
					if window.Flags&imgui.WindowFlagsChildWindow != 0 {
						child_window = " ChildWindow"
					}
					if window.Flags&imgui.WindowFlagsChildMenu != 0 {
						child_menu = " ChildMenu"
					}
				}

				im.BulletText("PopupID: %08x, Window: '%s'%s%s", im.OpenPopupStack[i].PopupId, window_name, child_window, child_menu)
			}
			im.TreePop()
		}
		if im.TreeNode("Internal state") {
			input_source_names := []string{"None", "Mouse", "Nav", "NavKeyboard", "NavGamepad"}
			assert(len(input_source_names) == int(imgui.InputSourceCOUNT))
			if im.HoveredWindow != nil {
				im.Text("HoveredWindow: '%s'", im.HoveredWindow.Name)
			} else {
				im.Text("HoveredWindow: '%s'", "NULL")
			}
			if im.HoveredRootWindow != nil {
				im.Text("HoveredRootWindow: '%s'", im.HoveredRootWindow.Name)
			} else {
				im.Text("HoveredRootWindow: '%s'", "NULL")
			}
			// Data is "in-flight" so depending on when the Metrics window is called we may see current frame information or not
			im.Text("ActiveId: 0x%08X/0x%08X (%.2f sec), ActiveIdSource: %s", im.ActiveId, im.ActiveIdPreviousFrame, im.ActiveIdTimer, input_source_names[im.ActiveIdSource])
			if im.MovingWindow != nil {
				im.Text("ActiveIdWindow: '%s'", im.ActiveIdWindow.Name)
			} else {
				im.Text("ActiveIdWindow: '%s'", "NULL")
			}
			if im.MovingWindow != nil {
				im.Text("MovingWindow: '%s'", im.MovingWindow.Name)
			} else {
				im.Text("MovingWindow: '%s'", "NULL")
			}
			if im.NavWindow != nil {
				im.Text("NavWindow: '%s'", im.NavWindow.Name)
			} else {
				im.Text("NavWindow: '%s'", "NULL")
			}
			im.Text("NavId: 0x%08X, NavLayer: %d", im.NavId, im.NavLayer)
			im.Text("NavInputSource: %s", input_source_names[im.NavInputSource])
			im.Text("NavActive: %v, NavVisible: %v", im.IO.NavActive, im.IO.NavVisible)
			im.Text("NavActivateId: 0x%08X, NavInputId: 0x%08X", im.NavActivateId, im.NavInputId)
			im.Text("NavDisableHighlight: %v, NavDisableMouseHover: %v", im.NavDisableHighlight, im.NavDisableMouseHover)
			im.Text("DragDrop: %v, SourceId = 0x%08X", im.DragDropActive, im.DragDropPayload.SourceId)
			im.TreePop()
		}
	}
	im.End()
}

func assert(x bool) {
	if !x {
		panic("assert failed")
	}
}
