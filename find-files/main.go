package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	matchFiles []string
	wg         = sync.WaitGroup{}
	lock       = sync.Mutex{}
)

func findFiles(root, filename string) {
	fmt.Println("Searching in", root)
	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock()
			matchFiles = append(matchFiles, filepath.Join(root, file.Name()))
			lock.Unlock()
		}
		if file.IsDir() {
			wg.Add(1)
			go findFiles(filepath.Join(root, file.Name()), filename)
		}
	}
	wg.Done()
}
func main() {
	startTime := time.Now()
	wg.Add(1)
	go findFiles("D:/desktop/attack_research/source", "readme")
	wg.Wait()
	spendTime := time.Since(startTime)
	for _, file := range matchFiles {
		fmt.Println("Matched", file)
	}
	fmt.Println("Spend time: ", spendTime) //单线程264.ms   多线程100ms
}
