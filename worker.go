package loadtest

import (
	"sync"
	"time"
)

type Worker struct {
	doneSigChain chan bool
	id           uint32
	stats        *workerStats
	runner       TestRunner
	workersGroup *sync.WaitGroup
}

func (w *Worker) run() {
	defer w.workersGroup.Done()

	for {
		select {
		case <-w.doneSigChain:
			//fmt.Println("Getting done signal", w.id)
			w.stats.endTest()
			return
		default:
			start := curTime()
			isOk := w.runner.Run(w.id)
			runTime := curTime() - start
			w.stats.affectStat(isOk, time.Duration(runTime))
		}
	}
}
