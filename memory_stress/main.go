package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("START MEMORY STRESS TEST")

	duration := 60 //s
	startTime := time.Now()

	for time.Since(startTime).Seconds() < float64(duration) {
		loadMemory()
	}

	fmt.Println("END MEMORY STRESS TEST")
}

func loadMemory() {
	const size = 100 * 1024 * 1024 // 100 MB
	for {
		_ = make([]byte, size)
		_ = make([]byte, 0)
	}
}
