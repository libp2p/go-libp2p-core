package metrics

import (
	"fmt"
	"math"
	"time"
)

func ExampleMeter() {
	meter := new(Meter)
	t := time.NewTicker(100 * time.Millisecond)
	for i := 0; i < 100; i++ {
		<-t.C
		meter.Mark(30)
	}

	// Get the current rate. This will be accurate *now* but not after we
	// sleep (because we calculate it using EWMA).
	rate := meter.Snapshot().Rate

	// Sleep 2 seconds to allow the total to catch up. We snapshot every
	// second so the total may not yet be accurate.
	time.Sleep(2 * time.Second)

	// Get the current total.
	total := meter.Snapshot().Total

	fmt.Printf("%d (%d/s)\n", total, roundTens(rate))
	// Output: 3000 (300/s)
}

func roundTens(x float64) int64 {
	return int64(math.Floor(x/10+0.5)) * 10
}
