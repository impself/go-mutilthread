package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *boid) calcAcceleration() Vector2D {
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius)
	avgVelocity := Vector2D{0, 0}
	count := 0.0
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherboid := boidsMap[int(i)][int(j)]; otherboid != -1 && otherboid != b.id {
				if dist := boids[otherboid].position.Distance(b.position); dist < viewRadius {
					count++
					avgVelocity = avgVelocity.Add(boids[otherboid].velocity)
				}
			}
		}
	}
	accel := Vector2D{0, 0}
	fmt.Println(count)
	if count > 0 {
		avgVelocity = avgVelocity.DivisionV(count)
		accel = avgVelocity.Substract(b.velocity).MultiplyV(adjRate)
	}
	return accel
}

func (b *boid) moveOne() {
	b.velocity = b.velocity.Add(b.calcAcceleration().limit(-1, 1))
	boidsMap[int(b.position.x)][int(b.position.y)] = -1   //移动前的bird赋值-1
	b.position = b.position.Add(b.velocity)               //更新位置
	boidsMap[int(b.position.x)][int(b.position.y)] = b.id //移动后更新bird位置索引
	next := b.position.Add(b.velocity)                    //下一次移动时会不会触碰边界，如果有，则反弹
	if next.x >= screenWidth || next.x < 0 {
		b.velocity = Vector2D{-b.velocity.x, b.velocity.y}
	}
	if next.y >= screenHeight || next.y < 0 {
		b.velocity = Vector2D{b.velocity.x, -b.velocity.y}
	}
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
