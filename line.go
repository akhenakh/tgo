package tgo

/*
#include "tg.h"
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// Line is a LineString element. It is a view into a parent Geom unless it was
// produced by Clone or Copy, in which case it owns its data.
type Line struct {
	cl *C.struct_tg_line
}

// Clone returns an independent, reference-counted copy of the line. It is
// cheap.
func (l *Line) Clone() *Line {
	cl := C.tg_line_clone(l.cl)
	if cl == nil {
		return nil
	}

	out := &Line{cl: cl}
	runtime.SetFinalizer(out, (*Line).free)

	return out
}

// Copy returns a deep copy of the line.
func (l *Line) Copy() *Line {
	cl := C.tg_line_copy(l.cl)
	if cl == nil {
		return nil
	}

	out := &Line{cl: cl}
	runtime.SetFinalizer(out, (*Line).free)

	return out
}

// AsGeom returns the line as a Geom without cloning. The result is a view that
// must not outlive the line.
func (l *Line) AsGeom() *Geom {
	return &Geom{cg: (*C.struct_tg_geom)(unsafe.Pointer(l.cl))}
}

// AsWKT returns the Well-Known Text representation of the line.
func (l *Line) AsWKT() string {
	return l.AsGeom().AsWKT()
}

// String implements fmt.Stringer and returns the line as Well-Known Text.
func (l *Line) String() string {
	return l.AsWKT()
}

// MemSize returns the allocation size of the line.
func (l *Line) MemSize() int {
	return int(C.tg_line_memsize(l.cl))
}

// Rect returns the minimum bounding rectangle of the line.
func (l *Line) Rect() Rect {
	return goRect(C.tg_line_rect(l.cl))
}

// NumPoints returns the number of points in the line.
func (l *Line) NumPoints() int {
	return int(C.tg_line_num_points(l.cl))
}

// PointAt returns the point at index.
func (l *Line) PointAt(index int) (Point, bool) {
	if index < 0 || index >= l.NumPoints() {
		return Point{}, false
	}
	return goPoint(C.tg_line_point_at(l.cl, C.int(index))), true
}

// Points returns a copy of the line's points.
func (l *Line) Points() []Point {
	n := l.NumPoints()
	if n == 0 {
		return nil
	}

	src := unsafe.Slice((*C.struct_tg_point)(C.tg_line_points(l.cl)), n)
	out := make([]Point, n)
	for i := range src {
		out[i] = goPoint(src[i])
	}

	return out
}

// NumSegments returns the number of segments in the line.
func (l *Line) NumSegments() int {
	return int(C.tg_line_num_segments(l.cl))
}

// SegmentAt returns the segment at index.
func (l *Line) SegmentAt(index int) (Segment, bool) {
	if index < 0 || index >= l.NumSegments() {
		return Segment{}, false
	}
	return goSegment(C.tg_line_segment_at(l.cl, C.int(index))), true
}

// Clockwise reports whether the line winds clockwise.
func (l *Line) Clockwise() bool {
	return bool(C.tg_line_clockwise(l.cl))
}

// Length returns the length of the line.
func (l *Line) Length() float64 {
	return float64(C.tg_line_length(l.cl))
}

// IndexSpread returns the number of segments grouped per index node, or zero
// when the line is not indexed.
func (l *Line) IndexSpread() int {
	return int(C.tg_line_index_spread(l.cl))
}

// IndexNumLevels returns the number of index levels.
func (l *Line) IndexNumLevels() int {
	return int(C.tg_line_index_num_levels(l.cl))
}

// IndexLevelNumRects returns the number of rectangles at an index level.
func (l *Line) IndexLevelNumRects(level int) int {
	return int(C.tg_line_index_level_num_rects(l.cl, C.int(level)))
}

// IndexLevelRect returns the rectangle at the given level and position.
func (l *Line) IndexLevelRect(level, rect int) Rect {
	return goRect(C.tg_line_index_level_rect(l.cl, C.int(level), C.int(rect)))
}

func (l *Line) free() {
	C.tg_line_free(l.cl)
}
