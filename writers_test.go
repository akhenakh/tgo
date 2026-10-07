package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriters(t *testing.T) {
	const wkt = "POLYGON((0 0,0 1,1 1,1 0,0 0))"
	g := geomFromWKT(t, wkt)

	require.Equal(t, wkt, g.AsWKT())
	require.Equal(t, wkt, g.String())

	tests := []struct {
		name   string
		encode func(*Geom) []byte
		decode func([]byte) (*Geom, error)
	}{
		{
			"geojson",
			func(g *Geom) []byte { return []byte(g.AsGeoJSON()) },
			UnmarshalGeoJSON,
		},
		{
			"wkb",
			func(g *Geom) []byte { return g.AsWKB() },
			UnmarshalWKB,
		},
		{
			"hex",
			func(g *Geom) []byte { return []byte(g.AsHex()) },
			func(b []byte) (*Geom, error) { return UnmarshalHex(string(b)) },
		},
		{
			"geobin",
			func(g *Geom) []byte { return g.AsGeoBin() },
			UnmarshalGeoBin,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := tt.encode(g)
			require.NotEmpty(t, encoded)

			got, err := tt.decode(encoded)
			require.NoError(t, err)
			require.Equal(t, wkt, got.AsWKT())
		})
	}
}

func TestGeoBinPeek(t *testing.T) {
	g := geomFromWKT(t, "POINT(1.5 2.5)")
	gb := g.AsGeoBin()
	require.NotEmpty(t, gb)

	require.Equal(t, Point{X: 1.5, Y: 2.5}, GeoBinPoint(gb))
	require.Equal(t, Rect{Min: Point{1.5, 2.5}, Max: Point{1.5, 2.5}}, GeoBinRect(gb))

	min, max, dims := GeoBinFullRect(gb)
	require.Equal(t, 2, dims)
	require.Equal(t, [4]float64{1.5, 2.5, 0, 0}, min)
	require.Equal(t, [4]float64{1.5, 2.5, 0, 0}, max)
}

func TestGeoBinInvalid(t *testing.T) {
	require.Equal(t, Rect{}, GeoBinRect(nil))
	require.Equal(t, Point{}, GeoBinPoint(nil))
	_, _, dims := GeoBinFullRect(nil)
	require.Equal(t, 0, dims)
}
