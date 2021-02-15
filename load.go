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

type FinishStrategyConst func(uint32) FinishTestStrategy

type LoadTest struct {
	Constructor         TestRunnerConstructor
	Concurrent          uint32
	DurationInSeconds   uint32
	FinishStrategyConst FinishStrategyConst
}

func (t *LoadTest) Start() *ResultStats {
	waitGroup := sync.WaitGroup{}

	var workers = make([]Worker, t.Concurrent)
	var allStats = make([]*workerStats, t.Concurrent)
	finishStrategyConstructor := t.FinishStrategyConst
	if finishStrategyConstructor == nil {
		finishStrategyConstructor = GetOneChanStrategy
	}
	finishStrategy := finishStrategyConstructor(t.Concurrent)

	for i := uint32(0); i < t.Concurrent; i++ {
		waitGroup.Add(1)
		stats := initStats()
		doneChan := finishStrategy.getFinishChan()
		worker := Worker{
			workersGroup: &waitGroup,
			id:           i + 1,
			stats:        stats,
			runner:       t.Constructor(),
			doneSigChain: doneChan,
		}
		workers[i] = worker
		allStats[i] = stats
		go worker.run()
	}

	time.Sleep(time.Duration(t.DurationInSeconds) * time.Second)

	finishStrategy.finishTesting()

	fmt.Println("Send done flag to all workers")
	waitGroup.Wait()
	fmt.Println("All workers finished")

	return constructTotalStats(allStats)
}
