package tgo

/*
#cgo LDFLAGS: -lm
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Poly represents a polygon geometry.
type Poly struct {
	cp *C.struct_tg_poly
}

// MultiPoly represents a multi-polygon geometry object.
type MultiPoly struct {
	cg *C.struct_tg_geom
}

// AsGeom returns a Geom of the polygon without cloning.
func (p *Poly) AsGeom() *Geom {
	cg := (*C.struct_tg_geom)(unsafe.Pointer(p.cp))
	return &Geom{
		cg: cg,
	}
}

// IsClockWise returns true if the polygon is clock wise.
func (p *Poly) IsClockWise() bool {
	return bool(C.tg_poly_clockwise(p.cp))
}

// AsWKT returns the representation of the poly as WKT.
func (p *Poly) AsWKT() string {
	return p.AsGeom().AsWKT()
}

// HolesCount returns the holes count.
func (p *Poly) HolesCount() int {
	return int(C.tg_poly_num_holes(p.cp))
}

// Exterior returns the exterior Ring of the poly.
func (p *Poly) Exterior() *Ring {
	cr := C.tg_poly_exterior(p.cp)
	return &Ring{
		cr: cr,
	}
}

// AsGeom returns a Geom of the multipolygon without cloning.
func (mp *MultiPoly) AsGeom() *Geom {
	return &Geom{
		cg: mp.cg,
	}
}

// AsWKT returns the representation of the multipoly as WKT.
func (mp *MultiPoly) AsWKT() string {
	return mp.AsGeom().AsWKT()
}

// PolygonsCount returns the count of polygons in the multipoly.
func (mp *MultiPoly) PolygonsCount() int {
	return int(C.tg_geom_num_polys(mp.cg))
}

// PolygonAt returns the Poly at index, true if it's applicable.
func (mp *MultiPoly) PolygonAt(index int) (*Poly, bool) {
	if index < 0 || index+1 > mp.PolygonsCount() {
		return nil, false
	}

	cp := C.tg_geom_poly_at(mp.cg, C.int(index))
	return &Poly{cp: cp}, true
}
