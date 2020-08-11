package main

import (
	"fmt"
	"github.com/name212/loadtest"
	"math/rand"
	"time"
)

func main() {
	test := loadtest.LoadTest{
		Constructor: func() loadtest.TestRunner {
			return &TestForRunner{}
		},
		Concurrent:        50000,
		DurationInSeconds: 10,
	}

	stats := test.Start()

	fmt.Printf("\nTest resuls\n\n")
	fmt.Printf("Runs %d: ok runs %d (%f%%)\n", stats.RunsCount, stats.OkRuns, stats.OkPercent)
	fmt.Printf(
		"Duration: Avg=%fms\tMin=%fms\tMax=%fms",
		float64(stats.AvgTime/time.Millisecond),
		float64(stats.MinTime/time.Millisecond),
		float64(stats.MaxTime/time.Millisecond),
	)
}

type TestForRunner struct{}

func (t *TestForRunner) Run(workerId uint32) bool {
	min := int64(50)
	max := int64(1000)
	timeToSleep := rand.Int63n(max-min) + min
	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	isOk := rand.Int31n(2) > 0

	fmt.Println("Worker ", workerId, " isOk = ", isOk, " after ", timeToSleep, "ms")
	return isOk
}
