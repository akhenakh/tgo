package tgo

/*
#include "tg.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// The c* helpers copy Go values into C-allocated memory so the buffers stay
// alive for the duration of a C call. Each returns a cleanup function that
// frees the buffer; calling it is always safe, including for empty slices.

func cPoints(pts []Point) (*C.struct_tg_point, func()) {
	if len(pts) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(pts)) * C.size_t(unsafe.Sizeof(C.struct_tg_point{})))
	s := unsafe.Slice((*C.struct_tg_point)(p), len(pts))
	for i, pt := range pts {
		s[i] = cPoint(pt)
	}

	return &s[0], func() { C.free(p) }
}

func cDoubles(vals []float64) (*C.double, func()) {
	if len(vals) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(vals)) * C.size_t(unsafe.Sizeof(C.double(0))))
	s := unsafe.Slice((*C.double)(p), len(vals))
	for i, v := range vals {
		s[i] = C.double(v)
	}

	return &s[0], func() { C.free(p) }
}

func cLines(lines []*C.struct_tg_line) (**C.struct_tg_line, func()) {
	if len(lines) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(lines)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	s := unsafe.Slice((**C.struct_tg_line)(p), len(lines))
	copy(s, lines)

	return &s[0], func() { C.free(p) }
}

func cRings(rings []*C.struct_tg_ring) (**C.struct_tg_ring, func()) {
	if len(rings) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(rings)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	s := unsafe.Slice((**C.struct_tg_ring)(p), len(rings))
	copy(s, rings)

	return &s[0], func() { C.free(p) }
}

func cPolys(polys []*C.struct_tg_poly) (**C.struct_tg_poly, func()) {
	if len(polys) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(polys)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	s := unsafe.Slice((**C.struct_tg_poly)(p), len(polys))
	copy(s, polys)

	return &s[0], func() { C.free(p) }
}

func cGeoms(geoms []*C.struct_tg_geom) (**C.struct_tg_geom, func()) {
	if len(geoms) == 0 {
		return nil, func() {}
	}

	p := C.malloc(C.size_t(len(geoms)) * C.size_t(unsafe.Sizeof(uintptr(0))))
	s := unsafe.Slice((**C.struct_tg_geom)(p), len(geoms))
	copy(s, geoms)

	return &s[0], func() { C.free(p) }
}
