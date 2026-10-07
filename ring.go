package tgo

/*
#include "tg.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Ring is a closed ring element. It is a view into a parent Geom or Poly
// unless it was produced by Clone or Copy, in which case it owns its data.
type Ring struct {
	cr *C.struct_tg_ring
}

// Clone returns an independent, reference-counted copy of the ring. It is
// cheap.
func (r *Ring) Clone() *Ring {
	cr := C.tg_ring_clone(r.cr)
	if cr == nil {
		return nil
	}

	out := &Ring{cr: cr}
	runtime.SetFinalizer(out, (*Ring).free)

	return out
}

// Copy returns a deep copy of the ring.
func (r *Ring) Copy() *Ring {
	cr := C.tg_ring_copy(r.cr)
	if cr == nil {
		return nil
	}

	out := &Ring{cr: cr}
	runtime.SetFinalizer(out, (*Ring).free)

	return out
}

// AsGeom converts the ring to a geometry without cloning. The result is a view
// that must not outlive the ring.
func (r *Ring) AsGeom() *Geom {
	return &Geom{cg: (*C.struct_tg_geom)(unsafe.Pointer(r.cr))}
}

// AsPoly converts the ring to a polygon without cloning. The result is a view
// that must not outlive the ring.
func (r *Ring) AsPoly() *Poly {
	return &Poly{cp: (*C.struct_tg_poly)(unsafe.Pointer(r.cr))}
}

// AsWKT returns the Well-Known Text representation of the ring.
func (r *Ring) AsWKT() string {
	return r.AsGeom().AsWKT()
}

// String implements fmt.Stringer and returns the ring as Well-Known Text.
func (r *Ring) String() string {
	return r.AsWKT()
}

// MemSize returns the allocation size of the ring.
func (r *Ring) MemSize() int {
	return int(C.tg_ring_memsize(r.cr))
}

// Rect returns the minimum bounding rectangle of the ring.
func (r *Ring) Rect() Rect {
	return goRect(C.tg_ring_rect(r.cr))
}

// NumPoints returns the number of points in the ring.
func (r *Ring) NumPoints() int {
	return int(C.tg_ring_num_points(r.cr))
}

// PointAt returns the point at index.
func (r *Ring) PointAt(index int) (Point, bool) {
	if index < 0 || index >= r.NumPoints() {
		return Point{}, false
	}
	return goPoint(C.tg_ring_point_at(r.cr, C.int(index))), true
}

// Points returns a copy of the ring's points.
func (r *Ring) Points() []Point {
	n := r.NumPoints()
	if n == 0 {
		return nil
	}

	src := unsafe.Slice((*C.struct_tg_point)(C.tg_ring_points(r.cr)), n)
	out := make([]Point, n)
	for i := range src {
		out[i] = goPoint(src[i])
	}

	return out
}

// NumSegments returns the number of segments in the ring.
func (r *Ring) NumSegments() int {
	return int(C.tg_ring_num_segments(r.cr))
}

// SegmentAt returns the segment at index.
func (r *Ring) SegmentAt(index int) (Segment, bool) {
	if index < 0 || index >= r.NumSegments() {
		return Segment{}, false
	}
	return goSegment(C.tg_ring_segment_at(r.cr, C.int(index))), true
}

// Convex reports whether the ring is convex.
func (r *Ring) Convex() bool {
	return bool(C.tg_ring_convex(r.cr))
}

// Clockwise reports whether the ring winds clockwise.
func (r *Ring) Clockwise() bool {
	return bool(C.tg_ring_clockwise(r.cr))
}

// Area returns the area of the ring.
func (r *Ring) Area() float64 {
	return float64(C.tg_ring_area(r.cr))
}

// Perimeter returns the perimeter of the ring.
func (r *Ring) Perimeter() float64 {
	return float64(C.tg_ring_perimeter(r.cr))
}

// IndexSpread returns the number of segments grouped per index node, or zero
// when the ring is not indexed.
func (r *Ring) IndexSpread() int {
	return int(C.tg_ring_index_spread(r.cr))
}

// IndexNumLevels returns the number of index levels.
func (r *Ring) IndexNumLevels() int {
	return int(C.tg_ring_index_num_levels(r.cr))
}

// IndexLevelNumRects returns the number of rectangles at an index level.
func (r *Ring) IndexLevelNumRects(level int) int {
	return int(C.tg_ring_index_level_num_rects(r.cr, C.int(level)))
}

// IndexLevelRect returns the rectangle at the given level and position.
func (r *Ring) IndexLevelRect(level, rect int) Rect {
	return goRect(C.tg_ring_index_level_rect(r.cr, C.int(level), C.int(rect)))
}

func (r *Ring) free() {
	C.tg_ring_free(r.cr)
}
