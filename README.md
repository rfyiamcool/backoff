# backoff

backoff policy, avoid `thundering herd` effect. the backoff lib is concurrent-safe.

## Install:

```
go get github.com/rfyiamcool/backoff
```

## Desc

**How to calculate the delay time**

```
float64(b.MinDelay) * math.Pow(b.Factor, b.Attempts)
```

## Usage:

**Simple example 1**

factor = 2

```
b := NewBackOff(
    WithMinDelay(100*time.Millisecond),
    WithMaxDelay(10*time.Second),
    WithFactor(2),
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
```

factor = 1.7

```
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
```

**Jitter example**

enabling Jitter adds some randomization to the backoff durations.

```
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
```