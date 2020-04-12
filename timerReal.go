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
	countTimer := 1000
	var x []int64
	for ind := 0; ind < countTimer; ind++ {
		k := float64(ind * ind)
		x = append(x, int64(k))
	}
	longTimer := time.Since(start)
	(*t).minPeriodNano = int64(longTimer) / int64(countTimer)
	fmt.Printf(" => Init timer - star: %s  end - %s  period - %v ns\n", start, longTimer, (*t).minPeriodNano)
	return 0
}
