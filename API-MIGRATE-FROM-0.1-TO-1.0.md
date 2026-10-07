# Migrating tgo from v0.1 to v1.0

v1.0 completes the binding of the [tidwall/tg](https://github.com/tidwall/tg) C
API and makes a few naming changes to be consistently Go-ish. The module path is
unchanged (`github.com/akhenakh/tgo`); the Go language requirement moves from
1.21 to **1.23**.

Most of the change is additive. The breaking changes are all renames, listed
below.

## Breaking changes

### Renamed constants

`GeomType` values are now prefixed with `Type`, and `IndexType` values with
`Index`.

| v0.1 | v1.0 |
| --- | --- |
| `Point` | `TypePoint` |
| `LineString` | `TypeLineString` |
| `Polygon` | `TypePolygon` |
| `MultiPoint` | `TypeMultiPoint` |
| `MultiLineString` | `TypeMultiLineString` |
| `MultiPolygon` | `TypeMultiPolygon` |
| `GeometryCollection` | `TypeGeometryCollection` |
| `None` | `IndexNone` |
| `Natural` | `IndexNatural` |
| `YStripes` | `IndexYStripes` |
| — | `IndexDefault` (new) |

The numeric values are unchanged; `IndexDefault` is the new zero value.

### Renamed methods

| v0.1 | v1.0 |
| --- | --- |
| `*Geom.AsText()` | `*Geom.AsWKT()` |
| `*Poly.AsText()` | `*Poly.AsWKT()` |
| `*Ring.AsText()` | `*Ring.AsWKT()` |
| `*MultiPoly.AsText()` | `*MultiPoly.AsWKT()` |
| `*Poly.IsClockWise()` | `*Poly.Clockwise()` |
| `*Poly.HolesCount()` | `*Poly.NumHoles()` |
| `*MultiPoly.PolygonsCount()` | `*MultiPoly.NumPolygons()` |

Every type that had `AsText` now also implements `String()` (same value as
`AsWKT`).

```go
// v0.1
if g.Type() == tgo.Polygon { ... }
s := g.AsText()
p, _ := g.AsPoly()
n := p.HolesCount()

// v1.0
if g.Type() == tgo.TypePolygon { ... }
s := g.AsWKT()
p, _ := g.AsPoly()
n := p.NumHoles()
```

### Behavior changes

- `UnmarshalWKT("")` / `UnmarshalGeoJSON(nil)` now return an `"empty data"`
  error.
- Parsers use the length-aware C entry points and free failed parses; a parse
  error no longer leaks the error geometry.
- The vendored C library was updated from upstream `efbb5b7` to `ea9c0e8`. This
  is a behavior-only change (bug fixes), with no public C API change.

## Additions

### Value types

`Point`, `Rect`, `Segment` (plain Go values), `LineString`, `Polygon`, with
helpers:

```go
p.Rect(); p.In(r)
s.Rect(); s.Intersects(other)
r.Center(); r.Union(o); r.Extend(p); r.Intersects(o); r.Contains(p)
```

### Parsers and writers

- `UnmarshalHex`, `UnmarshalGeoBin` (and their `...AndIndex` forms).
- `GeoBinRect`, `GeoBinPoint`, `GeoBinFullRect`.
- `(*Geom).AsGeoJSON`, `AsWKB`, `AsHex`, `AsGeoBin`.

### Constructors and copying

- `NewPoint`, `NewLineString`, `NewPolygon`, `NewMultiPoint`,
  `NewMultiLineString`, `NewMultiPolygon`, `NewGeometryCollection`, each with
  `Z`/`M`/`ZM` and `Empty` variants where the C API provides them,
  `NewError`, and the `...AndIndex` forms for the indexed types.
- `Clone` (cheap, reference-counted) and `Copy` (deep) on `Geom`, `Line`,
  `Ring` and `Poly`.

### Accessors

On `*Geom`: `TypeString`, `Rect`, `FullRect`, `IsEmpty`, `IsFeature`,
`IsFeatureCollection`, `Err`, `Point`, `NumPoints`, `PointAt`, `AsLine`,
`NumLines`, `LineAt`, `NumPolys`, `PolyAt`, `NumGeometries`, `GeometryAt`,
`Dims`, `HasZ`, `HasM`, `Z`, `M`, `NumExtraCoords`, `ExtraCoords`,
`IntersectsRect`, `IntersectsXY`.

The existing predicates are also available as methods (`g1.Contains(g2)`,
`g.Intersects(other)`, ...).

### Typed elements

- New `Line` type with `Rect`, `MemSize`, `NumPoints`, `Points`, `PointAt`,
  `NumSegments`, `SegmentAt`, `Clockwise`, `Length`, index introspection,
  `AsGeom`, `AsWKT`, `Clone`, `Copy`.
- `Ring` gains `Clone`, `Copy`, `MemSize`, `Rect`, `NumPoints`, `Points`,
  `PointAt`, `NumSegments`, `SegmentAt`, `Convex`, `Clockwise` and index
  introspection.
- `Poly` gains `Clone`, `Copy`, `MemSize`, `HoleAt` and `Rect`.

### Searches and nearest

- `(*Geom).Search(rect) iter.Seq2[int, *Geom]`.
- `(*Ring).SearchLine`, `(*Ring).SearchRing`, `(*Line).SearchLine`, returning
  `iter.Seq[SegmentPair]`.
- `Nearest` and `NearestK` on `Ring` and `Line`, returning `NearestSegment`.

### Environment

`SetIndex`, `SetIndexSpread`, `SetPrintFixedFloats`.

## Not wrapped

`tg_env_set_allocator` is intentionally not exposed: installing a C allocator
from Go is unsafe with cgo. The C library's default allocator is used.

## Memory model

The Go wrappers do not copy geometry data. Accessors return **views** into the
C-owned memory, and the data is released when the owning `*Geom` (or a `Clone`)
is garbage-collected. A view must therefore not outlive its source; call
`Clone()` (cheap, reference-counted) or `Copy()` (deep) when you need to keep it
around.

See the README for details.
