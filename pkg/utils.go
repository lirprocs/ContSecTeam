package pkg

import (
	"math/rand"
	"time"
)

func SleepRandom() {
	d := time.Duration(100+rand.Intn(400)) * time.Millisecond
	time.Sleep(d)
}

func ShouldFail() bool {
	return rand.Intn(100) < 20
}

func Backoff(attempt int) time.Duration {
	base := time.Duration(1<<attempt) * 100 * time.Millisecond
	jitter := time.Duration(rand.Intn(100)) * time.Millisecond
	return base + jitter
}
