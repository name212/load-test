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
		Concurrent:        100000,
		DurationInSeconds: 25,
	}

	stats := test.Start()

	fmt.Println(stats)
}

type TestForRunner struct{}

func (t *TestForRunner) Run(workerId uint32) bool {
	min := int64(50)
	max := int64(1000)
	timeToSleep := rand.Int63n(max-min) + min
	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	isOk := rand.Int31n(2) > 0

	//fmt.Println("Worker ", workerId, " isOk = ", isOk, " after ", timeToSleep, "ms")
	return isOk
}
