package tgo

/*
#include "tg.h"
*/
import "C"

// Point is a position in two-dimensional space.
type Point struct {
	X, Y float64
}

// Rect is an axis-aligned rectangle defined by its minimum and maximum points.
type Rect struct {
	Min, Max Point
}

// Segment is a line segment joining two points.
type Segment struct {
	A, B Point
}

// LineString is a series of connected points.
type LineString []Point

// Polygon describes a polygon by its exterior ring and optional holes. It is
// the plain-value form used by the geometry constructors.
type Polygon struct {
	Exterior LineString
	Holes    []LineString
}

func cPoint(p Point) C.struct_tg_point {
	return C.struct_tg_point{C.double(p.X), C.double(p.Y)}
}

func goPoint(p C.struct_tg_point) Point {
	return Point{X: float64(p.x), Y: float64(p.y)}
}

func cRect(r Rect) C.struct_tg_rect {
	return C.struct_tg_rect{cPoint(r.Min), cPoint(r.Max)}
}

func goRect(r C.struct_tg_rect) Rect {
	return Rect{Min: goPoint(r.min), Max: goPoint(r.max)}
}

func cSegment(s Segment) C.struct_tg_segment {
	return C.struct_tg_segment{cPoint(s.A), cPoint(s.B)}
}

func goSegment(s C.struct_tg_segment) Segment {
	return Segment{A: goPoint(s.a), B: goPoint(s.b)}
}

// Rect returns the rectangle that contains the point.
func (p Point) Rect() Rect {
	return goRect(C.tg_point_rect(cPoint(p)))
}

// In reports whether the point is inside or on the boundary of r.
func (p Point) In(r Rect) bool {
	return bool(C.tg_point_intersects_rect(cPoint(p), cRect(r)))
}

// Rect returns the minimum bounding rectangle of the segment.
func (s Segment) Rect() Rect {
	return goRect(C.tg_segment_rect(cSegment(s)))
}

// Intersects reports whether the segment intersects o.
func (s Segment) Intersects(o Segment) bool {
	return bool(C.tg_segment_intersects_segment(cSegment(s), cSegment(o)))
}

// Union returns the smallest rectangle that contains both r and o.
func (r Rect) Union(o Rect) Rect {
	return goRect(C.tg_rect_expand(cRect(r), cRect(o)))
}

// Extend returns the smallest rectangle that contains r and p.
func (r Rect) Extend(p Point) Rect {
	return goRect(C.tg_rect_expand_point(cRect(r), cPoint(p)))
}

// Center returns the center point of the rectangle.
func (r Rect) Center() Point {
	return goPoint(C.tg_rect_center(cRect(r)))
}

// Intersects reports whether r and o intersect.
func (r Rect) Intersects(o Rect) bool {
	return bool(C.tg_rect_intersects_rect(cRect(r), cRect(o)))
}

// Contains reports whether p is inside or on the boundary of r.
func (r Rect) Contains(p Point) bool {
	return bool(C.tg_rect_intersects_point(cRect(r), cPoint(p)))
}
