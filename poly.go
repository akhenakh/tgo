package tgo

/*
#include "tg.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Poly is a polygon element. It is a view into a parent Geom unless it was
// produced by Clone or Copy, in which case it owns its data.
type Poly struct {
	cp *C.struct_tg_poly
}

// MultiPoly is a view of a MultiPolygon geometry.
type MultiPoly struct {
	cg *C.struct_tg_geom
}

// Clone returns an independent, reference-counted copy of the polygon. It is
// cheap.
func (p *Poly) Clone() *Poly {
	cp := C.tg_poly_clone(p.cp)
	if cp == nil {
		return nil
	}

	out := &Poly{cp: cp}
	runtime.SetFinalizer(out, (*Poly).free)

	return out
}

// Copy returns a deep copy of the polygon.
func (p *Poly) Copy() *Poly {
	cp := C.tg_poly_copy(p.cp)
	if cp == nil {
		return nil
	}

	out := &Poly{cp: cp}
	runtime.SetFinalizer(out, (*Poly).free)

	return out
}

// AsGeom returns a Geom of the polygon without cloning. The result is a view
// that must not outlive the polygon.
func (p *Poly) AsGeom() *Geom {
	return &Geom{cg: (*C.struct_tg_geom)(unsafe.Pointer(p.cp))}
}

// AsWKT returns the representation of the poly as WKT.
func (p *Poly) AsWKT() string {
	return p.AsGeom().AsWKT()
}

// String implements fmt.Stringer and returns the poly as Well-Known Text.
func (p *Poly) String() string {
	return p.AsWKT()
}

// MemSize returns the allocation size of the polygon.
func (p *Poly) MemSize() int {
	return int(C.tg_poly_memsize(p.cp))
}

// Clockwise reports whether the polygon exterior winds clockwise.
func (p *Poly) Clockwise() bool {
	return bool(C.tg_poly_clockwise(p.cp))
}

// NumHoles returns the number of holes in the polygon.
func (p *Poly) NumHoles() int {
	return int(C.tg_poly_num_holes(p.cp))
}

// HoleAt returns the hole ring at index. The result is a view that must not
// outlive the polygon.
func (p *Poly) HoleAt(index int) (*Ring, bool) {
	if index < 0 || index >= p.NumHoles() {
		return nil, false
	}

	cr := C.tg_poly_hole_at(p.cp, C.int(index))
	if cr == nil {
		return nil, false
	}

	return &Ring{cr: cr}, true
}

// Exterior returns the exterior Ring of the poly. The result is a view that
// must not outlive the polygon.
func (p *Poly) Exterior() *Ring {
	return &Ring{cr: C.tg_poly_exterior(p.cp)}
}

// Rect returns the minimum bounding rectangle of the polygon.
func (p *Poly) Rect() Rect {
	return goRect(C.tg_poly_rect(p.cp))
}

func (p *Poly) free() {
	C.tg_poly_free(p.cp)
}

// AsGeom returns a Geom of the multipolygon without cloning. The result is a
// view that must not outlive the source geometry.
func (mp *MultiPoly) AsGeom() *Geom {
	return &Geom{cg: mp.cg}
}

// AsWKT returns the representation of the multipoly as WKT.
func (mp *MultiPoly) AsWKT() string {
	return mp.AsGeom().AsWKT()
}

// String implements fmt.Stringer and returns the multipoly as Well-Known Text.
func (mp *MultiPoly) String() string {
	return mp.AsWKT()
}

// NumPolygons returns the number of polygons in the multipoly.
func (mp *MultiPoly) NumPolygons() int {
	return int(C.tg_geom_num_polys(mp.cg))
}

// PolygonAt returns the Poly at index. The result is a view that must not
// outlive the source geometry.
func (mp *MultiPoly) PolygonAt(index int) (*Poly, bool) {
	if index < 0 || index >= mp.NumPolygons() {
		return nil, false
	}

	cp := C.tg_geom_poly_at(mp.cg, C.int(index))
	if cp == nil {
		return nil, false
	}

	return &Poly{cp: cp}, true
}
