package main

import (
	"fmt"
	"time"
)

type TimerR struct {
	minPeriodNano int64
}

func (t *TimerR) timerDelayNano(delay int64) {
	countTimer := (delay * 18) / t.minPeriodNano
	//	countTimer := (delay * 10) / t.minPeriodNano
	var x int64
	var ind int64
	for ind = 1; ind < countTimer; ind++ {
		k := float64(ind*ind) / float64(ind)
		x = x * int64(k) * (x + 123)
	}

}

func (t *TimerR) InitTimer() int {
	(*t).minPeriodNano = 1
	countTimer := 1000000000
	start := time.Now()
	t.timerDelayNano(int64(countTimer))
	longTimer := time.Since(start)
	(*t).minPeriodNano = int64(longTimer) / int64(countTimer)
	fmt.Printf(" => Init timer - star: %s  end - %s  period - %v ns\n", start, longTimer, (*t).minPeriodNano)
	return 0
}

/*
const loops = 500
const workAmt = 250 // 250us
var workAdj int

func (t *TimerR) InitTimer() int {
	doCalibrateWork()
	doChecks()
	doTests()
	return 0
}

// Calibrate doWork() to be about 1us of work by adjusting workAdj.
//
func doCalibrateWork() {
	workAdj = 100 // Initial guess
	now := time.Now()
	doWork(workAmt * loops)
	dur := time.Now().Sub(now) + (time.Microsecond * 1)
	workAdj = (workAdj * workAmt * int(time.Microsecond)) / (int(dur) / loops)
}

// Run the checks.
//
func doChecks() {
	fmt.Println()
	fmt.Println("Checks")
	fmt.Println()

	// Check time.Now resolution.
	now := time.Now()
	then := now
	for now == then {
		then = time.Now()
	}
	dur := then.Sub(now).Seconds()
	fmt.Printf("  time.Now resolution:  %.3fms\n", dur*1e3)

	// Check doWork(~workAmt).
	now = time.Now()
	doWork(workAmt * loops)
	dur = time.Now().Sub(now).Seconds() / loops
	fmt.Printf("  Work(~%dus):         %.0fus\n", workAmt, dur*1e6)
}

// Do the tests.
// Check time.Sleep and time.After, at various intervals (1, 0.5ms, 1ms, 1.5ms, 2ms),
// both with and without doWork.
//
func doTests() {
	fmt.Println()
	fmt.Println("Tests")

	fmt.Println("                  0          1        0.5ms       1ms       1.5ms       2ms")
	fmt.Println("               -------    -------    -------    -------    -------    -------")
	for _, test := range []string{"Sleep", "After"} {
		for _, work := range []bool{false, true} {
			if !work {
				fmt.Printf("  time.%s:  ", test)
			} else {
				fmt.Print("  ... + work:  ")
			}
			for _, i := range []time.Duration{0, 1, 500 * time.Microsecond, 1 * time.Millisecond,
				 1500 * time.Microsecond, 2 * time.Millisecond} {
				now := time.Now()
				for j := 0; j < loops; j++ {
					if work {
						_ = doWork(workAmt)
					}
					if test == "Sleep" {
						// This is the time.Sleep test.
						time.Sleep(i)
					} else {
						// This is the time.After test.
						ch := time.After(i)
						<-ch
					}
				}
				dur := time.Now().Sub(now).Seconds() / loops
				fmt.Printf("%.3fms    ", dur*1e3)
			}
			fmt.Println("")
		}
	}
}

// Work is calibrated (by adjusting workAdj) to be about 1 us of work.
//
func doWork(us int) (junk int) {
	var x int
	for i := 0; i < us*workAdj; i++ {
		x = (i + 1234) / (i + 12)
		x = x + x + x
	}
	return x
}
*/
