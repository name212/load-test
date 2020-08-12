package loadtest

import (
	"fmt"
	"sync"
	"time"
)

type TestRunner interface {
	Run(workerId uint32) bool
}

type TestRunnerConstructor func() TestRunner

type LoadTest struct {
	Constructor       TestRunnerConstructor
	Concurrent        uint32
	DurationInSeconds uint32
}

func (t *LoadTest) Start() *ResultStats {
	waitGroup := sync.WaitGroup{}
	doneChans := make([]chan bool, t.Concurrent)

	var workers = make([]Worker, t.Concurrent)
	var allStats = make([]*workerStats, t.Concurrent)

	for i := uint32(0); i < t.Concurrent; i++ {
		waitGroup.Add(1)
		stats := initStats()
		doneChan := make(chan bool, 1)
		worker := Worker{
			workersGroup: &waitGroup,
			id:           i + 1,
			stats:        stats,
			runner:       t.Constructor(),
			doneSigChain: doneChan,
		}
		workers[i] = worker
		allStats[i] = stats
		doneChans[i] = doneChan
		go worker.run()
	}

	time.Sleep(time.Duration(t.DurationInSeconds) * time.Second)

	for _, doneChan := range doneChans {
		doneChan <- true
	}
	fmt.Println("Send done flag to all workers")
	waitGroup.Wait()
	fmt.Println("All workers finished")
	return constructTotalStats(allStats)
}
