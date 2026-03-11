package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Level4 struct {
	ui *ebitenui.UI
	rootContainer *widget.Container 
	loadDialog  *widget.Container
	saveDialog  *widget.Container
    saveInput   *widget.TextInput
	showingSave bool
	showingLoad bool
	Obstacles []*Solid
	Balls     []Ball
	TmpShape  *Solid
	Mode      string
	IsPaused  bool
}

func (l *Level4) HandleDrawingMode() error {
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

func (l *Level4) HandleBallCreation() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if l.Mode == "normal" {
			xMouse, yMouse := ebiten.CursorPosition()
			newBall := GetNewBall(xMouse, yMouse)
			l.Balls = append(l.Balls, newBall)
		}
	}
}


func (l *Level4) HandleQuit(){
	
	if l.Mode == "normal" && inpututil.IsKeyJustPressed(ebiten.KeyEscape) && !l.showingSave && !l.showingLoad{
		os.Exit(0)
	}
}


func (l *Level4) HandleSave(){
	if l.Mode == "normal" && inpututil.IsKeyJustPressed(ebiten.KeyS) && !l.showingSave{
		l.showSaveDialog()
	}
	if l.Mode == "normal" && inpututil.IsKeyJustPressed(ebiten.KeyEscape) && l.showingSave{
		l.hideSaveDialog()
	}
}


func (l *Level4) HandleLoad(){
	if inpututil.IsKeyJustPressed(ebiten.KeyL) && !l.showingLoad && !l.showingSave {
    	l.showLoadDialog()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) && l.showingLoad {
		l.hideLoadDialog()
	}
}

func (l *Level4) HandleCollision(dt float32) {
	// now we need to update the balls
	for ballIndex := range len(l.Balls) {
		l.Balls[ballIndex].Update(dt)
	}

	// collision with obstacles
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
						aToBallEdge := Edge{A: edge.A, B: l.Balls[ballIndex].Position}
						aToBallEdge.Compute()
						l.Balls[ballIndex].BounceBack(aToBallEdge.Normalized)
						continue
					}

				} else if distanceAlongEdge > edge.Length {
					if DistanceSquared(l.Balls[ballIndex].Position, edge.B) < l.Balls[ballIndex].Radius*l.Balls[ballIndex].Radius {
						// collision!
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

	// collisions with other balls
	for ballAIndex := range len(l.Balls) {
		for ballBIndex := ballAIndex; ballBIndex < len(l.Balls); ballBIndex++ {
			if ballAIndex == ballBIndex {
				continue
			}
			ballToBallEdge := Edge{A: l.Balls[ballAIndex].Position, B: l.Balls[ballBIndex].Position}
			ballToBallEdge.Compute()
			overlap := (l.Balls[ballAIndex].Radius + l.Balls[ballBIndex].Radius) - ballToBallEdge.Length
			if overlap > 0 {
				// Collision!
				// I need to move the balls by half the overlap along their normal
				halfOverlap := overlap / 2
				overlapVector := ScaleVector(ballToBallEdge.Normalized, halfOverlap)
				l.Balls[ballAIndex].Position.Sub(overlapVector)
				l.Balls[ballBIndex].Position.Add(overlapVector)

				relativeVelocity := SubVectors(l.Balls[ballBIndex].Speed, l.Balls[ballAIndex].Speed)
				approachSpeed := relativeVelocity.Dot(ballToBallEdge.Normal)
				if approachSpeed > 0 {
					// the balls are moving towards from each other!
					ballASpeed := l.Balls[ballAIndex].Speed
					l.Balls[ballAIndex].Speed = l.Balls[ballBIndex].Speed
					l.Balls[ballBIndex].Speed = ballASpeed
				}
			}
		}
	}

}

func (l *Level4) Update() error {
	l.HandleQuit()
	l.ui.Update() 
	l.HandleSave()
	l.HandleLoad()
	
	if l.showingSave{
		return nil
	}
	
	l.HandleBallCreation()

	err := l.HandleDrawingMode()
	if err != nil {
		return err
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		l.IsPaused = !l.IsPaused
	}
	// substeps are used to avoid having balls "teleporting" through walls
	nbSubsteps := 10
	dt := float32(1) / float32(nbSubsteps)
	if l.IsPaused {
		dt = 0
	}
	for _ = range nbSubsteps {
		l.HandleCollision(dt)
	}

	return nil
}

func (l *Level4) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 128, G: 128, B: 128})
	for _, obst := range l.Obstacles {
		obst.Draw(screen)
	}
	if l.TmpShape != nil {
		l.TmpShape.Draw(screen)
	}
	for _, ball := range l.Balls {
		ball.Draw(screen)
	}
	l.ui.Draw(screen) 
	ebitenutil.DebugPrint(screen, fmt.Sprintf("mode: %s\nnb balls: %d", l.Mode, len(l.Balls)))
}

func (l *Level4) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func Level4Init() *Level4 {

	game := Level4{
		Mode: "normal",
	}

	game.rootContainer = widget.NewContainer(
        widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
    )

    game.ui = &ebitenui.UI{Container: game.rootContainer}

	game.buildSaveDialog()

	return &game

}
