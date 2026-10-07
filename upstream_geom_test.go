package tgo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFullRectUpstream reuses the table from upstream tests/test_rect.c
// (test_rect_fullrect).
func TestFullRectUpstream(t *testing.T) {
	tests := []struct {
		wkt      string
		dims     int
		min, max [4]float64
	}{
		{"POINT(10 20)", 2, [4]float64{10, 20}, [4]float64{10, 20}},
		{"POINT(10 20 30)", 3, [4]float64{10, 20, 30}, [4]float64{10, 20, 30}},
		{"POINT(10 20 30 40)", 4, [4]float64{10, 20, 30, 40}, [4]float64{10, 20, 30, 40}},
		{"LINESTRING(10 20,30 40)", 2, [4]float64{10, 20}, [4]float64{30, 40}},
		{"LINESTRING(10 20 30,40 50 60)", 3, [4]float64{10, 20, 30}, [4]float64{40, 50, 60}},
		{"LINESTRING(40 50 60,10 20 30)", 3, [4]float64{10, 20, 30}, [4]float64{40, 50, 60}},
		{"LINESTRING(10 20 30 40,50 60 70 80)", 4, [4]float64{10, 20, 30, 40}, [4]float64{50, 60, 70, 80}},
		{"LINESTRING(50 60 70 80,10 20 30 40)", 4, [4]float64{10, 20, 30, 40}, [4]float64{50, 60, 70, 80}},
		{"GEOMETRYCOLLECTION(POINT(10 20))", 2, [4]float64{10, 20}, [4]float64{10, 20}},
		{"GEOMETRYCOLLECTION(POINT(10 20),POINT(30 40))", 2, [4]float64{10, 20}, [4]float64{30, 40}},
		{"GEOMETRYCOLLECTION(POINT(10 20 30),POINT(40 50 60))", 3, [4]float64{10, 20, 30}, [4]float64{40, 50, 60}},
		{"GEOMETRYCOLLECTION(POINT(10 20 30 40),POINT(50 60 70 80))", 4, [4]float64{10, 20, 30, 40}, [4]float64{50, 60, 70, 80}},
		{"GEOMETRYCOLLECTION(POINT(40 50 60),POINT(10 20 30))", 3, [4]float64{10, 20, 30}, [4]float64{40, 50, 60}},
		{"GEOMETRYCOLLECTION(POINT(50 60 70 80),POINT(10 20 30 40))", 4, [4]float64{10, 20, 30, 40}, [4]float64{50, 60, 70, 80}},
		{"GEOMETRYCOLLECTION(POINT(10 20),POINT(30 40 50))", 3, [4]float64{10, 20, 50}, [4]float64{30, 40, 50}},
		{"GEOMETRYCOLLECTION(POINT(10 20),POINT(30 40 50 60))", 4, [4]float64{10, 20, 50, 60}, [4]float64{30, 40, 50, 60}},
		{"GEOMETRYCOLLECTION(POINT(10 20),POINT(30 40 50 60),POINT(70 80 90))", 4, [4]float64{10, 20, 50, 60}, [4]float64{70, 80, 90, 60}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g := geomFromWKT(t, tt.wkt)
			min, max, dims := g.FullRect()
			require.Equal(t, tt.dims, dims)
			require.Equal(t, tt.min[:tt.dims], min[:tt.dims])
			require.Equal(t, tt.max[:tt.dims], max[:tt.dims])
		})
	}

	// A nil geometry reports zero dimensions (upstream: tg_geom_fullrect(0,...)).
	var g *Geom
	_, _, dims := g.FullRect()
	require.Equal(t, 0, dims)
}

// TestRectUpstream reuses cases from upstream tests/test_rect.c.
func TestRectUpstream(t *testing.T) {
	r := Rect{Min: Point{0, 0}, Max: Point{10, 10}}

	require.Equal(t, Point{5, 5}, r.Center())
	require.Equal(t, Point{0, 0}, (Rect{}).Center())

	intersects := []struct {
		o    Rect
		want bool
	}{
		{Rect{Point{0, 0}, Point{10, 10}}, true},
		{Rect{Point{2, 2}, Point{8, 8}}, true},
		{Rect{Point{-1, 0}, Point{10, 10}}, true},
		{Rect{Point{0, -1}, Point{10, 10}}, true},
		{Rect{Point{0, 0}, Point{11, 10}}, true},
		{Rect{Point{0, 0}, Point{10, 11}}, true},
		{Rect{Point{11, 0}, Point{21, 10}}, false},
		{Rect{Point{0, 11}, Point{10, 21}}, false},
		{Rect{Point{11, 11}, Point{21, 21}}, false},
		{Rect{Point{-11, 11}, Point{1, 21}}, false},
		{Rect{Point{-11, -11}, Point{-1, -1}}, false},
	}
	for _, tt := range intersects {
		require.Equal(t, tt.want, r.Intersects(tt.o))
	}

	require.True(t, r.Contains(Point{5, 5}))
	require.False(t, r.Contains(Point{15, 15}))
	for x := 0; x <= 10; x++ {
		for y := 0; y <= 10; y++ {
			require.True(t, r.Contains(Point{float64(x), float64(y)}))
		}
	}
	require.False(t, r.Contains(Point{-15, -15}))
	require.False(t, r.Contains(Point{-15, 5}))
	require.False(t, r.Contains(Point{-15, 15}))
	require.False(t, r.Contains(Point{0, -15}))
}

// TestLineUpstream reuses cases from upstream tests/test_line.c.
func TestLineUpstream(t *testing.T) {
	clockwise := []struct {
		wkt  string
		want bool
	}{
		{"LINESTRING(0 0,0 10,10 10,10 0,0 0)", true},
		{"LINESTRING(0 0,10 0,10 10,0 10,0 0)", false},
		{"LINESTRING(0 0,0 10,10 10)", true},
		{"LINESTRING(0 0,10 0,10 10)", false},
	}
	for _, tt := range clockwise {
		l, ok := geomFromWKT(t, tt.wkt).AsLine()
		require.True(t, ok)
		require.Equal(t, tt.want, l.Clockwise())
	}

	// u1: {0,10},{0,0},{10,0},{10,10}, length 30.
	l, ok := geomFromWKT(t, "LINESTRING(0 10,0 0,10 0,10 10)").AsLine()
	require.True(t, ok)
	require.Equal(t, 30.0, l.Length())
	require.Equal(t, []Point{{0, 10}, {0, 0}, {10, 0}, {10, 10}}, l.Points())

	// A clone shares memory, a copy does not.
	cl := l.Clone()
	require.Equal(t, l.Points(), cl.Points())
	cp := l.Copy()
	require.Equal(t, l.Points(), cp.Points())
}

// TestPolyRingUpstream reuses cases from upstream tests/test_poly.c and
// tests/test_ring.c.
func TestPolyRingUpstream(t *testing.T) {
	cw := geomFromWKT(t, "POLYGON((0 0,0 4,4 4,4 0,0 0))")
	p, ok := cw.AsPoly()
	require.True(t, ok)
	require.True(t, p.Clockwise())
	require.Equal(t, 0, p.NumHoles())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{4, 4}}, p.Rect())

	ccw := geomFromWKT(t, "POLYGON((0 0,4 0,4 4,0 4,0 0))")
	p2, ok := ccw.AsPoly()
	require.True(t, ok)
	require.False(t, p2.Clockwise())

	// Octagon is convex; an L-shape is not.
	oct := geomFromWKT(t, "POLYGON((3 0,7 0,10 3,10 7,7 10,3 10,0 7,0 3,3 0))")
	po, ok := oct.AsPoly()
	require.True(t, ok)
	require.True(t, po.Exterior().Convex())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{10, 10}}, po.Rect())

	lshape := geomFromWKT(t, "POLYGON((0 0,0 4,2 4,2 2,4 2,4 0,0 0))")
	pl, ok := lshape.AsPoly()
	require.True(t, ok)
	require.False(t, pl.Exterior().Convex())

	// Poly with a hole: bounds from the exterior, hole accessible by index.
	holey := geomFromWKT(t, "POLYGON((0 0,0 10,10 10,10 0,0 0),(4 4,4 6,6 6,6 4,4 4))")
	ph, ok := holey.AsPoly()
	require.True(t, ok)
	require.Equal(t, 1, ph.NumHoles())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{10, 10}}, ph.Rect())
	hole, ok := ph.HoleAt(0)
	require.True(t, ok)
	require.Equal(t, Rect{Min: Point{4, 4}, Max: Point{6, 6}}, hole.Rect())
	_, ok = ph.HoleAt(-1)
	require.False(t, ok)
	_, ok = ph.HoleAt(1)
	require.False(t, ok)
}

// TestIndexUpstream reuses the indexing expectations from upstream
// tests/test_index.c: an index is only built once there are enough points.
func TestIndexUpstream(t *testing.T) {
	small := geomFromWKT(t, "LINESTRING(0 0,1 1,2 0,3 1)")
	ls, _ := small.AsLine()
	require.Equal(t, 0, ls.IndexSpread(), "small lines are not indexed")

	pts := make([]Point, 64)
	for i := range pts {
		pts[i] = Point{X: float64(i), Y: float64(i % 2)}
	}
	g, err := NewLineString(pts)
	require.NoError(t, err)
	l, ok := g.AsLine()
	require.True(t, ok)
	require.Equal(t, 16, l.IndexSpread(), "default index spread")
	require.GreaterOrEqual(t, l.IndexNumLevels(), 1)
	require.Positive(t, l.IndexLevelNumRects(0))
}
