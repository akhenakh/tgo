package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPoint_RectAndIn(t *testing.T) {
	p := Point{X: 3, Y: 4}

	r := p.Rect()
	require.Equal(t, Rect{Min: p, Max: p}, r)
	require.True(t, p.In(r))
	require.False(t, Point{X: 3.1, Y: 4}.In(r))
}

func TestRect_Methods(t *testing.T) {
	a := Rect{Min: Point{0, 0}, Max: Point{2, 2}}
	b := Rect{Min: Point{1, 1}, Max: Point{3, 3}}

	require.True(t, a.Intersects(b))
	require.False(t, a.Intersects(Rect{Min: Point{5, 5}, Max: Point{6, 6}}))
	require.True(t, a.Intersects(Rect{Min: Point{2, 0}, Max: Point{4, 2}}), "touching edges intersect")

	require.True(t, a.Contains(Point{2, 2}), "boundary is contained")
	require.False(t, a.Contains(Point{2.1, 2}))

	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{3, 3}}, a.Union(b))
	require.Equal(t, Rect{Min: Point{0, 0}, Max: Point{4, 2}}, a.Extend(Point{4, 2}))
	require.Equal(t, Point{1, 1}, a.Center())
}

func TestSegment_Methods(t *testing.T) {
	seg := Segment{A: Point{1, 1}, B: Point{3, 3}}
	require.Equal(t, Rect{Min: Point{1, 1}, Max: Point{3, 3}}, seg.Rect())

	tests := []struct {
		name string
		a, b Segment
		want bool
	}{
		{"share endpoint", Segment{Point{0, 0}, Point{0, 1}}, Segment{Point{0, 0}, Point{1, 0}}, true},
		{"cross", Segment{Point{0, 0}, Point{1, 1}}, Segment{Point{1, 0}, Point{0, 1}}, true},
		{"collinear overlap", Segment{Point{0, 0}, Point{2, 0}}, Segment{Point{1, 0}, Point{3, 0}}, true},
		{"collinear disjoint", Segment{Point{0, 0}, Point{1, 0}}, Segment{Point{2, 0}, Point{3, 0}}, false},
		{"parallel disjoint", Segment{Point{0, 0}, Point{0, 1}}, Segment{Point{1, 0}, Point{1, 1}}, false},
		{"touch at endpoint", Segment{Point{0, 0}, Point{1, 1}}, Segment{Point{1, 1}, Point{2, 0}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.a.Intersects(tt.b))
			require.Equal(t, tt.want, tt.b.Intersects(tt.a), "symmetric")
		})
	}
}
