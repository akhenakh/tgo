package tgo

/*
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

func (r *Ring) AsText() string {
	return r.AsGeom().AsText()
}
