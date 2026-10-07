package tgo

/*
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// extra coordinate kinds for the shared constructor helpers.
const (
	extraNone = iota
	extraZ
	extraM
	extraZM
)

// checkExtraCoords validates the length of an extra-coordinate slice against
// the number of positions in the geometry.
func checkExtraCoords(coords []float64, npoints, kind int) error {
	var want int
	switch kind {
	case extraZ, extraM:
		want = npoints
	case extraZM:
		want = 2 * npoints
	default:
		return nil
	}
	if len(coords) != want {
		return fmt.Errorf("tgo: expected %d extra coordinates, got %d", want, len(coords))
	}
	return nil
}

func polygonNumPoints(p Polygon) int {
	n := len(p.Exterior)
	for _, h := range p.Holes {
		n += len(h)
	}
	return n
}

func polygonsNumPoints(polys []Polygon) int {
	n := 0
	for _, p := range polys {
		n += polygonNumPoints(p)
	}
	return n
}

func multiLineNumPoints(lines []LineString) int {
	n := 0
	for _, l := range lines {
		n += len(l)
	}
	return n
}

// buildLine creates a C line from points. The returned line is owned by the
// caller and must be freed with tg_line_free.
func buildLine(points []Point, ix IndexType) *C.struct_tg_line {
	cpts, free := cPoints(points)
	defer free()
	return C.tg_line_new_ix(cpts, C.int(len(points)), C.enum_tg_index(ix))
}

// buildRing creates a C ring from points. The returned ring is owned by the
// caller and must be freed with tg_ring_free.
func buildRing(points []Point, ix IndexType) *C.struct_tg_ring {
	cpts, free := cPoints(points)
	defer free()
	return C.tg_ring_new_ix(cpts, C.int(len(points)), C.enum_tg_index(ix))
}

// buildPoly creates a C polygon by cloning its rings. The returned poly is
// owned by the caller and must be freed with tg_poly_free.
func buildPoly(p Polygon, ix IndexType) *C.struct_tg_poly {
	ext := buildRing(p.Exterior, ix)
	defer C.tg_ring_free(ext)

	holes := make([]*C.struct_tg_ring, len(p.Holes))
	for i, h := range p.Holes {
		holes[i] = buildRing(h, ix)
	}
	defer func() {
		for _, h := range holes {
			C.tg_ring_free(h)
		}
	}()

	choles, freeHoles := cRings(holes)
	defer freeHoles()

	return C.tg_poly_new(ext, choles, C.int(len(holes)))
}

// Clone returns a new geometry that shares the same underlying data using
// reference counting. It is cheap. It returns nil on allocation failure.
func (g *Geom) Clone() *Geom {
	if g == nil || g.cg == nil {
		return nil
	}

	cg := C.tg_geom_clone(g.cg)
	if cg == nil {
		return nil
	}

	ng := &Geom{cg: cg}
	runtime.SetFinalizer(ng, (*Geom).free)

	return ng
}

// Copy returns a deep copy of the geometry. It returns nil on allocation
// failure.
func (g *Geom) Copy() *Geom {
	if g == nil || g.cg == nil {
		return nil
	}

	cg := C.tg_geom_copy(g.cg)
	if cg == nil {
		return nil
	}

	ng := &Geom{cg: cg}
	runtime.SetFinalizer(ng, (*Geom).free)

	return ng
}

// NewPoint creates a Point geometry.
func NewPoint(p Point) (*Geom, error) {
	return wrapGeom(C.tg_geom_new_point(cPoint(p)))
}

// NewPointZ creates a Point geometry with a Z coordinate.
func NewPointZ(p Point, z float64) (*Geom, error) {
	return wrapGeom(C.tg_geom_new_point_z(cPoint(p), C.double(z)))
}

// NewPointM creates a Point geometry with an M coordinate.
func NewPointM(p Point, m float64) (*Geom, error) {
	return wrapGeom(C.tg_geom_new_point_m(cPoint(p), C.double(m)))
}

// NewPointZM creates a Point geometry with Z and M coordinates.
func NewPointZM(p Point, z, m float64) (*Geom, error) {
	return wrapGeom(C.tg_geom_new_point_zm(cPoint(p), C.double(z), C.double(m)))
}

// NewPointEmpty creates an empty Point geometry.
func NewPointEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_point_empty())
}

// NewLineString creates a LineString geometry.
func NewLineString(points []Point) (*Geom, error) {
	return newLineString(points, IndexDefault, nil, extraNone)
}

// NewLineStringAndIndex creates a LineString geometry using the provided
// indexing option.
func NewLineStringAndIndex(points []Point, ix IndexType) (*Geom, error) {
	return newLineString(points, ix, nil, extraNone)
}

// NewLineStringZ creates a LineString geometry with Z coordinates. zs holds
// one value per point.
func NewLineStringZ(points []Point, zs []float64) (*Geom, error) {
	return newLineString(points, IndexDefault, zs, extraZ)
}

// NewLineStringM creates a LineString geometry with M coordinates. ms holds
// one value per point.
func NewLineStringM(points []Point, ms []float64) (*Geom, error) {
	return newLineString(points, IndexDefault, ms, extraM)
}

// NewLineStringZM creates a LineString geometry with Z and M coordinates.
// coords holds interleaved Z and M values, two per point.
func NewLineStringZM(points []Point, coords []float64) (*Geom, error) {
	return newLineString(points, IndexDefault, coords, extraZM)
}

// NewLineStringEmpty creates an empty LineString geometry.
func NewLineStringEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_linestring_empty())
}

func newLineString(points []Point, ix IndexType, coords []float64, kind int) (*Geom, error) {
	if err := checkExtraCoords(coords, len(points), kind); err != nil {
		return nil, err
	}

	line := buildLine(points, ix)
	defer C.tg_line_free(line)

	cc, free := cDoubles(coords)
	defer free()

	switch kind {
	case extraZ:
		return wrapGeom(C.tg_geom_new_linestring_z(line, cc, C.int(len(coords))))
	case extraM:
		return wrapGeom(C.tg_geom_new_linestring_m(line, cc, C.int(len(coords))))
	case extraZM:
		return wrapGeom(C.tg_geom_new_linestring_zm(line, cc, C.int(len(coords))))
	default:
		return wrapGeom(C.tg_geom_new_linestring(line))
	}
}

// NewPolygon creates a Polygon geometry from an exterior ring and optional
// holes.
func NewPolygon(p Polygon) (*Geom, error) {
	return newPolygon(p, IndexDefault, nil, extraNone)
}

// NewPolygonAndIndex creates a Polygon geometry using the provided indexing
// option.
func NewPolygonAndIndex(p Polygon, ix IndexType) (*Geom, error) {
	return newPolygon(p, ix, nil, extraNone)
}

// NewPolygonZ creates a Polygon geometry with Z coordinates. zs holds one
// value per position across the exterior and all holes.
func NewPolygonZ(p Polygon, zs []float64) (*Geom, error) {
	return newPolygon(p, IndexDefault, zs, extraZ)
}

// NewPolygonM creates a Polygon geometry with M coordinates. ms holds one
// value per position across the exterior and all holes.
func NewPolygonM(p Polygon, ms []float64) (*Geom, error) {
	return newPolygon(p, IndexDefault, ms, extraM)
}

// NewPolygonZM creates a Polygon geometry with Z and M coordinates. coords
// holds interleaved Z and M values, two per position.
func NewPolygonZM(p Polygon, coords []float64) (*Geom, error) {
	return newPolygon(p, IndexDefault, coords, extraZM)
}

// NewPolygonEmpty creates an empty Polygon geometry.
func NewPolygonEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_polygon_empty())
}

func newPolygon(p Polygon, ix IndexType, coords []float64, kind int) (*Geom, error) {
	if err := checkExtraCoords(coords, polygonNumPoints(p), kind); err != nil {
		return nil, err
	}

	poly := buildPoly(p, ix)
	defer C.tg_poly_free(poly)

	cc, free := cDoubles(coords)
	defer free()

	switch kind {
	case extraZ:
		return wrapGeom(C.tg_geom_new_polygon_z(poly, cc, C.int(len(coords))))
	case extraM:
		return wrapGeom(C.tg_geom_new_polygon_m(poly, cc, C.int(len(coords))))
	case extraZM:
		return wrapGeom(C.tg_geom_new_polygon_zm(poly, cc, C.int(len(coords))))
	default:
		return wrapGeom(C.tg_geom_new_polygon(poly))
	}
}

// NewMultiPoint creates a MultiPoint geometry.
func NewMultiPoint(points []Point) (*Geom, error) {
	return newMultiPoint(points, nil, extraNone)
}

// NewMultiPointZ creates a MultiPoint geometry with Z coordinates. zs holds
// one value per point.
func NewMultiPointZ(points []Point, zs []float64) (*Geom, error) {
	return newMultiPoint(points, zs, extraZ)
}

// NewMultiPointM creates a MultiPoint geometry with M coordinates. ms holds
// one value per point.
func NewMultiPointM(points []Point, ms []float64) (*Geom, error) {
	return newMultiPoint(points, ms, extraM)
}

// NewMultiPointZM creates a MultiPoint geometry with Z and M coordinates.
// coords holds interleaved Z and M values, two per point.
func NewMultiPointZM(points []Point, coords []float64) (*Geom, error) {
	return newMultiPoint(points, coords, extraZM)
}

// NewMultiPointEmpty creates an empty MultiPoint geometry.
func NewMultiPointEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_multipoint_empty())
}

func newMultiPoint(points []Point, coords []float64, kind int) (*Geom, error) {
	if err := checkExtraCoords(coords, len(points), kind); err != nil {
		return nil, err
	}

	cpts, freePts := cPoints(points)
	defer freePts()

	cc, freeC := cDoubles(coords)
	defer freeC()

	n := C.int(len(points))

	switch kind {
	case extraZ:
		return wrapGeom(C.tg_geom_new_multipoint_z(cpts, n, cc, C.int(len(coords))))
	case extraM:
		return wrapGeom(C.tg_geom_new_multipoint_m(cpts, n, cc, C.int(len(coords))))
	case extraZM:
		return wrapGeom(C.tg_geom_new_multipoint_zm(cpts, n, cc, C.int(len(coords))))
	default:
		return wrapGeom(C.tg_geom_new_multipoint(cpts, n))
	}
}

// NewMultiLineString creates a MultiLineString geometry.
func NewMultiLineString(lines []LineString) (*Geom, error) {
	return newMultiLineString(lines, nil, extraNone)
}

// NewMultiLineStringZ creates a MultiLineString geometry with Z coordinates.
// zs holds one value per position across all lines.
func NewMultiLineStringZ(lines []LineString, zs []float64) (*Geom, error) {
	return newMultiLineString(lines, zs, extraZ)
}

// NewMultiLineStringM creates a MultiLineString geometry with M coordinates.
// ms holds one value per position across all lines.
func NewMultiLineStringM(lines []LineString, ms []float64) (*Geom, error) {
	return newMultiLineString(lines, ms, extraM)
}

// NewMultiLineStringZM creates a MultiLineString geometry with Z and M
// coordinates. coords holds interleaved Z and M values, two per position.
func NewMultiLineStringZM(lines []LineString, coords []float64) (*Geom, error) {
	return newMultiLineString(lines, coords, extraZM)
}

// NewMultiLineStringEmpty creates an empty MultiLineString geometry.
func NewMultiLineStringEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_multilinestring_empty())
}

func newMultiLineString(lines []LineString, coords []float64, kind int) (*Geom, error) {
	if err := checkExtraCoords(coords, multiLineNumPoints(lines), kind); err != nil {
		return nil, err
	}

	clines := make([]*C.struct_tg_line, len(lines))
	for i, l := range lines {
		clines[i] = buildLine(l, IndexDefault)
	}
	defer func() {
		for _, l := range clines {
			C.tg_line_free(l)
		}
	}()

	cp, free := cLines(clines)
	defer free()

	cc, freeC := cDoubles(coords)
	defer freeC()

	n := C.int(len(lines))

	switch kind {
	case extraZ:
		return wrapGeom(C.tg_geom_new_multilinestring_z(cp, n, cc, C.int(len(coords))))
	case extraM:
		return wrapGeom(C.tg_geom_new_multilinestring_m(cp, n, cc, C.int(len(coords))))
	case extraZM:
		return wrapGeom(C.tg_geom_new_multilinestring_zm(cp, n, cc, C.int(len(coords))))
	default:
		return wrapGeom(C.tg_geom_new_multilinestring(cp, n))
	}
}

// NewMultiPolygon creates a MultiPolygon geometry.
func NewMultiPolygon(polys []Polygon) (*Geom, error) {
	return newMultiPolygon(polys, nil, extraNone)
}

// NewMultiPolygonZ creates a MultiPolygon geometry with Z coordinates. zs
// holds one value per position across all polygons.
func NewMultiPolygonZ(polys []Polygon, zs []float64) (*Geom, error) {
	return newMultiPolygon(polys, zs, extraZ)
}

// NewMultiPolygonM creates a MultiPolygon geometry with M coordinates. ms
// holds one value per position across all polygons.
func NewMultiPolygonM(polys []Polygon, ms []float64) (*Geom, error) {
	return newMultiPolygon(polys, ms, extraM)
}

// NewMultiPolygonZM creates a MultiPolygon geometry with Z and M coordinates.
// coords holds interleaved Z and M values, two per position.
func NewMultiPolygonZM(polys []Polygon, coords []float64) (*Geom, error) {
	return newMultiPolygon(polys, coords, extraZM)
}

// NewMultiPolygonEmpty creates an empty MultiPolygon geometry.
func NewMultiPolygonEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_multipolygon_empty())
}

func newMultiPolygon(polys []Polygon, coords []float64, kind int) (*Geom, error) {
	if err := checkExtraCoords(coords, polygonsNumPoints(polys), kind); err != nil {
		return nil, err
	}

	cpolys := make([]*C.struct_tg_poly, len(polys))
	for i, p := range polys {
		cpolys[i] = buildPoly(p, IndexDefault)
	}
	defer func() {
		for _, p := range cpolys {
			C.tg_poly_free(p)
		}
	}()

	cp, free := cPolys(cpolys)
	defer free()

	cc, freeC := cDoubles(coords)
	defer freeC()

	n := C.int(len(polys))

	switch kind {
	case extraZ:
		return wrapGeom(C.tg_geom_new_multipolygon_z(cp, n, cc, C.int(len(coords))))
	case extraM:
		return wrapGeom(C.tg_geom_new_multipolygon_m(cp, n, cc, C.int(len(coords))))
	case extraZM:
		return wrapGeom(C.tg_geom_new_multipolygon_zm(cp, n, cc, C.int(len(coords))))
	default:
		return wrapGeom(C.tg_geom_new_multipolygon(cp, n))
	}
}

// NewGeometryCollection creates a GeometryCollection from the given
// geometries. The inputs are cloned; the caller keeps ownership.
func NewGeometryCollection(geoms ...*Geom) (*Geom, error) {
	if len(geoms) == 0 {
		return NewGeometryCollectionEmpty()
	}

	cgeoms := make([]*C.struct_tg_geom, len(geoms))
	for i, g := range geoms {
		if g == nil || g.cg == nil {
			return nil, errors.New("tgo: nil geometry")
		}
		cgeoms[i] = g.cg
	}

	cp, free := cGeoms(cgeoms)
	defer free()

	return wrapGeom(C.tg_geom_new_geometrycollection(cp, C.int(len(cgeoms))))
}

// NewGeometryCollectionEmpty creates an empty GeometryCollection.
func NewGeometryCollectionEmpty() (*Geom, error) {
	return wrapGeom(C.tg_geom_new_geometrycollection_empty())
}

// NewError creates an error geometry carrying msg. It returns nil on
// allocation failure.
func NewError(msg string) *Geom {
	cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cmsg))

	cg := C.tg_geom_new_error(cmsg)
	if cg == nil {
		return nil
	}

	g := &Geom{cg: cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g
}
