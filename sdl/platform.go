package sdl

/*
#include "sdl.h"
*/
import "C"

func GetPlatform() string {
	return C.GoString(C.SDL_GetPlatform())
}
