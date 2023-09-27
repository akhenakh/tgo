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

type Geom struct {
	cg *C.struct_tg_geom
}

type GeomType uint8

const (
	Point              GeomType = iota + 1 // Point.
	LineString                             // LineString.
	Polygon                                // Polygon.
	MultiPoint                             // MultiPoint, collection of points.
	MultiLineString                        // MultiLineString, collection of linestrings.
	MultiPolygon                           // MultiPolygon, collection of polygons.
	GeometryCollection                     // GeometryCollection, collection of geometries.
)

type IndexType uint32

const (
	None IndexType = iota + 1
	Natural
	YStripes
)

func UnmarshalWKT(data string) (*Geom, error) {
	return UnmarshalWKTAndIndex(data, None)
}

func UnmarshalWKTAndIndex(data string, idxt IndexType) (*Geom, error) {
	cd := C.CString(data)
	defer C.free(unsafe.Pointer(cd))

	cg := C.tg_parse_wkt_ix(cd, C.enum_tg_index(idxt))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func UnmarshalWKB(data []byte) (*Geom, error) {
	return UnmarshalWKBAndIndex(data, None)
}

func UnmarshalWKBAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_wkb_ix((*C.uchar)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func UnmarshalGeoJSON(data []byte) (*Geom, error) {
	return UnmarshalGeoJSONAndIndex(data, None)
}

func UnmarshalGeoJSONAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_geojson_ix((*C.char)(cdata), C.enum_tg_index(idxt))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func Equals(g1, g2 *Geom) bool {
	return bool(C.tg_geom_equals(g1.cg, g2.cg))
}

func Intersects(g1, g2 *Geom) bool {
	return bool(C.tg_geom_intersects(g1.cg, g2.cg))
}

func Disjoint(g1, g2 *Geom) bool {
	return bool(C.tg_geom_disjoint(g1.cg, g2.cg))
}

func Contains(g1, g2 *Geom) bool {
	return bool(C.tg_geom_contains(g1.cg, g2.cg))
}

func Within(g1, g2 *Geom) bool {
	return bool(C.tg_geom_within(g1.cg, g2.cg))
}

func Covers(g1, g2 *Geom) bool {
	return bool(C.tg_geom_covers(g1.cg, g2.cg))
}

func CoveredBy(g1, g2 *Geom) bool {
	return bool(C.tg_geom_coveredby(g1.cg, g2.cg))
}

func Touches(g1, g2 *Geom) bool {
	return bool(C.tg_geom_touches(g1.cg, g2.cg))
}

// MemSize returns the allocation size of the geometry in the C world.
func (g *Geom) MemSize() int {
	return int(C.tg_geom_memsize(g.cg))
}

// Properties returns a string that represents any extra JSON from a parsed GeoJSON
// geometry. Such as the "id" or "properties" fields.
func (g *Geom) Properties() string {
	return C.GoString(C.tg_geom_extra_json(g.cg))
}

// func (g *Geom) RingSearch() {
// 	r := C.tg_ring_new(points, len(coords)/2)
// 	C.tg_ring_ring_search(g.cg,r, bool(*iter)(struct tg_segment aseg, int aidx, struct tg_segment bseg, int bidx, void *udata), void *udata);
// 	C.tg_ring_free(r)
// }

func (g *Geom) AsPoly() (*Poly, bool) {
	cp := C.tg_geom_poly(g.cg)

	if cp == nil {
		return nil, false
	}

	p := &Poly{cp: cp}
	return p, true
}

func (g *Geom) AsMultiPoly() (*MultiPoly, bool) {
	if g.Type() != MultiPolygon {
		return nil, false
	}
	count := C.tg_geom_num_polys(g.cg)
	if count < 1 {
		return nil, false
	}

	mp := &MultiPoly{cg: g.cg}
	return mp, true
}

// AsText returns geometry as WKT
func (g *Geom) AsText() string {
	if g.cg == nil {
		return ""
	}

	csz := C.tg_geom_wkt(g.cg, nil, C.size_t(0))

	cwkt := C.malloc(csz + 1)
	C.tg_geom_wkt(g.cg, (*C.char)(cwkt), csz+1)
	defer C.free(cwkt)
	return C.GoString((*C.char)(cwkt))
}

func (g *Geom) Type() GeomType {
	switch C.tg_geom_typeof(g.cg) {
	case C.TG_POINT:
		return Point
	case C.TG_LINESTRING:
		return LineString
	case C.TG_POLYGON:
		return Polygon
	case C.TG_MULTIPOINT:
		return MultiPoint
	case C.TG_MULTILINESTRING:
		return MultiLineString
	case C.TG_MULTIPOLYGON:
		return MultiPolygon
	case C.TG_GEOMETRYCOLLECTION:
		return GeometryCollection
	}

	return 0
}

func (g *Geom) free() {
	C.tg_geom_free(g.cg)
}
