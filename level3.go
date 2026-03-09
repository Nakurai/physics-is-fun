package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Drawable interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Level3 struct {
	Balls        []Ball
	Obstacles    []*Solid
	Mode         string
	TmpShape     *Solid
	NbCollisions int
	IsRecording  bool
}

func (l *Level3) HandleDrawingMode() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyControlLeft) {
		l.Mode = "drawing"
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyControlLeft) {
		if l.Mode == "drawing" {
			// we close the shape
			edge := l.TmpShape.GetLastEdge()
			if edge != nil {
				lastEdge := Edge{A: edge.B, B: l.TmpShape.Edges[0].A}
				l.TmpShape.Edges = append(l.TmpShape.Edges, lastEdge)
			}
			// we add the new shape to the obstacles and reinitialize the tmp shape
			l.TmpShape.ComputeEdges()
			l.Obstacles = append(l.Obstacles, l.TmpShape)
			l.TmpShape = nil
			l.Mode = "normal"
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if l.Mode == "drawing" {
			if l.TmpShape == nil {
				newSolid := Solid{Edges: []Edge{}, Corners: [4]Vector[float32]{}}
				l.TmpShape = &newSolid
			}
			xMouse, yMouse := ebiten.CursorPosition()
			edge := l.TmpShape.GetLastEdge()
			if edge == nil {
				// if there is no previous edge, then the only edge we know starts and ends at the current mouse position
				edge = &Edge{A: Vector[float32]{X: float32(xMouse), Y: float32(yMouse)}, B: Vector[float32]{X: float32(xMouse), Y: float32(yMouse)}}
				l.TmpShape.Edges = append(l.TmpShape.Edges, *edge)
			} else {
				// otherwise, the new edge starts where the last one ended and ends at the current mouse position
				newEdge := Edge{A: edge.B, B: Vector[float32]{X: float32(xMouse), Y: float32(yMouse)}}
				l.TmpShape.Edges = append(l.TmpShape.Edges, newEdge)
			}

		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if l.Mode == "drawing" {
			xMouse, yMouse := ebiten.CursorPosition()
			edge := l.TmpShape.GetLastEdge()
			if edge != nil {
				edge.B = Vector[float32]{X: float32(xMouse), Y: float32(yMouse)}
			}
		}
	}

	return nil
}

func (l *Level3) HandleRecording() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if l.IsRecording {
			l.StopRecording()
			l.IsRecording = false
		} else {
			widthWindow, heightWindow := ebiten.WindowSize()
			l.StartRecording(widthWindow, heightWindow)
			l.IsRecording = true
		}
	}

	return nil
}

func (l *Level3) Update() error {

	// if the user is currently drawing, we handle it
	err := l.HandleDrawingMode()
	if err != nil {
		return err
	}
	err = l.HandleRecording()
	if err != nil {
		return err
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if l.Mode == "normal" {
			xMouse, yMouse := ebiten.CursorPosition()
			newBall := Ball{
				Radius:   5,
				Position: Vector[float32]{X: float32(xMouse), Y: float32(yMouse)},
				Speed:    Vector[float32]{X: 200, Y: 200},
			}
			l.Balls = append(l.Balls, newBall)
		}
	}

	// substeps are used to avoid having balls "teleporting" through walls
	l.NbCollisions = 0
	nbSubsteps := 100
	for _ = range nbSubsteps {
		// now we need to update the balls
		for ballIndex := range len(l.Balls) {
			l.Balls[ballIndex].Update(float32(1) / float32(nbSubsteps))
		}

		for ballIndex := range len(l.Balls) {
			for _, obstacle := range l.Obstacles {
				for _, edge := range obstacle.Edges {
					edgeBallVector := SubVectors(l.Balls[ballIndex].Position, edge.A)

					tangentToNormal := TangentVector(edge.Normal, 1)
					distanceAlongEdge := edgeBallVector.Dot(tangentToNormal)

					// distance along the edge is giving us how far along the edge the ball is. Depending on where this is, we need to check if the ball is colliding or not
					if distanceAlongEdge < 0 {
						// here the ball is before the edge. All we need is to check the distance between the ball's position and the edge's point
						if DistanceSquared(l.Balls[ballIndex].Position, edge.A) < l.Balls[ballIndex].Radius*l.Balls[ballIndex].Radius {
							// collision!
							l.NbCollisions += 1
							aToBallEdge := Edge{A: edge.A, B: l.Balls[ballIndex].Position}
							aToBallEdge.Compute()
							l.Balls[ballIndex].BounceBack(aToBallEdge.Normalized)
							continue
						}

					} else if distanceAlongEdge > edge.Length {
						if DistanceSquared(l.Balls[ballIndex].Position, edge.B) < l.Balls[ballIndex].Radius*l.Balls[ballIndex].Radius {
							// collision!
							l.NbCollisions += 1
							bToBallEdge := Edge{A: edge.B, B: l.Balls[ballIndex].Position}
							bToBallEdge.Compute()
							l.Balls[ballIndex].BounceBack(bToBallEdge.Normalized)
							continue
						}

					} else {

						distanceToEdge := edgeBallVector.Dot(edge.Normal)
						absDistanceToEdge := distanceToEdge
						if absDistanceToEdge < 0 {
							absDistanceToEdge *= -1
						}
						if absDistanceToEdge < l.Balls[ballIndex].Radius {
							// collision!
							l.NbCollisions += 1
							bounceNormal := edge.Normal
							if distanceToEdge < 0 {
								bounceNormal = ScaleVector(bounceNormal, -1)
							}
							// replacing the ball at the border
							overlap := l.Balls[ballIndex].Radius - absDistanceToEdge
							l.Balls[ballIndex].Position.X += bounceNormal.X * overlap
							l.Balls[ballIndex].Position.Y += bounceNormal.Y * overlap

							// bouncing the ball back only if it is not already moving away
							relativeVelocity := l.Balls[ballIndex].Speed.Dot(bounceNormal)
							if relativeVelocity < 0 {
								// the ball is moving toward the ball!
								l.Balls[ballIndex].BounceBack(bounceNormal)
							}
							continue
						}
					}

				}
			}
		}
	}

	return nil
}

func (l *Level3) Draw(screen *ebiten.Image) {
	for _, obst := range l.Obstacles {
		obst.Draw(screen)
	}
	if l.TmpShape != nil {
		l.TmpShape.Draw(screen)
	}
	for _, ball := range l.Balls {
		ball.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("mode: %s\nrecording: %v\nnb balls: %d", l.Mode, l.IsRecording, len(l.Balls)))
}

func (l *Level3) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func Level3Init() *Level3 {
	g := Level3{}
	g.Mode = "normal"

	return &g
}
