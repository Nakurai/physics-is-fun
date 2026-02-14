package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	blockSize = 16
)

type Block struct {
	x, y float32
}
type Ball struct {
	radius         float32
	x, y           float32
	speedX, speedY float32
}

type Level1 struct {
	blocks                     []Block
	balls                      []Ball
	topLeftX, topLeftY         int
	topRightX, topRightY       int
	bottomRightX, bottomRightY int
	bottomLeftX, bottomLeftY   int
	minSpeed, maxSpeed         int
	minRadius, maxRadius       int
	nbBalls                    int
}

func (g *Level1) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.AddBall()
	}

	for i := range len(g.balls) {
		// checking if the ball is touching the inner sides of the fence
		isBallTouchingLeftSide := g.balls[i].x-g.balls[i].radius <= float32(g.topLeftX)
		isBallMovingTowardsLeft := g.balls[i].speedX < 0
		isBallTouchingRightSide := g.balls[i].x+g.balls[i].radius >= float32(g.topRightX)
		isBallMovingTowardsRight := g.balls[i].speedX > 0
		if (isBallTouchingLeftSide && isBallMovingTowardsLeft) || (isBallTouchingRightSide && isBallMovingTowardsRight) {
			g.balls[i].speedX = -g.balls[i].speedX
		}
		isBallTouchingTopSide := g.balls[i].y-g.balls[i].radius <= float32(g.topLeftY)
		isBallMovingTowardsTop := g.balls[i].speedY < 0
		isBallTouchingBottomSide := g.balls[i].y+g.balls[i].radius >= float32(g.bottomLeftY)
		isBallMovingTowardsBottom := g.balls[i].speedY > 0
		if (isBallTouchingTopSide && isBallMovingTowardsTop) || (isBallTouchingBottomSide && isBallMovingTowardsBottom) {
			g.balls[i].speedY = -g.balls[i].speedY
		}
	}
	// Move the balls
	for i := range len(g.balls) {
		g.balls[i].x += g.balls[i].speedX
		g.balls[i].y += g.balls[i].speedY
	}

	return nil
}

func (g *Level1) Draw(screen *ebiten.Image) {
	for _, block := range g.blocks {
		vector.FillRect(screen, block.x, block.y, blockSize, blockSize, color.RGBA{0, 0, 255, 255}, true)
	}
	for _, ball := range g.balls {
		vector.FillCircle(screen, ball.x, ball.y, ball.radius, color.RGBA{255, 0, 0, 255}, true)
	}
}

func (g *Level1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Level1) AddBall() {
	randRadius := float32(rand.Intn(g.maxRadius-g.minRadius+1) + g.minRadius)

	minXLeft := g.topLeftX + int(randRadius)
	maxXRight := g.topRightX - int(randRadius)
	minYTop := g.topLeftY + int(randRadius)
	maxYBottom := g.bottomLeftY - int(randRadius)

	xRand := float32(rand.Intn(maxXRight-minXLeft+1) + minXLeft)
	yRand := float32(rand.Intn(maxYBottom-minYTop+1) + minYTop)

	randSpeedX := float32(rand.Intn(g.maxSpeed-g.minSpeed+1) + g.minSpeed)
	randSpeedY := float32(rand.Intn(g.maxSpeed-g.minSpeed+1) + g.minSpeed)
	g.balls = append(g.balls, Ball{
		radius: randRadius,
		x:      xRand,
		y:      yRand,
		speedX: randSpeedX,
		speedY: randSpeedY,
	})
	g.nbBalls++
}

func Level1Init() *Level1 {
	g := Level1{}
	offset := 150
	fenceWidth := 50
	fenceHeight := 50
	for i := 0; i < fenceWidth; i++ {
		for j := 0; j < fenceHeight; j++ {
			if i == 0 || i == fenceWidth-1 || j == 0 || j == fenceHeight-1 {
				b := Block{
					x: float32(offset + j*blockSize),
					y: float32(offset + i*blockSize),
				}
				g.blocks = append(g.blocks, b)
			}
		}
	}

	g.topLeftX = offset + blockSize
	g.topLeftY = offset + blockSize
	g.topRightX = offset + blockSize*(fenceWidth-1)
	g.topRightY = offset + blockSize
	g.bottomRightX = offset + blockSize*(fenceWidth-1)
	g.bottomRightY = offset + blockSize*(fenceHeight-1)
	g.bottomLeftX = offset + blockSize
	g.bottomLeftY = offset + blockSize*(fenceHeight-1)
	// the most top left possible position for the ball's center is at (offset+blockSize+radius, offsetoffset+blockSize+radius)
	// the most top right possible position for the ball's center is at (offset+blockSize*(fencewidth-1)-radius, offsetoffset+blockSize+radius)
	// the most bottom left possible position for the ball's center is at (offset+blockSize+radius, offset+blockSize*(fenceHeight-1)-radius)
	// the most bottom right possible position for the ball's center is at (offset+blockSize*(fencewidth-1)-radius, offset+blockSize*(fenceHeight-1)-radius)
	g.minSpeed = 3
	g.maxSpeed = 10
	g.minRadius = 10
	g.maxRadius = 20
	g.nbBalls = 0
	nbBallsInit := 3
	for i := 0; i < nbBallsInit; i++ {
		g.AddBall()
	}
	return &g
}
