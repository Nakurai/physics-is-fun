package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Level2 struct {
	Balls []Ball
	Fence Boundary
}

func (l *Level2) Update() error {

	// checking if the user pressed either the left or right arrow
	if ebiten.IsKeyPressed(ebiten.KeyLeft){
		l.Fence.Rotate(-1)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight){
		l.Fence.Rotate(1)
	}

	for i := range len(l.Balls) {
		l.Balls[i].Update()
		// collision detection
		for _, edge := range(l.Fence.Edges){
			// compute vector from any point of th eedge to the ball's center
			edgeBallVector := SubVectors(l.Balls[i].Position, edge.A)
			
			// calculate the dot product of the edge's normal vector and the vector above
			distance := edgeBallVector.Dot(edge.Normal)

			// collision occurs if the ball has entered the edge
			if distance < l.Balls[i].Radius {
				overlap := l.Balls[i].Radius - distance
				l.Balls[i].Position.X += edge.Normal.X * overlap
				l.Balls[i].Position.Y += edge.Normal.Y * overlap
				l.Balls[i].BounceBack(edge.Normal)
			}

		}
	}


	return nil
}
func (l *Level2) Draw(screen *ebiten.Image) {
	l.Fence.Draw(screen)
	for _, edge := range l.Fence.Edges {
		scaledNormal := ScaleVector(edge.Normal, 25)         // 25 is arbitrary, just to display something on the screen
		normalPoint := AddVectors(edge.Middle, scaledNormal) // applying the scaled direction to the middle point to have our destination
		vector.StrokeLine(screen, edge.Middle.X, edge.Middle.Y, normalPoint.X, normalPoint.Y, 1, Red, true)
	}
	for i := range len(l.Balls) {
		l.Balls[i].Draw(screen)
	} 
}

func (g *Level2) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (l *Level2) AddBall(ball Ball) {

}

func Level2Init() *Level2 {
	g := Level2{}
	width, length := ebiten.WindowSize()
	g.Fence = Boundary{
		Center:  Vector[float32]{X: float32(width / 2)+20, Y: float32(length / 2)+20},
		Dimension: Vector[float32]{X: 500, Y: 500},
	}
	g.Fence.ComputeEdges()

	g.Balls = append(g.Balls, Ball{
		Radius: 25,
		Position: Vector[float32]{X: float32(width / 2), Y: float32(length / 2)},
		Speed: Vector[float32]{X: 5, Y: 2},
	})

	return &g
}
