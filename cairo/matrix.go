package cairo

/*
#include <cairo.h>
*/
import "C"
import "unsafe"

type Matrix struct {
	XX, YX float64
	XY, YY float64
	X0, Y0 float64
}

func (m *Matrix) InitIdentity() {
	C.cairo_matrix_init_identity((*C.cairo_matrix_t)(unsafe.Pointer(m)))
}

func (m *Matrix) InitScale(sx, sy float64) {
	C.cairo_matrix_init_scale((*C.cairo_matrix_t)(unsafe.Pointer(m)), C.double(sx), C.double(sy))
}
