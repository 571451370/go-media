package readline

/*

#include <stdlib.h>
#include <readline/readline.h>
#include <readline/history.h>
#include "go_readline.h"

*/
// #cgo LDFLAGS: -lreadline
import "C"

import (
	"sync"
	"unsafe"
)

type handles struct {
	sync.Mutex
	key map[int]func(int, int) int
}

func newHandles() *handles {
	return &handles{
		key: make(map[int]func(int, int) int),
	}
}

var (
	hands = newHandles()
)

//export goKeyCallback
func goKeyCallback(arg0, invokingKey C.int) C.int {
	hands.Lock()
	defer hands.Unlock()
	f := hands.key[int(invokingKey)]
	return C.int(f(int(arg0), int(invokingKey)))
}

func BindKey(key int, fun func(int, int) int) {
	hands.Lock()
	defer hands.Unlock()
	hands.key[key] = fun
	C.rl_bind_key(C.int(key), (*_Ctype_rl_command_func_t)(C.keyCallback))
}

func Read(prompt string) {
	cprompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cprompt))
	C.readline(cprompt)
}

func AddHistory(input string) {
	cinput := C.CString(input)
	defer C.free(unsafe.Pointer(cinput))
	C.add_history(cinput)
}

func Complete(ignore, invokingKey int) int {
	return int(C.rl_complete(C.int(ignore), C.int(invokingKey)))
}
