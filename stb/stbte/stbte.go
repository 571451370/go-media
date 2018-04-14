package stbte

/*
#include "stbte.h"
#define STB_TEXTEDIT_IMPLEMENTATION
#include "stb_textedit.h"
*/
import "C"
import (
	"sync"
	"unsafe"
)

type (
	State       C.STB_TexteditState
	TextEditRow C.StbTexteditRow
)

type String interface {
	Len() int
	GetChar(idx int) rune
	GetWidth(line_start_idx, char_idx int) float64
	LayoutRow(r *TextEditRow, line_start_idx int)
	MoveWordRight(idx int) int
	MoveWordLeft(idx int) int
	DeleteChars(pos, n int)
	InsertChars(pos int, new_text []byte) int
}

type context struct {
	sync.Mutex
	str String
}

var ctx context

func (s *State) Init(is_single_line bool) {
	C.stb_textedit_initialize_state((*C.STB_TexteditState)(s), truth(is_single_line))
}

func (s *State) Click(str String, x, y float64) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_click(nil, (*C.STB_TexteditState)(s), C.float(x), C.float(y))
}

func (s *State) Drag(str String, x, y float64) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_drag(nil, (*C.STB_TexteditState)(s), C.float(x), C.float(y))
}

func (s *State) Cut(str String) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_cut(nil, (*C.STB_TexteditState)(s))
}

func (s *State) Paste(str String, text []byte) int {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	return int(C.stb_textedit_paste(nil, (*C.STB_TexteditState)(s), (*C.char)(unsafe.Pointer(&text[0])), C.int(len(text))))
}

func (s *State) Key(str String, key int) {
	ctx.Lock()
	ctx.str = str
	defer ctx.Unlock()
	C.stb_textedit_key(nil, (*C.STB_TexteditState)(s), C.int(key))
}

func truth(cond bool) C.int {
	if cond {
		return 1
	}
	return 0
}

//export stringlen
func stringlen(unsafe.Pointer) int {
	p := ctx.str
	return p.Len()
}

//export getchar
func getchar(_ unsafe.Pointer, idx int) rune {
	p := ctx.str
	return p.GetChar(idx)
}

//export getwidth
func getwidth(_ unsafe.Pointer, line_start_idx, char_idx int) float64 {
	p := ctx.str
	return p.GetWidth(line_start_idx, char_idx)
}

//export layoutrow
func layoutrow(r *TextEditRow, _ unsafe.Pointer, line_start_idx int) {
	p := ctx.str
	p.LayoutRow(r, line_start_idx)
}

//export movewordright
func movewordright(_ unsafe.Pointer, idx int) int {
	p := ctx.str
	return p.MoveWordRight(idx)
}

//export movewordleft
func movewordleft(_ unsafe.Pointer, idx int) int {
	p := ctx.str
	return p.MoveWordLeft(idx)
}

//export insertchars
func insertchars(_ unsafe.Pointer, pos int, text *byte, new_text_len int) int {
	s := ((*[1 << 30]byte)(unsafe.Pointer(text)))[:new_text_len:new_text_len]
	p := ctx.str
	return p.InsertChars(pos, s)
}

//export deletechars
func deletechars(_ unsafe.Pointer, pos, n int) {
	p := ctx.str
	p.DeleteChars(pos, n)
}