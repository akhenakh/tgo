package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiPoly_PolygonsCount(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int
	}{
		{
			"multipolygon 2",
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)),((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20)))`,
			2,
		},
		{
			"multipolygon 1",
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)))`,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := UnmarshalWKT(tt.data)
			require.NoError(t, err)
			mp, valid := g.AsMultiPoly()
			require.True(t, valid)
			require.Equal(t, tt.data, mp.AsGeom().AsText())
			require.Equal(t, MultiPolygon, mp.AsGeom().Type())
			require.Equal(t, tt.want, mp.PolygonsCount())
		})
	}
}

func TestMultiPoly_PolygonAt(t *testing.T) {
	tests := []struct {
		name  string
		index int
		data  string
		want  string
	}{
		{
			"multipolygon 2",
			0,
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)),((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20)))`,
			`POLYGON((40 40,20 45,45 30,40 40))`,
		},
		{
			"multipolygon 2",
			1,
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)),((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20)))`,
			`POLYGON((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20))`,
		},
		{
			"multipolygon 2",
			2,
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)),((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20)))`,
			``,
		},
		{
			"multipolygon 2",
			-1,
			`MULTIPOLYGON(((40 40,20 45,45 30,40 40)),((20 35,45 20,30 5,10 10,10 30,20 35),(30 20,20 25,20 15,30 20)))`,
			``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := UnmarshalWKT(tt.data)
			require.NoError(t, err)
			mp, valid := g.AsMultiPoly()
			require.True(t, valid)
			p, valid := mp.PolygonAt(tt.index)
			if tt.want != "" && !valid {
				t.Fatalf("not expecting polygon but got some")
			}
			if tt.want == "" {
				return
			}

			require.Equal(t, tt.want, p.AsText())
		})
	}
}

func TestPoly_HolesCount(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int
	}{
		{
			"polygon 1 hole",
			`POLYGON ((20 35, 45 20, 30 5, 10 10, 10 30, 20 35), (30 20, 20 25, 20 15, 30 20))`,
			1,
		},
		{
			"polygon no hole",
			`POLYGON ((20 35, 45 20, 30 5, 10 10, 10 30, 20 35))`,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := UnmarshalWKT(tt.data)
			require.NoError(t, err)
			p, valid := g.AsPoly()
			require.True(t, valid)
			require.Equal(t, tt.want, p.HolesCount())
		})
	}
}
