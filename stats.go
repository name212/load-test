package loadtest

import (
	"math"
	"time"
)

type BaseStats struct {
	StartTime int64
	EndTime   int64
	RunsCount uint64
	OkRuns    uint64
	MinTime   time.Duration
	MaxTime   time.Duration
}

type WorkerStats struct {
	BaseStats
	cumulativeTime time.Duration
}

type ResultStats struct {
	BaseStats
	AvgTime   time.Duration
	OkPercent float64
}

func (stats *WorkerStats) affectStat(isOk bool, runTime time.Duration) {
	stats.RunsCount += 1

	if isOk {
		stats.OkRuns += 1
	}

	stats.cumulativeTime += runTime
	if runTime > stats.MaxTime {
		stats.MaxTime = runTime
	}

	if runTime < stats.MinTime {
		stats.MinTime = runTime
	}

}

func (stats *WorkerStats) endTest() {
	stats.EndTime = curTime()
}

func initStats() *WorkerStats {
	stats := WorkerStats{
		cumulativeTime: 0,
	}

	stats.StartTime = curTime()
	stats.EndTime = 0
	stats.RunsCount = 0
	stats.OkRuns = 0
	stats.MinTime = math.MaxInt64
	stats.MaxTime = 0

	return &stats
}

func constructTotalStats(allStats []*WorkerStats) *ResultStats {
	totalStats := ResultStats{}

	totalTime := time.Duration(0)

	totalStats.MinTime = time.Duration(math.MaxInt64)
	totalStats.MaxTime = time.Duration(0)

	totalStats.StartTime = int64(math.MaxInt64)
	totalStats.EndTime = int64(0)

	for _, stats := range allStats {
		if stats.MinTime < totalStats.MinTime {
			totalStats.MinTime = stats.MinTime
		}

		if stats.MaxTime > totalStats.MaxTime {
			totalStats.MaxTime = stats.MaxTime
		}

		if stats.StartTime < totalStats.StartTime {
			totalStats.StartTime = stats.StartTime
		}

		if stats.EndTime > totalStats.EndTime {
			totalStats.EndTime = stats.EndTime
		}

		totalTime += stats.cumulativeTime
		totalStats.RunsCount += stats.RunsCount
		totalStats.OkRuns += stats.OkRuns
	}

	if totalStats.RunsCount > 0 {
		totalStats.OkPercent = float64(totalStats.OkRuns) / float64(totalStats.RunsCount)
		totalStats.OkPercent *= 100
	}

	avg := float64(totalTime) / float64(totalStats.RunsCount)
	totalStats.AvgTime = time.Duration(math.Round(avg))
	return &totalStats
}
