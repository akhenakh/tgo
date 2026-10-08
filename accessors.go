package tgo

/*
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TypeString returns the name of the geometry type.
func (g *Geom) TypeString() string {
	if g.cg == nil {
		return ""
	}
	return C.GoString(C.tg_geom_type_string(C.tg_geom_typeof(g.cg)))
}

// Rect returns the minimum bounding rectangle of the geometry.
func (g *Geom) Rect() Rect {
	return goRect(C.tg_geom_rect(g.cg))
}

// FullRect returns the minimum bounding rectangle on all dimensions and the
// number of dimensions. It returns dims == 0 for a nil geometry.
func (g *Geom) FullRect() (min, max [4]float64, dims int) {
	if g == nil || g.cg == nil {
		return min, max, 0
	}

	// The min/max arrays live in C memory: tg_geom_fullrect writes them, and C
	// must not write to Go memory.
	cmin := cAllocN[C.double](4)
	defer C.free(unsafe.Pointer(cmin))
	cmax := cAllocN[C.double](4)
	defer C.free(unsafe.Pointer(cmax))

	n := C.tg_geom_fullrect(g.cg, cmin, cmax)

	gmin := unsafe.Slice(cmin, 4)
	gmax := unsafe.Slice(cmax, 4)
	for i := range 4 {
		min[i] = float64(gmin[i])
		max[i] = float64(gmax[i])
	}

	return min, max, int(n)
}

// IsEmpty reports whether the geometry is empty.
func (g *Geom) IsEmpty() bool {
	return bool(C.tg_geom_is_empty(g.cg))
}

// IsFeature reports whether the geometry was parsed from a GeoJSON Feature.
func (g *Geom) IsFeature() bool {
	return bool(C.tg_geom_is_feature(g.cg))
}

// IsFeatureCollection reports whether the geometry was parsed from a GeoJSON
// FeatureCollection.
func (g *Geom) IsFeatureCollection() bool {
	return bool(C.tg_geom_is_featurecollection(g.cg))
}

// Err returns the parse error carried by the geometry, or an empty string.
func (g *Geom) Err() string {
	return C.GoString(C.tg_geom_error(g.cg))
}

// Point returns the geometry's point. For a non-point geometry it returns the
// center of the bounding rectangle.
func (g *Geom) Point() Point {
	return goPoint(C.tg_geom_point(g.cg))
}

// NumPoints returns the number of points in a MultiPoint geometry, or zero for
// any other type.
func (g *Geom) NumPoints() int {
	return int(C.tg_geom_num_points(g.cg))
}

// PointAt returns the point at index for a MultiPoint geometry.
func (g *Geom) PointAt(index int) (Point, bool) {
	if index < 0 || index >= g.NumPoints() {
		return Point{}, false
	}
	return goPoint(C.tg_geom_point_at(g.cg, C.int(index))), true
}

// NumLines returns the number of lines in a MultiLineString geometry, or zero
// for any other type.
func (g *Geom) NumLines() int {
	return int(C.tg_geom_num_lines(g.cg))
}

// LineAt returns the line at index for a MultiLineString geometry. The
// returned Line is a view into g and must not outlive it.
func (g *Geom) LineAt(index int) (*Line, bool) {
	if index < 0 || index >= g.NumLines() {
		return nil, false
	}

	cl := C.tg_geom_line_at(g.cg, C.int(index))
	if cl == nil {
		return nil, false
	}

	return &Line{cl: cl}, true
}

// NumPolys returns the number of polygons in a MultiPolygon geometry, or zero
// for any other type.
func (g *Geom) NumPolys() int {
	return int(C.tg_geom_num_polys(g.cg))
}

// PolyAt returns the polygon at index for a MultiPolygon geometry. The
// returned Poly is a view into g and must not outlive it.
func (g *Geom) PolyAt(index int) (*Poly, bool) {
	if index < 0 || index >= g.NumPolys() {
		return nil, false
	}

	cp := C.tg_geom_poly_at(g.cg, C.int(index))
	if cp == nil {
		return nil, false
	}

	return &Poly{cp: cp}, true
}

// NumGeometries returns the number of geometries in a GeometryCollection, or
// zero for any other type.
func (g *Geom) NumGeometries() int {
	return int(C.tg_geom_num_geometries(g.cg))
}

// GeometryAt returns the geometry at index for a GeometryCollection. The
// returned Geom is a view into g and must not outlive it.
func (g *Geom) GeometryAt(index int) (*Geom, bool) {
	if index < 0 || index >= g.NumGeometries() {
		return nil, false
	}

	cg := C.tg_geom_geometry_at(g.cg, C.int(index))
	if cg == nil {
		return nil, false
	}

	return &Geom{cg: cg}, true
}

// Dims returns the number of dimensions: 2, 3 (Z or M), or 4 (Z and M).
func (g *Geom) Dims() int {
	return int(C.tg_geom_dims(g.cg))
}

// HasZ reports whether the geometry has Z coordinates.
func (g *Geom) HasZ() bool {
	return bool(C.tg_geom_has_z(g.cg))
}

// HasM reports whether the geometry has M coordinates.
func (g *Geom) HasM() bool {
	return bool(C.tg_geom_has_m(g.cg))
}

// Z returns the Z coordinate of a Point geometry, or zero for any other type.
func (g *Geom) Z() float64 {
	return float64(C.tg_geom_z(g.cg))
}

// M returns the M coordinate of a Point geometry, or zero for any other type.
func (g *Geom) M() float64 {
	return float64(C.tg_geom_m(g.cg))
}

// NumExtraCoords returns the number of extra (Z/M) coordinates carried by the
// geometry. Points store their Z/M separately and return zero.
func (g *Geom) NumExtraCoords() int {
	return int(C.tg_geom_num_extra_coords(g.cg))
}

// ExtraCoords returns the extra (Z/M) coordinates carried by the geometry.
func (g *Geom) ExtraCoords() []float64 {
	n := g.NumExtraCoords()
	if n == 0 {
		return nil
	}

	ptr := C.tg_geom_extra_coords(g.cg)
	if ptr == nil {
		return nil
	}

	src := unsafe.Slice((*C.double)(ptr), n)
	out := make([]float64, n)
	for i, v := range src {
		out[i] = float64(v)
	}

	return out
}

// IntersectsRect reports whether the geometry intersects the rectangle.
func (g *Geom) IntersectsRect(r Rect) bool {
	return bool(C.tg_geom_intersects_rect(g.cg, cRect(r)))
}

// IntersectsXY reports whether the geometry intersects the point (x, y).
func (g *Geom) IntersectsXY(x, y float64) bool {
	return bool(C.tg_geom_intersects_xy(g.cg, C.double(x), C.double(y)))
}

// Equals reports whether g and o are equal.
func (g *Geom) Equals(o *Geom) bool {
	return bool(C.tg_geom_equals(g.cg, o.cg))
}

// Intersects reports whether g and o intersect.
func (g *Geom) Intersects(o *Geom) bool {
	return bool(C.tg_geom_intersects(g.cg, o.cg))
}

// Disjoint reports whether g and o are disjoint.
func (g *Geom) Disjoint(o *Geom) bool {
	return bool(C.tg_geom_disjoint(g.cg, o.cg))
}

// Contains reports whether g contains o.
func (g *Geom) Contains(o *Geom) bool {
	return bool(C.tg_geom_contains(g.cg, o.cg))
}

// Within reports whether g is within o.
func (g *Geom) Within(o *Geom) bool {
	return bool(C.tg_geom_within(g.cg, o.cg))
}

// Covers reports whether g covers o.
func (g *Geom) Covers(o *Geom) bool {
	return bool(C.tg_geom_covers(g.cg, o.cg))
}

// CoveredBy reports whether g is covered by o.
func (g *Geom) CoveredBy(o *Geom) bool {
	return bool(C.tg_geom_coveredby(g.cg, o.cg))
}

// Touches reports whether g touches o.
func (g *Geom) Touches(o *Geom) bool {
	return bool(C.tg_geom_touches(g.cg, o.cg))
}
