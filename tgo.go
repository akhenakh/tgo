package tgo

/*
#cgo LDFLAGS: -lm
#include "tg.h"
#include <stdlib.h>
#include <stdio.h>

#define MAX_RESPONSE_PER_PIP 8

struct pip_iter_properties_ctx {
	struct tg_point pip_point;
	char *properties[MAX_RESPONSE_PER_PIP];
	uint8_t count;
};

struct pip_iter_one_ctx {
	struct tg_point pip_point;
	struct tg_geom *geom;
};

bool pip_iter_properties(const struct tg_geom *child, int index, void *udata) {
	struct pip_iter_properties_ctx *ctx = udata;
	if (tg_geom_intersects_xy(child, ctx->pip_point.x, ctx->pip_point.y)) {
		ctx->properties[ctx->count] = (char*)tg_geom_extra_json(child);

		//printf("%d %s\n", index, ctx->properties[index]);
		ctx->count++;
		if (ctx->count >= MAX_RESPONSE_PER_PIP) {
			return false;
		}

	}
	return true;
}

bool pip_iter_one(const struct tg_geom *child, int index, void *udata) {
	struct pip_iter_one_ctx *ctx = udata;
	if (tg_geom_intersects_xy(child, ctx->pip_point.x, ctx->pip_point.y)) {
		ctx->geom = (struct tg_geom *)child;
		return true;
	}
	return true;
}

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
	None     IndexType = iota + 1 // no indexing available, or disabled
	Natural                       // indexing with natural ring order, for rings/lines
	YStripes                      // indexing using segment striping, rings only
)

// UnmarshalWKT parses geometries from a WKT representation.
// Using the Natural indexation.
func UnmarshalWKT(data string) (*Geom, error) {
	return UnmarshalWKTAndIndex(data, Natural)
}

// UnmarshalWKTAndIndex parses geometries from a WKT representation,
// and sets the indexation type.
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

// UnmarshalWKB parses geometries from a WKB representation.
// Using the Natural indexation.
func UnmarshalWKB(data []byte) (*Geom, error) {
	return UnmarshalWKBAndIndex(data, Natural)
}

// UnmarshalWKBAndIndex parses geometries from a WKB representation,
// and sets the indexation type.
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

// UnmarshalGeoJSON parses geometries from a GeoJSON representation.
// Using the Natural indexation.
func UnmarshalGeoJSON(data []byte) (*Geom, error) {
	return UnmarshalGeoJSONAndIndex(data, Natural)
}

// UnmarshalGeoJSONAndIndex parses geometries from a GeoJSON representation,
// and sets the indexation type.
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

// Parse data into a geometry by auto detecting the input type. The input data can be WKB, WKT, Hex, or GeoJSON.
func Parse(data []byte) (*Geom, error) {
	return ParseAndIndex(data, Natural)
}

// Parse data into a geometry by auto detecting the input type. The input data can be WKB, WKT, Hex, or GeoJSON.
func ParseAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_ix(unsafe.Pointer(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil

}

// Equals returns true if the two geometries are equal
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

// StabOne performs a Point in Polygon query using the point (x,y)
// and returns the first encountered child geometry
func (g *Geom) StabOne(x, y float64) *Geom {
	// creating a point
	p := C.struct_tg_point{C.double(x), C.double(y)}

	// creating a context for the iterator
	ctx := C.struct_pip_iter_one_ctx{pip_point: p}

	// calling the C func tg_geom_search
	// void tg_geom_search(const struct tg_geom *geom, struct tg_rect rect,
	//	bool (*iter)(const struct tg_geom *geom, int index, void *udata),
	//	void *udata);
	C.tg_geom_search(g.cg, C.tg_point_rect(p), (*[0]byte)(C.pip_iter_one), (unsafe.Pointer(&ctx)))
	if ctx.geom != nil {
		return &Geom{cg: ctx.geom}
	}

	return nil
}

// func (g *Geom) RingSearch() {
// 	r := C.tg_ring_new(points, len(coords)/2)
// 	C.tg_ring_ring_search(g.cg,r, bool(*iter)(struct tg_segment aseg, int aidx, struct tg_segment bseg, int bidx, void *udata), void *udata);
// 	C.tg_ring_free(r)
// }

// AsPoly returns a Poly of the geometry, returns false if not applicable.
func (g *Geom) AsPoly() (*Poly, bool) {
	cp := C.tg_geom_poly(g.cg)

	if cp == nil {
		return nil, false
	}

	p := &Poly{cp: cp}
	return p, true
}

// AsMultiPoly returns a MultiPoly of the geometry, returns false if not applicable.
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

// Type returns the geometry type.
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
