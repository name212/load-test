package loadtest

import (
	"time"
)

func curTime() int64 {
	return time.Now().UnixNano()
}
