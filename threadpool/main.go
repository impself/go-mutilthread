package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point struct {
	x int
	y int
}

const threadNum int = 16

var (
	str2Point = regexp.MustCompile(`\((\d*),(\d*)\)`)
	wg        = sync.WaitGroup{}
)

func findArea(inputChannle chan string) {
	for line := range inputChannle {
		var points []Point
		for _, p := range str2Point.FindAllStringSubmatch(line, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point{x: x, y: y})
		}
		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)]
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		//fmt.Println(math.Abs(area) / 2.0)
	}
	wg.Done()
}
func main() {
	absPath, _ := filepath.Abs("../threadpool")
	file, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	text := string(file)
	inputChannel := make(chan string, 1000)
	wg.Add(threadNum) //增加8个线程
	for i := 0; i < threadNum; i++ {
		go findArea(inputChannel)
	}
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		inputChannel <- line
	}
	close(inputChannel)
	wg.Wait()
	sinceTime := time.Since(start)
	fmt.Println(sinceTime)
}
