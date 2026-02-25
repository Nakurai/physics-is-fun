package main

import "golang.org/x/exp/constraints"

type Vector[T constraints.Integer | constraints.Float] struct {
	X T
	Y T
}

// Add updates the vector by adding v2
func (v *Vector[T]) Add(v2 Vector[T]) {
	v.X += v2.X
	v.Y += v2.Y
}

// Add updates the vector by multiplying v2
func (v *Vector[T]) Times(v2 Vector[T]) {
	v.X *= v2.X
	v.Y *= v2.Y
}

// Dot returns the dot product of the vector by v2
func (v *Vector[T]) Dot(v2 Vector[T]) T {
	return v.X*v2.X + v.Y*v2.Y
}

// SubVectors adds v2 from v1
func AddVectors[T constraints.Integer | constraints.Float](v1, v2 Vector[T]) Vector[T] {
	return Vector[T]{X: v1.X + v2.X, Y: v1.Y + v2.Y}
}

// SubVectors substracts v2 from v1 and return a new resulting vector
func SubVectors[T constraints.Integer | constraints.Float](v1, v2 Vector[T]) Vector[T] {
	return Vector[T]{X: v1.X - v2.X, Y: v1.Y - v2.Y}
}

// ScaleVector will multiply each vector's components by the scale factor and return the new resulting vector
func ScaleVector[T constraints.Integer | constraints.Float](v Vector[T], scale T) Vector[T] {
	return Vector[T]{X: v.X * scale, Y: v.Y * scale}
}
