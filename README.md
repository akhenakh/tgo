tgo
---
![Tests](https://github.com/akhenakh/tgo/actions/workflows/build.yml/badge.svg)

[![GoDoc](https://pkg.go.dev/badge/github.com/akhenakh/tgo)](https://pkg.go.dev/github.com/akhenakh/tgo)

[<img src="img/tgo.jpg">](https://github.com/akhenakh/tgo/)

Go bindings for [tidwall/tg](https://github.com/tidwall/tg) Geometry library for C - Fast point-in-polygon.

tgo compiles the self-contained `tg.c` amalgamation, no external dependencies are needed.

Requires Go 1.23+ and CGO (a C compiler).

## Usage

Simply `go get` this library with CGO enabled.

#### Parsing

```go
// WKT
g, _ := tgo.UnmarshalWKT("POLYGON((0 0,0 1,1 1,1 0,0 0))")

// GeoJSON
g, _ := tgo.UnmarshalGeoJSON([]byte(`{"type":"Point","coordinates":[1.5,2.5]}`))

// WKB
g, _ := tgo.UnmarshalWKB(data)

// Hex-encoded WKB or the compact geobin format
g, _ := tgo.UnmarshalHex("0101000000...")
g, _ := tgo.UnmarshalGeoBin(data)

// Auto-detect WKB, WKT, Hex or GeoJSON
g, _ := tgo.Parse(data)
```

Each parser also has an `...AndIndex` variant that selects the indexing
strategy (`tgo.IndexNone`, `tgo.IndexNatural`, `tgo.IndexYStripes`).

#### Writing

```go
s := g.AsWKT()        // WKT
s = g.AsGeoJSON()     // GeoJSON
s = g.AsHex()         // hex-encoded WKB
b := g.AsWKB()        // WKB
b = g.AsGeoBin()      // geobin
s = g.String()        // same as AsWKT, via fmt.Stringer
```

#### Constructors

```go
p, _ := tgo.NewPoint(tgo.Point{X: 1, Y: 2})
l, _ := tgo.NewLineString(tgo.LineString{{0, 0}, {1, 1}})
poly, _ := tgo.NewPolygon(tgo.Polygon{
	Exterior: tgo.LineString{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {0, 0}},
})
gc, _ := tgo.NewGeometryCollection(p, l)
```

Z/M and empty variants are available (`NewPointZ`, `NewLineStringZM`,
`NewPolygonEmpty`, ...).

#### Accessors

```go
g.Type()            // tgo.TypePolygon
g.TypeString()      // "Polygon"
g.Rect()            // bounding rectangle
g.IsEmpty()
g.Dims()            // 2, 3 (Z or M) or 4 (Z and M)
g.HasZ(); g.HasM()
g.NumPolys()
if p, ok := g.AsPoly(); ok { ... }
```

#### Typed elements

`AsPoly`, `AsLine` and `Poly.Exterior` return `*Poly`, `*Line` and `*Ring`
views into the source geometry. Views must not outlive their source; call
`Clone` (cheap, reference-counted) or `Copy` (deep) to keep one independently.

```go
p, _ := g.AsPoly()
r := p.Exterior()
r.Area(); r.Perimeter(); r.Convex(); r.Clockwise()
r.NumPoints(); r.PointAt(0); r.Points(); r.SegmentAt(0)
```

#### Predicates

```go
tgo.Intersects(g1, g2) // also Equals, Disjoint, Contains, Within,
                       // Covers, CoveredBy, Touches
g1.Contains(g2)        // the same operations are also methods
g.IntersectsXY(2, 48)
g.IntersectsRect(r)
```

#### Point in Polygon on large FeatureCollections

```go
// load your collection using UnmarshalGeoJSON
found := g.StabOne(2, 48)
if found != nil {
	fmt.Println(found.Properties())
}
// Output: {"properties":{ "ADMIN": "France", "ISO_A2": "FR", "ISO_A3": "FRA" }}

// Or iterate all matches with an index-aware search.
for index, child := range g.Search(tgo.Rect{Min: tgo.Point{X: 2, Y: 48}, Max: tgo.Point{X: 2, Y: 48}}) {
	_ = index
	_ = child.Properties()
}
```

#### Searches and nearest segment

```go
for pair := range ringA.SearchRing(ringB) { /* pair.A, pair.B */ }
for pair := range ring.SearchLine(line)    { /* ... */ }

nearest, ok := ring.Nearest(x, y)
k := ring.NearestK(x, y, 5) // 5 nearest segments, nearest first
```

#### Value types

`Point`, `Rect` and `Segment` are plain Go values with pure helpers:

```go
p := tgo.Point{X: 1, Y: 2}
p.Rect(); p.In(r)
tgo.Rect{Min: a, Max: b}.Center(); .Union(other); .Extend(p); .Intersects(r)
seg.Intersects(other); seg.Rect()
```

## Tests

Some tests are borrowed from [simplefeatures](https://github.com/peterstace/simplefeatures).
