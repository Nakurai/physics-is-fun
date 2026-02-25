package main

import "math"

type Edge struct {
	A      Vector[float32]
	B      Vector[float32]
	Middle Vector[float32] // middle is a point in space
	Normal Vector[float32] // the normal vector of this edge
}

// Middle returns a point in space. It is the middle of the edge
func (e Edge) GetMiddle() Vector[float32] {
	return Vector[float32]{X: (e.A.X + e.B.X) / 2, Y: (e.A.Y + e.B.Y) / 2}
}

// Normal returns the normal vector.
// a normal vector is a *direction*
// it is up to the caller to apply this direction to a specific point in space to get coordonates
func (e Edge) GetNormal() Vector[float32] {

	// the slope of the edge is (dx,dy)
	slope := SubVectors(e.B, e.A)

	// we calculate the length so we can normalize below
	// length is sqrt( dx^2 + dy^2 ) according to the pythagoeran theorem
	dx2 := math.Pow(float64(slope.X), 2)
	dy2 := math.Pow(float64(slope.Y), 2)
	length := float32(math.Sqrt(dx2 + dy2))

	// we normalize the original vector right away so we don't have to worry about it later
	// to normalize a vector: (dx/L , dy/L)
	normalized := Vector[float32]{X: slope.X / length, Y: slope.Y / length}

	// the normal formula is (-dy,dx) OR (dy,-dx)
	normal := Vector[float32]{X: -normalized.Y, Y: normalized.X}

	return normal
}

// Compute calculate the middle point and the normal vector and store the result
func (e *Edge) Compute() {
	e.Middle = e.GetMiddle()
	e.Normal = e.GetNormal()
}
