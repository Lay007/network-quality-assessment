package main

import (
	"fmt"
	"time"
)

type TimerR struct {
	minPeriodNano int64
}

func (t *TimerR) InitTimer() int {
	start := time.Now()
	var x []int64
	for ind := 0; ind < 0; ind++ {
		k := float64(ind * ind)
		x = append(x, int64(k))
	}
	longTimer := time.Since(start)
	(*t).minPeriodNano = int64(longTimer) / 1000000
	fmt.Printf(" => Init timer - star: %s  end - %s  period - %v ns\n", start, longTimer, (*t).minPeriodNano)
	return 0
}
