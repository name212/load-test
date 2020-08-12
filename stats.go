package loadtest

import (
	"fmt"
	"math"
	"time"
)

type BaseStats struct {
	StartTime time.Time
	EndTime   time.Time
	RunsCount uint64
	OkRuns    uint64
	MinTime   time.Duration
	MaxTime   time.Duration
}

type workerStats struct {
	BaseStats
	cumulativeTime time.Duration
}

type ResultStats struct {
	BaseStats
	AvgTime   time.Duration
	OkPercent float64
}

func (stats *workerStats) affectStat(isOk bool, runTime time.Duration) {
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

func (stats *workerStats) endTest() {
	stats.EndTime = time.Now()
}

func initStats() *workerStats {
	stats := workerStats{
		cumulativeTime: 0,
	}

	stats.StartTime = time.Now()
	stats.EndTime = time.Now()
	stats.RunsCount = 0
	stats.OkRuns = 0
	stats.MinTime = math.MaxInt64
	stats.MaxTime = 0

	return &stats
}

func constructTotalStats(allStats []*workerStats) *ResultStats {
	totalStats := ResultStats{}

	totalTime := time.Duration(0)

	totalStats.MinTime = time.Duration(math.MaxInt64)
	totalStats.MaxTime = time.Duration(0)

	totalStats.StartTime = time.Unix(0, 0)
	totalStats.EndTime = time.Now()

	for _, stats := range allStats {
		if stats.MinTime < totalStats.MinTime {
			totalStats.MinTime = stats.MinTime
		}

		if stats.MaxTime > totalStats.MaxTime {
			totalStats.MaxTime = stats.MaxTime
		}

		if totalStats.StartTime.Sub(stats.StartTime) < 0 {
			totalStats.StartTime = stats.StartTime
		}

		if stats.EndTime.Sub(totalStats.EndTime) > 0 {
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

func (s *ResultStats) String() string {
	fmStr := `
*** Test results ***

Runs %d: ok runs %d (%.2f%%)
Duration: 
	Avg = %s	
	Min = %s	
	Max = %s

Start time: %s
End time  : %s
Duration  : %s
`
	return fmt.Sprintf(fmStr,
		s.RunsCount, s.OkRuns, s.OkPercent,
		s.AvgTime,
		s.MinTime,
		s.MaxTime,

		s.StartTime,
		s.EndTime,
		s.EndTime.Sub(s.StartTime),
	)
}
