package backoff

import (
	"context"
	"math"
	"math/rand"
	"time"
)

var (
	defaultFactor   float64 = 2
	defaultJitter           = false
	defaultMinDelay         = 100 * time.Millisecond
	defaultMaxDelay         = 2 * time.Second
)

type Backoff struct {
	Attempts float64
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
			Attempts: 0,
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
	dur := float64(b.MinDelay) * math.Pow(b.Factor, b.Attempts)
	if b.Jitter == true {
		dur = rand.Float64()*(dur-float64(b.MinDelay)) + float64(b.MinDelay)
	}
	if dur > float64(b.MaxDelay) {
		return b.MaxDelay
	}

	b.Attempts++
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
	b.Attempts = 0
}
