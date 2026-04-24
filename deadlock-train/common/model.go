package common

import "sync"

type Train struct {
	Id          int
	TrainLength int
	Front       int
}
type Intersection struct {
	Id       int        //ID
	Mutex    sync.Mutex //加锁
	LockedBy int        //被哪个火车抢占
}
type Crossing struct {
	Position     int
	Intersection *Intersection
}
