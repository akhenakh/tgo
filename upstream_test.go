package tgo

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// These tests lock in behavior fixed upstream in tidwall/tg after the version
// tgo previously vendored (upstream efbb5b7, 2025-04-30).

// TestUnmarshalWKBEWKB covers PostGIS EWKB Z/M/ZM type-word flags. Before the
// upstream fix (4b99b5a) the flag bits were masked away and 3D/4D coordinates
// were silently misread as 2D.
func TestUnmarshalWKBEWKB(t *testing.T) {
	tests := []struct {
		name   string
		hex    string
		asText string
	}{
		{
			"EWKB Z point",
			"0101000080000000000000f03f00000000000000400000000000000840",
			"POINT(1 2 3)",
		},
		{
			"EWKB M point",
			"0101000040000000000000f03f00000000000000400000000000000840",
			"POINT M(1 2 3)",
		},
		{
			"EWKB ZM point",
			"01010000c0000000000000f03f000000000000004000000000000008400000000000001040",
			"POINT(1 2 3 4)",
		},
		{
			"EWKB SRID + Z point",
			"01010000a0e6100000000000000000f03f00000000000000400000000000000840",
			"POINT(1 2 3)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := hex.DecodeString(tt.hex)
			require.NoError(t, err)

			g, err := UnmarshalWKB(data)
			require.NoError(t, err)
			require.Equal(t, tt.asText, g.AsText())
		})
	}
}

// TestEqualsEmpty covers empty geometries being equal to each other,
// including across types (upstream ab64219).
func TestEqualsEmpty(t *testing.T) {
	tests := []struct {
		a, b string
	}{
		{"POINT EMPTY", "POINT EMPTY"},
		{"LINESTRING EMPTY", "LINESTRING EMPTY"},
		{"POLYGON EMPTY", "POLYGON EMPTY"},
		{"POINT EMPTY", "LINESTRING EMPTY"},
	}

	for _, tt := range tests {
		t.Run(tt.a+" == "+tt.b, func(t *testing.T) {
			g1 := geomFromWKT(t, tt.a)
			g2 := geomFromWKT(t, tt.b)

			require.True(t, Equals(g1, g2))
			require.False(t, Intersects(g1, g2))
			require.True(t, Disjoint(g1, g2))
		})
	}
}

// TestRelationsUpstreamFixes exercises relation fixes for boundary and
// collection edge cases (upstream 9df9f9e, 8d294a6, d187b7e, 78a0ce3).
func TestRelationsUpstreamFixes(t *testing.T) {
	tests := []struct {
		name                                string
		a, b                                string
		contains, within, covers, coveredby bool
		touches                             bool
	}{
		{
			name:    "line touches multipolygon boundary",
			a:       "LINESTRING(0 1, 3 1)",
			b:       "MULTIPOLYGON(((0 0,1 0,1 1,0 1,0 0)),((1 1,2 1,2 2,1 2,1 1)))",
			touches: true,
		},
		{
			name:    "line touches polygon boundary without shared vertices",
			a:       "POLYGON((0 0,2 0,2 2,0 2,0 0))",
			b:       "LINESTRING(1 3,3 1)",
			touches: true,
		},
		{
			name:      "collection with boundary point covered by polygon",
			a:         "GEOMETRYCOLLECTION(POINT EMPTY, POINT(1 1))",
			b:         "POLYGON((0 0,2 0,2 2,0 2,0 0))",
			within:    true,
			coveredby: true,
		},
		{
			name: "collection of nested empties is disjoint",
			a:    "GEOMETRYCOLLECTION(POINT EMPTY, POINT EMPTY)",
			b:    "POLYGON((0 0,2 0,2 2,0 2,0 0))",
			// all false
		},
		{
			name:    "collection covers boundary point but does not contain it",
			a:       "GEOMETRYCOLLECTION(POINT(0 0), LINESTRING(0 0,0 1))",
			b:       "POINT(0 0)",
			covers:  true,
			touches: true,
		},
		{
			name:      "single point collection equals and contains point",
			a:         "GEOMETRYCOLLECTION(POINT(0 0))",
			b:         "POINT(0 0)",
			contains:  true,
			within:    true,
			covers:    true,
			coveredby: true,
		},
		{
			name:     "multilinestring contains shared endpoint point",
			a:        "MULTILINESTRING((0 0,0 1),(0 0,1 0))",
			b:        "POINT(0 0)",
			contains: true,
			covers:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g1 := geomFromWKT(t, tt.a)
			g2 := geomFromWKT(t, tt.b)

			require.Equal(t, tt.contains, Contains(g1, g2), "Contains")
			require.Equal(t, tt.within, Within(g1, g2), "Within")
			require.Equal(t, tt.covers, Covers(g1, g2), "Covers")
			require.Equal(t, tt.coveredby, CoveredBy(g1, g2), "CoveredBy")
			require.Equal(t, tt.touches, Touches(g1, g2), "Touches")
		})
	}
}
