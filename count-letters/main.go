package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var lock = sync.Mutex{}

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	rsp, _ := http.Get(url)
	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	for i := 0; i < 20; i++ {
		for _, b := range body {
			c := strings.ToLower(string(b))
			index := strings.Index(allLetters, c)
			if index >= 0 {
				atomic.AddInt32(&frequency[index], 1)
			}
		}
	}
	wg.Done()

}
func main() {
	var frequency [26]int32
	var wg = sync.WaitGroup{}
	start := time.Now()

	for i := 1000; i <= 1200; i++ {
		wg.Add(1)
		go countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}
	wg.Wait()
	end := time.Since(start)
	fmt.Printf("processing took %s\n", end)
	for i, f := range frequency {
		fmt.Printf("%s -> %d\n", string(allLetters[i]), f)
	}
}
