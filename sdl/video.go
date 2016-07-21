package sdl

/*
#include "sdl.h"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

type DisplayMode struct {
	Format     uint32
	W, H       int32
	Rate       int32
	driverdata *byte
}

type (
	Window               C.SDL_Window
	WindowFlags          C.SDL_WindowFlags
	WindowEventID        C.SDL_WindowEventID
	GLattr               C.SDL_GLattr
	GLprofile            C.SDL_GLprofile
	GLcontextFlag        C.SDL_GLcontextFlag
	GLcontextReleaseFlag C.SDL_GLcontextReleaseFlag
)

const (
	WINDOW_FULLSCREEN         WindowFlags = C.SDL_WINDOW_FULLSCREEN
	WINDOW_OPENGL             WindowFlags = C.SDL_WINDOW_OPENGL
	WINDOW_SHOWN              WindowFlags = C.SDL_WINDOW_SHOWN
	WINDOW_HIDDEN             WindowFlags = C.SDL_WINDOW_HIDDEN
	WINDOW_BORDERLESS         WindowFlags = C.SDL_WINDOW_BORDERLESS
	WINDOW_RESIZABLE          WindowFlags = C.SDL_WINDOW_RESIZABLE
	WINDOW_MINIMIZED          WindowFlags = C.SDL_WINDOW_MINIMIZED
	WINDOW_MAXIMIZED          WindowFlags = C.SDL_WINDOW_MAXIMIZED
	WINDOW_INPUT_GRABBED      WindowFlags = C.SDL_WINDOW_INPUT_GRABBED
	WINDOW_INPUT_FOCUS        WindowFlags = C.SDL_WINDOW_INPUT_FOCUS
	WINDOW_MOUSE_FOCUS        WindowFlags = C.SDL_WINDOW_MOUSE_FOCUS
	WINDOW_FULLSCREEN_DESKTOP WindowFlags = C.SDL_WINDOW_FULLSCREEN_DESKTOP
	WINDOW_FOREIGN            WindowFlags = C.SDL_WINDOW_FOREIGN
	WINDOW_ALLOW_HIGHDPI      WindowFlags = C.SDL_WINDOW_ALLOW_HIGHDPI
	WINDOW_MOUSE_CAPTURE      WindowFlags = C.SDL_WINDOW_MOUSE_CAPTURE
)

const (
	WINDOWPOS_UNDEFINED = C.SDL_WINDOWPOS_UNDEFINED
	WINDOWPOS_CENTERED  = C.SDL_WINDOWPOS_CENTERED
)

const (
	WINDOWEVENT_NONE         WindowEventID = C.SDL_WINDOWEVENT_NONE
	WINDOWEVENT_SHOWN        WindowEventID = C.SDL_WINDOWEVENT_SHOWN
	WINDOWEVENT_HIDDEN       WindowEventID = C.SDL_WINDOWEVENT_HIDDEN
	WINDOWEVENT_EXPOSED      WindowEventID = C.SDL_WINDOWEVENT_EXPOSED
	WINDOWEVENT_MOVED        WindowEventID = C.SDL_WINDOWEVENT_MOVED
	WINDOWEVENT_RESIZED      WindowEventID = C.SDL_WINDOWEVENT_RESIZED
	WINDOWEVENT_SIZE_CHANGED WindowEventID = C.SDL_WINDOWEVENT_SIZE_CHANGED
	WINDOWEVENT_MINIMIZED    WindowEventID = C.SDL_WINDOWEVENT_MINIMIZED
	WINDOWEVENT_MAXIMIZED    WindowEventID = C.SDL_WINDOWEVENT_MAXIMIZED
	WINDOWEVENT_RESTORED     WindowEventID = C.SDL_WINDOWEVENT_RESTORED
	WINDOWEVENT_ENTER        WindowEventID = C.SDL_WINDOWEVENT_ENTER
	WINDOWEVENT_LEAVE        WindowEventID = C.SDL_WINDOWEVENT_LEAVE
	WINDOWEVENT_FOCUS_GAINED WindowEventID = C.SDL_WINDOWEVENT_FOCUS_GAINED
	WINDOWEVENT_FOCUS_LOST   WindowEventID = C.SDL_WINDOWEVENT_FOCUS_LOST
	WINDOWEVENT_CLOSE        WindowEventID = C.SDL_WINDOWEVENT_CLOSE
)

const (
	GL_RED_SIZE                   GLattr = C.SDL_GL_RED_SIZE
	GL_GREEN_SIZE                 GLattr = C.SDL_GL_GREEN_SIZE
	GL_BLUE_SIZE                  GLattr = C.SDL_GL_BLUE_SIZE
	GL_ALPHA_SIZE                 GLattr = C.SDL_GL_ALPHA_SIZE
	GL_BUFFER_SIZE                GLattr = C.SDL_GL_BUFFER_SIZE
	GL_DOUBLEBUFFER               GLattr = C.SDL_GL_DOUBLEBUFFER
	GL_DEPTH_SIZE                 GLattr = C.SDL_GL_DEPTH_SIZE
	GL_STENCIL_SIZE               GLattr = C.SDL_GL_STENCIL_SIZE
	GL_ACCUM_RED_SIZE             GLattr = C.SDL_GL_ACCUM_RED_SIZE
	GL_ACCUM_GREEN_SIZE           GLattr = C.SDL_GL_ACCUM_GREEN_SIZE
	GL_ACCUM_BLUE_SIZE            GLattr = C.SDL_GL_ACCUM_BLUE_SIZE
	GL_ACCUM_ALPHA_SIZE           GLattr = C.SDL_GL_ACCUM_ALPHA_SIZE
	GL_STEREO                     GLattr = C.SDL_GL_STEREO
	GL_MULTISAMPLEBUFFERS         GLattr = C.SDL_GL_MULTISAMPLEBUFFERS
	GL_MULTISAMPLESAMPLES         GLattr = C.SDL_GL_MULTISAMPLESAMPLES
	GL_ACCELERATED_VISUAL         GLattr = C.SDL_GL_ACCELERATED_VISUAL
	GL_RETAINED_BACKING           GLattr = C.SDL_GL_RETAINED_BACKING
	GL_CONTEXT_MAJOR_VERSION      GLattr = C.SDL_GL_CONTEXT_MAJOR_VERSION
	GL_CONTEXT_MINOR_VERSION      GLattr = C.SDL_GL_CONTEXT_MINOR_VERSION
	GL_CONTEXT_EGL                GLattr = C.SDL_GL_CONTEXT_EGL
	GL_CONTEXT_FLAGS              GLattr = C.SDL_GL_CONTEXT_FLAGS
	GL_CONTEXT_PROFILE_MASK       GLattr = C.SDL_GL_CONTEXT_PROFILE_MASK
	GL_SHARE_WITH_CURRENT_CONTEXT GLattr = C.SDL_GL_SHARE_WITH_CURRENT_CONTEXT
	GL_FRAMEBUFFER_SRGB_CAPABLE   GLattr = C.SDL_GL_FRAMEBUFFER_SRGB_CAPABLE
	GL_CONTEXT_RELEASE_BEHAVIOR   GLattr = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR
)

const (
	GL_CONTEXT_PROFILE_CORE          GLprofile = C.SDL_GL_CONTEXT_PROFILE_CORE
	GL_CONTEXT_PROFILE_COMPATIBILITY GLprofile = C.SDL_GL_CONTEXT_PROFILE_COMPATIBILITY
	GL_CONTEXT_PROFILE_ES            GLprofile = C.SDL_GL_CONTEXT_PROFILE_ES
)

const (
	GL_CONTEXT_DEBUG_FLAG              GLcontextFlag = C.SDL_GL_CONTEXT_DEBUG_FLAG
	GL_CONTEXT_FORWARD_COMPATIBLE_FLAG GLcontextFlag = C.SDL_GL_CONTEXT_FORWARD_COMPATIBLE_FLAG
	GL_CONTEXT_ROBUST_ACCESS_FLAG      GLcontextFlag = C.SDL_GL_CONTEXT_ROBUST_ACCESS_FLAG
	GL_CONTEXT_RESET_ISOLATION_FLAG    GLcontextFlag = C.SDL_GL_CONTEXT_RESET_ISOLATION_FLAG
)

const (
	GL_CONTEXT_RELEASE_BEHAVIOR_NONE  GLcontextReleaseFlag = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR_NONE
	GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH GLcontextReleaseFlag = C.SDL_GL_CONTEXT_RELEASE_BEHAVIOR_FLUSH
)

func GetNumVideoDrivers() int {
	return int(C.SDL_GetNumVideoDrivers())
}

func GetVideoDriver(index int) string {
	return C.GoString(C.SDL_GetVideoDriver(C.int(index)))
}

func VideoInit(driverName string) error {
	cdriverName := C.CString(driverName)
	defer C.free(unsafe.Pointer(cdriverName))
	return ek(C.SDL_VideoInit(cdriverName))
}

func VideoQuit() {
	C.SDL_VideoQuit()
}

func GetCurrentVideoDriver() string {
	return C.GoString(C.SDL_GetCurrentVideoDriver())
}

func GetNumVideoDisplays() int {
	return int(C.SDL_GetNumVideoDisplays())
}

func GetDisplayName(displayIndex int) string {
	return C.GoString(C.SDL_GetDisplayName(C.int(displayIndex)))
}

func GetDisplayBounds(displayIndex int) (Rect, error) {
	var cr C.SDL_Rect
	rc := C.SDL_GetDisplayBounds(C.int(displayIndex), &cr)
	if rc < 0 {
		return Rect{}, GetError()
	}
	return Rect{int32(cr.x), int32(cr.y), int32(cr.w), int32(cr.h)}, nil
}

func CreateWindow(title string, x, y, w, h int, flags WindowFlags) (*Window, error) {
	return nil, GetError()
}

func (w *Window) Flags() WindowFlags {
	return WindowFlags(C.SDL_GetWindowFlags(w))
}

func (w *Window) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.SDL_SetWindowTitle(w, ctitle)
}

func (w *Window) Title() string {
	return C.GoString(C.SDL_GetWindowTitle(w))
}

func GetGrabbedWindow() *Window {
	return (*Window)(C.SDL_GetGrabbedWindow())
}

func (w *Window) SetBrightness(brightness float32) {
	C.SDL_SetWindowBrightness(w, C.float(brightness))
}

func (w *Window) Brightness() float32 {
	return float32(C.SDL_GetWindowBrightness(w))
}

func (w *Window) SetSize(width, height int) {
	C.SDL_SetWindowSize(w, C.int(width), C.int(height))
}

func (w *Window) Size() (width, height int) {
	var cw, ch C.int
	C.SDL_GetWindowSize(w, &cw, &ch)
	return int(cw), int(ch)
}

func (w *Window) SetMinimumSize(minWidth, minHeight int) {
	C.SDL_SetWindowMinimumSize(w, C.int(minWidth), C.int(minHeight))
}

func (w *Window) MinimumSize() (minWidth, minHeight int) {
	var mw, mh C.int
	C.SDL_GetWindowMinimumSize(w, &mw, &mh)
	return int(mw), int(mh)
}

func (w *Window) SetMaximumSize(maxWidth, maxHeight int) {
	C.SDL_SetWindowMaximumSize(w, C.int(maxWidth), C.int(maxHeight))
}

func (w *Window) MaximumSize() (maxWidth, maxHeight int) {
	var mw, mh C.int
	C.SDL_GetWindowMaximumSize(w, &mw, &mh)
	return int(mw), int(mh)
}

func (w *Window) SetBordered(bordered bool) {
	C.SDL_SetWindowBordered(w, truth(bordered))
}

func (w *Window) Show() {
	C.SDL_ShowWindow(w)
}

func (w *Window) Hide() {
	C.SDL_HideWindow(w)
}

func (w *Window) Raise() {
	C.SDL_RaiseWindow(w)
}

func (w *Window) Maximize() {
	C.SDL_MaximizeWindow(w)
}

func IsScreenSaverEnabled() bool {
	return C.SDL_IsScreenSaverEnabled() != 0
}

func EnableScreenSaver() {
	C.SDL_EnableScreenSaver()
}

func DisableScreenSaver() {
	C.SDL_DisableScreenSaver()
}

func GLSetAttribute(attr GLattr, value int) error {
	return ek(C.SDL_GL_SetAttribute(C.SDL_GLattr(attr), C.int(value)))
}

func GLGetAttribute(attr GLattr) (int, error) {
	var cvalue C.int
	rc := C.SDL_GL_GetAttribute(C.SDL_GLattr(attr), &cvalue)
	return int(cvalue), ek(rc)
}

func (t *Texture) Update(rect *Rect, pixels interface{}, pitch int) error {
	return ek(C.SDL_UpdateTexture(t, (*C.SDL_Rect)(unsafe.Pointer(rect)),
		unsafe.Pointer(reflect.ValueOf(pixels).Pointer()), C.int(pitch)))
}

func (t *Texture) SetBlendMode(mode BlendMode) error {
	return ek(C.SDL_SetTextureBlendMode(t, C.SDL_BlendMode(mode)))
}

func (s *Surface) SetBlendMode(blendMode BlendMode) error {
	return ek(C.SDL_SetSurfaceBlendMode(s, C.SDL_BlendMode(blendMode)))
}

func (w *Window) SetIcon(icon *Surface) {
	C.SDL_SetWindowIcon(w, icon)
}

func (w *Window) SetFullscreen(flags uint32) error {
	return ek(C.SDL_SetWindowFullscreen(w, C.Uint32(flags)))
}
