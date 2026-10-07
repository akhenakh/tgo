package tgo

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestWKTRoundTrip reuses the normalization cases from upstream
// tests/test_wkt.c: parse the input then check the WKT written back.
func TestWKTRoundTrip(t *testing.T) {
	tests := []struct{ in, want string }{
		// Basic syntax and number normalization.
		{"POINT ZM EMPTY", "POINT EMPTY"},
		{" POINT ( -.5 -000 )", "POINT(-0.5 0)"},
		{"POINT( 5 \t 1 \r 3 \n 4 )", "POINT(5 1 3 4)"},
		{" LINESTRINGZM EMPTY", "LINESTRING EMPTY"},
		{" LINESTRINGM (1 \t 2 \r 3,3 \n 4 5)", "LINESTRING M(1 2 3,3 4 5)"},
		{" MULTIPOINT ( ( 1  2 ) , ( 3  4 ) ) ", "MULTIPOINT(1 2,3 4)"},
		{" MULTIPOINT ( ) ", "MULTIPOINT EMPTY"},
		{" MULTIPOINT ( 1 2 3 ) ", "MULTIPOINT(1 2 3)"},
		{" MULTIPOINT (-.5 000.5e4) ", "MULTIPOINT(-0.5 5e3)"},
		{" MULTIPOINT ((-1 -1)) ", "MULTIPOINT(-1 -1)"},
		{" MULTIPOINT (1e5 2000e-2) ", "MULTIPOINT(1e5 20)"},
		{" MULTIPOINT (500000000000000e-10 1) ", "MULTIPOINT(5e4 1)"},
		{" MULTIPOINT ( (1 2) , (2 3)) ", "MULTIPOINT(1 2,2 3)"},
		{" POLYGON ( (1 2, 2 3, 1 2) , (2 3, 4 5, 2 3)) ", "POLYGON((1 2,2 3,1 2),(2 3,4 5,2 3))"},
		{" POLYGON ( (1 2 3, 2 3 4, 1 2 3) , (2 3 4, 4 5 6, 2 3 4)) ", "POLYGON((1 2 3,2 3 4,1 2 3),(2 3 4,4 5 6,2 3 4))"},
		{" POLYGON Z( (1 2 3, 2 3 4, 1 2 3) , (2 3 4, 4 5 6, 2 3 4)) ", "POLYGON((1 2 3,2 3 4,1 2 3),(2 3 4,4 5 6,2 3 4))"},
		{" POLYGON M( (1 2 3, 2 3 4, 1 2 3) , (2 3 4, 4 5 6, 2 3 4)) ", "POLYGON M((1 2 3,2 3 4,1 2 3),(2 3 4,4 5 6,2 3 4))"},
		{" POLYGON ( (1 2 3 4, 2 3 4 5, 1 2 3 4) , (2 3 4 5, 4 5 6 7, 2 3 4 5)) ", "POLYGON((1 2 3 4,2 3 4 5,1 2 3 4),(2 3 4 5,4 5 6 7,2 3 4 5))"},
		{" POLYGON ZM( (1 2 3 4, 2 3 4 5, 1 2 3 4) , (2 3 4 5, 4 5 6 7, 2 3 4 5)) ", "POLYGON((1 2 3 4,2 3 4 5,1 2 3 4),(2 3 4 5,4 5 6 7,2 3 4 5))"},
		{" MULTIPOINT Z (1 2 3) ", "MULTIPOINT(1 2 3)"},
		{" MULTIPOINT M (1 2 3) ", "MULTIPOINT M(1 2 3)"},
		{" MULTIPOINT (1 2 3 4) ", "MULTIPOINT(1 2 3 4)"},
		{" MULTIPOINT ZM (1 2 3 4) ", "MULTIPOINT(1 2 3 4)"},
		{"GEOMETRYCOLLECTION ZM(POINT(1 2 3))", "GEOMETRYCOLLECTION(POINT(1 2 3))"},
		{"PoinT Zm eMpTy", "POINT EMPTY"},
		{"lInEsTrInG zM EmPtY", "LINESTRING EMPTY"},

		// Dimension specifiers.
		{"POINTZ(1 2 3)", "POINT(1 2 3)"},
		{"POINT Z(1 2 3)", "POINT(1 2 3)"},
		{"POINTM(1 2 3)", "POINT M(1 2 3)"},
		{"POINTZM(1 2 3 4)", "POINT(1 2 3 4)"},
		{"POINT ZM(1 2 3 4)", "POINT(1 2 3 4)"},
		{"LINESTRING Z(1 2 4,5 6 7)", "LINESTRING(1 2 4,5 6 7)"},
		{"LINESTRING ZM(1 2 4 5,6 7 8 9)", "LINESTRING(1 2 4 5,6 7 8 9)"},
		{"POLYGON Z((1 2 3,3 4 5,1 2 3),(5 6 7,7 8 9,5 6 7))", "POLYGON((1 2 3,3 4 5,1 2 3),(5 6 7,7 8 9,5 6 7))"},

		// Empty geometry variants.
		{"POINT EMPTY", "POINT EMPTY"},
		{"POINT Z EMPTY", "POINT EMPTY"},
		{"POINT M EMPTY", "POINT EMPTY"},
		{"LINESTRING EMPTY", "LINESTRING EMPTY"},
		{"LINESTRING M EMPTY", "LINESTRING EMPTY"},
		{"POLYGON EMPTY", "POLYGON EMPTY"},
		{"POLYGON ZM EMPTY", "POLYGON EMPTY"},
		{"MULTIPOINT EMPTY", "MULTIPOINT EMPTY"},
		{"MULTIPOINT Z EMPTY", "MULTIPOINT EMPTY"},
		{"MULTILINESTRING EMPTY", "MULTILINESTRING EMPTY"},
		{"MULTILINESTRING Z EMPTY", "MULTILINESTRING EMPTY"},
		{"MULTIPOLYGON EMPTY", "MULTIPOLYGON EMPTY"},
		{"MULTIPOLYGON M EMPTY", "MULTIPOLYGON EMPTY"},
		{"GEOMETRYCOLLECTION EMPTY", "GEOMETRYCOLLECTION EMPTY"},
		{"GEOMETRYCOLLECTION Z EMPTY", "GEOMETRYCOLLECTION EMPTY"},

		// Mixed dimensions / nested collections.
		{"MULTILINESTRING((1 2,5 6))", "MULTILINESTRING((1 2,5 6))"},
		{"MULTILINESTRING((1 2 3,5 6 7))", "MULTILINESTRING((1 2 3,5 6 7))"},
		{"MULTILINESTRING M((1 2 3,5 6 7))", "MULTILINESTRING M((1 2 3,5 6 7))"},
		{"MULTILINESTRING((1 2 3 4,5 6 7 8))", "MULTILINESTRING((1 2 3 4,5 6 7 8))"},
		{"MULTIPOLYGON(((1 2,2 3,1 2)))", "MULTIPOLYGON(((1 2,2 3,1 2)))"},
		{"MULTIPOLYGON(((1 2 3,2 3 4,1 2 3)))", "MULTIPOLYGON(((1 2 3,2 3 4,1 2 3)))"},
		{"MULTIPOLYGON M(((1 2 3,2 3 4,1 2 3)))", "MULTIPOLYGON M(((1 2 3,2 3 4,1 2 3)))"},
		{"GEOMETRYCOLLECTION(POINT(1 2))", "GEOMETRYCOLLECTION(POINT(1 2))"},
		{"GEOMETRYCOLLECTION(POINT(1 2 3))", "GEOMETRYCOLLECTION(POINT(1 2 3))"},
		{"GEOMETRYCOLLECTION(POINT M(1 2 3))", "GEOMETRYCOLLECTION(POINT M(1 2 3))"},
		{"GEOMETRYCOLLECTION(POINT(1 2 3 4))", "GEOMETRYCOLLECTION(POINT(1 2 3 4))"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			g, err := UnmarshalWKT(tt.in)
			require.NoError(t, err)
			require.Equal(t, tt.want, g.AsWKT())
		})
	}
}

// TestWKBAndHexRoundTrip reuses the geometry list from upstream
// tests/test_wkb.c, round-tripping each through hex and WKB.
func TestWKBAndHexRoundTrip(t *testing.T) {
	wkts := []string{
		"POINT EMPTY",
		"POINT(1 2)",
		"POINT(1 2 3)",
		"POINT M(1 2 3)",
		"POINT(1 2 3 4)",
		"LINESTRING EMPTY",
		"LINESTRING(1 2,3 4)",
		"LINESTRING(1 2 3,4 5 6)",
		"LINESTRING M(1 2 3,4 5 6)",
		"LINESTRING(1 2 3 4,5 6 7 8)",
		"POLYGON EMPTY",
		"POLYGON((1 2,3 4,1 2))",
		"POLYGON((1 2,3 4,1 2),(2 3,4 5,2 3))",
		"POLYGON((1 2 3,3 4 5,1 2 3))",
		"POLYGON M((1 2 3,3 4 5,1 2 3))",
		"POLYGON((1 2 3 4,3 4 5 6,1 2 3 4))",
		"MULTIPOINT EMPTY",
		"MULTIPOINT(1 2)",
		"MULTIPOINT(1 2 3)",
		"MULTIPOINT M(1 2 3)",
		"MULTIPOINT(1 2 3 4)",
		"MULTILINESTRING EMPTY",
		"MULTILINESTRING((1 2,3 4))",
		"MULTILINESTRING((1 2 3,3 4 5))",
		"MULTILINESTRING M((1 2 3,3 4 5))",
		"MULTILINESTRING((1 2 3 4,3 4 5 6))",
		"MULTIPOLYGON EMPTY",
		"MULTIPOLYGON(((1 2,3 4,1 2)))",
		"MULTIPOLYGON(((1 2 3,3 4 5,1 2 3)))",
		"MULTIPOLYGON M(((1 2 3,3 4 5,1 2 3)))",
		"MULTIPOLYGON(((1 2 3 4,3 4 5 6,1 2 3 4)))",
		"GEOMETRYCOLLECTION EMPTY",
		"GEOMETRYCOLLECTION(POINT EMPTY)",
		"GEOMETRYCOLLECTION(POINT(1 2 3))",
		"GEOMETRYCOLLECTION(POINT M(1 2 3))",
		"GEOMETRYCOLLECTION(MULTIPOINT(1 2 3 4))",
	}

	for _, wkt := range wkts {
		t.Run(wkt, func(t *testing.T) {
			g := geomFromWKT(t, wkt)

			fromHex, err := UnmarshalHex(g.AsHex())
			require.NoError(t, err)
			require.Equal(t, wkt, fromHex.AsWKT())

			fromWKB, err := UnmarshalWKB(g.AsWKB())
			require.NoError(t, err)
			require.Equal(t, wkt, fromWKB.AsWKT())
		})
	}
}

// TestGeoJSONToWKT reuses the normalization cases from upstream
// tests/test_wkb.c that parse GeoJSON.
func TestGeoJSONToWKT(t *testing.T) {
	tests := []struct{ in, want string }{
		{
			`{"type":"Polygon","coordinates":[[[-112,33,3],[-111,33,4],[-111,32,5],[-112,32,6],[-112,33,3]],[[1,1,0],[1,2,37],[2,2,1],[1,2,10],[1,1,0]]],"hello":[1,2,3]}`,
			"POLYGON((-112 33 3,-111 33 4,-111 32 5,-112 32 6,-112 33 3),(1 1 0,1 2 37,2 2 1,1 2 10,1 1 0))",
		},
		{
			`{"type":"LineString","coordinates":[[1,2],[3,4],[1,2]],"hello":[1,2,3]}`,
			"LINESTRING(1 2,3 4,1 2)",
		},
		{`{"type":"Point","coordinates":[1,2],"a":1}`, "POINT(1 2)"},
		{`{"type":"Point","coordinates":[1,2,3],"a":1}`, "POINT(1 2 3)"},
		{`{"type":"Point","coordinates":[1,2,3,4],"a":1}`, "POINT(1 2 3 4)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			g, err := UnmarshalGeoJSON([]byte(tt.in))
			require.NoError(t, err)
			require.Equal(t, tt.want, g.AsWKT())
		})
	}
}
