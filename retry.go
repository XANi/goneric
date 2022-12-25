package goneric

import "time"

// Retry retries function that returns error up to n times
func Retry[T any](n int, f func() (T, error)) (T, error) {
	for i := 1; i < n; i++ {
		out, err := f()
		if err == nil {
			return out, err
		}
	}
	return f()
}

// RetryAfter retries function till it returns without error, first after min_interval,
// then with increasing intervals up to max_interval with last retry happening near total_timeout
// Intended use is to be able to say "retry for 10 minutes, at the very least every minute, but not shorter than 10 seconds to account for TCP retransmissions"
func RetryAfter[T any](
	min_interval,
	max_interval,
	total_timeout time.Duration,
	f func() (T, error),
) (T, error) {
	finish := time.Now().Add(total_timeout)
	interval := min_interval
	var commandDuration time.Duration
	for {
		start := time.Now()
		out, err := f()
		if err == nil {
			return out, err
		}
		commandDuration = time.Now().Sub(start)
		time.Sleep(interval)
		interval = Min(max_interval, (min_interval * 3 / 2))
		if start.Add(interval).Add(commandDuration).After(finish) {
			// if we're near finish we will just try to run command at last second
			ttf := finish.Sub(time.Now())
			time.Sleep(ttf)
			break
		}

	}
	return f()
}

// Try every function till one returns without error
func Try[T any](f ...func() (T, error)) (out T, err error) {
	for _, ff := range f {
		out, err = ff()
		if err == nil {
			return
		}
	}
	return
}
