package backoff

import (
	"context"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

var (
	defaultFactor   float64 = 2
	defaultJitter           = false
	defaultMinDelay         = 100 * time.Millisecond
	defaultMaxDelay         = 2 * time.Second
)

type Backoff struct {
	attempts uint64
	Factor   float64

	//Jitter eases contention by randomizing backoff steps
	Jitter bool

	// Min and Max are the minimum and maximum values of the backoff control
	MinDelay time.Duration
	MaxDelay time.Duration
}

type BackoffOption func(*Backoff)

func WithMinDelay(d time.Duration) BackoffOption {
	return func(b *Backoff) {
		b.MinDelay = d
	}
}

func WithMaxDelay(d time.Duration) BackoffOption {
	return func(b *Backoff) {
		b.MaxDelay = d
	}
}

func WithJitterFlag(f bool) BackoffOption {
	return func(b *Backoff) {
		b.Jitter = f
	}
}

func WithFactor(v float64) BackoffOption {
	return func(b *Backoff) {
		b.Factor = v
	}
}

func NewBackOff(opts ...BackoffOption) *Backoff {
	var (
		bo = &Backoff{
			attempts: 0,
			Factor:   defaultFactor,
			Jitter:   defaultJitter,
			MinDelay: defaultMinDelay,
			MaxDelay: defaultMaxDelay,
		}
	)

	for _, option := range opts {
		option(bo)
	}

	return bo
}

func (b *Backoff) beRevise() {
	if b.MinDelay == 0 {
		b.MinDelay = defaultMinDelay
	}
	if b.MaxDelay == 0 {
		b.MaxDelay = defaultMaxDelay
	}
	if b.Factor == 0 {
		b.Factor = defaultFactor
	}
}

func (b *Backoff) Duration() time.Duration {
	dur := float64(b.MinDelay) * math.Pow(b.Factor, float64(b.attempts))
	if b.Jitter == true {
		dur = rand.Float64()*(dur-float64(b.MinDelay)) + float64(b.MinDelay)
	}
	if dur > float64(b.MaxDelay) {
		return b.MaxDelay
	}

	atomic.AddUint64(&b.attempts, 1)
	return time.Duration(dur)
}

// Sleep
func (b *Backoff) Sleep() {
	time.Sleep(b.Duration())
}

// SleepCtx
func (b *Backoff) SleepCtx(ctx context.Context) {
	var timer = time.NewTimer(b.Duration())

	select {
	case <-timer.C:
		return

	case <-ctx.Done():
		timer.Stop()
		return
	}
}

//Resets the current value of the counter back to Min
func (b *Backoff) Reset() {
	atomic.StoreUint64(&b.attempts, 0)
}

// Attempts
func (b *Backoff) Attempts() uint64 {
	return atomic.LoadUint64(&b.attempts)
}

// AttemptsInt
func (b *Backoff) AttemptsInt() int {
	return int(atomic.LoadUint64(&b.attempts))
}
