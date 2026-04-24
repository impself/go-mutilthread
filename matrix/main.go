package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MatrixSize  = 300  // 矩阵维度 300x300
	totalRounds = 100  // 循环计算100轮
)

var (
	matrixA = [MatrixSize][MatrixSize]int{} // 输入矩阵A
	matrixB = [MatrixSize][MatrixSize]int{} // 输入矩阵B
	result  = [MatrixSize][MatrixSize]int{} // 结果矩阵 C = A * B
)

// GenerateRandomMatrix 用 [-5, 4] 范围的随机整数填充矩阵
func GenerateRandomMatrix(matrix *[MatrixSize][MatrixSize]int) {
	for row := 0; row < MatrixSize; row++ {
		for col := 0; col < MatrixSize; col++ {
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

// workOutRow 计算结果矩阵的第row行: result[row][*] = matrixA[row][*] * matrixB[*][*]
func workOutRow(row int) {
	for col := 0; col < MatrixSize; col++ {
		for i := 0; i < MatrixSize; i++ {
			result[row][col] += matrixA[row][i] * matrixB[i][col]
		}
	}
}

// cleanResult 将结果矩阵所有元素清零（+= 累加前必须清零）
func cleanResult(matrix *[MatrixSize][MatrixSize]int) {
	for row := 0; row < MatrixSize; row++ {
		for col := 0; col < MatrixSize; col++ {
			matrix[row][col] = 0
		}
	}
}

func main() {
	fmt.Println("start...")
	start := time.Now()

	for i := 0; i < totalRounds; i++ {
		// 每轮：清零结果 → 生成随机输入 → 按行并发计算矩阵乘法
		cleanResult(&result)
		GenerateRandomMatrix(&matrixA)
		GenerateRandomMatrix(&matrixB)

		// 启动 MatrixSize 个 goroutine，每个负责计算一行
		var wg sync.WaitGroup
		wg.Add(MatrixSize)
		for row := 0; row < MatrixSize; row++ {
			go func(r int) {
				workOutRow(r)
				wg.Done()
			}(row)
		}
		wg.Wait() // 等待所有行计算完成再进入下一轮
	}

	since := time.Since(start)
	fmt.Println("花费时间:  ", since)
	fmt.Println("done...")
}
