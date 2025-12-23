package interfaces

// DelayFn returns a delay in seconds given the current retry count.
type DelayFn func(currenRetries int64) (delay int64)

var (
	// ExponentialBackoffDelayFn returns 2^(retries-1) seconds.
	ExponentialBackoffDelayFn DelayFn = func(currenRetries int64) int64 {
		return 2 << (currenRetries - 1)
	}

	// LinearDelayFn returns retries seconds (1s, 2s, 3s, …).
	LinearDelayFn DelayFn = func(currenRetries int64) int64 {
		return currenRetries
	}

	NoDelayFn DelayFn = func(_ int64) int64 { return 0 }

	DefaultDelayFn DelayFn = LinearDelayFn
)
