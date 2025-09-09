package pkg

import (
	"math"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	tests := []struct {
		attempt int
		min     time.Duration
		max     time.Duration
	}{
		{0, 100 * time.Millisecond, 200 * time.Millisecond},
		{1, 200 * time.Millisecond, 300 * time.Millisecond},
		{2, 400 * time.Millisecond, 500 * time.Millisecond},
		{3, 800 * time.Millisecond, 900 * time.Millisecond},
		{4, 1600 * time.Millisecond, 1700 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				duration := Backoff(tt.attempt)
				if duration < tt.min || duration > tt.max {
					t.Errorf("Backoff(%d) = %v, want between %v and %v",
						tt.attempt, duration, tt.min, tt.max)
				}
			}
		})
	}
}

func TestBackoffExponentialGrowth(t *testing.T) {
	durations := make([]time.Duration, 5)
	for i := 0; i < 5; i++ {
		durations[i] = Backoff(i)
	}

	for i := 1; i < len(durations); i++ {
		if durations[i] <= durations[i-1] {
			t.Errorf("Backoff should grow: %v <= %v", durations[i], durations[i-1])
		}
	}
}

func TestSleepRandom(t *testing.T) {
	start := time.Now()
	SleepRandom()
	elapsed := time.Since(start)

	if elapsed < 100*time.Millisecond || elapsed > 500*time.Millisecond {
		t.Errorf("SleepRandom() took %v, want between 100ms and 500ms", elapsed)
	}
}

func TestShouldFail(t *testing.T) {
	total := 10000
	failures := 0

	for i := 0; i < total; i++ {
		if ShouldFail() {
			failures++
		}
	}

	failureRate := float64(failures) / float64(total) * 100
	expectedRate := 20.0
	tolerance := 2.0

	if math.Abs(failureRate-expectedRate) > tolerance {
		t.Errorf("ShouldFail() failure rate: %.2f%%, want %.2f%% ±%.2f%%",
			failureRate, expectedRate, tolerance)
	}
}
