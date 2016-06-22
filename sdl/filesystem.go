package sdl

/*
#include "sdl.h"
*/
import "C"
import "unsafe"

func BasePath() string {
	return C.GoString(C.SDL_GetBasePath())
}

func GetPrefPath(org, app string) string {
	corg := C.CString(org)
	capp := C.CString(app)
	defer C.free(unsafe.Pointer(corg))
	defer C.free(unsafe.Pointer(capp))
	return C.GoString(C.SDL_GetPrefPath(corg, capp))
}
