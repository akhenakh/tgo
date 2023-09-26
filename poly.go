package tgo

/*
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Poly struct {
	cp *C.struct_tg_poly
}

type MultiPoly struct {
	cg *C.struct_tg_geom
}

func (p *Poly) AsGeom() *Geom {
	cg := (*C.struct_tg_geom)(unsafe.Pointer(p.cp))
	return &Geom{
		cg: cg,
	}
}

func (p *Poly) AsText() string {
	return p.AsGeom().AsText()
}

func (p *Poly) HolesCount() int {
	return int(C.tg_poly_num_holes(p.cp))
}

func (mp *MultiPoly) AsGeom() *Geom {
	return &Geom{
		cg: mp.cg,
	}
}

func (mp *MultiPoly) AsText() string {
	return mp.AsGeom().AsText()
}

func (mp *MultiPoly) PolygonsCount() int {
	return int(C.tg_geom_num_polys(mp.cg))
}

func (mp *MultiPoly) PolygonAt(index int) (*Poly, bool) {
	if index < 0 || index+1 > mp.PolygonsCount() {
		return nil, false
	}

	cp := C.tg_geom_poly_at(mp.cg, C.int(index))
	return &Poly{cp: cp}, true
}
