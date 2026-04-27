package hierarchy

import (
	"sort"
	"time"

	. "github.com/inmself/deadlock-train/common"
)

// lockIntersectionsInDistance 锁定指定范围内的所有交叉口
// 参数说明：
//
//	id: 列车编号，用于标记锁的所有者
//	reserveStart: 预留范围起点（车头当前位置）
//	resverEnd: 预留范围终点（车头位置 + 车长）
//	crossings: 所有交叉口信息列表
func lockIntersectionsInDistance(id, reserveStart, resverEnd int, crossings []*Crossing) {
	// 收集需要加锁的交叉口（范围内且未被当前列车锁定的）
	var intersectionsToLock []*Intersection
	for _, crossing := range crossings {
		if reserveStart <= crossing.Position && resverEnd >= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}
	// 按交叉口ID从小到大排序，确保所有列车都按统一顺序加锁
	// 这样可以预防死锁（避免循环等待）
	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].Id < intersectionsToLock[j].Id
	})
	// 依次加锁各个交叉口
	for _, it := range intersectionsToLock {
		it.Mutex.Lock()
		it.LockedBy = id
		time.Sleep(time.Millisecond)
	}
}

// MoveTrain 移动列车，模拟列车在轨道上行进
// 参数说明：
//
//	train: 列车对象，包含位置、长度等信息
//	distance: 目标移动距离
//	crossings: 所有交叉口信息列表
func MoveTrain(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1
		for _, cross := range crossings {
			// 当列车车头到达交叉口时，预锁后续需要的交叉口
			if train.Front == cross.Position {
				lockIntersectionsInDistance(train.Id, cross.Position, cross.Position+train.TrainLength, crossings)
			}
			// 当列车车尾离开交叉口时，解锁
			back := train.Front - train.TrainLength
			if back == cross.Position {
				cross.Intersection.Mutex.Unlock()
				cross.Intersection.LockedBy = -1
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
}
