package sdl

/*
#include "sdl.h"
*/
import "C"
import "unsafe"

type HintPriority C.SDL_HintPriority

const (
	HINT_DEFAULT  HintPriority = C.SDL_HINT_DEFAULT
	HINT_NORMAL   HintPriority = C.SDL_HINT_NORMAL
	HINT_OVERRIDE HintPriority = C.SDL_HINT_OVERRIDE
)

const (
	HINT_FRAMEBUFFER_ACCELERATION        = C.SDL_HINT_FRAMEBUFFER_ACCELERATION
	HINT_RENDER_DRIVER                   = C.SDL_HINT_RENDER_DRIVER
	HINT_RENDER_OPENGL_SHADERS           = C.SDL_HINT_RENDER_OPENGL_SHADERS
	HINT_RENDER_DIRECT3D_THREADSAFE      = C.SDL_HINT_RENDER_DIRECT3D_THREADSAFE
	HINT_RENDER_DIRECT3D11_DEBUG         = C.SDL_HINT_RENDER_DIRECT3D11_DEBUG
	HINT_RENDER_SCALE_QUALITY            = C.SDL_HINT_RENDER_SCALE_QUALITY
	HINT_XINPUT_ENABLED                  = C.SDL_HINT_XINPUT_ENABLED
	HINT_XINPUT_USE_OLD_JOYSTICK_MAPPING = C.SDL_HINT_XINPUT_USE_OLD_JOYSTICK_MAPPING
	HINT_GAMECONTROLLERCONFIG            = C.SDL_HINT_GAMECONTROLLERCONFIG
	HINT_ALLOW_TOPMOST                   = C.SDL_HINT_ALLOW_TOPMOST
)

func SetHint(name, value string, prio HintPriority) bool {
	cname := C.CString(name)
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cname))
	defer C.free(unsafe.Pointer(cvalue))
	return C.SDL_SetHintWithPriority(cname, cvalue, C.SDL_HintPriority(prio)) != 0
}

func GetHint(name string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.GoString(C.SDL_GetHint(cname))
}

func ClearHints() {
	C.SDL_ClearHints()
}
