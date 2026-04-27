package main

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten"
	. "github.com/inmself/deadlock-train/arbitrator"
	. "github.com/inmself/deadlock-train/common"
)

var (
	trains        [4]*Train
	intersections [4]*Intersection
)

const TrainLength = 70

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		DrawTracks(screen)
		DrawIntersections(screen)
		DrawTrains(screen)

	}
	return nil
}

func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &Train{Id: i, TrainLength: TrainLength, Front: 0}
	}
	for i := 0; i < 4; i++ {
		intersections[i] = &Intersection{Id: i, Mutex: sync.Mutex{}, LockedBy: -1}
	}
	for i := 0; i < 4; i++ {
		go MoveTrain(trains[i], 300, []*Crossing{{Position: 125, Intersection: intersections[i]},
			{Position: 175, Intersection: intersections[(i+1)%4]}})
	}

	if err := ebiten.Run(update, 320, 320, 3, "Trains in a box "); err != nil {
		log.Fatal("err")
	}
}
