package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Ball struct {
	Radius   float32
	Position Vector[float32]
	Speed    Vector[float32]
	Color color.RGBA
}

func (b *Ball) Update(dt float32) {
	scaledUpdate := ScaleVector(b.Speed, dt)
	b.Position.Add(scaledUpdate)
}

// BounceBack reverses the velocity along the perpendicular axis from the edge, leaving the parallell axis intact
// the formula is: new velocity = velocity - 2 (velocity . unit normal) * unit normal
func (b *Ball) BounceBack(unitNormal Vector[float32]) {
	dotVelocity := b.Speed.X*unitNormal.X + b.Speed.Y*unitNormal.Y
	b.Speed.X = b.Speed.X - 2 * dotVelocity * unitNormal.X
	b.Speed.Y = b.Speed.Y - 2 * dotVelocity * unitNormal.Y
}

func (b *Ball) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, b.Position.X, b.Position.Y, b.Radius, b.Color, true)
}