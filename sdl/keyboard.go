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

func GetKeyboardState() []uint8 {
	var numkeys C.int
	state := C.SDL_GetKeyboardState(&numkeys)
	return C.GoBytes(unsafe.Pointer(state), numkeys)
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
