package sdl

/*
#include "sdl.h"
*/
import "C"
import "unsafe"

type (
	Scancode C.SDL_Scancode
	Keycode  C.SDL_Keycode
)

type Keysym struct {
	Scancode Scancode
	Sym      Keycode
	Mod      uint16
}

func GetKeyName(key Keycode) string {
	return C.GoString(C.SDL_GetKeyName(C.SDL_Keycode(key)))
}

func GetKeyFromName(name string) Keycode {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Keycode(C.SDL_GetKeyFromName(cname))
}

func StartTextInput() {
	C.SDL_StartTextInput()
}

func IsTextInputActive() bool {
	return C.SDL_IsTextInputActive() != 0
}

func StopTextInput() {
	C.SDL_StopTextInput()
}
