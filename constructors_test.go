package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPoint(t *testing.T) {
	g, err := NewPoint(Point{X: 1.5, Y: 2.5})
	require.NoError(t, err)

	require.Equal(t, TypePoint, g.Type())
	require.Equal(t, "POINT(1.5 2.5)", g.AsWKT())
	require.Equal(t, Point{X: 1.5, Y: 2.5}, g.Point())
	require.Equal(t, 2, g.Dims())
	require.False(t, g.IsEmpty())
}

func TestNewPointDims(t *testing.T) {
	g, err := NewPointZ(Point{X: 1, Y: 2}, 3)
	require.NoError(t, err)
	require.True(t, g.HasZ())
	require.False(t, g.HasM())
	require.Equal(t, 3, g.Dims())
	require.Equal(t, 3.0, g.Z())

	g, err = NewPointM(Point{X: 1, Y: 2}, 4)
	require.NoError(t, err)
	require.False(t, g.HasZ())
	require.True(t, g.HasM())
	require.Equal(t, 4.0, g.M())

	g, err = NewPointZM(Point{X: 1, Y: 2}, 3, 4)
	require.NoError(t, err)
	require.True(t, g.HasZ())
	require.True(t, g.HasM())
	require.Equal(t, 4, g.Dims())
}

func TestNewLineString(t *testing.T) {
	g, err := NewLineString(LineString{{0, 0}, {1, 1}, {2, 0}})
	require.NoError(t, err)

	require.Equal(t, TypeLineString, g.Type())
	require.Equal(t, "LINESTRING(0 0,1 1,2 0)", g.AsWKT())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{2, 1}}, g.Rect())
}

func TestNewLineStringExtraCoords(t *testing.T) {
	points := LineString{{0, 0}, {1, 1}}

	g, err := NewLineStringZ(points, []float64{1, 2})
	require.NoError(t, err)
	require.True(t, g.HasZ())
	require.Equal(t, []float64{1, 2}, g.ExtraCoords())

	_, err = NewLineStringZ(points, []float64{1})
	require.Error(t, err)

	g, err = NewLineStringZM(points, []float64{1, 2, 3, 4})
	require.NoError(t, err)
	require.True(t, g.HasZ())
	require.True(t, g.HasM())
	require.Equal(t, []float64{1, 2, 3, 4}, g.ExtraCoords())

	_, err = NewLineStringZM(points, []float64{1, 2})
	require.Error(t, err)
}

func TestNewPolygon(t *testing.T) {
	g, err := NewPolygon(Polygon{
		Exterior: LineString{{0, 0}, {0, 4}, {4, 4}, {4, 0}, {0, 0}},
		Holes:    []LineString{{{1, 1}, {1, 2}, {2, 2}, {2, 1}, {1, 1}}},
	})
	require.NoError(t, err)

	require.Equal(t, TypePolygon, g.Type())
	require.Equal(t, "POLYGON((0 0,0 4,4 4,4 0,0 0),(1 1,1 2,2 2,2 1,1 1))", g.AsWKT())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{4, 4}}, g.Rect())
}

func TestNewMultiGeometries(t *testing.T) {
	mp, err := NewMultiPoint([]Point{{1, 2}, {3, 4}})
	require.NoError(t, err)
	require.Equal(t, TypeMultiPoint, mp.Type())
	require.Equal(t, 2, mp.NumPoints())
	p, ok := mp.PointAt(1)
	require.True(t, ok)
	require.Equal(t, Point{3, 4}, p)
	_, ok = mp.PointAt(2)
	require.False(t, ok)

	ml, err := NewMultiLineString([]LineString{{{0, 0}, {1, 1}}, {{2, 2}, {3, 3}}})
	require.NoError(t, err)
	require.Equal(t, TypeMultiLineString, ml.Type())
	require.Equal(t, "MULTILINESTRING((0 0,1 1),(2 2,3 3))", ml.AsWKT())

	mpoly, err := NewMultiPolygon([]Polygon{
		{Exterior: LineString{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}}},
		{Exterior: LineString{{2, 2}, {2, 3}, {3, 3}, {3, 2}, {2, 2}}},
	})
	require.NoError(t, err)
	require.Equal(t, TypeMultiPolygon, mpoly.Type())
	require.Equal(t, 2, mpoly.NumPolys())
	pol, ok := mpoly.PolyAt(0)
	require.True(t, ok)
	require.Equal(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))", pol.AsWKT())
}

func TestNewGeometryCollection(t *testing.T) {
	pt, err := NewPoint(Point{1, 2})
	require.NoError(t, err)
	ls, err := NewLineString(LineString{{0, 0}, {1, 1}})
	require.NoError(t, err)

	gc, err := NewGeometryCollection(pt, ls)
	require.NoError(t, err)
	require.Equal(t, TypeGeometryCollection, gc.Type())
	require.Equal(t, 2, gc.NumGeometries())

	child, ok := gc.GeometryAt(0)
	require.True(t, ok)
	require.Equal(t, TypePoint, child.Type())

	_, ok = gc.GeometryAt(2)
	require.False(t, ok)
}

func TestNewEmpty(t *testing.T) {
	tests := []struct {
		name string
		ctor func() (*Geom, error)
		want GeomType
	}{
		{"point", NewPointEmpty, TypePoint},
		{"linestring", NewLineStringEmpty, TypeLineString},
		{"polygon", NewPolygonEmpty, TypePolygon},
		{"multipoint", NewMultiPointEmpty, TypeMultiPoint},
		{"multilinestring", NewMultiLineStringEmpty, TypeMultiLineString},
		{"multipolygon", NewMultiPolygonEmpty, TypeMultiPolygon},
		{"geometrycollection", NewGeometryCollectionEmpty, TypeGeometryCollection},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := tt.ctor()
			require.NoError(t, err)
			require.Equal(t, tt.want, g.Type())
			require.True(t, g.IsEmpty())
		})
	}
}

func TestCloneCopy(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 1,1 1,1 0,0 0))")

	cl := g.Clone()
	require.NotNil(t, cl)
	require.NotSame(t, g, cl)
	require.Equal(t, g.AsWKT(), cl.AsWKT())

	cp := g.Copy()
	require.NotNil(t, cp)
	require.Equal(t, g.AsWKT(), cp.AsWKT())
}

func TestNewError(t *testing.T) {
	g := NewError("boom")
	require.NotNil(t, g)
	require.Equal(t, "boom", g.Err())
}

func TestAccessorsFullRect(t *testing.T) {
	g := geomFromWKT(t, "POLYGON((0 0,0 2,3 2,3 0,0 0))")

	min, max, dims := g.FullRect()
	require.Equal(t, 2, dims)
	require.Equal(t, [4]float64{0, 0, 0, 0}, min)
	require.Equal(t, [4]float64{3, 2, 0, 0}, max)

	require.Equal(t, "Polygon", g.TypeString())
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{3, 2}}, g.Rect())
}
