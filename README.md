tgo
---
![Tests](https://github.com/akhenakh/tgo/actions/workflows/build.yml/badge.svg)

[![GoDoc](https://pkg.go.dev/badge/github.com/akhenakh/tgo)](https://pkg.go.dev/github.com/akhenakh/tgo)


Go bindings for [tidwall/tg](https://github.com/tidwall/tg) Geometry library for C - Fast point-in-polygon 

This is partial but functional, tg is a very small self contained C library, tgo compiles tg, no external dependencies needed.

## Usage

Simply go get this library with CGO enabled (you'll need a C compiler).

#### Read from WKT
```go
// Unmarshal from WKT
input := "POLYGON((0 0,0 1,1 1,1 0,0 0))"
g, _ := tgo.UnmarshalWKT(input)

// Marshal to WKT
output := g.AsText()
fmt.Println(output) // Prints: POLYGON((0 0,0 1,1 1,1 0,0 0))
```

#### Read from GeoJSON
```go
input := `{"type":"Feature","properties":{},"geometry":{"coordinates":[-79.20159897229003,43.636785010689835],"type":"Point"}}`
g, _ := tgo.UnmarshalGeoJSON(input)
```

#### Read from WKB
```go
input := []byte{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 248, 63, 0, 0, 0, 0, 0, 0, 4, 64}
g, _ := tgo.UnmarshalWKB(input)
```

#### Intersects
```go
if Intersects(g1, g2) {
	fmt.Println("Intersects")
}
```

#### Point in Polygon on large FeatureCollections
```go

// load your collection using UnmarshalGeoJSON
found := g.StabOne(2, 48)
if found != nil {
	fmt.Println(found.Properties())
}
// Output: {"properties":{ "ADMIN": "France", "ISO_A2": "FR", "ISO_A3": "FRA" }}
```

#### Types

```go
input := "POLYGON((0 0,0 1,1 1,1 0,0 0))"
g, _ := tgo.UnmarshalWKT(input)

if g.Types() == tgo.Polygon() {
	p, _ := g.AsPoly()
	p.HolesCount()
}
```

## Tests

Some tests are borrowed from [simplefeatures](https://github.com/peterstace/simplefeatures).
