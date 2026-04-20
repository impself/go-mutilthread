package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 500
	viewRadius                = 13
	adjRate                   = 0.015
)

var (
	green    = color.RGBA{10, 255, 50, 255}
	boids    [boidCount]*boid
	boidsMap [screenWidth + 1][screenHeight + 1]int //所有boid的位置索引  填充值为boid.id
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, boid := range boids {
		screen.Set(int(boid.position.x+1), int(boid.position.y), green)
		screen.Set(int(boid.position.x-1), int(boid.position.y), green)
		screen.Set(int(boid.position.x), int(boid.position.y+1), green)
		screen.Set(int(boid.position.x), int(boid.position.y-1), green)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	for i, row := range boidsMap {
		for j := range row {
			boidsMap[i][j] = -1
		}
	} //预先填充所有二维数组为-1
	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Boids")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
