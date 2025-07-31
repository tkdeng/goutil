package goutil

import (
	"sync"
	"time"
)

// Interval runs a callback in an event loop (similar to JavaScript)
type Interval struct {
	cb      func() bool
	ms      time.Duration
	lastRun int64
	id      uintptr
	mu      sync.Mutex
	stop    bool
}

var intervals []*Interval
var intervalMU sync.Mutex

func init() {
	go func() {
		for {
			time.Sleep(time.Millisecond)
			now := time.Now().UnixMilli()

			intervalMU.Lock()

			for _, interval := range intervals {
				if interval.stop {
					intervals = append(intervals[:interval.id], intervals[interval.id+1:]...)
					continue
				}

				if now-interval.lastRun < interval.ms.Milliseconds() {
					continue
				}
				interval.lastRun = now

				go func() {
					interval.mu.Lock()
					defer interval.mu.Unlock()

					if !interval.stop {
						runAgain := interval.cb()
						if !runAgain {
							interval.stop = true
						}
					}
				}()
			}

			intervalMU.Unlock()
		}
	}()
}

// SetInterval runs a callback in an event loop (similar to JavaScript).
//
// This may be useful for concurrently running low priority tasks without creating multiple goroutines.
//
// @cb return:
//   - true: continue interval
//   - false: break interval
//
// @ms: time in nanoseconds (minimum of 1 millisecond)
func SetInterval(cb func() bool, ms time.Duration) *Interval {
	intervalMU.Lock()
	defer intervalMU.Unlock()

	interval := Interval{
		cb: cb,
		ms: ms,
		id: uintptr(len(intervals)),
	}

	intervals = append(intervals, &interval)

	return &interval
}

// Clear stops the interval and removes it from the queue
func (interval *Interval) Clear() {
	interval.stop = true
}
