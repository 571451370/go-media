// +build ignore
package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
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

	AppLayout struct {
		Selected int
	}

	MenuOptions struct {
		Enabled    bool
		Float      float64
		Check      bool
		ComboItems int
	}

	AppLog struct {
		LastTime float64
		Log      ExampleAppLog
	}

	AppLongText struct {
		TestType     int
		Log          strings.Builder
		Lines        int
		DummyMembers [8]float64
	}

	Widgets struct {
		Basic struct {
			Clicked int
			Check   bool
			E       int

			Arr []float64

			ItemsCurrent int

			Input struct {
				Str0       []byte
				I0         int
				F0, D0, F1 float64
				Vec4a      [4]float64
			}

			Drag struct {
				I1, I2 int
				F1, F2 float64
			}

			Slider struct {
				I1     int
				F1, F2 float64
				Angle  float64
			}

			ColorEdit struct {
				Col1, Col2 color.RGBA
			}

			ListBox struct {
				ItemCurrent int
			}
		}

		Trees struct {
			AlignLabelWithCurrentXPosition bool
		}

		CollapsingHeader struct {
			ClosableGroup bool
		}

		Text struct {
			WordWrapping struct {
				WrapWidth float64
			}
			UTF8 struct {
				Buf []byte
			}
		}

		Images struct {
			PressedCount int
		}

		Combo struct {
			Flags        uint
			ItemCurrent  int
			ItemCurrent2 int
			ItemCurrent3 int
			ItemCurrent4 int
		}

		Selectables struct {
			Basic struct {
				Selection [5]bool
			}
			Single struct {
				Selected int
			}
			Multiple struct {
				Selection [5]bool
			}
			Rendering struct {
				Selected [3]bool
			}
			Columns struct {
				Selected [16]bool
			}
			Grid struct {
				Selected [16]bool
			}
		}

		FilteredTextInput struct {
			Buf1    []byte
			Buf2    []byte
			Buf3    []byte
			Buf4    []byte
			Buf5    []byte
			Buf6    []byte
			Bufpass []byte
		}

		MultilineTextInput struct {
			ReadOnly bool
			Text     []byte
		}

		Plots struct {
			Animate      bool
			Values       [90]float64
			ValuesOffset int
			RefreshTime  float64
			Phase        float64

			FuncType     int
			DisplayCount int

			Progress    float64
			ProgressDir float64
		}

		ColorPicker struct {
			Color              color.RGBA
			AlphaPreview       bool
			AlphaHalfPreview   bool
			OptionsMenu        bool
			Hdr                bool
			SavedPaletteInited bool
			SavedPalette       [32]color.RGBA
			BackupColor        color.RGBA

			Alpha       bool
			AlphaBar    bool
			SidePreview bool
			RefColor    bool
			RefColorV   color.RGBA
			InputsMode  int
			PickerMode  int
		}

		Range struct {
			Begin  float64
			End    float64
			BeginI int
			EndI   int
		}

		MultiComponents struct {
			Vec4f [4]float64
			Vec4i [4]int
		}

		VerticalSliders struct {
			Spacing  float64
			IntValue int
			Values   []float64
			Values2  []float64
		}
	}

	Layout struct {
		ChildRegion struct {
			DisableMouseWheel bool
			DisableMenu       bool
			Line              int
		}
		WidgetsWidth struct {
			F float64
		}
		Horizontal struct {
			C1, C2, C3, C4 bool
			F0, F1, F2     float64
			Item           int
			Selection      [4]int
		}

		Scrolling struct {
			Track      bool
			TrackLine  int
			ScrollToPx int
		}

		HorizontalScrolling struct {
			Lines int
		}

		Clipping struct {
			Size   f64.Vec2
			Offset f64.Vec2
		}
	}

	PopupsModal struct {
		Popups struct {
			SelectedFish int
		}
		ContextMenus struct {
			Value float64
			Name  []byte
		}
		Modals struct {
			DontAskMeNextTime bool
			Item              int
			Color             color.RGBA
		}
	}

	Columns struct {
		MixedItems struct {
			Foo float64
			Bar float64
		}
		Borders struct {
			Horizontal bool
			Vertical   bool
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

	Input struct {
		Tabbing struct {
			Buf []byte
		}
		Focus struct {
			Buf []byte
			F3  [3]float64
		}
		FocusHovered struct {
			EmbedAllInsideAChildWindow bool
		}
	}

	AppCustomRendering struct {
		Sz         float64
		Col        color.RGBA
		AddingLine bool
		Points     []f64.Vec2
	}

	StyleEditor struct {
		Init          bool
		RefSavedStyle imgui.Style
		Colors        struct {
			OutputDest         int
			OutputOnlyModified bool
			Filter             imgui.TextFilter
			AlphaFlags         int
		}
		Fonts struct {
			WindowScale float64
		}
	}

	StyleSelector struct {
		StyleIdx int
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

	ui.ClearColor = f64.Vec4{0.45, 0.55, 0.60, 1.00}.ToRGBA()

	ui.Widgets.Basic.Arr = []float64{0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2}

	ui.MenuOptions.Enabled = true
	ui.MenuOptions.Float = 0.5

	ui.MetricsWindow.ShowClipRects = true

	ui.StyleSelector.StyleIdx = -1
	ui.StyleEditor.Colors.OutputDest = 0
	ui.StyleEditor.Colors.OutputOnlyModified = true
	ui.StyleEditor.Fonts.WindowScale = 1.0
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
	ShowSimpleWindow()
	ShowAnotherWindow()
	ShowDemoWindow()
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

func ShowSimpleWindow() {
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

func ShowAnotherWindow() {
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

func ShowDemoWindow() {
	if !ui.ShowDemoWindow {
		return
	}
	im.SetNextWindowPos(f64.Vec2{650, 20}, imgui.CondFirstUseEver, f64.Vec2{0, 0})
	p_open := &ui.ShowDemoWindow

	// Demonstrate the various window flags. Typically you would just use the default.
	if ui.ShowAppMainMenuBar {
		ShowExampleAppMainMenuBar()
	}
	if ui.ShowAppConsole {
		ShowExampleAppConsole()
	}
	if ui.ShowAppLog {
		ShowExampleAppLog()
	}
	if ui.ShowAppLayout {
		ShowExampleAppLayout()
	}
	if ui.ShowAppPropertyEditor {
		ShowExampleAppPropertyEditor()
	}
	if ui.ShowAppLongText {
		ShowExampleAppLongText()
	}
	if ui.ShowAppAutoResize {
		ShowExampleAppAutoResize()
	}
	if ui.ShowAppConstrainedResize {
		ShowExampleAppConstrainedResize()
	}
	if ui.ShowAppFixedOverlay {
		ShowExampleAppFixedOverlay()
	}
	if ui.ShowAppWindowTitles {
		ShowExampleAppWindowTitles()
	}
	if ui.ShowAppCustomRendering {
		ShowExampleAppCustomRendering()
	}

	if ui.ShowAppMetrics {
		ShowMetricsWindow()
	}
	if ui.ShowAppStyleEditor {
		im.BeginEx("Style Editor", &ui.ShowAppStyleEditor, 0)
		ShowStyleEditor(nil)
		im.End()
	}
	if ui.ShowAppAbout {
		ShowAppAbout()
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
			ShowExampleMenuFile()
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
		ShowUserGuide()
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
			ShowStyleEditor(nil)
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
			clicked := &ui.Widgets.Basic.Clicked
			if im.Button("Button") {
				*clicked++
			}
			if *clicked != 0 {
				im.SameLine()
				im.Text("Thanks for clicking me!")
			}
			check := &ui.Widgets.Basic.Check
			im.Checkbox("checkbox", check)

			e := &ui.Widgets.Basic.E
			im.RadioButtonEx("radio a", e, 0)
			im.SameLine()
			im.RadioButtonEx("radio b", e, 1)
			im.SameLine()
			im.RadioButtonEx("radio c", e, 2)

			// Color buttons, demonstrate using PushID() to add unique identifier in the ID stack, and changing style.
			for i := 0; i < 7; i++ {
				if i > 0 {
					im.SameLine()
				}
				im.PushID(imgui.ID(i))
				im.PushStyleColor(imgui.ColButton, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.6, 0.6}))
				im.PushStyleColor(imgui.ColButtonHovered, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.7, 0.7}))
				im.PushStyleColor(imgui.ColButtonActive, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.8, 0.8}))
				im.Button("Click")
				im.PopStyleColorN(3)
				im.PopID()
			}
			// Arrow buttons
			spacing := im.GetStyle().ItemInnerSpacing.X
			if im.ArrowButton("##left", imgui.DirLeft) {
			}
			im.SameLineEx(0.0, spacing)
			if im.ArrowButton("##left", imgui.DirRight) {
			}

			im.Text("Hover over me")
			if im.IsItemHovered() {
				im.SetTooltip("I am a tooltip")
			}

			im.SameLine()
			im.Text("- or me")
			if im.IsItemHovered() {
				im.BeginTooltip()
				im.Text("I am a fancy tooltip")
				arr := ui.Widgets.Basic.Arr
				im.PlotLines("Curve", arr)
				im.EndTooltip()
			}

			im.Separator()
			im.LabelText("label", "Value")
			{
				// Using the _simplified_ one-liner Combo() api here
				items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO"}
				items_current := &ui.Widgets.Basic.ItemsCurrent
				im.ComboString("combo", items_current, items)
				im.SameLine()
				ShowHelpMarker("Refer to the \"Combo\" section below for an explanation of the full BeginCombo/EndCombo API, and demonstration of various flags.\n")
			}

			{
				str0 := ui.Widgets.Basic.Input.Str0
				i0 := &ui.Widgets.Basic.Input.I0
				im.InputText("input text", str0)
				im.SameLine()
				ShowHelpMarker("Hold SHIFT or use mouse to select text.\nCTRL+Left/Right to word jump.\nCTRL+A or double-click to select all.\nCTRL+X,CTRL+C,CTRL+V clipboard.\nCTRL+Z,CTRL+Y undo/redo.\nESCAPE to revert.\n")

				im.InputInt("input int", i0)
				im.SameLine()
				ShowHelpMarker("You can apply arithmetic operators +,*,/ on numerical values.\n  e.g. [ 100 ], input '*2', result becomes [ 200 ]\nUse +- to subtract.\n")

				f0 := &ui.Widgets.Basic.Input.F0
				im.InputFloatEx("input float", f0, 0.01, 1.0, "", 1)

				d0 := &ui.Widgets.Basic.Input.D0
				im.InputFloatEx("input double", d0, 0.01, 1.0, "%.6f", 1)

				f1 := &ui.Widgets.Basic.Input.F1
				im.InputFloatEx("input scientific", f1, 0.0, 0.0, "%e", 1)
				im.SameLine()
				ShowHelpMarker("You can input value using the scientific notation,\n  e.g. \"1e+8\" becomes \"100000000\".\n")

				vec4a := ui.Widgets.Basic.Input.Vec4a[:]
				im.InputFloatN("input float3", vec4a)
			}

			{
				i1 := &ui.Widgets.Basic.Drag.I1
				im.DragInt("drag int", i1)
				im.SameLine()
				ShowHelpMarker("Click and drag to edit value.\nHold SHIFT/ALT for faster/slower edit.\nDouble-click or CTRL+click to input value.")

				i2 := &ui.Widgets.Basic.Drag.I2
				im.DragIntEx("drag int 0..100", i2, 1, 0, 100, "%.0f%%")

				f1 := &ui.Widgets.Basic.Drag.F1
				f2 := &ui.Widgets.Basic.Drag.F2
				im.DragFloatEx("drag float", f1, 0.005, 0, 0, "", 1)
				im.DragFloatEx("drag small float", f2, 0.0001, 0.0, 0.0, "%.06f ns", 1)
			}

			{
				i1 := &ui.Widgets.Basic.Slider.I1
				im.SliderInt("slider int", i1, -1, 3)
				im.SameLine()
				ShowHelpMarker("CTRL+click to input value.")

				f1 := &ui.Widgets.Basic.Slider.F1
				f2 := &ui.Widgets.Basic.Slider.F2
				angle := &ui.Widgets.Basic.Slider.Angle
				im.SliderFloatEx("slider float", f1, 0.0, 1.0, "ratio = %.3f", 1)
				im.SliderFloatEx("slider log float", f2, -10.0, 10.0, "%.4f", 3)
				im.SliderAngle("slider angle", angle)
			}

			{
				col1 := &ui.Widgets.Basic.ColorEdit.Col1
				col2 := &ui.Widgets.Basic.ColorEdit.Col2
				im.ColorEdit3("color1", col1)
				im.SameLine()
				ShowHelpMarker("Click on the colored square to open a color picker.\nRight-click on the colored square to show options.\nCTRL+click on individual component to input value.\n")

				im.ColorEdit4("color 2", col2)
			}

			{
				// List box
				listbox_items := []string{"Apple", "Banana", "Cherry", "Kiwi", "Mango", "Orange", "Pineapple", "Strawberry", "Watermelon"}
				listbox_item_current := &ui.Widgets.Basic.ListBox.ItemCurrent
				im.ListBoxEx("listbox\n(single select)", listbox_item_current, listbox_items, 4)
			}

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
				ShowHelpMarker("This is a more standard looking tree with selectable nodes.\nClick to select, CTRL+Click to toggle, click on arrows or double-click to open.")
				align_label_with_current_x_position := &ui.Widgets.Trees.AlignLabelWithCurrentXPosition
				im.Checkbox("Align label with current X position)", align_label_with_current_x_position)
				im.Text("Hello!")
				if *align_label_with_current_x_position {
					im.UnindentEx(im.GetTreeNodeToLabelSpacing())
				}
				im.TreePop()
			}

			im.TreePop()
		}

		if im.TreeNode("Collapsing Headers") {
			closable_group := &ui.Widgets.CollapsingHeader.ClosableGroup
			im.Checkbox("Enable extra group", closable_group)
			if im.CollapsingHeader("Header") {
				im.Text("IsItemHovered: %v", im.IsItemHovered())
				for i := 0; i < 5; i++ {
					im.Text("Some content %d", i)
				}
			}
			if im.CollapsingHeaderOpen("Header with a close button", closable_group) {
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
				ShowHelpMarker("The TextDisabled color is stored in ImGuiStyle.")
				im.TreePop()
			}

			if im.TreeNode("Word Wrapping") {
				// Using shortcut. You can use PushTextWrapPos()/PopTextWrapPos() for more flexibility.
				im.TextWrapped("This text should automatically wrap on the edge of the window. The current implementation for text wrapping follows simple rules suitable for English and possibly other languages.")
				im.Spacing()

				wrap_width := &ui.Widgets.Text.WordWrapping.WrapWidth
				im.SliderFloatEx("Wrap width", wrap_width, -20, 600, "%.0f", 1)

				im.Text("Test paragraph 1:")
				pos := im.GetCursorScreenPos()
				im.GetWindowDrawList().AddRectFilled(
					f64.Vec2{pos.X + *wrap_width, pos.Y},
					f64.Vec2{pos.X + *wrap_width + 10, pos.Y + im.GetTextLineHeight()},
					color.RGBA{255, 0, 255, 255},
				)
				im.PushTextWrapPos(im.GetCursorPos().X + *wrap_width)
				im.Text("The lazy dog is a good dog. This paragraph is made to fit within %.0f pixels. Testing a 1 character word. The quick brown fox jumps over the lazy dog.", *wrap_width)
				im.GetWindowDrawList().AddRect(im.GetItemRectMin(), im.GetItemRectMax(), color.RGBA{255, 255, 0, 255})
				im.PopTextWrapPos()

				im.Text("Test paragraph 2:")
				pos = im.GetCursorScreenPos()
				im.GetWindowDrawList().AddRectFilled(
					f64.Vec2{pos.X + *wrap_width, pos.Y},
					f64.Vec2{pos.X + *wrap_width + 10, pos.Y + im.GetTextLineHeight()},
					color.RGBA{255, 0, 255, 255},
				)
				im.PushTextWrapPos(im.GetCursorPos().X + *wrap_width)
				im.Text("aaaaaaaa bbbbbbbb, c cccccccc,dddddddd. d eeeeeeee   ffffffff. gggggggg!hhhhhhhh")
				im.GetWindowDrawList().AddRect(im.GetItemRectMin(), im.GetItemRectMax(), color.RGBA{255, 255, 0, 255})
				im.PopTextWrapPos()

				im.TreePop()
			}

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
			if len(ui.Widgets.Text.UTF8.Buf) == 0 {
				// "nihongo"
				ui.Widgets.Text.UTF8.Buf = []byte("\xe6\x97\xa5\xe6\x9c\xac\xe8\xaa\x9e")
			}
			im.InputText("UTF-8 input", ui.Widgets.Text.UTF8.Buf)
			im.TreePop()
		}

		if im.TreeNode("Images") {
			io := im.GetIO()
			im.TextWrapped("Below we are displaying the font texture (which is the only texture we have access to in this demo). Use the 'ImTextureID' type as storage to pass pointers or identifier to your own texture data. Hover the texture for a zoomed view!")
			// Here we are grabbing the font texture because that's the only one we have access to inside the demo code.
			// Remember that ImTextureID is just storage for whatever you want it to be, it is essentially a value that will be passed to the render function inside the ImDrawCmd structure.
			// If you use one of the default imgui_impl_XXXX.cpp renderer, they all have comments at the top of their file to specify what they expect to be stored in ImTextureID.
			// (for example, the imgui_impl_dx11.cpp renderer expect a 'ID3D11ShaderResourceView*' pointer. The imgui_impl_glfw_gl3.cpp renderer expect a GLuint OpenGL texture identifier etc.)
			// If you decided that ImTextureID = MyEngineTexture*, then you can pass your MyEngineTexture* pointers to ImGui::Image(), and gather width/height through your own functions, etc.
			// Using ShowMetricsWindow() as a "debugger" to inspect the draw data that are being passed to your render will help you debug issues if you are confused about this.

			// Consider using the lower-level ImDrawList::AddImage() API, via ImGui::GetWindowDrawList()->AddImage().
			my_tex_id := io.Fonts.TexID
			my_tex_w := float64(io.Fonts.TexWidth)
			my_tex_h := float64(io.Fonts.TexHeight)

			im.Text("%.0fx%.0f", my_tex_w, my_tex_h)
			pos := im.GetCursorScreenPos()
			im.Image(my_tex_id, f64.Vec2{my_tex_w, my_tex_h}, f64.Vec2{0, 0}, f64.Vec2{1, 1}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 128})
			if im.IsItemHovered() {
				im.BeginTooltip()
				region_sz := 32.0
				region_x := io.MousePos.X - pos.X - region_sz*0.5
				if region_x < 0.0 {
					region_x = 0.0
				} else if region_x > my_tex_w-region_sz {
					region_x = my_tex_w - region_sz
				}
				region_y := io.MousePos.Y - pos.Y - region_sz*0.5
				if region_y < 0.0 {
					region_y = 0.0
				} else if region_y > my_tex_h-region_sz {
					region_y = my_tex_h - region_sz
				}
				zoom := 4.0
				im.Text("Min: (%.2f, %.2f)", region_x, region_y)
				im.Text("Max: (%.2f, %.2f)", region_x+region_sz, region_y+region_sz)
				uv0 := f64.Vec2{(region_x) / my_tex_w, (region_y) / my_tex_h}
				uv1 := f64.Vec2{(region_x + region_sz) / my_tex_w, (region_y + region_sz) / my_tex_h}
				im.Image(my_tex_id, f64.Vec2{region_sz * zoom, region_sz * zoom}, uv0, uv1, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 128})
				im.EndTooltip()
			}
			im.TextWrapped("And now some textured buttons..")
			pressed_count := &ui.Widgets.Images.PressedCount
			for i := 0; i < 8; i++ {
				im.PushID(imgui.ID(i))
				frame_padding := -1 + i // -1 uses default padding
				if im.ImageButtonEx(my_tex_id, f64.Vec2{32, 32}, f64.Vec2{0, 0}, f64.Vec2{32.0 / my_tex_w, 32 / my_tex_h}, frame_padding, color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}) {
					*pressed_count += 1
				}
				im.PopID()
				im.SameLine()
			}
			im.NewLine()
			im.Text("Pressed %d times.", *pressed_count)
			im.TreePop()

		}

		if im.TreeNode("Combo") {
			// Expose flags as checkbox for the demo
			flags := &ui.Widgets.Combo.Flags
			im.CheckboxFlags("ImGuiComboFlags_PopupAlignLeft", flags, uint(imgui.ComboFlagsPopupAlignLeft))
			if im.CheckboxFlags("ImGuiComboFlags_NoArrowButton", flags, uint(imgui.ComboFlagsNoArrowButton)) {
				// Clear the other flag, as we cannot combine both
				*flags &^= uint(imgui.ComboFlagsNoPreview)
			}
			if im.CheckboxFlags("ImGuiComboFlags_NoPreview", flags, uint(imgui.ComboFlagsNoPreview)) {
				// Clear the other flag, as we cannot combine both
				*flags &^= uint(imgui.ComboFlagsNoArrowButton)
			}

			// General BeginCombo() API, you have full control over your selection data and display type.
			// (your selection data could be an index, a pointer to the object, an id for the object, a flag stored in the object itself, etc.)
			items := []string{"AAAA", "BBBB", "CCCC", "DDDD", "EEEE", "FFFF", "GGGG", "HHHH", "IIII", "JJJJ", "KKKK", "LLLLLLL", "MMMM", "OOOOOOO"}
			items_current := &ui.Widgets.Combo.ItemCurrent
			if im.BeginComboEx("combo 1", items[*items_current], imgui.ComboFlags(*flags)) {
				for n := range items {
					is_selected := items[*items_current] == items[n]
					if im.SelectableEx(items[n], is_selected, 0, f64.Vec2{}) {
						*items_current = n
					}
					if is_selected {
						// Set the initial focus when opening the combo (scrolling + for keyboard navigation support in the upcoming navigation branch)
						im.SetItemDefaultFocus()
					}
				}
				im.EndCombo()
			}

			// Simplified one-liner Combo() API, using values packed in a single constant string
			item_current_2 := &ui.Widgets.Combo.ItemCurrent2
			im.ComboString("combo 2 (one liner)", item_current_2, []string{"aaaa", "bbbb", "cccc", "dddd", "eeee"})

			// Simplified one-liner Combo() using an array of const char*
			// If the selection isn't within 0..count, Combo won't display a preview
			item_current_3 := &ui.Widgets.Combo.ItemCurrent3
			im.ComboString("combo 3 (array)", item_current_3, items)

			// Simplified one-liner Combo() using an accessor function
			item_current_4 := &ui.Widgets.Combo.ItemCurrent4
			item_getter := func(idx int) (string, bool) {
				return items[idx], true
			}
			im.ComboItem("combo 4 (function", item_current_4, item_getter, len(items))

			im.TreePop()
		}

		if im.TreeNode("Selectables") {
			// Selectable() has 2 overloads:
			// - The one taking "bool selected" as a read-only selection information. When Selectable() has been clicked is returns true and you can alter selection state accordingly.
			// - The one taking "bool* p_selected" as a read-write selection information (convenient in some cases)
			// The earlier is more flexible, as in real application your selection may be stored in a different manner (in flags within objects, as an external list, etc).
			if im.TreeNode("Basic") {
				selection := ui.Widgets.Selectables.Basic.Selection[:]
				im.SelectableOpen("1. I am selectable", &selection[0])
				im.SelectableOpen("2. I am selectable", &selection[1])
				im.Text("3. I am not selectable")
				im.SelectableOpen("4. I am selectable", &selection[3])
				if im.SelectableEx("5. I am double clickable", selection[4], imgui.SelectableFlagsAllowDoubleClick, f64.Vec2{}) {
					if im.IsMouseDoubleClicked(0) {
						selection[4] = !selection[4]
					}
				}
				im.TreePop()
			}

			if im.TreeNode("Selection State: Single Selection") {
				selected := &ui.Widgets.Selectables.Single.Selected
				for n := 0; n < 5; n++ {
					buf := fmt.Sprintf("Object %d", n)
					if im.SelectableEx(buf, *selected == n, 0, f64.Vec2{}) {
						*selected = n
					}
				}
				im.TreePop()
			}

			if im.TreeNode("Selection State: Multiple Selection") {
				ShowHelpMarker("Hold CTRL and click to select multiple items.")
				selection := ui.Widgets.Selectables.Multiple.Selection[:]
				for n := 0; n < 5; n++ {
					buf := fmt.Sprintf("Object %d", n)
					if im.SelectableEx(buf, selection[n], 0, f64.Vec2{}) {
						// Clear selection when CTRL is not held
						if !im.GetIO().KeyCtrl {
							for i := range selection {
								selection[i] = false
							}
						}
						selection[n] = !selection[n]
					}
				}
				im.TreePop()
			}

			if im.TreeNode("Rendering more text into the same line") {
				// Using the Selectable() override that takes "bool* p_selected" parameter and toggle your booleans automatically.
				selected := ui.Widgets.Selectables.Rendering.Selected[:]
				im.SelectableOpen("main.c", &selected[0])
				im.SameLineEx(300, -1)
				im.Text(" 2,345 bytes")
				im.SelectableOpen("Hello.cpp", &selected[1])
				im.SameLineEx(300, -1)
				im.Text("12,345 bytes")
				im.SelectableOpen("Hello.h", &selected[2])
				im.SameLineEx(300, -1)
				im.Text(" 2,345 bytes")
				im.TreePop()
			}

			if im.TreeNode("In columns") {
				im.Columns(3, "", false)
				selected := ui.Widgets.Selectables.Columns.Selected[:]
				for i := range selected {
					label := fmt.Sprintf("Item %d", i)
					if im.SelectableOpen(label, &selected[i]) {
					}
					im.NextColumn()
				}
				im.Columns(1, "", true)
				im.TreePop()
			}

			if im.TreeNode("Grid") {
				selected := ui.Widgets.Selectables.Grid.Selected[:]
				for i := range selected {
					im.PushID(imgui.ID(i))
					if im.SelectableOpenEx("Sailor", &selected[i], 0, f64.Vec2{50, 50}) {
						x := i % 4
						y := i / 4
						if x > 0 {
							selected[i-1] = !selected[i-1]
						}
						if x < 3 {
							selected[i+1] = !selected[i+1]
						}
						if y > 0 {
							selected[i-4] = !selected[i-4]
						}
						if y < 3 {
							selected[i+4] = !selected[i+4]
						}
					}
					if i%4 < 3 {
						im.SameLine()
					}
					im.PopID()
				}
				im.TreePop()
			}

			im.TreePop()
		}

		if im.TreeNode("Filtered Text Input") {
			buf1 := ui.Widgets.FilteredTextInput.Buf1
			im.InputText("default", buf1)
			buf2 := ui.Widgets.FilteredTextInput.Buf2
			im.InputTextEx("decimal", buf2, f64.Vec2{}, imgui.InputTextFlagsCharsDecimal, nil)
			buf3 := ui.Widgets.FilteredTextInput.Buf3
			im.InputTextEx("hexadecimal", buf3, f64.Vec2{}, imgui.InputTextFlagsCharsHexadecimal|imgui.InputTextFlagsCharsUppercase, nil)
			buf4 := ui.Widgets.FilteredTextInput.Buf4
			im.InputTextEx("uppercase", buf4, f64.Vec2{}, imgui.InputTextFlagsCharsUppercase, nil)
			buf5 := ui.Widgets.FilteredTextInput.Buf5
			im.InputTextEx("no blank", buf5, f64.Vec2{}, imgui.InputTextFlagsCharsNoBlank, nil)
			im.TreePop()

			FilterLetters := func(data *imgui.TextEditCallbackData) int {
				return 0
			}
			buf6 := ui.Widgets.FilteredTextInput.Buf6
			im.InputTextEx("\"imgui\" letters", buf6, f64.Vec2{}, imgui.InputTextFlagsCallbackCharFilter, FilterLetters)

			bufpass := ui.Widgets.FilteredTextInput.Bufpass
			im.Text("Password input")
			im.InputTextEx("password", bufpass, f64.Vec2{}, imgui.InputTextFlagsPassword|imgui.InputTextFlagsCharsNoBlank, nil)
			im.SameLine()
			ShowHelpMarker("Display all characters as '*'.\nDisable clipboard cut and copy.\nDisable logging.\n")
			im.InputTextEx("password (clear)", bufpass, f64.Vec2{}, imgui.InputTextFlagsCharsNoBlank, nil)

			im.TreePop()
		}

		if im.TreeNode("Multi-line Text Input") {
			read_only := &ui.Widgets.MultilineTextInput.ReadOnly
			text := ui.Widgets.MultilineTextInput.Text

			im.PushStyleVar(imgui.StyleVarFramePadding, f64.Vec2{0, 0})
			im.Checkbox("Read-only", read_only)
			im.PopStyleVar()
			flags := imgui.InputTextFlagsAllowTabInput
			if *read_only {
				flags |= imgui.InputTextFlagsReadOnly
			}
			im.InputTextMultiline("##source", text, f64.Vec2{-1.0, im.GetTextLineHeight() * 16}, flags, nil)
			im.TreePop()
		}

		if im.TreeNode("Plots widgets") {
			animate := &ui.Widgets.Plots.Animate
			im.Checkbox("Animate", animate)

			arr := []float64{0.6, 0.1, 1.0, 0.5, 0.92, 0.1, 0.2}
			im.PlotLines("Frame Times", arr)

			// Create a dummy array of contiguous float values to plot
			// Tip: If your float aren't contiguous but part of a structure, you can pass a pointer to your first float and the sizeof() of your structure in the Stride parameter.
			values := ui.Widgets.Plots.Values[:]
			values_offset := &ui.Widgets.Plots.ValuesOffset
			refresh_time := &ui.Widgets.Plots.RefreshTime
			if !*animate || *refresh_time == 0.0 {
				*refresh_time = im.GetTime()
			}
			// Create dummy data at fixed 60 hz rate for the demo
			for *refresh_time < im.GetTime() {
				phase := &ui.Widgets.Plots.Phase
				values[*values_offset] = math.Cos(*phase)
				*values_offset = (*values_offset + 1) % len(values)
				*phase += 0.10 * float64(*values_offset)
				*refresh_time += 1.0 / 60.0
			}
			im.PlotLinesEx("Lines", values, *values_offset, "avg 0.0", -1.0, 1.0, f64.Vec2{0, 80})
			im.PlotHistogramEx("Histogram", arr, 0, "", 0.0, 1.0, f64.Vec2{0, 80})

			// Use functions to generate output
			// FIXME: This is rather awkward because current plot API only pass in indices. We probably want an API passing floats and user provide sample rate/count.
			func_type := &ui.Widgets.Plots.FuncType
			display_count := &ui.Widgets.Plots.DisplayCount
			im.Separator()
			im.PushItemWidth(100)
			im.ComboString("func", func_type, []string{"Sin", "Saw"})
			im.PopItemWidth()
			im.SameLine()
			im.SliderInt("Sample count", display_count, 1, 400)
			var fun func(idx int) float64
			if *func_type == 0 {
				fun = func(i int) float64 { return math.Sin(float64(i) * 0.1) }
			} else {
				fun = func(i int) float64 {
					if i&1 != 0 {
						return 1
					}
					return -1
				}
			}
			im.PlotLinesItemEx("Lines", fun, *display_count, 0, "", -1, 1, f64.Vec2{0, 80})
			im.PlotHistogramItemEx("Histogram", fun, *display_count, 0, "", -1.0, 1.0, f64.Vec2{0, 80})
			im.Separator()

			// Animate a simple progress bar
			progress := &ui.Widgets.Plots.Progress
			progress_dir := &ui.Widgets.Plots.ProgressDir
			if *animate {
				*progress += *progress_dir * 0.4 * im.GetIO().DeltaTime
				if *progress >= +1.1 {
					*progress = +1.1
					*progress_dir *= -1.0
				}
				if *progress <= -0.1 {
					*progress = -0.1
					*progress_dir *= -1.0
				}
			}
			// Typically we would use ImVec2(-1.0f,0.0f) to use all available width, or ImVec2(width,0.0f) for a specified width. ImVec2(0.0f,0.0f) uses ItemWidth.
			im.ProgressBarEx(*progress, f64.Vec2{0.0, 0.0}, "")
			im.SameLineEx(0.0, im.GetStyle().ItemInnerSpacing.X)
			im.Text("Progress Bar")

			progress_saturated := f64.Saturate(*progress)
			buf := fmt.Sprintf("%d/%d", int(progress_saturated*1753), 1753)
			im.ProgressBarEx(*progress, f64.Vec2{0, 0}, buf)
			im.TreePop()
		}

		if im.TreeNode("Color/Picker Widgets") {
			color_ := &ui.Widgets.ColorPicker.Color

			alpha_preview := &ui.Widgets.ColorPicker.AlphaPreview
			alpha_half_preview := &ui.Widgets.ColorPicker.AlphaHalfPreview
			options_menu := &ui.Widgets.ColorPicker.OptionsMenu
			hdr := &ui.Widgets.ColorPicker.Hdr
			im.Checkbox("With Alpha Preview", alpha_preview)
			im.Checkbox("With Half Alpha Preview", alpha_half_preview)
			im.Checkbox("With Options Menu", options_menu)
			im.SameLine()
			ShowHelpMarker("Right-click on the individual color widget to show options.")
			im.Checkbox("With HDR", hdr)
			im.SameLine()
			ShowHelpMarker("Currently all this does is to lift the 0..1 limits on dragging widgets.")
			var misc_flags imgui.ColorEditFlags
			if *hdr {
				misc_flags |= imgui.ColorEditFlagsHDR
			}
			if *alpha_half_preview {
				misc_flags |= imgui.ColorEditFlagsAlphaPreviewHalf
			}
			if *alpha_preview {
				misc_flags |= imgui.ColorEditFlagsAlphaPreview
			}
			if *options_menu {
				misc_flags |= imgui.ColorEditFlagsNoOptions
			}

			im.Text("Color widget:")
			im.SameLine()
			ShowHelpMarker("Click on the colored square to open a color picker.\nCTRL+click on individual component to input value.\n")
			im.ColorEdit3Ex("MyColor##1", color_, misc_flags)

			im.Text("Color widget HSV with Alpha:")
			im.ColorEdit4Ex("MyColor##2", color_, imgui.ColorEditFlagsHSV|misc_flags)

			im.Text("Color widget with Float Display:")
			im.ColorEdit4Ex("MyColor##2f", color_, imgui.ColorEditFlagsFloat|misc_flags)

			im.Text("Color button with Picker:")
			im.SameLine()
			ShowHelpMarker("With the ImGuiColorEditFlags_NoInputs flag you can hide all the slider/text inputs.\nWith the ImGuiColorEditFlags_NoLabel flag you can pass a non-empty label which will only be used for the tooltip and picker popup.")
			im.ColorEdit4Ex("MyColor##3", color_, imgui.ColorEditFlagsNoInputs|imgui.ColorEditFlagsNoLabel|misc_flags)

			im.Text("Color button with Custom Picker Popup:")

			// Generate a dummy palette
			saved_palette_inited := &ui.Widgets.ColorPicker.SavedPaletteInited
			saved_palette := ui.Widgets.ColorPicker.SavedPalette[:]
			if !*saved_palette_inited {
				for n := range saved_palette {
					saved_palette[n] = chroma.HSV2RGB(chroma.HSV{float64(n) / 31.0 * 360, 0.8, 0.8})
				}
				*saved_palette_inited = true
			}

			backup_color := &ui.Widgets.ColorPicker.BackupColor
			open_popup := im.ColorButtonEx("MyColor##3b", *color_, misc_flags, f64.Vec2{})
			im.SameLine()
			if im.Button("Palette") {
				open_popup = true
			}
			if open_popup {
				im.OpenPopup("mypicker")
				*backup_color = *color_
			}
			if im.BeginPopup("mypicker") {
				// FIXME: Adding a drag and drop example here would be perfect!
				im.Text("MY CUSTOM COLOR PICKER WITH AN AMAZING PALETTE!")
				im.Separator()
				im.ColorPicker4Ex("##picker", color_, misc_flags|imgui.ColorEditFlagsNoSidePreview|imgui.ColorEditFlagsNoSmallPreview, nil)
				im.SameLine()
				im.BeginGroup()
				im.Text("Current")
				im.ColorButtonEx("##current", *color_, imgui.ColorEditFlagsNoPicker|imgui.ColorEditFlagsAlphaPreviewHalf, f64.Vec2{60, 40})
				im.Text("Previous")
				if (im.ColorButtonEx("##previous", *backup_color, imgui.ColorEditFlagsNoPicker|imgui.ColorEditFlagsAlphaPreviewHalf, f64.Vec2{60, 40})) {
					*color_ = *backup_color
				}
				im.Separator()
				im.Text("Palette")
				for n := range saved_palette {
					im.PushID(imgui.ID(n))
					if n%8 != 0 {
						im.SameLineEx(0.0, im.GetStyle().ItemSpacing.Y)
					}
					if (im.ColorButtonEx("##palette", saved_palette[n], imgui.ColorEditFlagsNoAlpha|imgui.ColorEditFlagsNoPicker|imgui.ColorEditFlagsNoTooltip, f64.Vec2{20, 20})) {
						*color_ = saved_palette[n]
					}

					if im.BeginDragDropTarget() {
						// TODO
						im.EndDragDropTarget()
					}
					im.PopID()
				}
				im.EndGroup()
				im.EndPopup()
			}
			im.Text("Color button only:")
			im.ColorButtonEx("MyColor##3c", *color_, misc_flags, f64.Vec2{80, 80})

			im.Text("Color picker:")
			alpha := &ui.Widgets.ColorPicker.Alpha
			alpha_bar := &ui.Widgets.ColorPicker.AlphaBar
			side_preview := &ui.Widgets.ColorPicker.SidePreview
			ref_color := &ui.Widgets.ColorPicker.RefColor
			ref_color_v := &ui.Widgets.ColorPicker.RefColorV
			inputs_mode := &ui.Widgets.ColorPicker.InputsMode
			picker_mode := &ui.Widgets.ColorPicker.PickerMode
			im.Checkbox("With Alpha", alpha)
			im.Checkbox("With Alpha Bar", alpha_bar)
			im.Checkbox("With Side Preview", side_preview)
			if *side_preview {
				im.SameLine()
				im.Checkbox("With Ref Color", ref_color)
				if *ref_color {
					im.SameLine()
					im.ColorEdit4Ex("##RefColor", ref_color_v, imgui.ColorEditFlagsNoInputs|misc_flags)
				}
			}
			im.ComboString("Inputs Mode", inputs_mode, []string{"All Inputs", "No Inputs", "RGB Input", "HSV Input", "HEX Input"})
			im.ComboString("Picker Mode", picker_mode, []string{"Auto/Current", "Hue bar + SV rect", "Hue wheel + SV triangle"})
			im.SameLine()
			ShowHelpMarker("User can right-click the picker to change mode.")
			flags := misc_flags
			if !*alpha {
				// This is by default if you call ColorPicker3() instead of ColorPicker4()
				flags |= imgui.ColorEditFlagsNoAlpha
			}
			if *alpha_bar {
				flags |= imgui.ColorEditFlagsAlphaBar
			}
			if !*side_preview {
				flags |= imgui.ColorEditFlagsNoSidePreview
			}
			if *picker_mode == 1 {
				flags |= imgui.ColorEditFlagsPickerHueBar
			}
			if *picker_mode == 2 {
				flags |= imgui.ColorEditFlagsPickerHueWheel
			}
			if *inputs_mode == 1 {
				flags |= imgui.ColorEditFlagsNoInputs
			}
			if *inputs_mode == 2 {
				flags |= imgui.ColorEditFlagsRGB
			}
			if *inputs_mode == 3 {
				flags |= imgui.ColorEditFlagsHSV
			}
			if *inputs_mode == 4 {
				flags |= imgui.ColorEditFlagsHEX
			}
			if *ref_color {
				im.ColorPicker4Ex("MyColor##4", color_, flags, ref_color_v)
			} else {
				im.ColorPicker4Ex("MyColor##4", color_, flags, nil)
			}

			im.Text("Programmatically set defaults/options:")
			im.SameLine()
			ShowHelpMarker("SetColorEditOptions() is designed to allow you to set boot-time default.\nWe don't have Push/Pop functions because you can force options on a per-widget basis if needed, and the user can change non-forced ones with the options menu.\nWe don't have a getter to avoid encouraging you to persistently save values that aren't forward-compatible.")

			if im.Button("Uint8 + HSV") {
				im.SetColorEditOptions(imgui.ColorEditFlagsUint8 | imgui.ColorEditFlagsHSV)
			}
			im.SameLine()
			if im.Button("Float + HDR") {
				im.SetColorEditOptions(imgui.ColorEditFlagsFloat | imgui.ColorEditFlagsRGB)
			}

			im.TreePop()
		}

		if im.TreeNode("Range Widgets") {
			begin := &ui.Widgets.Range.Begin
			end := &ui.Widgets.Range.End
			begin_i := &ui.Widgets.Range.BeginI
			end_i := &ui.Widgets.Range.EndI
			im.DragFloatRange2Ex("range", begin, end, 0.25, 0.0, 100.0, "Min: %.1f %%", "Max: %.1f %%", 1)
			im.DragIntRange2Ex("range int (no bounds)", begin_i, end_i, 5, 0, 0, "Min: %.0f units", "Max: %.0f units")
			im.TreePop()
		}

		if im.TreeNode("Multi-component Widgets") {
			vec4f := ui.Widgets.MultiComponents.Vec4f[:]
			vec4i := ui.Widgets.MultiComponents.Vec4i[:]

			im.InputFloat2("input float2", vec4f)
			im.DragFloat2Ex("drag float2", vec4f, 0.01, 0.0, 1.0, "", 1)
			im.SliderFloat2("slider float 2", vec4f, 0.0, 1.0)
			im.DragInt2Ex("drag int2", vec4i, 1, 0, 255, "")
			im.InputInt2("input int2", vec4i)
			im.SliderInt2("slider int2", vec4i, 0, 255)
			im.Spacing()

			im.InputFloat3("input float3", vec4f)
			im.DragFloat3Ex("drag float3", vec4f, 0.01, 0.0, 1.0, "", 1.0)
			im.SliderFloat3Ex("slider float3", vec4f, 0.0, 1.0, "", 1.0)
			im.DragInt3Ex("drag int3", vec4i, 1, 0, 255, "")
			im.InputInt3("input int3", vec4i)
			im.SliderInt3("slider int3", vec4i, 0, 255)
			im.Spacing()

			im.InputFloat4("input float4", vec4f)
			im.DragFloat4Ex("drag float4", vec4f, 0.01, 0.0, 1.0, "", 1)
			im.SliderFloat4Ex("slider float4", vec4f, 0.0, 1.0, "", 1)
			im.InputInt4("input int4", vec4i)
			im.DragInt4Ex("drag int4", vec4i, 1, 0, 255, "")
			im.SliderInt4("slider int4", vec4i, 0, 255)

			im.TreePop()
		}

		if im.TreeNode("Vertical Sliders") {
			const spacing = 4
			im.PushStyleVar(imgui.StyleVarItemSpacing, f64.Vec2{spacing, spacing})

			int_value := &ui.Widgets.VerticalSliders.IntValue
			im.VSliderInt("##int", f64.Vec2{18, 160}, int_value, 0, 5)
			im.SameLine()

			values := []float64{0.0, 0.60, 0.35, 0.9, 0.70, 0.20, 0.0}
			im.PushStringID("set1")
			for i := 0; i < 7; i++ {
				if i > 0 {
					im.SameLine()
				}
				im.PushID(imgui.ID(i))
				im.PushStyleColor(imgui.ColFrameBg, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.6, 0.5}))
				im.PushStyleColor(imgui.ColFrameBgActive, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.7, 0.5}))
				im.PushStyleColor(imgui.ColSliderGrab, chroma.HSV2RGB(chroma.HSV{float64(i) / 7.0 * 360, 0.9, 0.5}))
				im.VSliderFloatEx("##v", f64.Vec2{18, 160}, &values[i], 0.0, 1.0, "", 1.0)
				im.PopStyleColorN(4)
				im.PopID()
			}
			im.PopID()

			im.SameLine()
			im.PushStringID("set2")
			values2 := []float64{0.20, 0.80, 0.40, 0.25}
			const rows = 3
			small_slider_size := f64.Vec2{18, (160 - (rows-1)*spacing) / rows}
			for nx := 0; nx < 4; nx++ {
				if nx > 0 {
					im.SameLine()
				}
				im.BeginGroup()
				for ny := 0; ny < rows; ny++ {
					im.PushID(imgui.ID(nx*rows + ny))
					im.VSliderFloatEx("##v", small_slider_size, &values2[nx], 0.0, 1.0, "", 1.0)
					if im.IsItemActive() || im.IsItemHovered() {
						im.SetTooltip("%.3f", values2[nx])
					}
					im.PopID()
				}
				im.EndGroup()
			}
			im.PopID()

			im.SameLine()
			im.PushStringID("set3")
			for i := 0; i < 4; i++ {
				if i > 0 {
					im.SameLine()
				}
				im.PushID(imgui.ID(i))
				im.PushStyleVar(imgui.StyleVarGrabMinSize, 40.0)
				im.VSliderFloatEx("##v", f64.Vec2{40, 160}, &values[i], 0.0, 1.0, "%.2f\nsec", 1.0)
				im.PopStyleVar()
				im.PopID()
			}
			im.PopID()
			im.PopStyleVar()
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Layout") {
		if im.TreeNode("Child regions") {
			disable_mouse_wheel := &ui.Layout.ChildRegion.DisableMouseWheel
			disable_menu := &ui.Layout.ChildRegion.DisableMenu
			line := &ui.Layout.ChildRegion.Line

			goto_line := im.Button("Goto")
			im.Checkbox("Disable Mouse Wheel", disable_mouse_wheel)
			im.Checkbox("Disable Menu", disable_menu)
			im.Button("Goto")
			im.SameLine()
			im.PushItemWidth(100)
			im.InputIntEx("##Line", line, 0, 0, imgui.InputTextFlagsEnterReturnsTrue)
			im.PopItemWidth()

			// Child 1: no border, enable horizontal scrollbar
			{
				window_flags := imgui.WindowFlagsHorizontalScrollbar
				if *disable_mouse_wheel {
					window_flags |= imgui.WindowFlagsNoScrollWithMouse
				}
				im.BeginChildEx("Child1", f64.Vec2{im.GetWindowContentRegionWidth() * 0.5, 300}, false, window_flags)
				for i := 0; i < 100; i++ {
					im.Text("%04d: scrollable region", i)
					if goto_line && *line == i {
						im.SetScrollHere()
					}
				}
				if goto_line && *line >= 100 {
					im.SetScrollHere()
				}
				im.EndChild()
			}
			im.SameLine()

			// Child 2: rounded border
			{
				im.PushStyleVar(imgui.StyleVarChildRounding, 5.0)
				var window_flags imgui.WindowFlags
				if *disable_mouse_wheel {
					window_flags |= imgui.WindowFlagsNoScrollWithMouse
				}
				if *disable_menu {
					window_flags |= imgui.WindowFlagsMenuBar
				}
				im.BeginChildEx("Child2", f64.Vec2{0, 300}, true, window_flags)
				if !*disable_menu && im.BeginMenuBar() {
					if im.BeginMenu("Menu") {
						ShowExampleMenuFile()
						im.EndMenu()
					}
					im.EndMenuBar()
				}
				im.Columns(2, "", true)
				for i := 0; i < 100; i++ {
					if i == 50 {
						im.NextColumn()
					}
					buf := fmt.Sprintf("%08x", i*5731)
					im.ButtonEx(buf, f64.Vec2{-1, 0}, 0)
				}
				im.EndChild()
				im.PopStyleVar()
			}
			im.TreePop()
		}

		if im.TreeNode("Widgets Width") {
			f := &ui.Layout.WidgetsWidth.F
			im.Text("PushItemWidth(100)")
			im.SameLine()
			ShowHelpMarker("Fixed width.")
			im.PushItemWidth(100)
			im.DragFloat("float##1", f)
			im.PopItemWidth()

			im.Text("PushItemWidth(GetWindowWidth() * 0.5f)")
			im.SameLine()
			ShowHelpMarker("Half of window width.")
			im.PushItemWidth(im.GetWindowWidth() * 0.5)
			im.DragFloat("float##2", f)
			im.PopItemWidth()

			im.Text("PushItemWidth(GetContentRegionAvailWidth() * 0.5f)")
			im.SameLine()
			ShowHelpMarker("Half of available width.\n(~ right-cursor_pos)\n(works within a column set)")
			im.PushItemWidth(im.GetContentRegionAvailWidth() * 0.5)
			im.DragFloat("float##3", f)
			im.PopItemWidth()

			im.Text("PushItemWidth(-100)")
			im.SameLine()
			ShowHelpMarker("Align to right edge minus 100")
			im.PushItemWidth(-100)
			im.DragFloat("float##4", f)
			im.PopItemWidth()

			im.Text("PushItemWidth(-1)")
			im.SameLine()
			ShowHelpMarker("Align to right edge")
			im.PushItemWidth(-1)
			im.DragFloat("float##5", f)
			im.PopItemWidth()

			im.TreePop()
		}

		if im.TreeNode("Basic Horizontal Layout") {
			im.TextWrapped("(Use ImGui::SameLine() to keep adding items to the right of the preceding item)")

			// Text
			im.Text("Two items: Hello")
			im.SameLine()
			im.TextColored(color.RGBA{255, 255, 0, 255}, "Sailor")

			// Adjust spacing
			im.Text("More spacing: Hello")
			im.SameLineEx(0, 20)
			im.TextColored(color.RGBA{255, 255, 0, 255}, "Sailor")

			// Button
			im.AlignTextToFramePadding()
			im.Text("Normal buttons")
			im.SameLine()
			im.Button("Banana")
			im.SameLine()
			im.Button("Apple")
			im.SameLine()
			im.Button("Corniflower")

			// Button
			im.Text("Small buttons")
			im.SameLine()
			im.SmallButton("Like this one")
			im.SameLine()
			im.Text("can fit within a text block.")

			// Aligned to arbitrary position. Easy/cheap column.
			im.Text("Aligned")
			im.SameLineEx(150, -1)
			im.Text("x=150")
			im.SameLineEx(300, -1)
			im.Text("x=300")
			im.Text("Aligned")
			im.SameLineEx(150, -1)
			im.SmallButton("x=150")
			im.SameLineEx(300, -1)
			im.SmallButton("x=300")

			// Checkbox
			c1 := &ui.Layout.Horizontal.C1
			c2 := &ui.Layout.Horizontal.C2
			c3 := &ui.Layout.Horizontal.C3
			c4 := &ui.Layout.Horizontal.C4

			im.Checkbox("My", c1)
			im.SameLine()
			im.Checkbox("Tailor", c2)
			im.SameLine()
			im.Checkbox("Is", c3)
			im.SameLine()
			im.Checkbox("Rich", c4)

			// Various
			f0 := &ui.Layout.Horizontal.F0
			f1 := &ui.Layout.Horizontal.F1
			f2 := &ui.Layout.Horizontal.F2
			im.PushItemWidth(80)
			items := []string{"AAAA", "BBBB", "CCCC", "DDDD"}
			item := &ui.Layout.Horizontal.Item
			im.ComboString("Combo", item, items)
			im.SameLine()
			im.SliderFloat("X", f0, 0.0, 5.0)
			im.SameLine()
			im.SliderFloat("Y", f1, 0.0, 5.0)
			im.SameLine()
			im.SliderFloat("Z", f2, 0.0, 5.0)
			im.PopItemWidth()

			im.PushItemWidth(80)
			im.Text("Lists:")
			selection := ui.Layout.Horizontal.Selection[:]
			for i := 0; i < 4; i++ {
				if i > 0 {
					im.SameLine()
				}
				im.PushID(imgui.ID(i))
				im.ListBox("", &selection[i], items)
				im.PopID()
			}
			im.PopItemWidth()

			// Dummy
			sz := f64.Vec2{30, 30}
			im.ButtonEx("A", sz, 0)
			im.SameLine()
			im.Dummy(sz)
			im.SameLine()
			im.ButtonEx("B", sz, 0)

			im.TreePop()
		}

		if im.TreeNode("Groups") {
			im.TextWrapped("(Using ImGui::BeginGroup()/EndGroup() to layout items. BeginGroup() basically locks the horizontal position. EndGroup() bundles the whole group so that you can use functions such as IsItemHovered() on it.)")
			im.BeginGroup()
			{
				im.BeginGroup()
				im.Button("AAA")
				im.SameLine()
				im.Button("BBB")
				im.SameLine()
				im.BeginGroup()
				im.Button("CCC")
				im.Button("DDD")
				im.EndGroup()
				im.SameLine()
				im.Button("EEE")
				im.EndGroup()
				if im.IsItemHovered() {
					im.SetTooltip("First group hovered")
				}
			}
			// Capture the group size and create widgets using the same size
			size := im.GetItemRectSize()
			values := []float64{0.5, 0.20, 0.80, 0.60, 0.25}
			im.PlotHistogramEx("##values", values, 0, "", 0.0, 1.0, size)

			im.ButtonEx("ACTION", f64.Vec2{(size.X - im.GetStyle().ItemSpacing.X) * 0.5, size.Y}, 0)
			im.SameLine()
			im.ButtonEx("REACTION", f64.Vec2{(size.X - im.GetStyle().ItemSpacing.X) * 0.5, size.Y}, 0)
			im.EndGroup()
			im.SameLine()

			im.ButtonEx("LEVERAGE\nBUZZWORD", size, 0)
			im.SameLine()

			if im.ListBoxHeader("List", size) {
				im.SelectableEx("Selected", true, 0, f64.Vec2{})
				im.SelectableEx("Not Selected", false, 0, f64.Vec2{})
				im.ListBoxFooter()
			}

			im.TreePop()
		}

		if im.TreeNode("Text Baseline Alignment") {
			im.TextWrapped("(This is testing the vertical alignment that occurs on text to keep it at the same baseline as widgets. Lines only composed of text or \"small\" widgets fit in less vertical spaces than lines with normal widgets)")

			im.Text("One\nTwo\nThree")
			im.SameLine()
			im.Text("Hello\nWorld")
			im.SameLine()
			im.Text("Banana")

			im.Text("Banana")
			im.SameLine()
			im.Text("Hello\nWorld")
			im.SameLine()
			im.Text("One\nTwo\nThree")

			im.Button("HOP##1")
			im.SameLine()
			im.Text("Banana")
			im.SameLine()
			im.Text("Hello\nWorld")
			im.SameLine()
			im.Text("Banana")

			im.Button("HOP##2")
			im.SameLine()
			im.Text("Hello\nWorld")
			im.SameLine()
			im.Text("Banana")

			im.Button("TEST##1")
			im.SameLine()
			im.Text("TEST")
			im.SameLine()
			im.SmallButton("TEST##2")

			// If your line starts with text, call this to align it to upcoming widgets.
			im.AlignTextToFramePadding()
			im.Text("Text aligned to Widget")
			im.SameLine()
			im.Button("Widget##1")
			im.SameLine()
			im.Text("Widget")
			im.SameLine()
			im.SmallButton("Widget##2")
			im.SameLine()
			im.Button("Widget##3")

			// Tree
			spacing := im.GetStyle().ItemInnerSpacing.X
			im.Button("Button##1")
			im.SameLineEx(0.0, spacing)
			if im.TreeNode("Node##1") {
				// Dummy tree data
				for i := 0; i < 6; i++ {
					im.BulletText("Item %d..", i)
					im.TreePop()
				}
			}

			// Vertically align text node a bit lower so it'll be vertically centered with upcoming widget. Otherwise you can use SmallButton (smaller fit).
			im.AlignTextToFramePadding()

			// Common mistake to avoid: if we want to SameLine after TreeNode we need to do it before we add child content.
			node_open := im.TreeNode("Node##2")
			im.SameLineEx(0.0, spacing)
			im.Button("Button##2")

			if node_open {
				// Dummy tree data
				for i := 0; i < 6; i++ {
					im.BulletText("Item %d..", i)
					im.TreePop()
				}
			}

			// Bullet
			im.Button("Button##3")
			im.SameLineEx(0.0, spacing)
			im.BulletText("Bullet text")

			im.AlignTextToFramePadding()
			im.BulletText("Node")
			im.SameLineEx(0.0, spacing)
			im.Button("Button##4")

			im.TreePop()
		}

		if im.TreeNode("Scrolling") {
			im.TextWrapped("(Use SetScrollHere() or SetScrollFromPosY() to scroll to a given position.)")
			track := &ui.Layout.Scrolling.Track
			track_line := &ui.Layout.Scrolling.TrackLine
			scroll_to_px := &ui.Layout.Scrolling.ScrollToPx
			im.Checkbox("Track", track)
			im.PushItemWidth(100)
			im.SameLineEx(130, -1)
			if im.DragIntEx("##line", track_line, 0.25, 0, 99, "Line = %.0f") {
				*track = true
			}
			scroll_to := im.Button("Scroll To Pos")
			im.SameLineEx(130, -1)
			if im.DragIntEx("##pos_y", scroll_to_px, 1.00, 0, 9999, "Y = %.0f px") {
				scroll_to = true
			}
			im.PopItemWidth()
			if scroll_to {
				*track = false
			}

			for i := 0; i < 5; i++ {
				if i > 0 {
					im.SameLine()
				}
				im.BeginGroup()
				switch i {
				case 0:
					im.Text("Top")
				case 1:
					im.Text("25%")
				case 2:
					im.Text("Center")
				case 3:
					im.Text("75%")
				default:
					im.Text("Bottom")
				}
				im.BeginChildIDEx(im.GetID(imgui.ID(i)), f64.Vec2{im.GetWindowWidth() * 0.17, 200.0}, true, 0)
				if scroll_to {
					im.SetScrollFromPosY(im.GetCursorStartPos().Y+float64(*scroll_to_px), float64(i)*0.25)
				}
				for line := 0; line < 100; line++ {
					if *track && line == *track_line {
						im.TextColored(color.RGBA{255, 255, 0, 255}, "Line %d", line)
						// 0.0f:top, 0.5f:center, 1.0f:bottom
						im.SetScrollHereEx(float64(i) * 0.25)
					} else {
						im.Text("Line %d", line)
					}
				}

				scroll_y := im.GetScrollY()
				scroll_max_y := im.GetScrollMaxY()
				im.EndChild()
				im.Text("%.0f/%0.f", scroll_y, scroll_max_y)
				im.EndGroup()
			}

			im.TreePop()
		}

		if im.TreeNode("Horizontal Scrolling") {
			im.Bullet()
			im.TextWrapped("Horizontal scrolling for a window has to be enabled explicitly via the ImGuiWindowFlags_HorizontalScrollbar flag.")
			im.Bullet()
			im.TextWrapped("You may want to explicitly specify content width by calling SetNextWindowContentWidth() before Begin().")

			lines := &ui.Layout.HorizontalScrolling.Lines
			im.SliderInt("Lines", lines, 1, 15)
			im.PushStyleVar(imgui.StyleVarFrameRounding, 3.0)
			im.PushStyleVar(imgui.StyleVarFramePadding, f64.Vec2{2.0, 1.0})
			im.BeginChildEx("scrolling", f64.Vec2{0, im.GetFrameHeightWithSpacing()*7 + 30}, true, imgui.WindowFlagsHorizontalScrollbar)
			for line := 0; line < *lines; line++ {
				// Display random stuff (for the sake of this trivial demo we are using basic Button+SameLine. If you want to create your own time line for a real application you may be better off
				// manipulating the cursor position yourself, aka using SetCursorPos/SetCursorScreenPos to position the widgets yourself. You may also want to use the lower-level ImDrawList API)
				num_buttons := 10
				if *lines&1 != 0 {
					num_buttons += line * 9
				} else {
					num_buttons += line * 3
				}
				for n := 0; n < num_buttons; n++ {
					if n > 0 {
						im.SameLine()
					}
					im.PushID(imgui.ID(n + line*1000))
					num_buf := fmt.Sprint(n)
					var label string
					if n%15 == 0 {
						label = "FizzBuzz"
					} else if n%3 == 0 {
						label = "Fizz"
					} else if n%5 == 0 {
						label = "Buzz"
					} else {
						label = num_buf
					}
					hue := float64(n) * 0.05 * 360
					im.PushStyleColor(imgui.ColButton, chroma.HSV2RGB(chroma.HSV{hue, 0.6, 0.6}))
					im.PushStyleColor(imgui.ColButtonHovered, chroma.HSV2RGB(chroma.HSV{hue, 0.7, 0.7}))
					im.PushStyleColor(imgui.ColButtonActive, chroma.HSV2RGB(chroma.HSV{hue, 0.8, 0.8}))
					im.ButtonEx(label, f64.Vec2{40.0 + math.Sin(float64(line+n))*20.0, 0.0}, 0)
					im.PopStyleColorN(3)
					im.PopID()
				}
				scroll_x := im.GetScrollX()
				scroll_max_x := im.GetScrollMaxX()
				im.EndChild()
				im.PopStyleVarN(2)
				scroll_x_delta := 0.0
				im.SmallButton("<<")
				if im.IsItemActive() {
					scroll_x_delta = -im.GetIO().DeltaTime * 1000.0
				}
				im.SameLine()
				im.Text("Scroll from code")
				im.SameLine()
				im.SmallButton(">>")
				if im.IsItemActive() {
					scroll_x_delta = +im.GetIO().DeltaTime * 1000.0
				}
				im.SameLine()
				im.Text("%.0f/%.0f", scroll_x, scroll_max_x)
				if scroll_x_delta != 0.0 {
					// Demonstrate a trick: you can use Begin to set yourself in the context of another window (here we are already out of your child window)
					im.BeginChild("scrolling")
					im.SetScrollX(im.GetScrollX() + scroll_x_delta)
					im.End()
				}
			}
			im.TreePop()
		}

		if im.TreeNode("Clipping") {
			size := &ui.Layout.Clipping.Size
			offset := &ui.Layout.Clipping.Offset
			im.TextWrapped("On a per-widget basis we are occasionally clipping text CPU-side if it won't fit in its frame. Otherwise we are doing coarser clipping + passing a scissor rectangle to the renderer. The system is designed to try minimizing both execution and CPU/GPU rendering cost.")
			im.DragV2Ex("size", size, 0.5, 0.0, 200.0, "%.0f", 1)
			im.TextWrapped("(Click and drag)")
			pos := im.GetCursorScreenPos()
			clip_rect := f64.Vec4{pos.X, pos.Y, pos.X + size.X, pos.Y + size.Y}
			im.InvisibleButton("##dummy", *size)
			if im.IsItemActive() && im.IsMouseDragging() {
				offset.X += im.GetIO().MouseDelta.X
				offset.Y += im.GetIO().MouseDelta.Y
			}
			im.GetWindowDrawList().AddRectFilled(pos, f64.Vec2{pos.X + size.X, pos.Y + size.Y}, color.RGBA{90, 90, 120, 255})
			im.GetWindowDrawList().AddTextEx(
				im.GetFont(), im.GetFontSize()*2.0,
				f64.Vec2{pos.X + offset.X, pos.Y + offset.Y}, color.RGBA{255, 255, 255, 255},
				"Line 1 hello\nLine 2 clip me!", 0.0, &clip_rect,
			)
			im.TreePop()
		}
	}

	if im.CollapsingHeader("Popups & Modal windows") {
		if im.TreeNode("Popups") {
			im.TextWrapped("When a popup is active, it inhibits interacting with windows that are behind the popup. Clicking outside the popup closes it.")

			selected_fish := &ui.PopupsModal.Popups.SelectedFish
			names := []string{"Bream", "Haddock", "Mackerel", "Pollock", "Tilefish"}
			toggles := []bool{true, false, false, false, false}

			// Simple selection popup
			// (If you want to show the current selection inside the Button itself, you may want to build a string using the "###" operator to preserve a constant ID with a variable label)
			if im.Button("Select..") {
				im.OpenPopup("select")
			}
			im.SameLine()
			if *selected_fish == -1 {
				im.TextUnformatted("<None>")
			} else {
				im.TextUnformatted(names[*selected_fish])
			}
			if im.BeginPopup("select") {
				im.Text("Aquarium")
				im.Separator()
				for i := range names {
					if im.Selectable(names[i]) {
						*selected_fish = i
					}
				}
				im.EndPopup()
			}
			// Showing a menu with toggles
			if im.Button("Toggle..") {
				im.OpenPopup("toggle")
			}
			if im.BeginPopup("toggle") {
				for i := range names {
					im.MenuItemSelect(names[i], "", &toggles[i])
				}
				if im.BeginMenu("Sub-Menu") {
					im.MenuItem("Click me")
					im.EndMenu()
				}

				im.Separator()
				im.Text("Tooltip here")
				if im.IsItemHovered() {
					im.SetTooltip("I am a tooltip over a popup")
				}

				if im.Button("Stacked Popup") {
					im.OpenPopup("another popup")
				}

				if im.BeginPopup("another popup") {
					for i := range names {
						im.MenuItemSelect(names[i], "", &toggles[i])
					}
					if im.BeginMenu("Sub-menu") {
						im.MenuItem("Click me")
						im.EndMenu()
					}
					im.EndPopup()
				}
				im.EndPopup()
			}

			if im.Button("Popup Menu..") {
				im.OpenPopup("FilePopup")
			}
			if im.BeginPopup("FilePopup") {
				ShowExampleMenuFile()
				im.EndPopup()
			}

			im.TreePop()
		}

		if im.TreeNode("Context menus") {
			// BeginPopupContextItem() is a helper to provide common/simple popup behavior of essentially doing:
			//    if (IsItemHovered() && IsMouseClicked(0))
			//       OpenPopup(id);
			//    return BeginPopup(id);
			// For more advanced uses you may want to replicate and cuztomize this code. This the comments inside BeginPopupContextItem() implementation.

			value := &ui.PopupsModal.ContextMenus.Value
			im.Text("Value = %.3f (<-- right-click here)", *value)
			if im.BeginPopupContextItemEx("item context menu", 1) {
				if im.Selectable("Set to zero") {
					*value = 0.0
				}
				if im.Selectable("Set to PI") {
					*value = 3.1415
				}
				im.PushItemWidth(-1)
				im.DragFloatEx("##Value", value, 0.1, 0.0, 0.0, "", 1)
				im.PopItemWidth()
				im.EndPopup()
			}
			name := ui.PopupsModal.ContextMenus.Name
			buf := fmt.Sprintf("Button: %s###Button", name)
			im.Button(buf)
			// When used after an item that has an ID (here the Button), we can skip providing an ID to BeginPopupContextItem().
			if im.BeginPopupContextItem() {
				im.Text("Edit name:")
				im.InputText("##edit", name)
				if im.Button("Close") {
					im.CloseCurrentPopup()
				}
				im.EndPopup()
			}
			im.SameLine()
			im.Text("(<-- right-click here)")
			im.TreePop()
		}

		if im.TreeNode("Modals") {
			im.TextWrapped("Modal windows are like popups but the user cannot close them by clicking outside the window.")
			if im.Button("Delete..") {
				im.OpenPopup("Delete?")
			}

			if im.BeginPopupModalEx("Delete?", nil, imgui.WindowFlagsAlwaysAutoResize) {
				im.Text("All those beautiful files will be deleted.\nThis operation cannot be undone!\n\n")
				im.Separator()

				//static int dummy_i = 0;
				//ImGui::Combo("Combo", &dummy_i, "Delete\0Delete harder\0");

				dont_ask_me_next_time := &ui.PopupsModal.Modals.DontAskMeNextTime
				im.PushStyleVar(imgui.StyleVarFramePadding, f64.Vec2{0, 0})
				im.Checkbox("Don't ask me next time", dont_ask_me_next_time)
				im.PopStyleVar()

				if im.ButtonEx("OK", f64.Vec2{120, 0}, 0) {
					im.CloseCurrentPopup()
				}
				im.SetItemDefaultFocus()
				im.SameLine()
				if im.ButtonEx("Cancel", f64.Vec2{120, 0}, 0) {
					im.CloseCurrentPopup()
				}
				im.EndPopup()
			}

			if im.Button("Stacked modals..") {
				im.OpenPopup("Stacked 1")
			}
			if im.BeginPopupModal("Stacked 1") {
				im.Text("Hello from Stacked The First\nUsing style.Colors[ImGuiCol_ModalWindowDarkening] for darkening.")
				item := &ui.PopupsModal.Modals.Item
				im.ComboString("Combo", item, []string{"aaaa", "bbbb", "cccc", "dddd", "eeee"})
				color := &ui.PopupsModal.Modals.Color
				// This is to test behavior of stacked regular popups over a modal
				im.ColorEdit4("color", color)
				if im.Button("Add another modal..") {
					im.OpenPopup("Stacked 2")
				}
				if im.BeginPopupModal("Stacked 2") {
					im.Text("Hello from Stacked The Second!")
					if im.Button("Close") {
						im.CloseCurrentPopup()
					}
					im.EndPopup()
				}

				if im.Button("Close") {
					im.CloseCurrentPopup()
				}
				im.EndPopup()
			}

			im.TreePop()
		}

		if im.TreeNode("Menus inside a regular window") {
			im.TextWrapped("Below we are testing adding menu items to a regular window. It's rather unusual but should work!")
			im.Separator()
			// NB: As a quirk in this very specific example, we want to differentiate the parent of this menu from the parent of the various popup menus above.
			// To do so we are encloding the items in a PushID()/PopID() block to make them two different menusets. If we don't, opening any popup above and hovering our menu here
			// would open it. This is because once a menu is active, we allow to switch to a sibling menu by just hovering on it, which is the desired behavior for regular menus.
			im.PushStringID("foo")
			im.MenuItemSelect("Menu item", "CTRL+M", nil)
			if im.BeginMenu("Menu inside a regular window") {
				ShowExampleMenuFile()
				im.EndMenu()
			}
			im.PopID()
			im.Separator()
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

		// Create multiple items in a same cell before switching to next column
		if im.TreeNode("Mixed items") {
			im.Columns(3, "mixed", true)
			im.Separator()

			im.Text("Hello")
			im.Button("Banana")
			im.NextColumn()

			im.Text("ImGui")
			im.Button("Apple")
			foo := &ui.Columns.MixedItems.Foo
			im.InputFloatEx("red", foo, 0.05, 0, "%.3f", 1)
			im.Text("An extra line here.")
			im.NextColumn()

			im.Text("Sailor")
			im.Button("Corniflower")
			bar := &ui.Columns.MixedItems.Bar
			im.InputFloatEx("blue", bar, 0.05, 0, "%.3f", 1)
			im.NextColumn()

			if im.CollapsingHeader("Category A") {
				im.Text("Blah blah blah")
			}
			im.NextColumn()
			if im.CollapsingHeader("Category B") {
				im.Text("Blah blah blah")
			}
			im.NextColumn()
			if im.CollapsingHeader("Category C") {
				im.Text("Blah blah blah")
			}
			im.NextColumn()

			im.Columns(1, "", true)
			im.Separator()
			im.TreePop()
		}

		// Word wrapping
		if im.TreeNode("Word-wrapping") {
			im.Columns(2, "word-wrapping", true)
			im.Separator()
			im.TextWrapped("The quick brown fox jumps over the lazy dog.")
			im.TextWrapped("Hello Left")
			im.NextColumn()
			im.TextWrapped("The quick brown fox jumps over the lazy dog.")
			im.TextWrapped("Hello Right")
			im.Columns(1, "", true)
			im.Separator()
			im.TreePop()
		}

		if im.TreeNode("Borders") {
			// NB: Future columns API should allow automatic horizontal borders.
			h_borders := &ui.Columns.Borders.Horizontal
			v_borders := &ui.Columns.Borders.Vertical
			im.Checkbox("horizontal", h_borders)
			im.SameLine()
			im.Checkbox("vertical", v_borders)
			im.Columns(4, "", *v_borders)
			for i := 0; i < 4*3; i++ {
				if *h_borders && im.GetColumnIndex() == 0 {
					im.Separator()
				}
				im.Text("%c%c%c", 'a'+i, 'a'+i, 'a'+i)
				im.Text("Width %.2f\nOffset %.2f", im.GetColumnWidth(), im.GetColumnOffset(-1))
				im.NextColumn()
			}
			im.Columns(1, "", true)
			if *h_borders {
				im.Separator()
			}
			im.TreePop()
		}

		if im.TreeNode("Horizontal Scrolling") {
			im.SetNextWindowContentSize(f64.Vec2{1500.0, 0.0})
			im.BeginChildEx("##ScrollingRegion", f64.Vec2{0, im.GetFontSize() * 20}, false, imgui.WindowFlagsHorizontalScrollbar)
			im.Columns(10, "", true)
			ITEMS_COUNT := 2000
			// Also demonstrate using the clipper for large list
			var clipper imgui.ListClipper
			clipper.Init(im, ITEMS_COUNT, -1)
			for clipper.Step() {
				for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
					for j := 0; j < 10; j++ {
						im.Text("Line %d Column %d...", i, j)
						im.NextColumn()
					}
				}
			}
			im.Columns(1, "", true)
			im.EndChild()
			im.TreePop()
		}

		node_open := im.TreeNode("Tree within single cell")
		im.SameLine()
		ShowHelpMarker("NB: Tree node must be poped before ending the cell. There's no storage of state per-cell.")
		if node_open {
			im.Columns(2, "tree items", true)
			im.Separator()
			if im.TreeNode("Hello") {
				im.BulletText("Sailor")
				im.TreePop()
			}
			im.NextColumn()
			if im.TreeNode("Bonjour") {
				im.BulletText("Marin")
				im.TreePop()
			}
			im.NextColumn()
			im.Columns(1, "", true)
			im.Separator()
			im.TreePop()
		}
		im.PopID()
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
		im.SameLine()
		ShowHelpMarker("Instruct ImGui to render a mouse cursor for you in software. Note that a mouse cursor rendered via your application GPU rendering path will feel more laggy than hardware cursor, but will be more in sync with your other visuals.\n\nSome desktop applications may use both kinds of cursors (e.g. enable software cursor only when resizing/dragging something).")

		im.CheckboxFlags("io.ConfigFlags: NavEnableGamepad", (*uint)(&io.ConfigFlags), uint(imgui.ConfigFlagsNavEnableGamepad))
		im.CheckboxFlags("io.ConfigFlags: NavEnableKeyboard", (*uint)(&io.ConfigFlags), uint(imgui.ConfigFlagsNavEnableKeyboard))
		im.CheckboxFlags("io.ConfigFlags: NavEnableSetMousePos", (*uint)(&io.ConfigFlags), uint(imgui.ConfigFlagsNavEnableSetMousePos))
		im.SameLine()
		ShowHelpMarker("Instruct navigation to move the mouse cursor. See comment for ImGuiConfigFlags_NavEnableSetMousePos.")
		im.CheckboxFlags("io.ConfigFlags: NoMouseCursorChange", (*uint)(&io.ConfigFlags), uint(imgui.ConfigFlagsNoMouseCursorChange))
		im.SameLine()
		ShowHelpMarker("Instruct back-end to not alter mouse cursor shape and visibility.")

		if im.TreeNode("Keyboard, Mouse & Navigation State") {
			if im.IsMousePosValid() {
				im.Text("Mouse pos: (%g, %g)", io.MousePos.X, io.MousePos.Y)
			} else {
				im.Text("Mouse pos: <INVALID>")
			}
			im.Text("Mouse delta: (%g, %g)", io.MouseDelta.X, io.MouseDelta.Y)
			im.Text("Mouse down:")
			for i := range io.MouseDown {
				if io.MouseDownDuration[i] >= 0.0 {
					im.SameLine()
					im.Text("b%d (%.02f secs)", i, io.MouseDownDuration[i])
				}
			}
			im.Text("Mouse clicked:")
			for i := range io.MouseDown {
				if im.IsMouseClicked(i, false) {
					im.SameLine()
					im.Text("b%d", i)
				}
			}
			im.Text("Mouse dbl-clicked:")
			for i := range io.MouseDown {
				if im.IsMouseDoubleClicked(i) {
					im.SameLine()
					im.Text("b%d", i)
				}
			}
			im.Text("Mouse released:")
			for i := range io.MouseDown {
				if im.IsMouseReleased(i) {
					im.SameLine()
					im.Text("b%d", i)
				}
			}
			im.Text("Mouse wheel: %.1f", io.MouseWheel)

			im.Text("Keys down:")
			for i := range io.KeysDown {
				if io.KeysDownDuration[i] >= 0.0 {
					im.SameLine()
					im.Text("%d (%.02f secs)", i, io.KeysDownDuration[i])
				}
			}
			im.Text("Keys pressed:")
			for i := range io.KeysDown {
				if im.IsKeyPressed(i, true) {
					im.SameLine()
					im.Text("%d", i)
				}
			}
			im.Text("Keys release:")
			for i := range io.KeysDown {
				if im.IsKeyReleased(i) {
					im.SameLine()
					im.Text("%d", i)
				}
			}
			var ctrl, shift, alt, super string
			if io.KeyCtrl {
				ctrl = "CTRL"
			}
			if io.KeyShift {
				shift = "SHIFT"
			}
			if io.KeyAlt {
				alt = "ALT"
			}
			if io.KeySuper {
				super = "SUPER"
			}
			im.Text("Keys mods: %s%s%s%s", ctrl, shift, alt, super)

			im.Text("NavInputs down:")
			for i := range io.NavInputs {
				if io.NavInputs[i] > 0.0 {
					im.SameLine()
					im.Text("[%d] %.2f", i, io.NavInputs[i])
				}
			}
			im.Text("NavInputs pressed:")
			for i := range io.NavInputs {
				if io.NavInputsDownDuration[i] == 0.0 {
					im.SameLine()
					im.Text("[%d]", i)
				}
			}
			im.Text("NavInputs duration:")
			for i := range io.NavInputs {
				if io.NavInputsDownDuration[i] >= 0.0 {
					im.SameLine()
					im.Text("[%d] %.2f", i, io.NavInputsDownDuration[i])
				}
			}

			im.Button("Hovering me sets the\nkeyboard capture flag")
			if im.IsItemHovered() {
				im.CaptureKeyboardFromApp(true)
			}
			im.SameLine()
			im.Button("Holding me clears the\nthe keyboard capture flag")
			if im.IsItemActive() {
				im.CaptureKeyboardFromApp(false)
			}

			im.TreePop()
		}

		if im.TreeNode("Tabbing") {
			im.Text("Use TAB/SHIFT+TAB to cycle through keyboard editable fields.")
			buf := ui.Input.Tabbing.Buf
			im.InputText("1", buf)
			im.InputText("2", buf)
			im.InputText("3", buf)
			im.PushAllowKeyboardFocus(false)
			im.InputText("4 (tab skip)", buf)
			im.PopAllowKeyboardFocus()
			im.InputText("5", buf)
			im.TreePop()
		}

		if im.TreeNode("Focus from code") {
			focus_1 := im.Button("Focus on 1")
			im.SameLine()
			focus_2 := im.Button("Focus on 2")
			im.SameLine()
			focus_3 := im.Button("Focus on 3")
			has_focus := 0
			buf := ui.Input.Focus.Buf

			if focus_1 {
				im.SetKeyboardFocusHere()
			}
			im.InputText("1", buf)
			if im.IsItemActive() {
				has_focus = 1
			}

			if focus_2 {
				im.SetKeyboardFocusHere()
			}
			im.InputText("2", buf)
			if im.IsItemActive() {
				has_focus = 2
			}

			im.PushAllowKeyboardFocus(false)
			if focus_3 {
				im.SetKeyboardFocusHere()
			}
			im.InputText("3 (tab skip)", buf)
			if im.IsItemActive() {
				has_focus = 3
			}
			im.PopAllowKeyboardFocus()

			if has_focus != 0 {
				im.Text("Item with focus: %d", has_focus)
			} else {
				im.Text("Item with focus: <none>")
			}

			// Use >= 0 parameter to SetKeyboardFocusHere() to focus an upcoming item
			f3 := ui.Input.Focus.F3[:]
			focus_ahead := -1
			if im.Button("Focus on X") {
				focus_ahead = 0
			}
			im.SameLine()
			if im.Button("Focus on Y") {
				focus_ahead = 1
			}
			im.SameLine()
			if im.Button("Focus on Z") {
				focus_ahead = 2
			}
			if focus_ahead != -1 {
				im.SetKeyboardFocusHereEx(focus_ahead)
			}
			im.SliderFloat3("Float3", f3, 0.0, 1.0)

			im.TextWrapped("NB: Cursor & selection are preserved when refocusing last used item in code.")
			im.TreePop()
		}

		if im.TreeNode("Focused & Hovered Test") {
			embed_all_inside_a_child_window := &ui.Input.FocusHovered.EmbedAllInsideAChildWindow
			im.Checkbox("Embed everything inside a child window (for additional testing)", embed_all_inside_a_child_window)
			if *embed_all_inside_a_child_window {
				im.BeginChildEx("embeddingchild", f64.Vec2{0, im.GetFontSize() * 25}, true, 0)
			}

			// Testing IsWindowFocused() function with its various flags (note that the flags can be combined)
			im.BulletText(
				"IsWindowFocused() = %v\n"+
					"IsWindowFocused(_ChildWindows) = %v\n"+
					"IsWindowFocused(_ChildWindows|_RootWindow) = %v\n"+
					"IsWindowFocused(_RootWindow) = %v\n"+
					"IsWindowFocused(_AnyWindow) = %v\n",
				im.IsWindowFocused(),
				im.IsWindowFocusedEx(imgui.FocusedFlagsChildWindows),
				im.IsWindowFocusedEx(imgui.FocusedFlagsChildWindows|imgui.FocusedFlagsRootWindow),
				im.IsWindowFocusedEx(imgui.FocusedFlagsRootWindow),
				im.IsWindowFocusedEx(imgui.FocusedFlagsAnyWindow))

			// Testing IsWindowHovered() function with its various flags (note that the flags can be combined)
			im.BulletText(
				"IsWindowHovered() = %v\n"+
					"IsWindowHovered(_AllowWhenBlockedByPopup) = %v\n"+
					"IsWindowHovered(_AllowWhenBlockedByActiveItem) = %v\n"+
					"IsWindowHovered(_ChildWindows) = %v\n"+
					"IsWindowHovered(_ChildWindows|_RootWindow) = %v\n"+
					"IsWindowHovered(_RootWindow) = %v\n"+
					"IsWindowHovered(_AnyWindow) = %v\n",
				im.IsWindowHovered(),
				im.IsWindowHoveredEx(imgui.HoveredFlagsAllowWhenBlockedByPopup),
				im.IsWindowHoveredEx(imgui.HoveredFlagsAllowWhenBlockedByActiveItem),
				im.IsWindowHoveredEx(imgui.HoveredFlagsChildWindows),
				im.IsWindowHoveredEx(imgui.HoveredFlagsChildWindows|imgui.HoveredFlagsRootWindow),
				im.IsWindowHoveredEx(imgui.HoveredFlagsRootWindow),
				im.IsWindowHoveredEx(imgui.HoveredFlagsAnyWindow))

			// Testing IsItemHovered() function (because BulletText is an item itself and that would affect the output of IsItemHovered, we pass all lines in a single items to shorten the code)
			im.Button("ITEM")
			im.BulletText(
				"IsItemHovered() = %v\n"+
					"IsItemHovered(_AllowWhenBlockedByPopup) = %v\n"+
					"IsItemHovered(_AllowWhenBlockedByActiveItem) = %v\n"+
					"IsItemHovered(_AllowWhenOverlapped) = %v\n"+
					"IsItemhovered(_RectOnly) = %v\n",
				im.IsItemHovered(),
				im.IsItemHoveredEx(imgui.HoveredFlagsAllowWhenBlockedByPopup),
				im.IsItemHoveredEx(imgui.HoveredFlagsAllowWhenBlockedByActiveItem),
				im.IsItemHoveredEx(imgui.HoveredFlagsAllowWhenOverlapped),
				im.IsItemHoveredEx(imgui.HoveredFlagsRectOnly))

			im.BeginChildEx("child", f64.Vec2{0, 50}, true, 0)
			im.Text("This is another child window for testing IsWindowHovered() flags.")
			im.EndChild()

			if *embed_all_inside_a_child_window {
				im.EndChild()
			}

			im.TreePop()
		}

		if im.TreeNode("Dragging") {
			im.TextWrapped("You can use im.GetMouseDragDelta(0) to query for the dragged amount on any widget.")
			for button := 0; button < 3; button++ {
				im.Text("IsMouseDragging(%d):\n  w/ default threshold: %v,\n  w/ zero threshold: %v\n  w/ large threshold: %v",
					button, im.IsMouseDraggingEx(button, -1.0), im.IsMouseDraggingEx(button, 0.0), im.IsMouseDraggingEx(button, 20.0))
			}
			im.Button("Drag Me")
			if im.IsItemActive() {
				// Draw a line between the button and the mouse cursor
				draw_list := im.GetWindowDrawList()
				draw_list.PushClipRectFullScreen()
				draw_list.AddLineEx(io.MouseClickedPos[0], io.MousePos, im.GetColorFromStyle(imgui.ColButton), 4.0)
				draw_list.PopClipRect()

				// Drag operations gets "unlocked" when the mouse has moved past a certain threshold (the default threshold is stored in io.MouseDragThreshold)
				// You can request a lower or higher threshold using the second parameter of IsMouseDragging() and GetMouseDragDelta()
				value_raw := im.GetMouseDragDelta(0, 0.0)
				value_with_lock_threshold := im.GetMouseDragDelta(0, -1)
				mouse_delta := io.MouseDelta
				im.SameLine()
				im.Text("Raw (%.1f, %.1f), WithLockThresold (%.1f, %.1f), MouseDelta (%.1f, %.1f)",
					value_raw.X, value_raw.Y, value_with_lock_threshold.X, value_with_lock_threshold.Y, mouse_delta.X, mouse_delta.Y)
			}
			im.TreePop()
		}

		if im.TreeNode("Mouse cursors") {
			mouse_cursors_names := []string{"Arrow", "TextInput", "Move", "ResizeNS", "ResizeEW", "ResizeNESW", "ResizeNWSE"}
			assert(len(mouse_cursors_names) == int(imgui.MouseCursorCOUNT))
			im.Text("Current mouse cursor = %d: %s", im.GetMouseCursor(), mouse_cursors_names[im.GetMouseCursor()])
			im.Text("Hover to see mouse cursors:")
			im.SameLine()
			ShowHelpMarker("Your application can render a different mouse cursor based on what im.GetMouseCursor() returns. If software cursor rendering (io.MouseDrawCursor) is set ImGui will draw the right cursor for you, otherwise your backend needs to handle it.")
			for i := 0; i < int(imgui.MouseCursorCOUNT); i++ {
				label := fmt.Sprintf("Mouse cursor %d: %s", i, mouse_cursors_names[i])
				im.Bullet()
				im.SelectableEx(label, false, 0, f64.Vec2{})
				if im.IsItemHovered() || im.IsItemFocused() {
					im.SetMouseCursor(imgui.MouseCursor(i))
				}
			}
			im.TreePop()
		}
	}

	im.End()
}

func ShowUserGuide() {
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

func ShowExampleAppMainMenuBar() {
	if im.BeginMainMenuBar() {
		if im.BeginMenu("File") {
			ShowExampleMenuFile()
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

func ShowExampleAppConsole() {
}

// Demonstrate creating a simple log window with basic filtering.
func ShowExampleAppLog() {
	// Demo: add random items (unless Ctrl is held)
	p_open := &ui.ShowAppLog
	log := &ui.AppLog.Log
	last_time := ui.AppLog.LastTime
	time := im.GetTime()
	if time-last_time >= 0.20 && !im.GetIO().KeyCtrl {
		random_words := []string{"system", "info", "warning", "error", "fatal", "notice", "log"}
		log.AddLog("[%s] Hello, time is %.1f, frame count is %d\n", random_words[rand.Intn(len(random_words))], time, im.GetFrameCount())
		last_time = time
	}
	log.Draw("Example: Log", p_open)
}

// Demonstrate create a window with multiple child windows.
func ShowExampleAppLayout() {
	p_open := &ui.ShowAppLayout
	im.SetNextWindowSize(f64.Vec2{500, 440}, imgui.CondFirstUseEver)
	if im.BeginEx("Example: Layout", p_open, imgui.WindowFlagsMenuBar) {
		if im.BeginMenuBar() {
			if im.BeginMenu("File") {
				if im.MenuItem("Close") {
					*p_open = false
				}
				im.EndMenu()
			}
			im.EndMenuBar()
		}

		// left
		selected := &ui.AppLayout.Selected
		im.BeginChildEx("left pane", f64.Vec2{150, 0}, true, 0)
		for i := 0; i < 100; i++ {
			label := fmt.Sprintf("MyObject %v", i)
			if im.SelectableEx(label, *selected == i, 0, f64.Vec2{}) {
				*selected = i
			}
		}
		im.EndChild()
		im.SameLine()

		// right
		im.BeginGroup()

		// Leave room for 1 line below us
		im.BeginChildEx("item view", f64.Vec2{0, -im.GetFrameHeightWithSpacing()}, true, 0)
		im.Text("MyObject: %d", *selected)
		im.Separator()
		im.TextWrapped("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. ")
		im.EndChild()

		if im.Button("Revert") {
		}
		im.SameLine()
		if im.Button("Save") {
		}

		im.EndGroup()
	}
	im.End()
}

func ShowExampleAppPropertyEditor() {
	p_open := &ui.ShowAppPropertyEditor
	im.SetNextWindowSize(f64.Vec2{430, 450}, imgui.CondFirstUseEver)
	if !im.BeginEx("Example: Property editor", p_open, 0) {
		im.End()
		return
	}

	ShowHelpMarker("This example shows how you may implement a property editor using two columns.\nAll objects/fields data are dummies here.\nRemember that in many simple cases, you can use ImGui::SameLine(xxx) to position\nyour cursor horizontally instead of using the Columns() API.")

	im.PushStyleVar(imgui.StyleVarFramePadding, f64.Vec2{2, 2})
	im.Columns(2, "", true)
	im.Separator()

	var ShowDummyObject func(prefix string, uid int)
	ShowDummyObject = func(prefix string, uid int) {
		// Use object uid as identifier. Most commonly you could also use the object pointer as a base ID.
		im.PushID(imgui.ID(uid))
		// Text and Tree nodes are less high than regular widgets, here we add vertical spacing to make the tree lines equal high.
		im.AlignTextToFramePadding()

		node_open := im.TreeNodeStringID("Object", "%s_%u", prefix, uid)
		im.NextColumn()
		im.AlignTextToFramePadding()
		im.Text("my sailor is rich")
		im.NextColumn()
		if node_open {
			dummy_members := ui.AppLongText.DummyMembers[:]
			for i := 0; i < 8; i++ {
				im.PushID(imgui.ID(i))
				if i < 2 {
					ShowDummyObject("Child", 424242)
				} else {
					im.AlignTextToFramePadding()
					label := fmt.Sprintf("Field_%d", i)
					im.Bullet()
					im.Selectable(label)
					im.NextColumn()
					im.PushItemWidth(-1)
					if i >= 5 {
						im.InputFloat("##value", &dummy_members[i], 1.0)
					} else {
						im.DragFloatEx("##value", &dummy_members[i], 0.01, 0, 0, "", 1)
					}
					im.PopItemWidth()
					im.NextColumn()
				}
				im.PopID()
			}
			im.TreePop()
		}
		im.PopID()
	}

	// Iterate dummy objects with dummy members (all the same data)
	for obj_i := 0; obj_i < 3; obj_i++ {
		ShowDummyObject("Object", obj_i)
	}

	im.Columns(1, "", true)
	im.Separator()
	im.PopStyleVar()
	im.End()
}

// Demonstrate/test rendering huge amount of text, and the incidence of clipping.
func ShowExampleAppLongText() {
	p_open := &ui.ShowAppLongText
	test_type := &ui.AppLongText.TestType
	lines := &ui.AppLongText.Lines
	log := &ui.AppLongText.Log

	im.SetNextWindowSize(f64.Vec2{520, 600}, imgui.CondFirstUseEver)
	if !im.BeginEx("Example: Long text display", p_open, 0) {
		im.End()
		return
	}
	im.Text("Printing unusually long amount of text.")
	im.ComboString("Test type", test_type, []string{"Single call to TextUnformatted()", "Multiple calls to Text(), clipped manually", "Multiple calls to Text(), not clipped (slow)"})
	im.Text("Buffer contents: %d lines, %d bytes", *lines, log.Len())
	if im.Button("Clear") {
		log.Reset()
		*lines = 0
	}
	im.SameLine()
	if im.Button("Add 1000 lines") {
		for i := 0; i < 1000; i++ {
			fmt.Fprintf(log, "%d The quick brown fox jumps over the lazy dog\n", *lines+i)
		}
		*lines += 1000
	}
	im.BeginChild("Log")
	switch *test_type {
	case 0:
		// Single call to TextUnformatted() with a big buffer
		im.TextUnformatted(log.String())
	case 1:
		// Multiple calls to Text(), manually coarsely clipped - demonstrate how to use the ImGuiListClipper helper.
		im.PushStyleVar(imgui.StyleVarItemSpacing, f64.Vec2{0, 0})
		var clipper imgui.ListClipper
		clipper.Init(im, *lines, -1)
		for clipper.Step() {
			for i := clipper.DisplayStart; i < clipper.DisplayEnd; i++ {
				im.Text("%d The quick brown fox jumps over the lazy dog", i)
			}
		}
		im.PopStyleVar()
	case 2:
		// Multiple calls to Text(), not clipped (slow)
		im.PushStyleVar(imgui.StyleVarItemSpacing, f64.Vec2{0, 0})
		for i := 0; i < *lines; i++ {
			im.Text("%d The quick brown fox jumps over the lazy dog", i)
		}
		im.PopStyleVar()
	}
	im.EndChild()
	im.End()
}

// Demonstrate creating a window which gets auto-resized according to its content.
func ShowExampleAppAutoResize() {
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
func ShowExampleAppConstrainedResize() {
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
func ShowExampleAppFixedOverlay() {
	const DISTANCE = 10.0

	corner := &ui.AppFixedOverlay.Corner
	p_open := &ui.ShowAppFixedOverlay

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
	if *corner != -1 {
		im.SetNextWindowPos(window_pos, imgui.CondAlways, window_pos_pivot)
	}
	// Transparent background
	im.SetNextWindowBgAlpha(0.3)
	window_flags := imgui.WindowFlagsNoTitleBar | imgui.WindowFlagsNoResize | imgui.WindowFlagsAlwaysAutoResize | imgui.WindowFlagsNoMove | imgui.WindowFlagsNoSavedSettings | imgui.WindowFlagsNoFocusOnAppearing | imgui.WindowFlagsNoNav
	if *corner != -1 {
		window_flags |= imgui.WindowFlagsNoMove
	}
	if im.BeginEx("Example: Fixed Overlay", p_open, window_flags) {
		im.Text("Simple overlay\nin the corner of the screen.\n(right-click to change position)")
		im.Separator()
		if im.IsMousePosValid() {
			im.Text("Mouse Position: (%.1f,%.1f)", im.GetIO().MousePos.X, im.GetIO().MousePos.Y)
		} else {
			im.Text("Mouse Position: <invalid>")
		}
		if im.BeginPopupContextWindow() {
			if im.MenuItemEx("Custom", "", *corner == -1, true) {
				*corner = -1
			}
			if im.MenuItemEx("Top-left", "", *corner == 0, true) {
				*corner = 0
			}
			if im.MenuItemEx("Top-right", "", *corner == 1, true) {
				*corner = 1
			}
			if im.MenuItemEx("Bottom-left", "", *corner == 2, true) {
				*corner = 2
			}
			if im.MenuItemEx("Bottom-right", "", *corner == 3, true) {
				*corner = 3
			}
			if im.MenuItem("Close") {
				*p_open = false
			}
			im.EndPopup()
		}
		im.End()
	}
}

// Demonstrate using "##" and "###" in identifiers to manipulate ID generation.
// This apply to regular items as well. Read FAQ section "How can I have multiple widgets with the same label? Can I have widget without a label? (Yes). A primer on the purpose of labels/IDs." for details.
func ShowExampleAppWindowTitles() {
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

// Demonstrate using the low-level ImDrawList to draw custom shapes.
func ShowExampleAppCustomRendering() {
	p_open := &ui.ShowAppCustomRendering

	im.SetNextWindowSize(f64.Vec2{350, 560}, imgui.CondFirstUseEver)
	if !im.BeginEx("Example: Custom rendering", p_open, 0) {
		im.End()
		return
	}

	// Tip: If you do a lot of custom rendering, you probably want to use your own geometrical types and benefit of overloaded operators, etc.
	// Define IM_VEC2_CLASS_EXTRA in imconfig.h to create implicit conversions between your types and ImVec2/ImVec4.
	// ImGui defines overloaded operators but they are internal to imgui.cpp and not exposed outside (to avoid messing with your types)
	// In this example we are not using the maths operators!
	draw_list := im.GetWindowDrawList()

	// Primitives
	im.Text("Primitives")
	sz := &ui.AppCustomRendering.Sz
	col := &ui.AppCustomRendering.Col
	im.DragFloatEx("Size", sz, 0.2, 2.0, 72.0, "%.0f", 1)
	im.ColorEdit3("Color", col)
	{
		p := im.GetCursorScreenPos()
		x := p.X + 4.0
		y := p.Y + 4.0
		spacing := 8.0
		for n := 0; n < 2; n++ {
			thickness := 4.0
			if n == 0 {
				thickness = 1.0
			}
			draw_list.AddCircleEx(f64.Vec2{x + *sz*0.5, y + *sz*0.5}, *sz*0.5, *col, 20, thickness)
			x += *sz + spacing
			draw_list.AddRectEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, 0.0, imgui.DrawCornerFlagsAll, thickness)
			x += *sz + spacing
			draw_list.AddRectEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, 10.0, imgui.DrawCornerFlagsAll, thickness)
			x += *sz + spacing
			draw_list.AddRectEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, 10.0, imgui.DrawCornerFlagsTopLeft|imgui.DrawCornerFlagsBotRight, thickness)
			x += *sz + spacing
			draw_list.AddTriangleEx(f64.Vec2{x + *sz*0.5, y}, f64.Vec2{x + *sz, y + *sz - 0.5}, f64.Vec2{x, y + *sz - 0.5}, *col, thickness)
			x += *sz + spacing
			draw_list.AddLineEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y}, *col, thickness)
			x += *sz + spacing
			draw_list.AddLineEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, thickness)
			x += *sz + spacing
			draw_list.AddLineEx(f64.Vec2{x, y}, f64.Vec2{x, y + *sz}, *col, thickness)
			x += spacing
			draw_list.AddBezierCurve(f64.Vec2{x, y}, f64.Vec2{x + *sz*1.3, y + *sz*0.3}, f64.Vec2{x + *sz - *sz*1.3, y + *sz - *sz*0.3}, f64.Vec2{x + *sz, y + *sz}, *col, thickness)
			x = p.X + 4
			y += *sz + spacing
		}
		draw_list.AddCircleFilledEx(f64.Vec2{x + *sz*0.5, y + *sz*0.5}, *sz*0.5, *col, 32)
		x += *sz + spacing
		draw_list.AddRectFilled(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col)
		x += *sz + spacing
		draw_list.AddRectFilledEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, 10.0, 0)
		x += *sz + spacing
		draw_list.AddRectFilledEx(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz}, *col, 10.0, imgui.DrawCornerFlagsTopLeft|imgui.DrawCornerFlagsBotRight)
		x += *sz + spacing
		draw_list.AddTriangleFilled(f64.Vec2{x + *sz*0.5, y}, f64.Vec2{x + *sz, y + *sz - 0.5}, f64.Vec2{x, y + *sz - 0.5}, *col)
		x += *sz + spacing
		draw_list.AddRectFilledMultiColor(f64.Vec2{x, y}, f64.Vec2{x + *sz, y + *sz},
			color.RGBA{0, 0, 0, 255}, color.RGBA{255, 0, 0, 255}, color.RGBA{255, 255, 0, 255}, color.RGBA{0, 255, 0, 255})
		im.Dummy(f64.Vec2{(*sz + spacing) * 8, (*sz + spacing) * 3})
	}
	im.Separator()
	{
		points := &ui.AppCustomRendering.Points
		adding_line := &ui.AppCustomRendering.AddingLine
		im.Text("Canvas example")
		if im.Button("Clear") {
			*points = (*points)[:0]
		}
		if len(*points) >= 2 {
			im.SameLine()
			if im.Button("Undo") {
				*points = (*points)[len(*points)-2:]
			}
		}
		im.Text("Left-click and drag to add lines,\nRight-click to undo")

		// Here we are using InvisibleButton() as a convenience to 1) advance the cursor and 2) allows us to use IsItemHovered()
		// However you can draw directly and poll mouse/keyboard by yourself. You can manipulate the cursor using GetCursorPos() and SetCursorPos().
		// If you only use the ImDrawList API, you can notify the owner window of its extends by using SetCursorPos(max).

		// ImDrawList API uses screen coordinates!
		canvas_pos := im.GetCursorScreenPos()
		// Resize canvas to what's available
		canvas_size := im.GetContentRegionAvail()
		if canvas_size.X < 50.0 {
			canvas_size.X = 50.0
		}
		if canvas_size.Y < 50.0 {
			canvas_size.Y = 50.0
		}
		draw_list.AddRectFilledMultiColor(canvas_pos, f64.Vec2{canvas_pos.X + canvas_size.X, canvas_pos.Y + canvas_size.Y},
			color.RGBA{50, 50, 50, 255}, color.RGBA{50, 50, 60, 255}, color.RGBA{60, 60, 70, 255}, color.RGBA{50, 50, 60, 255})
		draw_list.AddRect(canvas_pos, f64.Vec2{canvas_pos.X + canvas_size.X, canvas_pos.Y + canvas_size.Y}, color.RGBA{255, 255, 255, 255})

		adding_preview := false
		im.InvisibleButton("canvas", canvas_size)
		mouse_pos_in_canvas := f64.Vec2{im.GetIO().MousePos.X - canvas_pos.X, im.GetIO().MousePos.Y - canvas_pos.Y}
		if *adding_line {
			adding_preview = true
			*points = append(*points, mouse_pos_in_canvas)
			if !im.IsMouseDown(0) {
				*adding_line = false
				adding_preview = false
			}
		}
		if im.IsItemHovered() {
			if !*adding_line && im.IsMouseClicked(0, true) {
				*points = append(*points, mouse_pos_in_canvas)
				*adding_line = true
			}
			if im.IsMouseClicked(1, true) && len(*points) > 0 {
				*adding_line, adding_preview = false, false
				*points = (*points)[:len(*points)-2]
			}
		}
		// clip lines within the canvas (if we resize it, etc.)
		draw_list.PushClipRectEx(canvas_pos, f64.Vec2{canvas_pos.X + canvas_size.X, canvas_pos.Y + canvas_size.Y}, true)
		for i := 0; i < len(*points)-1; i += 2 {
			draw_list.PushClipRectEx(canvas_pos, f64.Vec2{canvas_pos.X + canvas_size.X, canvas_pos.Y + canvas_size.Y}, true)
		}
		draw_list.PopClipRect()
		if adding_preview {
			*points = (*points)[:len(*points)-1]
		}
	}
	im.End()
}

func ShowAppAbout() {
	im.BeginEx("About Dear ImGui", &ui.ShowAppAbout, imgui.WindowFlagsAlwaysAutoResize)
	im.Text("Dear ImGui, %s", im.GetVersion())
	im.Separator()
	im.Text("By Omar Cornut and all dear imgui contributors.")
	im.Text("Dear ImGui is licensed under the MIT License, see LICENSE for more information.")
	im.End()
}

func ShowExampleMenuFile() {
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
				ShowExampleMenuFile()
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

func ShowStyleEditor(ref *imgui.Style) {
	style := im.GetStyle()
	ref_saved_style := &ui.StyleEditor.RefSavedStyle

	// Default to using internal storage as reference
	init := &ui.StyleEditor.Init
	if *init && ref == nil {
		*ref_saved_style = *style
	}
	*init = false
	if ref == nil {
		ref = ref_saved_style
	}

	im.PushItemWidth(im.GetWindowWidth() * 0.50)

	if ShowStyleSelector("Colors##Selector") {
		*ref_saved_style = *style
	}
	ShowFontSelector("Fonts##Selector")

	// Simplified Settings
	if im.SliderFloatEx("FrameRounding", &style.FrameRounding, 0.0, 12.0, "%.0f", 1.0) {
		// Make GrabRounding always the same as FrameRounding
		style.GrabRounding = style.FrameRounding
	}

	{
		window_border := style.WindowBorderSize > 0
		if im.Checkbox("WindowBorder", &window_border) {
			style.WindowBorderSize = 0
			if window_border {
				style.WindowBorderSize = 1
			}
		}
	}
	im.SameLine()

	{
		frame_border := style.FrameBorderSize > 0
		if im.Checkbox("FrameBorder", &frame_border) {
			style.FrameBorderSize = 0
			if frame_border {
				style.FrameBorderSize = 1
			}
		}
	}
	im.SameLine()

	{
		popup_border := style.PopupBorderSize > 0
		if im.Checkbox("PopupBorder", &popup_border) {
			style.PopupBorderSize = 0
			if popup_border {
				style.PopupBorderSize = 1
			}
		}
	}

	// Save/Revert button
	if im.Button("Save Ref") {
		*ref_saved_style = *style
		*ref = *ref_saved_style
	}
	im.SameLine()
	if im.Button("Revert Ref") {
		*style = *ref
	}
	im.SameLine()
	ShowHelpMarker("Save/Revert in local non-persistent storage. Default Colors definition are not affected. Use \"Export Colors\" below to save them somewhere.")

	if im.TreeNode("Rendering") {
		im.Checkbox("Anti-aliased lines", &style.AntiAliasedLines)
		im.SameLine()
		ShowHelpMarker("When disabling anti-aliasing lines, you'll probably want to disable borders in your style as well.")
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
		im.SliderV2Ex("WindowPadding", &style.WindowPadding, 0.0, 20.0, "%.0f", 1)
		im.SliderFloatEx("PopupRounding", &style.PopupRounding, 0.0, 16.0, "%.0f", 1)
		im.SliderV2Ex("FramePadding", &style.FramePadding, 0.0, 20.0, "%.0f", 1)
		im.SliderV2Ex("ItemSpacing", &style.ItemSpacing, 0.0, 20.0, "%.0f", 1)
		im.SliderV2Ex("ItemInnerSpacing", &style.ItemInnerSpacing, 0.0, 20.0, "%.0f", 1)
		im.SliderV2Ex("TouchExtraPadding", &style.TouchExtraPadding, 0.0, 10.0, "%.0f", 1)
		im.SliderFloatEx("IndentSpacing", &style.IndentSpacing, 0.0, 30.0, "%.0f", 1)
		im.SliderFloatEx("ScrollbarSize", &style.ScrollbarSize, 1.0, 20.0, "%.0f", 1)
		im.SliderFloatEx("GrabMinSize", &style.GrabMinSize, 1.0, 20.0, "%.0f", 1)
		im.Text("BorderSize")
		im.SliderFloatEx("WindowBorderSize", &style.WindowBorderSize, 0.0, 1.0, "%.0f", 1)
		im.SliderFloatEx("ChildBorderSize", &style.ChildBorderSize, 0.0, 1.0, "%.0f", 1)
		im.SliderFloatEx("PopupBorderSize", &style.PopupBorderSize, 0.0, 1.0, "%.0f", 1)
		im.SliderFloatEx("FrameBorderSize", &style.FrameBorderSize, 0.0, 1.0, "%.0f", 1)
		im.Text("Rounding")
		im.SliderFloatEx("WindowRounding", &style.WindowRounding, 0.0, 14.0, "%.0f", 1)
		im.SliderFloatEx("ChildRounding", &style.ChildRounding, 0.0, 16.0, "%.0f", 1)
		im.SliderFloatEx("FrameRounding", &style.FrameRounding, 0.0, 12.0, "%.0f", 1)
		im.SliderFloatEx("ScrollbarRounding", &style.ScrollbarRounding, 0.0, 12.0, "%.0f", 1)
		im.SliderFloatEx("GrabRounding", &style.GrabRounding, 0.0, 12.0, "%.0f", 1)
		im.Text("Alignment")
		im.SliderV2Ex("WindowTitleAlign", &style.WindowTitleAlign, 0.0, 1.0, "%.2f", 1)
		im.SliderV2Ex("ButtonTextAlign", &style.ButtonTextAlign, 0.0, 1.0, "%.2f", 1)
		im.SameLine()
		ShowHelpMarker("Alignment applies when a button is larger than its text content.")
		im.Text("Safe Area Padding")
		im.SameLine()
		ShowHelpMarker("Adjust if you cannot see the edges of your screen (e.g. on a TV where scaling has not been configured).")
		im.SliderV2Ex("DisplaySafeAreaPadding", &style.DisplaySafeAreaPadding, 0.0, 30.0, "%.0f", 1)
		im.TreePop()
	}

	if im.TreeNode("Colors") {
		output_dest := &ui.StyleEditor.Colors.OutputDest
		output_only_modified := &ui.StyleEditor.Colors.OutputOnlyModified
		if im.Button("Export Unsaved") {
			if *output_dest == 0 {
				im.LogToClipboard()
			} else {
				im.LogToTTY()
			}
			im.LogText("ImVec4* colors = ImGui::GetStyle().Colors;\n")
			for i := imgui.Col(0); i < imgui.ColCOUNT; i++ {
				col := style.Colors[i]
				name := im.GetStyleColorName(i)
				if !*output_only_modified || col != ref.Colors[i] {
					im.LogText("colors[ImGuiCol_%s]%*s= ImVec4(%.2ff, %.2ff, %.2ff, %.2ff);\n", name, 23-len(name), "", col.X, col.Y, col.Z, col.W)
				}
			}
			im.LogFinish()
		}
		im.SameLine()
		im.PushItemWidth(120)
		im.ComboString("##output_type", output_dest, []string{"To Clipboard", "To TTY"})
		im.PopItemWidth()
		im.SameLine()
		im.Checkbox("Only Modified Colors", output_only_modified)

		im.Text("Tip: Left-click on colored square to open color picker,\nRight-click to open edit options menu.")

		alpha_flags := &ui.StyleEditor.Colors.AlphaFlags
		im.RadioButtonEx("Opaque", alpha_flags, 0)
		im.SameLine()
		im.RadioButtonEx("Alpha", alpha_flags, int(imgui.ColorEditFlagsAlphaPreview))
		im.SameLine()
		im.RadioButtonEx("Both", alpha_flags, int(imgui.ColorEditFlagsAlphaPreviewHalf))
		im.SameLine()

		im.BeginChildEx("#colors", f64.Vec2{0, 300}, true, imgui.WindowFlagsAlwaysVerticalScrollbar|imgui.WindowFlagsAlwaysHorizontalScrollbar|imgui.WindowFlagsNavFlattened)
		im.PushItemWidth(-160)
		for i := imgui.Col(0); i < imgui.ColCOUNT; i++ {
			name := im.GetStyleColorName(i)
			im.PushID(imgui.ID(i))
			im.ColorEditV4Ex("##color", &style.Colors[i], imgui.ColorEditFlagsAlphaBar|imgui.ColorEditFlags(*alpha_flags))
			if style.Colors[i] != ref.Colors[i] {
				// Tips: in a real user application, you may want to merge and use an icon font into the main font, so instead of "Save"/"Revert" you'd use icons.
				// Read the FAQ and misc/fonts/README.txt about using icon fonts. It's really easy and super convenient!
				im.SameLineEx(0.0, style.ItemInnerSpacing.X)
				if im.Button("Save") {
					ref.Colors[i] = style.Colors[i]
				}
				im.SameLineEx(0.0, style.ItemInnerSpacing.X)
				if im.Button("Revert") {
					style.Colors[i] = ref.Colors[i]
				}
			}
			im.SameLineEx(0.0, style.ItemInnerSpacing.X)
			im.TextUnformatted(name)
			im.PopID()
		}
		im.PopItemWidth()
		im.EndChild()

		im.TreePop()
	}

	fonts_opened := im.TreeNodeStringID("Fonts", "Fonts (%d)", len(im.GetIO().Fonts.Fonts))
	if fonts_opened {
		atlas := im.GetIO().Fonts
		if im.TreeNodeStringIDEx("Atlas Texture", 0, "Atlas texture (%dx%d pixels)", atlas.TexWidth, atlas.TexHeight) {
			im.Image(atlas.TexID, f64.Vec2{float64(atlas.TexWidth), float64(atlas.TexHeight)}, f64.Vec2{0, 0}, f64.Vec2{1, 1}, color.RGBA{255, 255, 255, 255}, color.RGBA{255, 255, 255, 128})
			im.TreePop()
		}
		im.PushItemWidth(100)
		for i := range atlas.Fonts {
			font := atlas.Fonts[i]
			im.PushID(font.ID)
			font_name := ""
			if font.ConfigData != nil {
				font_name = font.ConfigData[0].Name
			}
			font_details_opened := im.TreeNodeIDEx(font.ID, 0, "Font %d: '%s', %.2f px, %d glyphs", i, font_name, font.FontSize, len(font.Glyphs))
			im.SameLine()
			if im.SmallButton("Set as default") {
				im.GetIO().FontDefault = font
			}
			if font_details_opened {
				im.PushFont(font)
				im.Text("The quick brown fox jumps over the lazy dog")
				im.PopFont()
				// Scale only this font
				im.DragFloatEx("Font scale", &font.Scale, 0.005, 0.3, 2.0, "%.1f", 1)
				im.InputFloatEx("Font offset", &font.DisplayOffset.Y, 1, 1, "", 1)
				im.SameLine()
				ShowHelpMarker("Note than the default embedded font is NOT meant to be scaled.\n\nFont are currently rendered into bitmaps at a given size at the time of building the atlas. You may oversample them to get some flexibility with scaling. You can also render at multiple sizes and select which one to use at runtime.\n\n(Glimmer of hope: the atlas system should hopefully be rewritten in the future to make scaling more natural and automatic.)")
				im.Text("Ascent: %f, Descent: %f, Height: %f", font.Ascent, font.Descent, font.Ascent-font.Descent)
				im.Text("Fallback character: '%c' (%d)", font.FallbackChar, font.FallbackChar)
				im.Text("Texture surface: %d pixels (approx) ~ %dx%d", font.MetricsTotalSurface,
					int(math.Sqrt(float64(font.MetricsTotalSurface))),
					int(math.Sqrt(float64(font.MetricsTotalSurface))))
				for config_i := 0; config_i < font.ConfigDataCount; config_i++ {
					cfg := &font.ConfigData[config_i]
					if cfg != nil {
						im.BulletText("Input %d: '%s', Oversample: (%d,%d), PixelSnapH: %v", config_i, cfg.Name, cfg.OversampleH, cfg.OversampleV, cfg.PixelSnapH)
					}
				}
				if im.TreeNodeStringIDEx("Glyphs", 0, "Glyphs (%d)", len(font.Glyphs)) {
					// Display all glyphs of the fonts in separate pages of 256 characters
					for base := 0; base < 0x10000; base += 256 {
						count := 0
						for n := 0; n < 256; n++ {
							if font.FindGlyphNoFallback(rune(base+n)) != nil {
								count += 1
							}
						}
						glyph_str := "glyph"
						if count > 0 {
							glyph_str = "glyphs"
						}
						if count > 0 && im.TreeNodeIDEx(imgui.ID(base), 0, "U+%04X..U+%04X (%d %s)", base, base+255, count, glyph_str) {
							cell_size := font.FontSize * 1
							cell_spacing := style.ItemSpacing.Y
							base_pos := im.GetCursorScreenPos()
							draw_list := im.GetWindowDrawList()
							for n := 0; n < 256; n++ {
								cell_p1 := f64.Vec2{
									base_pos.X + float64(n%16)*(cell_size+cell_spacing),
									base_pos.Y + float64(n/16)*(cell_size+cell_spacing),
								}
								cell_p2 := f64.Vec2{cell_p1.X + cell_size, cell_p1.Y + cell_size}
								glyph := font.FindGlyphNoFallback(rune(base + n))
								col := color.RGBA{255, 255, 255, 100}
								if glyph != nil {
									col = color.RGBA{255, 255, 255, 50}
								}
								draw_list.AddRect(cell_p1, cell_p2, col)
								// We use ImFont::RenderChar as a shortcut because we don't have UTF-8 conversion functions available to generate a string.
								font.RenderChar(draw_list, cell_size, cell_p1, im.GetColorFromStyle(imgui.ColText), rune(base+n))
								if glyph != nil && im.IsMouseHoveringRect(cell_p1, cell_p2) {
									im.BeginTooltip()
									im.Text("Codepoint: U+%04X", base+n)
									im.Separator()
									im.Text("AdvanceX: %.1f", glyph.AdvanceX)
									im.Text("Pos: (%.2f,%.2f)->(%.2f,%.2f)", glyph.X0, glyph.Y0, glyph.X1, glyph.Y1)
									im.Text("UV: (%.3f,%.3f)->(%.3f,%.3f)", glyph.U0, glyph.V0, glyph.U1, glyph.V1)
									im.EndTooltip()
								}
							}
							im.Dummy(f64.Vec2{(cell_size + cell_spacing) * 16, (cell_size + cell_spacing) * 16})
							im.TreePop()
						}
					}
					im.TreePop()
				}
				im.TreePop()
			}
			im.PopID()
		}
		window_scale := &ui.StyleEditor.Fonts.WindowScale
		im.SetWindowFontScale(*window_scale)
		im.TreePop()
	}

	im.PopItemWidth()
}

func ShowStyleSelector(label string) bool {
	style_idx := &ui.StyleSelector.StyleIdx
	if im.ComboString(label, style_idx, []string{"Classic", "Dark", "Light"}) {
		switch *style_idx {
		case 0:
			im.StyleColorsClassic(nil)
		case 1:
			im.StyleColorsDark(nil)
		case 2:
			im.StyleColorsLight(nil)
		}
		return true
	}
	return false
}

// Demo helper function to select among loaded fonts.
// Here we use the regular BeginCombo()/EndCombo() api which is more the more flexible one.
func ShowFontSelector(label string) {
	io := im.GetIO()
	font_current := im.GetFont()
	if im.BeginCombo(label, font_current.GetDebugName()) {
		for n := range io.Fonts.Fonts {
			if im.SelectableEx(io.Fonts.Fonts[n].GetDebugName(), io.Fonts.Fonts[n] == font_current, 0, f64.Vec2{}) {
				io.FontDefault = io.Fonts.Fonts[n]
			}
		}
		im.EndCombo()
	}
	im.SameLine()
	ShowHelpMarker(
		"- Load additional fonts with io.Fonts->AddFontFromFileTTF().\n" +
			"- The font atlas is built when calling io.Fonts->GetTexDataAsXXXX() or io.Fonts->Build().\n" +
			"- Read FAQ and documentation in misc/fonts/ for more details.\n" +
			"- If you need to add/remove fonts at runtime (e.g. for DPI change), do it before calling NewFrame().")
}

func ShowHelpMarker(desc string) {
	im.TextDisabled("(?)")
	if im.IsItemHovered() {
		im.BeginTooltip()
		im.PushTextWrapPos(im.GetFontSize() * 35.0)
		im.TextUnformatted(desc)
		im.PopTextWrapPos()
		im.EndTooltip()
	}
}

func ShowMetricsWindow() {
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

// Usage:
//  static ExampleAppLog my_log;
//  my_log.AddLog("Hello %d world\n", 123);
//  my_log.Draw("title");
type ExampleAppLog struct {
	Buf            strings.Builder
	Filter         imgui.TextFilter
	Lines          []string
	ScrollToBottom bool
}

func (c *ExampleAppLog) AddLog(format string, args ...interface{}) {
}

func (c *ExampleAppLog) Clear() {
}

func (c *ExampleAppLog) Draw(title string, p_open *bool) {
	im.SetNextWindowSize(f64.Vec2{500, 400}, imgui.CondFirstUseEver)
	im.BeginEx(title, p_open, 0)
	if im.Button("Clear") {
		c.Clear()
	}
	im.SameLine()
	copy := im.Button("Copy")
	im.SameLine()
	im.Separator()
	im.BeginChildEx("scrolling", f64.Vec2{0, 0}, false, imgui.WindowFlagsHorizontalScrollbar)
	if copy {
		im.LogToClipboard()
	}
	if c.Filter.IsActive() {
	} else {
	}

	if c.ScrollToBottom {
		im.SetScrollHereEx(1.0)
	}
	c.ScrollToBottom = false
	im.EndChild()
	im.End()
}
