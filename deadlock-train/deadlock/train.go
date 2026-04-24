package deadlock

import (
	"time"

	. "github.com/inmself/deadlock-train/common"
)

func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, cross := range crossings {
			if train.Front == cross.Position {
				cross.Intersection.Mutex.Lock()
				cross.Intersection.LockedBy = train.Id
			}
			back := train.Front - train.TrainLength
			if back == cross.Position {
				cross.Intersection.Mutex.Unlock()
				cross.Intersection.LockedBy = -1
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
