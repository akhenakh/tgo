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

func UnmarshalWKT(data string) (*Geom, error) {
	cd := C.CString(data)
	defer C.free(unsafe.Pointer(cd))

	cg := C.tg_parse_wkt(cd)
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func UnmarshalWKB(data []byte) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_wkb((*C.uchar)(cdata), C.size_t(len(data)))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func UnmarshalGeoJSON(data []byte) (*Geom, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	cdata := C.CBytes(data)
	defer C.free(cdata)

	cg := C.tg_parse_geojson((*C.char)(cdata))
	cerr := C.tg_geom_error(cg)
	if cerr != nil {
		return nil, fmt.Errorf("%s", C.GoString(cerr))
	}

	g := &Geom{cg}
	runtime.SetFinalizer(g, (*Geom).free)

	return g, nil
}

func Intersects(g1, g2 *Geom) bool {
	return bool(C.tg_geom_intersects(g1.cg, g2.cg))
}

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

func (g *Geom) free() {
	C.tg_geom_free(g.cg)
}
