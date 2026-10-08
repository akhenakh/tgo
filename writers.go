package tgo

/*
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// cStringOut calls fn twice: first to size the buffer, then to fill it, and
// returns the result as a Go string. fn must write a NUL-terminated string.
func cStringOut(fn func(dst *C.char, n C.size_t) C.size_t) string {
	n := fn(nil, 0)
	if n == 0 {
		return ""
	}

	buf := C.malloc(n + 1)
	defer C.free(buf)

	fn((*C.char)(buf), n+1)

	return C.GoString((*C.char)(buf))
}

// cBytesOut calls fn twice: first to size the buffer, then to fill it, and
// returns the result as a Go byte slice.
func cBytesOut(fn func(dst *C.uchar, n C.size_t) C.size_t) []byte {
	n := fn(nil, 0)
	if n == 0 {
		return nil
	}

	buf := C.malloc(n)
	defer C.free(buf)

	fn((*C.uchar)(buf), n)

	return C.GoBytes(buf, C.int(n))
}

// AsWKT returns the geometry as Well-Known Text.
func (g *Geom) AsWKT() string {
	if g.cg == nil {
		return ""
	}
	return cStringOut(func(dst *C.char, n C.size_t) C.size_t {
		return C.tg_geom_wkt(g.cg, dst, n)
	})
}

// AsGeoJSON returns the geometry as GeoJSON.
func (g *Geom) AsGeoJSON() string {
	if g.cg == nil {
		return ""
	}
	return cStringOut(func(dst *C.char, n C.size_t) C.size_t {
		return C.tg_geom_geojson(g.cg, dst, n)
	})
}

// AsHex returns the geometry as hex-encoded WKB.
func (g *Geom) AsHex() string {
	if g.cg == nil {
		return ""
	}
	return cStringOut(func(dst *C.char, n C.size_t) C.size_t {
		return C.tg_geom_hex(g.cg, dst, n)
	})
}

// AsWKB returns the geometry as Well-Known Binary.
func (g *Geom) AsWKB() []byte {
	if g.cg == nil {
		return nil
	}
	return cBytesOut(func(dst *C.uchar, n C.size_t) C.size_t {
		return C.tg_geom_wkb(g.cg, (*C.uint8_t)(dst), n)
	})
}

// AsGeoBin returns the geometry in the compact geobin format.
func (g *Geom) AsGeoBin() []byte {
	if g.cg == nil {
		return nil
	}
	return cBytesOut(func(dst *C.uchar, n C.size_t) C.size_t {
		return C.tg_geom_geobin(g.cg, (*C.uint8_t)(dst), n)
	})
}

// String implements fmt.Stringer and returns the geometry as Well-Known Text.
func (g *Geom) String() string {
	return g.AsWKT()
}

// GeoBinRect returns the minimum bounding rectangle stored in a geobin buffer.
func GeoBinRect(data []byte) Rect {
	if len(data) == 0 {
		return Rect{}
	}
	return goRect(C.tg_geobin_rect((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data))))
}

// GeoBinPoint returns the point stored in a geobin buffer.
func GeoBinPoint(data []byte) Point {
	if len(data) == 0 {
		return Point{}
	}
	return goPoint(C.tg_geobin_point((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data))))
}

// GeoBinFullRect returns the minimum bounding rectangle of a geobin buffer on
// all dimensions, along with the number of dimensions. It returns dims == 0
// for an invalid buffer.
func GeoBinFullRect(data []byte) (min, max [4]float64, dims int) {
	if len(data) == 0 {
		return min, max, 0
	}

	// The min/max arrays live in C memory: tg_geobin_fullrect writes them, and
	// C must not write to Go memory.
	cmin := cAllocN[C.double](4)
	defer C.free(unsafe.Pointer(cmin))
	cmax := cAllocN[C.double](4)
	defer C.free(unsafe.Pointer(cmax))

	n := C.tg_geobin_fullrect((*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data)), cmin, cmax)

	gmin := unsafe.Slice(cmin, 4)
	gmax := unsafe.Slice(cmax, 4)
	for i := range 4 {
		min[i] = float64(gmin[i])
		max[i] = float64(gmax[i])
	}

	return min, max, int(n)
}
