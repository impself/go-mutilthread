package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	MatrixSize  = 300 // 矩阵维度 300x300
	totalRounds = 100 // 循环计算100轮
)

var (
	matrixA      = [MatrixSize][MatrixSize]int{} // 输入矩阵A
	matrixB      = [MatrixSize][MatrixSize]int{} // 输入矩阵B
	result       = [MatrixSize][MatrixSize]int{} // 结果矩阵 C = A * B
	workStart    = NewBarrier(MatrixSize + 1)
	workComplete = NewBarrier(MatrixSize + 1)
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
	workStart.Wait()
	for col := 0; col < MatrixSize; col++ {
		for i := 0; i < MatrixSize; i++ {
			result[row][col] += matrixA[row][i] * matrixB[i][col]
		}
	}
	workComplete.Wait()
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
	// 启动 MatrixSize 个 goroutine，每个负责计算一行

	start := time.Now()

	var wg sync.WaitGroup
	for row := 0; row < MatrixSize; row++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := 0; i < totalRounds; i++ {
				workOutRow(r)
			}
		}(row)
	}

	for i := 0; i < totalRounds; i++ {
		// 每轮：清零结果 → 生成随机输入 → 按行并发计算矩阵乘法
		cleanResult(&result)
		GenerateRandomMatrix(&matrixA)
		GenerateRandomMatrix(&matrixB)
		workStart.Wait()
		workComplete.Wait()

	}
	wg.Wait()
	since := time.Since(start)
	fmt.Println("花费时间:  ", since)
	fmt.Println("done...")
}
