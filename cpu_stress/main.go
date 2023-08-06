package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	fmt.Println("START CPU STRESS TEST")

	duration := 60 //s
	startTime := time.Now()

	for time.Since(startTime).Seconds() < float64(duration) {
		calculateSin()
	}

	fmt.Println("END CPU STRESS TEST")
}

func calculateSin() {
	iterations := 10000000

	for i := 0; i < iterations; i++ {
		_ = math.Sin(float64(i))
	}
}
