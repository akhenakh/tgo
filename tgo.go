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
		return false; // stop: keep the first match
	}
	return true;
}

*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// Geom is a geometry. It owns its C data unless it is a view produced by an
// accessor, in which case it must not outlive the source geometry.
type Geom struct {
	cg *C.struct_tg_geom
}

// GeomType identifies the underlying type of a Geom.
type GeomType uint8

const (
	TypePoint              GeomType = iota + 1 // Point.
	TypeLineString                             // LineString.
	TypePolygon                                // Polygon.
	TypeMultiPoint                             // MultiPoint, collection of points.
	TypeMultiLineString                        // MultiLineString, collection of linestrings.
	TypeMultiPolygon                           // MultiPolygon, collection of polygons.
	TypeGeometryCollection                     // GeometryCollection, collection of geometries.
)

// IndexType is the indexing strategy used for rings, lines and polygons.
type IndexType uint32

// Indexing options. IndexDefault lets the library pick (IndexNatural unless
// changed via SetIndex).
const (
	IndexDefault  IndexType = iota // Library default.
	IndexNone                      // No indexing available, or disabled.
	IndexNatural                   // Indexing with natural ring order, for rings/lines.
	IndexYStripes                  // Indexing using segment striping, rings only.
)

// wrapGeom takes ownership of cg and wraps it as a *Geom with a finalizer.
// If cg represents a parse error it is freed and the error is returned.
func wrapGeom(cg *C.struct_tg_geom) (*Geom, error) {
	if cerr := C.tg_geom_error(cg); cerr != nil {
		C.tg_geom_free(cg)
		return nil, errors.New(C.GoString(cerr))
	}

	g := &Geom{cg: cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

// UnmarshalWKT parses a geometry from a WKT representation.
// Using the Natural indexation.
func UnmarshalWKT(data string) (*Geom, error) {
	return UnmarshalWKTAndIndex(data, IndexNatural)
}

// UnmarshalWKTAndIndex parses a geometry from a WKT representation,
// and sets the indexing type.
func UnmarshalWKTAndIndex(data string, idxt IndexType) (*Geom, error) {
	if data == "" {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes([]byte(data))
	defer C.free(cdata)

	cg := C.tg_parse_wktn_ix((*C.char)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// UnmarshalWKB parses a geometry from a WKB representation.
// Using the Natural indexation.
func UnmarshalWKB(data []byte) (*Geom, error) {
	return UnmarshalWKBAndIndex(data, IndexNatural)
}

// UnmarshalWKBAndIndex parses a geometry from a WKB representation,
// and sets the indexing type.
func UnmarshalWKBAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_wkb_ix((*C.uchar)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// UnmarshalGeoJSON parses a geometry from a GeoJSON representation.
// Using the Natural indexation.
func UnmarshalGeoJSON(data []byte) (*Geom, error) {
	return UnmarshalGeoJSONAndIndex(data, IndexNatural)
}

// UnmarshalGeoJSONAndIndex parses a geometry from a GeoJSON representation,
// and sets the indexing type.
func UnmarshalGeoJSONAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_geojsonn_ix((*C.char)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// UnmarshalHex parses a geometry from a hex-encoded WKB representation.
// Using the Natural indexation.
func UnmarshalHex(data string) (*Geom, error) {
	return UnmarshalHexAndIndex(data, IndexNatural)
}

// UnmarshalHexAndIndex parses a geometry from a hex-encoded WKB
// representation, and sets the indexing type.
func UnmarshalHexAndIndex(data string, idxt IndexType) (*Geom, error) {
	if data == "" {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes([]byte(data))
	defer C.free(cdata)

	cg := C.tg_parse_hexn_ix((*C.char)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// UnmarshalGeoBin parses a geometry from the geobin representation.
// Using the Natural indexation.
func UnmarshalGeoBin(data []byte) (*Geom, error) {
	return UnmarshalGeoBinAndIndex(data, IndexNatural)
}

// UnmarshalGeoBinAndIndex parses a geometry from the geobin representation,
// and sets the indexing type.
func UnmarshalGeoBinAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_geobin_ix((*C.uchar)(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// Parse parses data into a geometry by auto-detecting the input type. The
// input data can be WKB, WKT, Hex, or GeoJSON.
func Parse(data []byte) (*Geom, error) {
	return ParseAndIndex(data, IndexNatural)
}

// ParseAndIndex parses data into a geometry by auto-detecting the input type.
// The input data can be WKB, WKT, Hex, or GeoJSON.
func ParseAndIndex(data []byte, idxt IndexType) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_ix(unsafe.Pointer(cdata), C.size_t(len(data)), C.enum_tg_index(idxt))
	return wrapGeom(cg)
}

// Equals returns true if the two geometries are equal
func Equals(g1, g2 *Geom) bool {
	return bool(C.tg_geom_equals(g1.cg, g2.cg))
}

// Intersects reports whether two geometries intersect.
func Intersects(g1, g2 *Geom) bool {
	return bool(C.tg_geom_intersects(g1.cg, g2.cg))
}

// Disjoint reports whether two geometries are disjoint (do not intersect).
func Disjoint(g1, g2 *Geom) bool {
	return bool(C.tg_geom_disjoint(g1.cg, g2.cg))
}

// Contains reports whether g1 contains g2.
func Contains(g1, g2 *Geom) bool {
	return bool(C.tg_geom_contains(g1.cg, g2.cg))
}

// Within reports whether the geometry g1 is within the geometry g2.
func Within(g1, g2 *Geom) bool {
	return bool(C.tg_geom_within(g1.cg, g2.cg))
}

// Covers reports whether the first geometry covers the second geometry.
// A geometry g1 covers g2 if no point of g2 lies outside g1.
func Covers(g1, g2 *Geom) bool {
	return bool(C.tg_geom_covers(g1.cg, g2.cg))
}

// CoveredBy reports whether the geometry g1 is covered by geometry g2.
// A geometry g1 is covered by g2 if every point of g1 is also a point of g2.
func CoveredBy(g1, g2 *Geom) bool {
	return bool(C.tg_geom_coveredby(g1.cg, g2.cg))
}

// Touches reports whether two geometries touch each other.
// Two geometries are considered to touch if they share at least one point in their boundaries,
// but do not intersect in their interiors. The function returns true if the geometries touch, false otherwise.
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

	// creating a context for the iterator. It lives in C memory because
	// pip_iter_one writes ctx.geom from C, and C must not write to Go memory.
	ctx := cAlloc[C.struct_pip_iter_one_ctx]()
	defer C.free(unsafe.Pointer(ctx))
	ctx.pip_point = p

	// calling the C func tg_geom_search
	// void tg_geom_search(const struct tg_geom *geom, struct tg_rect rect,
	//	bool (*iter)(const struct tg_geom *geom, int index, void *udata),
	//	void *udata);
	C.tg_geom_search(g.cg, C.tg_point_rect(p), (*[0]byte)(C.pip_iter_one), (unsafe.Pointer(ctx)))
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

// AsLine returns a Line of the geometry, returns false if not applicable. The
// result is a view that must not outlive the geometry.
func (g *Geom) AsLine() (*Line, bool) {
	cl := C.tg_geom_line(g.cg)

	if cl == nil {
		return nil, false
	}

	return &Line{cl: cl}, true
}

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
	if g.Type() != TypeMultiPolygon {
		return nil, false
	}
	count := C.tg_geom_num_polys(g.cg)
	if count < 1 {
		return nil, false
	}

	mp := &MultiPoly{cg: g.cg}
	return mp, true
}

// Type returns the geometry type.
func (g *Geom) Type() GeomType {
	switch C.tg_geom_typeof(g.cg) {
	case C.TG_POINT:
		return TypePoint
	case C.TG_LINESTRING:
		return TypeLineString
	case C.TG_POLYGON:
		return TypePolygon
	case C.TG_MULTIPOINT:
		return TypeMultiPoint
	case C.TG_MULTILINESTRING:
		return TypeMultiLineString
	case C.TG_MULTIPOLYGON:
		return TypeMultiPolygon
	case C.TG_GEOMETRYCOLLECTION:
		return TypeGeometryCollection
	}

	return 0
}

func (g *Geom) free() {
	C.tg_geom_free(g.cg)
}
