package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeomSearch(t *testing.T) {
	p1 := geomFromWKT(t, "POINT(0 0)")
	p2 := geomFromWKT(t, "POINT(5 5)")
	p3 := geomFromWKT(t, "POINT(0 1)")
	gc, err := NewGeometryCollection(p1, p2, p3)
	require.NoError(t, err)

	var indices []int
	for i, child := range gc.Search(Rect{Min: Point{-1, -1}, Max: Point{1, 2}}) {
		indices = append(indices, i)
		require.Equal(t, TypePoint, child.Type())
	}
	require.Equal(t, []int{0, 2}, indices)

	// A non-collection geometry yields nothing.
	single := geomFromWKT(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))")
	count := 0
	for range single.Search(Rect{Min: Point{-1, -1}, Max: Point{2, 2}}) {
		count++
	}
	require.Equal(t, 0, count)
}

func TestSegmentSearches(t *testing.T) {
	r1 := geomFromWKT(t, "POLYGON((0 0,0 2,2 2,2 0,0 0))")
	p1, _ := r1.AsPoly()
	ring1 := p1.Exterior()

	r2 := geomFromWKT(t, "POLYGON((1 1,1 3,3 3,3 1,1 1))")
	p2, _ := r2.AsPoly()
	ring2 := p2.Exterior()

	var pairs []SegmentPair
	for pair := range ring1.SearchRing(ring2) {
		require.True(t, pair.A.Intersects(pair.B))
		pairs = append(pairs, pair)
	}
	require.NotEmpty(t, pairs)

	// Ring/line search.
	line := geomFromWKT(t, "LINESTRING(-1 1,3 1)")
	l, _ := line.AsLine()

	n := 0
	for pair := range ring1.SearchLine(l) {
		require.True(t, pair.A.Intersects(pair.B))
		n++
	}
	require.Positive(t, n)

	// Line/line search.
	a := geomFromWKT(t, "LINESTRING(0 0,2 2)")
	la, _ := a.AsLine()
	b := geomFromWKT(t, "LINESTRING(0 2,2 0)")
	lb, _ := b.AsLine()

	n = 0
	for pair := range la.SearchLine(lb) {
		require.True(t, pair.A.Intersects(pair.B))
		n++
	}
	require.Equal(t, 1, n)
}

func TestNearest(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 4,4 4,4 0,0 0))")
	p, _ := g.AsPoly()
	ring := p.Exterior()

	nearest, ok := ring.Nearest(2, 5)
	require.True(t, ok)
	require.Equal(t, 1.0, nearest.Distance)
	require.Equal(t, 4.0, nearest.Segment.A.Y)
	require.Equal(t, 4.0, nearest.Segment.B.Y)

	res := ring.NearestK(2, 5, 2)
	require.Len(t, res, 2)
	require.LessOrEqual(t, res[0].Distance, res[1].Distance)

	line := geomFromWKT(t, "LINESTRING(0 0,3 4,6 0)")
	l, _ := line.AsLine()
	nearest, ok = l.Nearest(3, 4)
	require.True(t, ok)
	require.Equal(t, 0.0, nearest.Distance)
}

func TestEnv(t *testing.T) {
	// Restore the library defaults so other tests are unaffected.
	SetIndex(IndexDefault)
	SetIndexSpread(16)
	SetPrintFixedFloats(false)

	SetIndex(IndexNatural)
	g := geomFromWKT(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))")
	require.Equal(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))", g.AsWKT())
}
