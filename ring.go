package tgo

/*
#cgo LDFLAGS: -lm
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Ring represents a ring geometry element.
type Ring struct {
	cr *C.struct_tg_ring
}

// AsGeom converts the ring to a geometric object.
func (r *Ring) AsGeom() *Geom {
	cg := (*C.struct_tg_geom)(unsafe.Pointer(r.cr))
	return &Geom{
		cg: cg,
	}
}

// AsPoly converts the ring to a polygon.
func (r *Ring) AsPoly() *Poly {
	cp := (*C.struct_tg_poly)(unsafe.Pointer(r.cr))
	return &Poly{
		cp: cp,
	}
}

// AsText returns the Well-Known Text representation of the ring.
func (r *Ring) AsText() string {
	return r.AsGeom().AsText()
}

// Area returns the area of the ring.
func (r *Ring) Area() float64 {
	return float64(C.tg_ring_area(r.cr))
}

// Perimeter returns the perimeter of the ring.
func (r *Ring) Perimeter() float64 {
	return float64(C.tg_ring_perimeter(r.cr))
}
