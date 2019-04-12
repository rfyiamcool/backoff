package backoff

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
	b := NewBackOff(
		WithMinDelay(100*time.Millisecond),
		WithMaxDelay(10*time.Second),
	)

	equals(t, b.Duration(), 100*time.Millisecond)
	equals(t, b.Duration(), 200*time.Millisecond)
	equals(t, b.Duration(), 400*time.Millisecond)
	for index := 0; index < 100; index++ {
		b.Duration()
	}

	// is max
	equals(t, b.Duration(), 10*time.Second)
	b.Reset()
	equals(t, b.Duration(), 100*time.Millisecond)
}

func Test2(t *testing.T) {

	b := NewBackOff(
		WithMinDelay(100*time.Millisecond),
		WithMaxDelay(10*time.Second),
		WithFactor(1.5),
	)

	equals(t, b.Duration(), 100*time.Millisecond)
	equals(t, b.Duration(), 150*time.Millisecond)
	equals(t, b.Duration(), 225*time.Millisecond)
	b.Reset()
	equals(t, b.Duration(), 100*time.Millisecond)
}

func Test3(t *testing.T) {

	b := NewBackOff(
		WithMinDelay(100*time.Millisecond),
		WithMaxDelay(10*time.Second),
		WithFactor(1.7),
	)

	equals(t, b.Duration(), 100*time.Nanosecond)
	equals(t, b.Duration(), 175*time.Nanosecond)
	equals(t, b.Duration(), 306*time.Nanosecond)
	b.Reset()
	equals(t, b.Duration(), 100*time.Nanosecond)
}

func TestJitter(t *testing.T) {
	b := NewBackOff(
		WithMinDelay(100*time.Millisecond),
		WithMaxDelay(10*time.Second),
		WithFactor(2),
		WithJitterFlag(true),
	)

	equals(t, b.Duration(), 100*time.Millisecond)
	between(t, b.Duration(), 100*time.Millisecond, 200*time.Millisecond)
	between(t, b.Duration(), 100*time.Millisecond, 400*time.Millisecond)
	b.Reset()
	equals(t, b.Duration(), 100*time.Millisecond)
}

func between(t *testing.T, actual, low, high time.Duration) {
	if actual < low {
		t.Fatalf("Got %s, Expecting >= %s", actual, low)
	}
	if actual > high {
		t.Fatalf("Got %s, Expecting <= %s", actual, high)
	}
}

func equals(t *testing.T, d1, d2 time.Duration) {
	if d1 != d2 {
		t.Fatalf("Got %s, Expecting %s", d1, d2)
	}
}
