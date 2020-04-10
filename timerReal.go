package serverSLA
/*
import (
	"fmt"
	"time"
)

type TimerR struct {
	minPeriodNano int64
}

func (t *TimerR) InitTimer() {
	fmt.Println(" => InitHi <=")
	start := time.Now()
	var x []int64
	for ind := 0; ind < 100000; ind++ {
		k := float64(ind * ind)
		x = append(x, int64(k))
	}
	longTimer := time.Since(start)
	(*t).minPeriodNano = int64(longTimer) / 100000
	fmt.Printf(" => Init timer - star: %s  end - %s  period - %v ns", start, longTimer, (*t).minPeriodNano)
}
*/