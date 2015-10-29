package main

import (
	"math"
	"testing"
	"time"

	a "github.com/bmizerany/assert"
)

const testWindowSize = time.Second * 600
const testMaxRestarts = 5

func TestWindow(t *testing.T) {
	ts := time.Now()
	w := makeWindow(ts, 5, 5)
	a.Equal(t, len(window(w, ts, testWindowSize)), 3)
}

func TestLimit(t *testing.T) {
	ts := time.Now()
	w := makeWindow(ts, 5, 2)
	m := &Metadata{
		Restarts: w,
	}
	a.Equal(t, limit(m, testWindowSize, testMaxRestarts), true)
	a.Equal(t, len(m.Restarts), len(w))
	w = makeWindow(ts, 5, 5)
	m.Restarts = w
	a.Equal(t, limit(m, testWindowSize, testMaxRestarts), false)
	a.Equal(t, len(m.Restarts), 4)
}

func makeWindow(ts time.Time, size int, exp int) []time.Time {
	var w []time.Time
	for i := 1; i <= size; i++ {
		w = append(w, time.Unix(ts.Unix()-int64(math.Pow(float64(exp), float64(i))), 0))
	}
	return w
}
