// +build linux darwin openbsd netbsd freebsd dragonfly solaris

package sdl

/*
#include "sdl.h"

#cgo pkg-config: sdl2
*/
import "C"
