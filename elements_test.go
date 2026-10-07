package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLine(t *testing.T) {
	g := geomFromWKT(t, "LINESTRING(0 0,3 4,6 0)")

	l, ok := g.AsLine()
	require.True(t, ok)

	require.Equal(t, 3, l.NumPoints())
	require.Equal(t, 2, l.NumSegments())
	require.Equal(t, 10.0, l.Length())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{6, 4}}, l.Rect())
	require.Positive(t, l.MemSize())
	require.Equal(t, g.AsWKT(), l.AsWKT())
	require.Equal(t, g.AsWKT(), l.String())

	p0, ok := l.PointAt(0)
	require.True(t, ok)
	require.Equal(t, Point{0, 0}, p0)
	_, ok = l.PointAt(3)
	require.False(t, ok)

	require.Equal(t, []Point{{0, 0}, {3, 4}, {6, 0}}, l.Points())

	seg, ok := l.SegmentAt(0)
	require.True(t, ok)
	require.Equal(t, Segment{A: Point{0, 0}, B: Point{3, 4}}, seg)

	require.Equal(t, TypeLineString, l.AsGeom().Type())
}

func TestRingElement(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 4,4 4,4 0,0 0))")

	p, ok := g.AsPoly()
	require.True(t, ok)

	r := p.Exterior()
	require.Equal(t, 16.0, r.Area())
	require.Equal(t, 16.0, r.Perimeter())
	require.True(t, r.Convex())
	require.True(t, r.Clockwise())
	require.Equal(t, 5, r.NumPoints(), "closed ring keeps the closing point")
	require.Equal(t, 4, r.NumSegments())
	require.Positive(t, r.MemSize())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{4, 4}}, r.Rect())

	_, ok = r.PointAt(0)
	require.True(t, ok)
	require.Len(t, r.Points(), r.NumPoints())

	require.Equal(t, g.AsWKT(), r.AsGeom().AsWKT())
}

func TestPolyElement(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 4,4 4,4 0,0 0),(1 1,1 2,2 2,2 1,1 1))")

	p, ok := g.AsPoly()
	require.True(t, ok)

	require.Equal(t, 1, p.NumHoles())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{4, 4}}, p.Rect())
	require.Positive(t, p.MemSize())

	hole, ok := p.HoleAt(0)
	require.True(t, ok)
	require.Equal(t, 1.0, hole.Area())
	_, ok = p.HoleAt(1)
	require.False(t, ok)
}

func TestElementCloneCopy(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 4,4 4,4 0,0 0))")
	p, _ := g.AsPoly()
	r := p.Exterior()

	rc := r.Clone()
	require.NotNil(t, rc)
	require.NotSame(t, r, rc)
	require.Equal(t, r.Area(), rc.Area())

	rp := r.Copy()
	require.NotNil(t, rp)
	require.Equal(t, r.Area(), rp.Area())

	pc := p.Clone()
	require.NotNil(t, pc)
	require.Equal(t, p.AsWKT(), pc.AsWKT())

	l := geomFromWKT(t, "LINESTRING(0 0,1 1)")
	line, _ := l.AsLine()
	lc := line.Clone()
	require.NotNil(t, lc)
	require.Equal(t, line.AsWKT(), lc.AsWKT())
}

func TestGeomLineAccessors(t *testing.T) {
	g := geomFromWKT(t, "MULTILINESTRING((0 0,1 1),(2 2,3 3))")
	require.Equal(t, 2, g.NumLines())

	line, ok := g.LineAt(1)
	require.True(t, ok)
	require.Equal(t, "LINESTRING(2 2,3 3)", line.AsWKT())

	_, ok = g.LineAt(2)
	require.False(t, ok)

	// AsLine is only applicable to LineString.
	_, ok = g.AsLine()
	require.False(t, ok)

	poly := geomFromWKT(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))")
	_, ok = poly.AsLine()
	require.False(t, ok)
}
