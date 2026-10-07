package tgo

/*
#include "tg.h"
#include <stdlib.h>

// The distance helpers below are exported by tg.c but are not declared in
// tg.h. tg_ring_nearest_segment() has no built-in distance calculation, so we
// reuse these to implement the concrete nearest queries.
double tg_point_distance_rect(struct tg_point p, struct tg_rect r);
double tg_point_distance_segment(struct tg_point p, struct tg_segment s);

struct tgo_geom_search_item {
	const struct tg_geom *geom;
	int index;
};

struct tgo_geom_search_ctx {
	struct tgo_geom_search_item *items;
	int len;
	int cap;
	int oom;
};

bool tgo_geom_search_collect(const struct tg_geom *child, int index,
	void *udata)
{
	struct tgo_geom_search_ctx *ctx = udata;
	if (ctx->len == ctx->cap) {
		int ncap = ctx->cap == 0 ? 16 : ctx->cap * 2;
		struct tgo_geom_search_item *items = realloc(ctx->items,
			(size_t)ncap * sizeof(struct tgo_geom_search_item));
		if (!items) {
			ctx->oom = 1;
			return false;
		}
		ctx->items = items;
		ctx->cap = ncap;
	}
	ctx->items[ctx->len].geom = child;
	ctx->items[ctx->len].index = index;
	ctx->len++;
	return true;
}

struct tgo_seg_pair {
	struct tg_segment a;
	int aidx;
	struct tg_segment b;
	int bidx;
};

struct tgo_seg_search_ctx {
	struct tgo_seg_pair *items;
	int len;
	int cap;
	int oom;
};

bool tgo_seg_collect(struct tg_segment aseg, int aidx, struct tg_segment bseg,
	int bidx, void *udata)
{
	struct tgo_seg_search_ctx *ctx = udata;
	if (ctx->len == ctx->cap) {
		int ncap = ctx->cap == 0 ? 16 : ctx->cap * 2;
		struct tgo_seg_pair *items = realloc(ctx->items,
			(size_t)ncap * sizeof(struct tgo_seg_pair));
		if (!items) {
			ctx->oom = 1;
			return false;
		}
		ctx->items = items;
		ctx->cap = ncap;
	}
	ctx->items[ctx->len].a = aseg;
	ctx->items[ctx->len].aidx = aidx;
	ctx->items[ctx->len].b = bseg;
	ctx->items[ctx->len].bidx = bidx;
	ctx->len++;
	return true;
}

struct tgo_nearest_result {
	struct tg_segment seg;
	double dist;
	int index;
};

struct tgo_nearest_ctx {
	struct tg_point point;
	struct tgo_nearest_result *results;
	int k;
	int cap;
	int len;
};

double tgo_nearest_rect_dist(struct tg_rect rect, int *more, void *udata) {
	(void)more;
	struct tgo_nearest_ctx *ctx = udata;
	return tg_point_distance_rect(ctx->point, rect);
}

double tgo_nearest_seg_dist(struct tg_segment seg, int *more, void *udata) {
	(void)more;
	struct tgo_nearest_ctx *ctx = udata;
	return tg_point_distance_segment(ctx->point, seg);
}

bool tgo_nearest_iter(struct tg_segment seg, double dist, int index,
	void *udata)
{
	struct tgo_nearest_ctx *ctx = udata;
	if (ctx->len < ctx->cap) {
		ctx->results[ctx->len].seg = seg;
		ctx->results[ctx->len].dist = dist;
		ctx->results[ctx->len].index = index;
		ctx->len++;
	}
	return ctx->k <= 0 || ctx->len < ctx->k;
}
*/
import "C"
import (
	"iter"
	"unsafe"
)

// SegmentPair is a pair of segments, with their indices, reported by the
// segment search iterators.
type SegmentPair struct {
	A      Segment
	AIndex int
	B      Segment
	BIndex int
}

// NearestSegment is a segment and its distance from a query point.
type NearestSegment struct {
	Segment  Segment
	Distance float64
	Index    int
}

// Search iterates over the child geometries of a Multi*/GeometryCollection
// (or GeoJSON FeatureCollection) that intersect r, from the index. The yielded
// geometries are views into g and must not outlive it. Search yields nothing
// for geometries that are not collections.
func (g *Geom) Search(r Rect) iter.Seq2[int, *Geom] {
	return func(yield func(int, *Geom) bool) {
		if g == nil || g.cg == nil {
			return
		}

		var ctx C.struct_tgo_geom_search_ctx
		C.tg_geom_search(g.cg, cRect(r), (*[0]byte)(C.tgo_geom_search_collect), unsafe.Pointer(&ctx))
		defer C.free(unsafe.Pointer(ctx.items))

		n := int(ctx.len)
		if n == 0 {
			return
		}

		items := unsafe.Slice(ctx.items, n)
		for i := 0; i < n; i++ {
			if !yield(int(items[i].index), &Geom{cg: items[i].geom}) {
				return
			}
		}
	}
}

func yieldSegmentPairs(ctx C.struct_tgo_seg_search_ctx, yield func(SegmentPair) bool) {
	n := int(ctx.len)
	if n == 0 {
		return
	}

	items := unsafe.Slice(ctx.items, n)
	for i := 0; i < n; i++ {
		ok := yield(SegmentPair{
			A:      goSegment(items[i].a),
			AIndex: int(items[i].aidx),
			B:      goSegment(items[i].b),
			BIndex: int(items[i].bidx),
		})
		if !ok {
			return
		}
	}
}

// SearchLine iterates over the intersecting segments of the ring and the line.
func (r *Ring) SearchLine(l *Line) iter.Seq[SegmentPair] {
	return func(yield func(SegmentPair) bool) {
		if r == nil || r.cr == nil || l == nil || l.cl == nil {
			return
		}
		var ctx C.struct_tgo_seg_search_ctx
		C.tg_ring_line_search(r.cr, l.cl, (*[0]byte)(C.tgo_seg_collect), unsafe.Pointer(&ctx))
		defer C.free(unsafe.Pointer(ctx.items))
		yieldSegmentPairs(ctx, yield)
	}
}

// SearchRing iterates over the intersecting segments of the two rings.
func (r *Ring) SearchRing(o *Ring) iter.Seq[SegmentPair] {
	return func(yield func(SegmentPair) bool) {
		if r == nil || r.cr == nil || o == nil || o.cr == nil {
			return
		}
		var ctx C.struct_tgo_seg_search_ctx
		C.tg_ring_ring_search(r.cr, o.cr, (*[0]byte)(C.tgo_seg_collect), unsafe.Pointer(&ctx))
		defer C.free(unsafe.Pointer(ctx.items))
		yieldSegmentPairs(ctx, yield)
	}
}

// SearchLine iterates over the intersecting segments of the two lines.
func (l *Line) SearchLine(o *Line) iter.Seq[SegmentPair] {
	return func(yield func(SegmentPair) bool) {
		if l == nil || l.cl == nil || o == nil || o.cl == nil {
			return
		}
		var ctx C.struct_tgo_seg_search_ctx
		C.tg_line_line_search(l.cl, o.cl, (*[0]byte)(C.tgo_seg_collect), unsafe.Pointer(&ctx))
		defer C.free(unsafe.Pointer(ctx.items))
		yieldSegmentPairs(ctx, yield)
	}
}

// nearestFrom returns up to k segments ordered from nearest to farthest. k <= 0
// returns every segment ordered by distance.
func nearestFrom(cr *C.struct_tg_ring, x, y float64, k int) []NearestSegment {
	n := int(C.tg_ring_num_segments(cr))
	if n == 0 {
		return nil
	}

	results := C.malloc(C.size_t(n) * C.size_t(unsafe.Sizeof(C.struct_tgo_nearest_result{})))
	if results == nil {
		return nil
	}
	defer C.free(results)

	ctx := C.struct_tgo_nearest_ctx{
		point:   cPoint(Point{X: x, Y: y}),
		results: (*C.struct_tgo_nearest_result)(results),
		k:       C.int(k),
		cap:     C.int(n),
	}

	C.tg_ring_nearest_segment(cr,
		(*[0]byte)(C.tgo_nearest_rect_dist),
		(*[0]byte)(C.tgo_nearest_seg_dist),
		(*[0]byte)(C.tgo_nearest_iter),
		unsafe.Pointer(&ctx))

	m := int(ctx.len)
	out := make([]NearestSegment, m)
	arr := unsafe.Slice(ctx.results, m)
	for i := 0; i < m; i++ {
		out[i] = NearestSegment{
			Segment:  goSegment(arr[i].seg),
			Distance: float64(arr[i].dist),
			Index:    int(arr[i].index),
		}
	}

	return out
}

// Nearest returns the segment of the ring nearest to (x, y).
func (r *Ring) Nearest(x, y float64) (NearestSegment, bool) {
	res := nearestFrom(r.cr, x, y, 1)
	if len(res) == 0 {
		return NearestSegment{}, false
	}
	return res[0], true
}

// NearestK returns the k segments of the ring nearest to (x, y), ordered from
// nearest to farthest.
func (r *Ring) NearestK(x, y float64, k int) []NearestSegment {
	return nearestFrom(r.cr, x, y, k)
}

// Nearest returns the segment of the line nearest to (x, y).
func (l *Line) Nearest(x, y float64) (NearestSegment, bool) {
	res := nearestFrom((*C.struct_tg_ring)(unsafe.Pointer(l.cl)), x, y, 1)
	if len(res) == 0 {
		return NearestSegment{}, false
	}
	return res[0], true
}

// NearestK returns the k segments of the line nearest to (x, y), ordered from
// nearest to farthest.
func (l *Line) NearestK(x, y float64, k int) []NearestSegment {
	return nearestFrom((*C.struct_tg_ring)(unsafe.Pointer(l.cl)), x, y, k)
}
