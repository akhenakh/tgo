package tgo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRing_Perimeter(t *testing.T) {
	tests := []struct {
		name string
		data string
		want float64
	}{
		{
			"polygon",
			`POLYGON((0 0, 40 0, 40 40, 0 40, 0 0))`,
			160,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := UnmarshalWKT(tt.data)
			require.NoError(t, err)
			p, valid := g.AsPoly()
			require.True(t, valid)
			r := p.Exterior()
			require.Equal(t, tt.want, r.Perimeter())
		})
	}
}

func TestRing_Area(t *testing.T) {
	tests := []struct {
		name string
		data string
		want float64
	}{
		{
			"polygon",
			`POLYGON((0 0, 40 0, 40 40, 0 40, 0 0))`,
			1600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := UnmarshalWKT(tt.data)
			require.NoError(t, err)
			p, valid := g.AsPoly()
			require.True(t, valid)
			r := p.Exterior()
			require.Equal(t, tt.want, r.Area())
		})
	}
}
