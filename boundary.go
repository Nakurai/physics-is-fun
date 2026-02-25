package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Boundary struct {
	Center    Vector[float32] // center of the boundary
	Dimension Vector[float32]
	Edges     []Edge
	Rotation  float64 // in radians
}

// Rotate will rotate the boundary by a certain number of degreee
// !! angleDeg is in degrees
func (b *Boundary) Rotate(angleDeg float64) {
	b.Rotation += angleDeg * math.Pi / 180
	b.ComputeEdges()
}

func (b *Boundary) Update() {
}

func (b *Boundary) Draw(screen *ebiten.Image) {
	for _, e := range b.Edges {
        vector.StrokeLine(screen, e.A.X, e.A.Y, e.B.X, e.B.Y, 1, Red, false)
    }
}

// rotateAround rotates point p around origin o by angle theta (radians).
func rotateAround(p, o Vector[float32], theta float64) Vector[float32] {
    sin, cos := math.Sincos(theta)
    dx := float64(p.X - o.X)
    dy := float64(p.Y - o.Y)
    return Vector[float32]{
        X: o.X + float32(dx*cos-dy*sin),
        Y: o.Y + float32(dx*sin+dy*cos),
    }
}

func (b *Boundary) ComputeEdges() {
	halfDim := ScaleVector(b.Dimension, 0.5)

	topLeft := Vector[float32]{X: b.Center.X - halfDim.X, Y: b.Center.Y - halfDim.Y}
	topLeft = rotateAround(topLeft, b.Center, b.Rotation)
	
	topRight := Vector[float32]{X: b.Center.X + halfDim.X, Y: b.Center.Y - halfDim.Y}
	topRight = rotateAround(topRight, b.Center, b.Rotation)
	
	bottomLeft := Vector[float32]{X: b.Center.X - halfDim.X, Y: b.Center.Y + halfDim.Y}
	bottomLeft = rotateAround(bottomLeft, b.Center, b.Rotation)
	
	bottomRight := Vector[float32]{X: b.Center.X + halfDim.X, Y: b.Center.Y + halfDim.Y}
	bottomRight = rotateAround(bottomRight, b.Center, b.Rotation)

	b.Edges = []Edge{
		{A: topLeft, B: topRight},
		{A: topRight, B: bottomRight},
		{A: bottomRight, B: bottomLeft},
		{A: bottomLeft, B: topLeft},
	}
	for i := range len(b.Edges) {
		b.Edges[i].Compute()
	}
}
