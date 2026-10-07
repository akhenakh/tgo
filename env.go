package tgo

/*
#include "tg.h"
*/
import "C"

// SetIndex sets the default indexing option used when constructing geometries.
// It should be called at program startup, before creating any geometry.
func SetIndex(ix IndexType) {
	C.tg_env_set_index(C.enum_tg_index(ix))
}

// SetIndexSpread sets the number of segments grouped per index node. It should
// be called at program startup, before creating any geometry.
func SetIndexSpread(spread int) {
	C.tg_env_set_index_spread(C.int(spread))
}

// SetPrintFixedFloats controls whether numbers are printed in fixed notation.
// It should be called at program startup.
func SetPrintFixedFloats(print bool) {
	C.tg_env_set_print_fixed_floats(C.bool(print))
}
