package tgo

/*
#cgo LDFLAGS: -lm
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Ring struct {
	cr *C.struct_tg_ring
}

func (r *Ring) AsGeom() *Geom {
	cg := (*C.struct_tg_geom)(unsafe.Pointer(r.cr))
	return &Geom{
		cg: cg,
	}
}

func (r *Ring) AsPoly() *Poly {
	cp := (*C.struct_tg_poly)(unsafe.Pointer(r.cr))
	return &Poly{
		cp: cp,
	}
}

func (r *Ring) AsText() string {
	return r.AsGeom().AsText()
}

func (r *Ring) Area() float64 {
	return float64(C.tg_ring_area(r.cr))
}

func (r *Ring) Perimeter() float64 {
	return float64(C.tg_ring_perimeter(r.cr))
}
