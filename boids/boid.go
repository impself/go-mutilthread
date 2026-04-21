package main

import (
	"math"
	"math/rand"
	"time"
)

// 想象在人群中走路：
// 对齐 = 大家都往同一个方向走
// 聚合 = 跟着大部队走，别掉队
// 分离 = 别tm撞到人身上！保持社交距离！
type boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgPosition, avgVelocity, seperation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0
	lock.RLock() //加读锁   可以多个同时读，单个写
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherboid := boidsMap[int(i)][int(j)]; otherboid != -1 && otherboid != b.id {
				if dist := boids[otherboid].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherboid].velocity)                                     //当前boid的范围内平均速度
					avgPosition = avgPosition.Add(boids[otherboid].position)                                     //当前boid的范围内平均位置
					seperation = seperation.Add(b.position.Substract(boids[otherboid].position).DivisionV(dist)) //当前boid的范围内，与其他boid的距离（距离越远权重越小）
				}
			}
		}
	}
	lock.RUnlock()
	accel := Vector2D{b.borderBounce(b.position.x, screenWidth),
		b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		avgVelocity = avgVelocity.DivisionV(count)                       //往邻居的中心靠拢
		avgPosition = avgPosition.DivisionV(count)                       //跟邻居的速度方向保持一致
		accelVel := avgVelocity.Substract(b.velocity).MultiplyV(adjRate) //b与中心点之间的速度向量（加上0.015权重，变化更平缓）
		accelPos := avgPosition.Substract(b.position).MultiplyV(adjRate) //b与中心点之间的位置向量
		accelSeperation := seperation.MultiplyV(adjRate)                 //分离产生的加速度（防止鸟群距离太近）
		accel = accel.Add(accelPos).Add(accelVel).Add(accelSeperation)
	}
	return accel
}

func (b *boid) borderBounce(pos, screen float64) float64 { //判断边界，防止碰墙
	if pos < viewRadius {
		return 1 / math.Max(pos, 1.0)
	} else if pos > screen-viewRadius {
		return 1 / math.Min(pos-screen, -1.0)
	}
	return 0
}

func (b *boid) moveOne() {
	calAcc := b.calcAcceleration() //加速度
	lock.Lock()                    //写锁，只能单个读写
	b.velocity = b.velocity.Add(calAcc.LimitMax(1))
	boidsMap[int(b.position.x)][int(b.position.y)] = -1   //移动前的bird赋值-1
	b.position = b.position.Add(b.velocity)               //更新位置
	boidsMap[int(b.position.x)][int(b.position.y)] = b.id //移动后更新bird位置索引
	lock.Unlock()
}

func (b *boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func createBoid(i int) {
	b := &boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight},
		velocity: Vector2D{rand.Float64()*2 - 1, rand.Float64()*2 - 1},
		id:       i,
	}
	boids[i] = b
	boidsMap[int(b.position.x)][int(b.position.y)] = b.id //写入位置信息
	go b.start()
}
